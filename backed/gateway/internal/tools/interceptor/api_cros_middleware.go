package interceptor

import (
	"net/http"
)

func CrosMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")         // 指明哪些请求源被允许访问资源，值可以为 "*"，"null"，或者单个源地址。
		w.Header().Set("Access-Control-Allow-Headers", "*")        //对于预请求来说，指明了哪些头信息可以用于实际的请求中。
		w.Header().Set("Access-Control-Allow-Methods", "*")        //对于预请求来说，哪些请求方式可以用于实际的请求。
		w.Header().Set("Access-Control-Expose-Headers", "*")       //对于预请求来说，指明哪些头信息可以安全的暴露给 CORS API 规范的 API
		w.Header().Set("Access-Control-Allow-Credentials", "true") //指明当请求中省略 creadentials 标识时响应是否暴露。对于预请求来说，它表明实际的请求中可以包含用户凭证。

		//放行所有OPTIONS方法
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		next(w, r)
	}
}
