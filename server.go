package main

import (
	"net/http"
	"path/filepath"
)

type Server struct {
	Config      *Config
	AssetServer http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "", "/":
		s.ServeUpload(w, r)
	default:
		s.AssetServer.ServeHTTP(w, r)
	}
}

func (s *Server) ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		s.serveUploadPost(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join(s.Config.AssetDir, "upload.html"))
}

func (s *Server) serveUploadPost(w http.ResponseWriter, r *http.Request) {
	// TODO: allow the user to upload content.
}
