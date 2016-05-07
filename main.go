package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/unixpickle/ratelimit"
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
		Config:   config,
		AssetFS:  http.Dir(config.AssetDir),
		AssetDir: config.AssetDir,
		Database: database,
		SessionStore: sessions.NewCookieStore(securecookie.GenerateRandomKey(16),
			securecookie.GenerateRandomKey(16)),
		HostNamer:   &ratelimit.HTTPRemoteNamer{},
		RateLimiter: ratelimit.NewTimeSliceLimiter(time.Minute*10, 20),
	}
	if err := http.ListenAndServe(":"+os.Args[2], server); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
