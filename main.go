package main

import (
	"fmt"
	"net/http"
	"os"
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
	assetServer := http.FileServer(http.Dir(config.AssetDir))
	server := &Server{
		Config:      config,
		AssetServer: assetServer,
		AssetDir:    config.AssetDir,
		Database:    database,
	}
	http.ListenAndServe(":"+os.Args[2], server)
}
