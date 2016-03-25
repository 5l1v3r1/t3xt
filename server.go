package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	BadRequestFilename    = "bad_request.html"
	InternalErrorFilename = "internal_error.html"
	NotFoundFilename      = "not_found.html"
	ViewFilename          = "view.html"
	ListFilename          = "list.html"
)

var viewPathRegexp = regexp.MustCompile("^/view/([a-f0-9]*)$")

const listingResultCount = 15

type Server struct {
	Config   *Config
	AssetFS  http.FileSystem
	AssetDir string
	Database *Database
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "", "/":
		s.serveUpload(w, r)
	case "/list":
		s.serveList(w, r)
	default:
		if match := viewPathRegexp.FindStringSubmatch(r.URL.Path); match != nil {
			s.serveView(w, r, match[1])
		} else {
			f, err := s.AssetFS.Open(r.URL.Path)
			if err == nil {
				defer f.Close()
				http.ServeContent(w, r, path.Base(r.URL.Path), time.Now(), f)
			} else {
				s.serveError(w, r, http.StatusNotFound, NotFoundFilename)
			}
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

func (s *Server) serveList(w http.ResponseWriter, r *http.Request) {
	entries, err := s.listEntriesForRequest(r)
	if err != nil {
		s.serveError(w, r, http.StatusBadRequest, BadRequestFilename)
		return
	} else if len(entries) == 0 {
		if r.FormValue("before") != "" {
			http.Redirect(w, r, "/list?after=0", http.StatusTemporaryRedirect)
		} else if r.FormValue("after") != "" {
			http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		return
	}

	listing := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		head, _ := s.Database.Head(entry.ID)
		listing[i] = map[string]interface{}{
			"id":       entry.ID,
			"secretId": entry.ShareID,
			"head":     head,
			"lines":    entry.LineCount,
			"postTime": entry.PostDate.Unix(),
		}
	}
	next, last := s.availableListDirections(entries)
	fullData := map[string]interface{}{
		"posts":   listing,
		"hasNext": next,
		"hasLast": last,
	}
	encodedData, err := json.Marshal(fullData)
	if err != nil {
		s.serveError(w, r, http.StatusInternalServerError, InternalErrorFilename)
		return
	}
	s.injectAndServe(w, r, string(encodedData), ListFilename)
}

func (s *Server) listEntriesForRequest(r *http.Request) ([]DatabaseEntry, error) {
	if beforeID := r.FormValue("before"); beforeID != "" {
		if idNum, err := strconv.Atoi(beforeID); err != nil {
			return nil, err
		} else if idNum < 0 {
			return nil, errors.New("bad index")
		} else {
			return s.Database.EntriesBefore(idNum, listingResultCount), nil
		}
	} else if afterID := r.FormValue("after"); afterID != "" {
		if idNum, err := strconv.Atoi(afterID); err != nil {
			return nil, err
		} else if idNum < 0 {
			return nil, errors.New("bad index")
		} else {
			return s.Database.EntriesAfter(idNum, listingResultCount), nil
		}
	} else {
		return s.Database.LatestEntries(listingResultCount), nil
	}
}

func (s *Server) availableListDirections(l []DatabaseEntry) (next, last bool) {
	lastID := l[0].ID
	firstID := l[len(l)-1].ID
	if firstID > 0 {
		last = len(s.Database.EntriesBefore(firstID-1, 1)) == 1
	}
	next = len(s.Database.EntriesAfter(lastID+1, 1)) == 1
	return
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
