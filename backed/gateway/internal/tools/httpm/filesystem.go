package httpm

import (
	"embed"
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"io"
	"net/http"
)

func FileSystem(fs *embed.FS, root string, routes []rest.Route) []rest.Route {
	if routes == nil {
		routes = make([]rest.Route, 0)
	}
	readDir, err := fs.ReadDir(root)
	if err != nil {
		return routes
	}
	for _, entry := range readDir {
		fileName := fmt.Sprintf("%s/%s", root, entry.Name())
		if entry.IsDir() {
			routes = FileSystem(fs, fileName, routes)
		} else {
			routes = append(routes, rest.Route{
				Method: http.MethodGet,
				Path:   "/" + fileName,
				Handler: func(writer http.ResponseWriter, request *http.Request) {
					file, err := fs.Open(fileName)
					if err != nil {
						writer.WriteHeader(http.StatusBadRequest)
						return
					}
					all, err := io.ReadAll(file)
					if err != nil {
						writer.WriteHeader(http.StatusBadRequest)
						return
					}
					_, _ = writer.Write(all)
				},
			})
		}
	}
	return routes
}
