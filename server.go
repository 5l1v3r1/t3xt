package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	BadRequestFilename    = "bad_request.html"
	InternalErrorFilename = "internal_error.html"
)

type Server struct {
	Config      *Config
	AssetServer http.Handler
	AssetDir    string
	Database    *Database
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "", "/":
		s.serveUpload(w, r)
	default:
		s.AssetServer.ServeHTTP(w, r)
	}
}

func (s *Server) serveUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		s.serveUploadPost(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join(s.Config.AssetDir, "upload.html"))
}

func (s *Server) serveUploadPost(w http.ResponseWriter, r *http.Request) {
	var info DatabaseEntry
	r.ParseForm()

	requiredFields := []string{"language", "code"}
	for _, field := range requiredFields {
		if res, ok := r.PostForm[field]; !ok || len(res) != 1 {
			s.serveError(w, r, http.StatusBadRequest, BadRequestFilename)
			return
		}
	}

	info.Language = r.PostFormValue("language")
	info.PostDate = time.Now()
	bodyReader := bytes.NewBufferString(r.PostFormValue("code"))
	if info, err := s.Database.CreateEntry(info, bodyReader); err != nil {
		s.serveError(w, r, http.StatusInternalServerError, InternalErrorFilename)
	} else {
		s.redirectPost(w, r, info)
	}
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, code int, file string) {
	reader, err := os.Open(filepath.Join(s.AssetDir, file))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	io.Copy(w, reader)
}

func (s *Server) redirectPost(w http.ResponseWriter, r *http.Request, info DatabaseEntry) {
	http.Redirect(w, r, info.ShareID, http.StatusTemporaryRedirect)
}
