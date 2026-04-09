package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"new-lb/store"
)

// Start launches the HTTP server on localhost:8080.
// webPath is the path to the web/index.html file.
func Start(s *store.Store, webPath string, dataDir string) error {
	h := NewHandler(s, dataDir)
	mux := http.NewServeMux()

	// API routes - handle /api/entries and /api/entries/{name}
	mux.HandleFunc("/api/entries/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/entries/")
		if name == "" {
			writeError(w, http.StatusBadRequest, "missing entry name")
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.HandleGetEntry(w, r, name)
		case http.MethodPatch:
			h.HandlePatchEntry(w, r, name)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.HandleGetEntries(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/collect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.HandleCollect(w, r)
	})

	mux.HandleFunc("/api/classify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.HandleClassify(w, r)
	})

	mux.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.HandleReport(w, r)
	})

	// Merge groups route
	mux.HandleFunc("/api/merge-groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.HandleGetMergeGroups(w, r)
	})

	// Snapshot routes - /api/snapshots/{name}/...
	mux.HandleFunc("/api/snapshots/", func(w http.ResponseWriter, r *http.Request) {
		// Extract name from path: /api/snapshots/{name} or /api/snapshots/{name}/load
		rest := strings.TrimPrefix(r.URL.Path, "/api/snapshots/")
		if rest == "" {
			writeError(w, http.StatusBadRequest, "missing snapshot name")
			return
		}
		// Check for /load suffix
		if name, ok := strings.CutSuffix(rest, "/load"); ok {
			if r.Method != http.MethodPost {
				writeError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			h.HandleLoadSnapshot(w, r, name)
			return
		}
		// Otherwise it's /api/snapshots/{name}
		name := rest
		switch r.Method {
		case http.MethodDelete:
			h.HandleDeleteSnapshot(w, r, name)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// Snapshot routes - /api/snapshots (list/create)
	mux.HandleFunc("/api/snapshots", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.HandleListSnapshots(w, r)
		case http.MethodPost:
			h.HandleSaveSnapshot(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// Serve web UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := os.ReadFile(webPath)
		if err != nil {
			log.Printf("[server] failed to read %s: %v", webPath, err)
			writeError(w, http.StatusInternalServerError, "failed to read index.html")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	addr := "127.0.0.1:8080"
	log.Printf("[server] starting on http://%s", addr)
	return http.ListenAndServe(addr, mux)
}
