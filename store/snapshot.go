package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"new-lb/model"
)

// Classification holds the user-editable classification fields for a single ALB.
type Classification struct {
	ALBName     string `json:"albName"`
	Solution    string `json:"solution"`
	Environment string `json:"environment"`
	Action      string `json:"action"`
	MergeTarget string `json:"mergeTarget"`
	MergedName  string `json:"mergedName,omitempty"`
	Note        string `json:"note"`
}

// Snapshot is a saved set of classifications with metadata.
type Snapshot struct {
	Name            string           `json:"name"`
	CreatedAt       string           `json:"createdAt"`
	EntryCount      int              `json:"entryCount"`
	Classifications []Classification `json:"classifications"`
}

// SnapshotMeta is the metadata-only view of a snapshot (no classifications).
type SnapshotMeta struct {
	Name       string `json:"name"`
	CreatedAt  string `json:"createdAt"`
	EntryCount int    `json:"entryCount"`
}

// invalidNameChars matches characters not allowed in snapshot names.
// Allowed: Korean, alphanumeric, -, _, ., space
// Blocked: / \ : * ? < > |
var invalidNameChars = regexp.MustCompile(`[/\\:*?<>|]`)

// validateSnapshotName checks the snapshot name for validity.
func validateSnapshotName(name string) error {
	if name == "" {
		return fmt.Errorf("snapshot name is empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("snapshot name exceeds 100 characters")
	}
	if invalidNameChars.MatchString(name) {
		return fmt.Errorf("snapshot name contains invalid characters (/\\:*?<>|)")
	}
	// Block path traversal sequences
	if strings.Contains(name, "..") {
		return fmt.Errorf("snapshot name contains invalid sequence '..'")
	}
	return nil
}

// snapshotsDir returns the snapshots directory path under dataDir.
func snapshotsDir(dataDir string) string {
	return filepath.Join(dataDir, "snapshots")
}

// SaveSnapshot extracts classification fields from entries and saves them as a named snapshot.
func SaveSnapshot(dataDir string, entries []model.Entry, name string) error {
	name = strings.TrimSpace(name)
	if err := validateSnapshotName(name); err != nil {
		return err
	}

	dir := snapshotsDir(dataDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create snapshots directory: %w", err)
	}

	classifications := make([]Classification, len(entries))
	for i, e := range entries {
		classifications[i] = Classification{
			ALBName:     e.ALBName,
			Solution:    e.Solution,
			Environment: e.Environment,
			Action:      e.Action,
			MergeTarget: e.MergeTarget,
			MergedName:  e.MergedName,
			Note:        e.Note,
		}
	}

	snap := Snapshot{
		Name:            name,
		CreatedAt:       time.Now().Format(time.RFC3339),
		EntryCount:      len(entries),
		Classifications: classifications,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	filePath := filepath.Join(dir, name+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("write snapshot file: %w", err)
	}

	log.Printf("[snapshot] saved snapshot %q with %d entries", name, len(entries))
	return nil
}

// ListSnapshots returns metadata for all snapshots, sorted by creation time (newest first).
func ListSnapshots(dataDir string) ([]SnapshotMeta, error) {
	dir := snapshotsDir(dataDir)
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SnapshotMeta{}, nil
		}
		return nil, fmt.Errorf("read snapshots directory: %w", err)
	}

	var metas []SnapshotMeta
	for _, de := range dirEntries {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".json") {
			continue
		}
		filePath := filepath.Join(dir, de.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[snapshot] warning: could not read %s: %v", filePath, err)
			continue
		}
		var snap Snapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			log.Printf("[snapshot] warning: could not parse %s: %v", filePath, err)
			continue
		}
		metas = append(metas, SnapshotMeta{
			Name:       snap.Name,
			CreatedAt:  snap.CreatedAt,
			EntryCount: snap.EntryCount,
		})
	}

	// Sort newest first
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].CreatedAt > metas[j].CreatedAt
	})

	if metas == nil {
		metas = []SnapshotMeta{}
	}
	return metas, nil
}

// LoadSnapshot reads and decodes a snapshot by name.
func LoadSnapshot(dataDir string, name string) (*Snapshot, error) {
	name = strings.TrimSpace(name)
	if err := validateSnapshotName(name); err != nil {
		return nil, err
	}

	filePath := filepath.Join(snapshotsDir(dataDir), name+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot not found: %s", name)
		}
		return nil, fmt.Errorf("read snapshot file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("parse snapshot: %w", err)
	}
	return &snap, nil
}

// DeleteSnapshot removes a snapshot file by name.
func DeleteSnapshot(dataDir string, name string) error {
	name = strings.TrimSpace(name)
	if err := validateSnapshotName(name); err != nil {
		return err
	}

	filePath := filepath.Join(snapshotsDir(dataDir), name+".json")
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot not found: %s", name)
		}
		return fmt.Errorf("delete snapshot file: %w", err)
	}

	log.Printf("[snapshot] deleted snapshot %q", name)
	return nil
}

// ApplySnapshot overwrites classification fields in entries using the snapshot data.
// Returns (matched, unmatched) counts where unmatched means the snapshot entry has no
// corresponding entry in the current entries list.
func ApplySnapshot(snap *Snapshot, entries []model.Entry) (matched int, unmatched int) {
	entryMap := make(map[string]*model.Entry, len(entries))
	for i := range entries {
		entryMap[entries[i].ALBName] = &entries[i]
	}

	for _, c := range snap.Classifications {
		entry, exists := entryMap[c.ALBName]
		if !exists {
			unmatched++
			continue
		}
		entry.Solution = c.Solution
		entry.Environment = c.Environment
		entry.Action = c.Action
		entry.MergeTarget = c.MergeTarget
		entry.MergedName = c.MergedName
		entry.Note = c.Note
		matched++
	}

	log.Printf("[snapshot] applied snapshot %q: matched=%d, unmatched=%d", snap.Name, matched, unmatched)
	return matched, unmatched
}
