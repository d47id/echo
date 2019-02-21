package lifecycle

import (
	"encoding/json"
	"net/http"
)

var (
	// Version is the semantic version assigned to the build
	Version = "unset"
	// BuildTime is the date and time that the build was created
	BuildTime = "unset"
	// Branch is the branch from which the build was created
	Branch = "unset"
	// Commit is the hash of the current commit at the time the build was created
	Commit = "unset"
)

// info is a serializable wrapper for the version info set at compile-time
type info struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
	Branch    string `json:"branch"`
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

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getInfo())
}
