package classifier

import (
	"log"
	"regexp"
	"strings"

	"new-lb/model"
)

// solutionKeywords maps solution names to their keywords for matching.
// Order matters: more specific keywords should be checked first.
var solutionKeywords = []struct {
	Solution string
	Keywords []string
}{
	// Specific/unique keywords first
	{"signage", []string{"signage", "scms", "swaiting", "tizenweb", "oss-cms", "oss-waiting", "product-8080", "waiting-808"}},
	{"cms", []string{"cms-elb", "dev-cms"}},
	{"aiagent", []string{"aiagent", "knowledge-graph"}},
	{"nserise", []string{"nseries", "nextpay-kiosk", "nextpay-nkds", "nextpay-norder", "nextpay-npos", "nextpay-order", "nextpay-shop", "socket-dev", "webservice"}},
	{"bss", []string{"bss-"}},
	{"ncount", []string{"ncount"}},
	{"srt", []string{"srt"}},
	{"wine", []string{"wine"}},
	{"ws2025", []string{"ws2025"}},
	{"dooh", []string{"dooh"}},
}

// suspendedProjects are projects that are stopped/pending confirmation.
var suspendedProjects = map[string]bool{
	"nserise": true,
	"bss":     true,
	"srt":     true,
	"ws2025":  true,
}

// patternRegex matches ALB names like alb-dev-signage, lb-prd-cms, elb-stg-xxx
var patternRegex = regexp.MustCompile(`^(alb|lb|elb)-?(dev|stg|prd)-(.+)$`)

// Classify runs the auto-classification pipeline on all entries.
// Returns the number of classified entries and the number of unknowns.
func Classify(entries []model.Entry) (classified int, unknown int) {
	log.Println("[classifier] running auto-classification...")

	for i := range entries {
		// Skip solution auto-classification if user already set it manually
		if entries[i].Solution == "" || entries[i].Solution == "unknown" {
			sol, env := classifyEntry(entries[i].ALBName)
			entries[i].Solution = sol
			entries[i].Environment = env
		}

		// Auto-detect environment if still empty or unknown (independent of solution)
		if entries[i].Environment == "" || entries[i].Environment == "unknown" {
			entries[i].Environment = detectEnvironment(entries[i].ALBName)
		}

		// Skip action auto-inference if user already set it manually
		if entries[i].Action == "" || entries[i].Action == "미정" {
			entries[i].Action = inferAction(&entries[i])
		}

		if entries[i].Solution == "unknown" {
			unknown++
		} else {
			classified++
		}

		log.Printf("[classifier] %s -> solution=%s, env=%s, action=%s",
			entries[i].ALBName, entries[i].Solution, entries[i].Environment, entries[i].Action)
	}

	log.Printf("[classifier] classification complete: %d classified, %d unknown", classified, unknown)
	return classified, unknown
}

// classifyEntry determines solution and environment for a single ALB name.
func classifyEntry(albName string) (solution, environment string) {
	lower := strings.ToLower(albName)

	// Step 1: Regex pattern matching (alb|lb|elb)-{env}-{solution})
	if m := patternRegex.FindStringSubmatch(lower); m != nil {
		env := m[2]
		rest := m[3]
		sol := matchSolutionByKeyword(rest)
		if sol == "" {
			sol = matchSolutionByKeyword(lower)
		}
		if sol != "" {
			return sol, env
		}
	}

	// Step 2: Keyword matching (primary strategy)
	solution = matchSolutionByKeyword(lower)
	if solution == "" {
		solution = "unknown"
	}

	environment = detectEnvironment(albName)
	return solution, environment
}

// matchSolutionByKeyword checks if the ALB name contains any known solution keyword.
func matchSolutionByKeyword(lowerName string) string {
	for _, sk := range solutionKeywords {
		for _, kw := range sk.Keywords {
			if strings.Contains(lowerName, kw) {
				return sk.Solution
			}
		}
	}
	return ""
}

// detectEnvironment determines the environment from the ALB name.
func detectEnvironment(albName string) string {
	lower := strings.ToLower(albName)
	switch {
	case strings.Contains(lower, "-dev") || strings.HasPrefix(lower, "dev-"):
		return "dev"
	case strings.Contains(lower, "staging") || strings.Contains(lower, "-stg"):
		return "stg"
	default:
		return "prd"
	}
}

// inferAction determines the recommended action for an entry.
func inferAction(entry *model.Entry) string {
	// 1. No target -> delete
	if entry.Status == "no_target" {
		return "삭제"
	}
	// 2. No record -> delete
	if entry.Status == "no_record" {
		return "삭제"
	}
	// 3. HTTPS error -> delete
	if entry.Status == "error" {
		return "삭제"
	}
	// 4. Suspended project -> pending
	if suspendedProjects[entry.Solution] {
		return "미정"
	}
	// 5. Normal -> keep
	return "유지"
}
