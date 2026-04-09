package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"new-lb/classifier"
	"new-lb/collector"
	"new-lb/model"
	"new-lb/store"
)

// Handler holds dependencies for API handlers.
type Handler struct {
	store   *store.Store
	dataDir string
}

// NewHandler creates a new Handler.
func NewHandler(s *store.Store, dataDir string) *Handler {
	return &Handler{store: s, dataDir: dataDir}
}

// HandleGetEntries handles GET /api/entries with optional query filters.
func (h *Handler) HandleGetEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Apply filters
	solution := r.URL.Query().Get("solution")
	status := r.URL.Query().Get("status")
	action := r.URL.Query().Get("action")
	env := r.URL.Query().Get("environment")
	search := r.URL.Query().Get("search")

	var filtered []model.Entry
	for _, e := range entries {
		if solution != "" && e.Solution != solution {
			continue
		}
		if status != "" && e.Status != status {
			continue
		}
		if action != "" && e.Action != action {
			continue
		}
		if env != "" && e.Environment != env {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(e.ALBName), strings.ToLower(search)) {
			continue
		}
		filtered = append(filtered, e)
	}

	writeJSON(w, http.StatusOK, filtered)
}

// HandleGetEntry handles GET /api/entries/{name}.
func (h *Handler) HandleGetEntry(w http.ResponseWriter, r *http.Request, name string) {
	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, e := range entries {
		if e.ALBName == name {
			writeJSON(w, http.StatusOK, e)
			return
		}
	}
	writeError(w, http.StatusNotFound, "entry not found: "+name)
}

// HandlePatchEntry handles PATCH /api/entries/{name}.
func (h *Handler) HandlePatchEntry(w http.ResponseWriter, r *http.Request, name string) {
	var updates map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	entry, err := h.store.UpdateEntry(name, updates)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, entry)
}

// HandleCollect handles POST /api/collect.
func (h *Handler) HandleCollect(w http.ResponseWriter, r *http.Request) {
	log.Println("[api] starting collection...")

	// Load existing entries to preserve user-modified fields
	oldEntries, err := h.store.Load()
	if err != nil {
		log.Printf("[api] warning: could not load existing entries: %v", err)
		oldEntries = nil
	}

	entries, err := collector.Run(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "collection failed: "+err.Error())
		return
	}

	// Merge user-modified fields from old entries
	if len(oldEntries) > 0 {
		entries = store.MergeEntries(oldEntries, entries)
	}

	if err := h.store.Save(entries); err != nil {
		writeError(w, http.StatusInternalServerError, "save failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "completed",
		"count":  len(entries),
	})
}

// HandleClassify handles POST /api/classify.
func (h *Handler) HandleClassify(w http.ResponseWriter, r *http.Request) {
	log.Println("[api] starting classification...")

	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load failed: "+err.Error())
		return
	}

	if len(entries) == 0 {
		writeError(w, http.StatusBadRequest, "no entries found, run collect first")
		return
	}

	classified, unknown := classifier.Classify(entries)

	if err := h.store.Save(entries); err != nil {
		writeError(w, http.StatusInternalServerError, "save failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"classified": classified,
		"unknown":    unknown,
	})
}

// HandleReport handles GET /api/report.
func (h *Handler) HandleReport(w http.ResponseWriter, r *http.Request) {
	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	report := model.Report{
		Total:         len(entries),
		BySolution:    make(map[string]int),
		ByStatus:      make(map[string]int),
		ByAction:      make(map[string]int),
		ByEnvironment: make(map[string]int),
	}

	for _, e := range entries {
		report.BySolution[e.Solution]++
		report.ByStatus[e.Status]++
		report.ByAction[e.Action]++
		report.ByEnvironment[e.Environment]++

		if e.Action == "삭제" {
			report.DeleteCandidates = append(report.DeleteCandidates, e.ALBName)
		}
	}

	if report.DeleteCandidates == nil {
		report.DeleteCandidates = []string{}
	}
	if report.MergeCandidates == nil {
		report.MergeCandidates = []model.MergePlan{}
	}

	writeJSON(w, http.StatusOK, report)
}

// HandleSaveSnapshot handles POST /api/snapshots — saves current classifications as a snapshot.
func (h *Handler) HandleSaveSnapshot(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load entries failed: "+err.Error())
		return
	}
	if len(entries) == 0 {
		writeError(w, http.StatusBadRequest, "no entries to save, run collect first")
		return
	}

	if err := store.SaveSnapshot(h.dataDir, entries, body.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "saved",
		"name":   body.Name,
		"count":  len(entries),
	})
}

// HandleListSnapshots handles GET /api/snapshots — returns snapshot metadata list.
func (h *Handler) HandleListSnapshots(w http.ResponseWriter, r *http.Request) {
	metas, err := store.ListSnapshots(h.dataDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, metas)
}

// HandleLoadSnapshot handles POST /api/snapshots/{name}/load — loads and applies a snapshot.
func (h *Handler) HandleLoadSnapshot(w http.ResponseWriter, r *http.Request, name string) {
	snap, err := store.LoadSnapshot(h.dataDir, name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load entries failed: "+err.Error())
		return
	}

	matched, unmatched := store.ApplySnapshot(snap, entries)

	if err := h.store.Save(entries); err != nil {
		writeError(w, http.StatusInternalServerError, "save entries failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "applied",
		"name":      name,
		"matched":   matched,
		"unmatched": unmatched,
	})
}

// HandleDeleteSnapshot handles DELETE /api/snapshots/{name}.
func (h *Handler) HandleDeleteSnapshot(w http.ResponseWriter, r *http.Request, name string) {
	if err := store.DeleteSnapshot(h.dataDir, name); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "deleted",
		"name":   name,
	})
}

// mergeGroupEntry is a single source entry in a merge group response.
type mergeGroupInfo struct {
	MergedName    string        `json:"mergedName"`
	Sources       []string      `json:"sources"`
	SourceEntries []model.Entry `json:"sourceEntries"`
}

// HandleGetMergeGroups handles GET /api/merge-groups.
// Groups entries by mergeTarget and returns the grouped result.
func (h *Handler) HandleGetMergeGroups(w http.ResponseWriter, r *http.Request) {
	entries, err := h.store.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Build a map of ALB name -> Entry for target lookup
	entryMap := make(map[string]*model.Entry, len(entries))
	for i := range entries {
		entryMap[entries[i].ALBName] = &entries[i]
	}

	// Group entries by mergeTarget
	groups := make(map[string]*mergeGroupInfo)
	for _, e := range entries {
		if e.MergeTarget == "" {
			continue
		}
		g, exists := groups[e.MergeTarget]
		if !exists {
			// Determine mergedName: use target's MergedName if available
			mergedName := ""
			if target, ok := entryMap[e.MergeTarget]; ok {
				mergedName = target.MergedName
			}
			g = &mergeGroupInfo{
				MergedName:    mergedName,
				Sources:       []string{},
				SourceEntries: []model.Entry{},
			}
			groups[e.MergeTarget] = g
		}
		g.Sources = append(g.Sources, e.ALBName)
		g.SourceEntries = append(g.SourceEntries, e)
	}

	writeJSON(w, http.StatusOK, groups)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
