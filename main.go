package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/unixpickle/ratelimit"
	"github.com/unixpickle/whichlang"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config.json> <port>\n", os.Args[0])
		os.Exit(1)
	}
	config, err := GetConfig(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load config:", err)
		os.Exit(1)
	}
	database, err := OpenDatabase(config.DbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open database:", err)
		os.Exit(1)
	}

	server := &Server{
		Config:     config,
		AssetFS:    http.Dir(config.AssetDir),
		AssetDir:   config.AssetDir,
		Database:   database,
		Classifier: readClassifier(config),
		SessionStore: sessions.NewCookieStore(securecookie.GenerateRandomKey(16),
			securecookie.GenerateRandomKey(16)),
		HostNamer:   &ratelimit.HTTPRemoteNamer{},
		RateLimiter: ratelimit.NewTimeSliceLimiter(time.Minute*10, 20),
	}
	handler := context.ClearHandler(server)
	if err := http.ListenAndServe(":"+os.Args[2], handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readClassifier(config *Config) whichlang.Classifier {
	if config.ClassifierType == "" {
		return nil
	}

	decoder, ok := whichlang.Decoders[config.ClassifierType]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown classifier type:", config.ClassifierType)
		os.Exit(1)
	}
	classifierData, err := ioutil.ReadFile(config.ClassifierPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load classifier:", err)
		os.Exit(1)
	}
	classifier, err := decoder(classifierData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load classifier:", err)
		os.Exit(1)
	}
	return classifier
}
