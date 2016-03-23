package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	BadRequestFilename    = "bad_request.html"
	InternalErrorFilename = "internal_error.html"
	NotFoundFilename      = "not_found.html"
	ViewFilename          = "view.html"
)

var viewPathRegexp = regexp.MustCompile("^/view/([a-f0-9]*)$")

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
		if match := viewPathRegexp.FindStringSubmatch(r.URL.Path); match != nil {
			s.serveView(w, r, match[1])
		} else {
			s.AssetServer.ServeHTTP(w, r)
		}
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
	info.PosterIP = ipAddressFromRequest(r)
	bodyReader := bytes.NewBufferString(r.PostFormValue("code"))
	if info, err := s.Database.CreateEntry(info, bodyReader); err != nil {
		s.serveError(w, r, http.StatusInternalServerError, InternalErrorFilename)
	} else {
		s.redirectPost(w, r, info)
	}
}

func (s *Server) redirectPost(w http.ResponseWriter, r *http.Request, info DatabaseEntry) {
	http.Redirect(w, r, "/view/"+info.ShareID, http.StatusTemporaryRedirect)
}

func (s *Server) serveView(w http.ResponseWriter, r *http.Request, shareID string) {
	entry, reader, err := s.Database.OpenEntry(shareID)
	if err != nil {
		s.serveError(w, r, http.StatusNotFound, NotFoundFilename)
		return
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		s.serveError(w, r, http.StatusInternalServerError, InternalErrorFilename)
		return
	}
	postData := map[string]interface{}{
		"content":  string(data),
		"postTime": entry.PostDate.Unix(),
		"language": entry.Language,
	}
	encodedPostData, err := json.Marshal(postData)
	if err != nil {
		s.serveError(w, r, http.StatusInternalServerError, InternalErrorFilename)
		return
	}
	s.injectAndServe(w, r, string(encodedPostData), ViewFilename)
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

// injectAndServe serves an HTML page, replacing a pre-determined part of JavaScript with the data
// passed to the data argument.
func (s *Server) injectAndServe(w http.ResponseWriter, r *http.Request, data, pageFilename string) {
	contents, err := ioutil.ReadFile(filepath.Join(s.AssetDir, pageFilename))
	if err != nil {
		s.serveError(w, r, http.StatusInternalServerError, InternalErrorFilename)
		return
	}
	contentStr := string(contents)
	contentStr = strings.Replace(contentStr, "/* SCRIPT INJECT */{}", data, 1)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(contentStr))
}

func ipAddressFromRequest(r *http.Request) string {
	if forwardHeader := r.Header.Get("X-Forwarded-For"); forwardHeader != "" {
		return strings.Split(forwardHeader, ", ")[0]
	}

	// r.RemoteAddr is either "IPv4Address:port" or "[IPv6Address]:port".
	if strings.HasPrefix(r.RemoteAddr, "[") {
		return strings.Split(r.RemoteAddr, "]")[0][1:]
	} else {
		return strings.Split(r.RemoteAddr, ":")[0]
	}
}
