package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"new-lb/model"
)

// Store handles JSON file persistence with mutex protection.
type Store struct {
	mu       sync.Mutex
	filePath string
}

// New creates a Store pointing to the given file path.
// It ensures the parent directory exists.
func New(filePath string) *Store {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[store] warning: could not create directory %s: %v", dir, err)
	}
	return &Store{filePath: filePath}
}

// Load reads entries from the JSON file.
// Returns an empty slice if the file does not exist.
func (s *Store) Load() ([]model.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *Store) loadLocked() ([]model.Entry, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Entry{}, nil
		}
		return nil, fmt.Errorf("read %s: %w", s.filePath, err)
	}
	if len(data) == 0 {
		return []model.Entry{}, nil
	}
	var entries []model.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		log.Printf("[store] warning: failed to parse %s, returning empty: %v", s.filePath, err)
		return []model.Entry{}, nil
	}
	return entries, nil
}

// Save writes entries to the JSON file atomically (write to tmp, then rename).
func (s *Store) Save(entries []model.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked(entries)
}

func (s *Store) saveLocked(entries []model.Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal entries: %w", err)
	}
	tmpPath := s.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmpPath, s.filePath); err != nil {
		return fmt.Errorf("rename temp to final: %w", err)
	}
	log.Printf("[store] saved %d entries to %s", len(entries), s.filePath)
	return nil
}

// MergeEntries preserves user-modified fields from oldEntries into newEntries.
// Matching is done by ALB name. For each new entry, if a matching old entry exists:
//   - solution: kept if old value is not "unknown"
//   - environment: kept if old value is not "" or "unknown"
//   - action: kept if old value is not "미정"
//   - mergeTarget: always kept from old
//   - mergedName: always kept from old
//   - note: always kept from old
func MergeEntries(oldEntries, newEntries []model.Entry) []model.Entry {
	oldMap := make(map[string]*model.Entry, len(oldEntries))
	for i := range oldEntries {
		oldMap[oldEntries[i].ALBName] = &oldEntries[i]
	}

	for i := range newEntries {
		old, exists := oldMap[newEntries[i].ALBName]
		if !exists {
			continue
		}
		if old.Solution != "unknown" {
			newEntries[i].Solution = old.Solution
		}
		if old.Environment != "" && old.Environment != "unknown" {
			newEntries[i].Environment = old.Environment
		}
		if old.Action != "미정" {
			newEntries[i].Action = old.Action
		}
		if old.MergeTarget != "" {
			newEntries[i].MergeTarget = old.MergeTarget
		}
		if old.MergedName != "" {
			newEntries[i].MergedName = old.MergedName
		}
		if old.Note != "" {
			newEntries[i].Note = old.Note
		}
	}

	log.Printf("[store] merged user data: %d old entries matched against %d new entries", len(oldEntries), len(newEntries))
	return newEntries
}

// UpdateEntry modifies fields of a single entry by ALB name.
// Supported fields: solution, action, note, mergeTarget, environment.
func (s *Store) UpdateEntry(name string, updates map[string]any) (*model.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := s.loadLocked()
	if err != nil {
		return nil, err
	}

	var found *model.Entry
	for i := range entries {
		if entries[i].ALBName == name {
			found = &entries[i]
			break
		}
	}
	if found == nil {
		return nil, fmt.Errorf("entry not found: %s", name)
	}

	if v, ok := updates["solution"]; ok {
		if s, ok := v.(string); ok {
			found.Solution = s
		}
	}
	if v, ok := updates["action"]; ok {
		if s, ok := v.(string); ok {
			found.Action = s
		}
	}
	if v, ok := updates["note"]; ok {
		if s, ok := v.(string); ok {
			found.Note = s
		}
	}
	if v, ok := updates["mergeTarget"]; ok {
		if s, ok := v.(string); ok {
			found.MergeTarget = s
		}
	}
	if v, ok := updates["environment"]; ok {
		if s, ok := v.(string); ok {
			found.Environment = s
		}
	}
	if v, ok := updates["mergedName"]; ok {
		if s, ok := v.(string); ok {
			found.MergedName = s
		}
	}

	if err := s.saveLocked(entries); err != nil {
		return nil, err
	}
	return found, nil
}
