package interceptor

import (
	"bytes"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/result"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const grpcPrefix = "Grpc-Metadata-"
const authKey = "Authorization"

type AuthorityMiddleware struct {
	svcCtx *svc.ServiceContext
	parser *jwt.Parser
}

func NewAuthorityMiddleware() *AuthorityMiddleware {
	return &AuthorityMiddleware{
		parser: jwt.NewParser(),
	}
}

type logResponseWriter struct {
	writer http.ResponseWriter
	code   int
	// 存响应数据
	buf *bytes.Buffer

	re *regexp.Regexp
}

func newLogResponseWriter(writer http.ResponseWriter, code int) *logResponseWriter {
	var buf bytes.Buffer
	return &logResponseWriter{
		writer: writer,
		code:   code,
		buf:    &buf,
		re:     regexp.MustCompile(`code\s*=\s*Code\((\d+)\)\s*desc\s*=\s*(.*)`),
	}
}

func (m *AuthorityMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.ignore(r.URL.Path) {
			tokenStr := r.Header.Get(authKey)
			claims, err := m.parseJwt(tokenStr)
			if err != nil {
				result.HttpResult(r, w, nil, err)
				return
			}
			for k, v := range claims {
				r.Header.Set(grpcPrefix+k, fmt.Sprintf("%v", v))
			}
		}

		//var dup io.ReadCloser
		//dup, r.Body, _ = drainBody(r.Body)
		lwr := newLogResponseWriter(w, http.StatusOK)
		next(lwr, r)
		//r.Body = dup
	}
}

func (m *AuthorityMiddleware) ignore(url string) bool {
	return strings.HasPrefix(url, "/p/")
}

func (m *AuthorityMiddleware) parseJwt(tokenStr string) (jwt.MapClaims, error) {
	token, err := m.parser.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.svcCtx.Config.JwtAuth.AccessSecret), nil
	})
	if err != nil {
		return nil, errorx.New(errorx.Token, "Token invalid, please log in again")
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims.VerifyExpiresAt(m.svcCtx.Config.JwtAuth.AccessExpire, false) {
		return nil, errorx.New(errorx.Token, "Token expire, please log in again")
	}
	return claims, nil
}

func (w *logResponseWriter) Write(bs []byte) (int, error) {
	// gateway error convent to json
	if w.code == http.StatusInternalServerError {
		matches := w.re.FindStringSubmatch(string(bs))
		if len(matches) >= 3 {
			c := matches[1]
			errMsg := matches[2]
			errCode, _ := strconv.ParseUint(c, 10, 32)
			toJson := errorx.New(uint32(errCode), errMsg).ToJson()
			w.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			bs = []byte(toJson)
		}
	}
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *logResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.code = code
	if code == http.StatusInternalServerError {
		code = http.StatusOK
	}
	w.writer.WriteHeader(code)
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
