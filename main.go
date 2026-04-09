package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"new-lb/classifier"
	"new-lb/collector"
	"new-lb/server"
	"new-lb/store"
)

func main() {
	mode := flag.String("mode", "serve", "Run mode: collect, classify, serve, all")
	dataDir := flag.String("data", "data", "Data directory for entries.json")
	flag.Parse()

	// Resolve data path
	dataPath := filepath.Join(*dataDir, "entries.json")
	s := store.New(dataPath)

	log.Printf("[main] mode=%s, data=%s", *mode, dataPath)

	switch *mode {
	case "collect":
		runCollect(s)

	case "classify":
		runClassify(s)

	case "serve":
		runServe(s, *dataDir)

	case "all":
		runCollect(s)
		runClassify(s)
		runServe(s, *dataDir)

	default:
		log.Fatalf("[main] unknown mode: %s (use: collect, classify, serve, all)", *mode)
	}
}

func runCollect(s *store.Store) {
	log.Println("[main] starting collection...")

	// Load existing entries to preserve user-modified fields
	oldEntries, err := s.Load()
	if err != nil {
		log.Printf("[main] warning: could not load existing entries: %v", err)
		oldEntries = nil
	}

	ctx := context.Background()
	entries, err := collector.Run(ctx)
	if err != nil {
		log.Fatalf("[main] collection failed: %v", err)
	}

	// Merge user-modified fields from old entries
	if len(oldEntries) > 0 {
		entries = store.MergeEntries(oldEntries, entries)
	}

	if err := s.Save(entries); err != nil {
		log.Fatalf("[main] save failed: %v", err)
	}
	log.Printf("[main] collection complete: %d entries saved", len(entries))
}

func runClassify(s *store.Store) {
	log.Println("[main] starting classification...")
	entries, err := s.Load()
	if err != nil {
		log.Fatalf("[main] load failed: %v", err)
	}
	if len(entries) == 0 {
		log.Println("[main] no entries to classify, run collect first")
		return
	}

	classified, unknown := classifier.Classify(entries)
	log.Printf("[main] classified=%d, unknown=%d", classified, unknown)

	// Run tikitaka for unknowns if running interactively
	if unknown > 0 && isTerminal() {
		classifier.TikiTaka(entries)
	}

	if err := s.Save(entries); err != nil {
		log.Fatalf("[main] save failed: %v", err)
	}
	log.Println("[main] classification results saved")
}

func runServe(s *store.Store, dataDir string) {
	log.Println("[main] starting web server...")
	// Determine web path relative to the executable's working directory
	webPath := filepath.Join("web", "index.html")
	if err := server.Start(s, webPath, dataDir); err != nil {
		log.Fatalf("[main] server failed: %v", err)
	}
}

// isTerminal checks if stdin is a terminal (for tikitaka interactive mode).
func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
