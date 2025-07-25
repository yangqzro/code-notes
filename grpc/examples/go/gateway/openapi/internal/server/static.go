package server

import (
	"net/http"
	"os"
)

func StaticServer(root string) http.Handler {
	return http.FileServerFS(os.DirFS(root))
}
