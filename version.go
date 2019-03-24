package main

import (
	"encoding/json"
	"net/http"
)

var (
	// Version should be overridden with the application's git tag or other
	// version identifier, if applicable
	Version string
	// BuildTime should be overridden with the date and time that the
	// application is built
	BuildTime = "unset"
	// Branch should be overridden with the name of the branch from which the
	// application is being built, if applicable
	Branch string
	// Commit should be overridden with the commit sha or other suitable
	// identifier
	Commit = "unset"
)

// info is a serializable wrapper for the version info set at compile-time
type info struct {
	Version   string `json:"version,omitempty"`
	BuildTime string `json:"buildTime"`
	Branch    string `json:"branch,omitempty"`
	Commit    string `json:"commit"`
}

// getInfo returns a populated VersionInfo
func getInfo() info {
	return info{
		Version:   Version,
		BuildTime: BuildTime,
		Branch:    Branch,
		Commit:    Commit,
	}
}

// versionHandler provides version info via http
func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getInfo())
}
