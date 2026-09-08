package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

//go:embed web
var webFS embed.FS

func serveWeb(addr string) {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Game of Life  |  open http://%s\n", addr)
	if err := http.ListenAndServe(addr, http.FileServer(http.FS(sub))); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
