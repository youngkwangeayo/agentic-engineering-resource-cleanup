package classifier

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"new-lb/model"
)

var solutions = []string{
	"signage", "cms", "nserise", "bss", "ncount",
	"srt", "aiagent", "wine", "ws2025", "dooh", "unknown",
}

var actions = []string{"유지", "합치기", "삭제", "미정"}

// TikiTaka runs an interactive CLI session for unclassified ALBs.
// It prompts the user to choose solution and action for each unknown entry.
func TikiTaka(entries []model.Entry) {
	var unknowns []int
	for i := range entries {
		if entries[i].Solution == "unknown" {
			unknowns = append(unknowns, i)
		}
	}

	if len(unknowns) == 0 {
		log.Println("[tikitaka] no unknown entries, skipping")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n=== Tikitaka: %d unknown ALBs to classify ===\n\n", len(unknowns))

	for seq, idx := range unknowns {
		e := &entries[idx]

		// Show ALB info
		recordNames := make([]string, 0, len(e.Records))
		for _, r := range e.Records {
			recordNames = append(recordNames, r.Name)
		}
		recordStr := "none"
		if len(recordNames) > 0 {
			recordStr = strings.Join(recordNames, ", ")
		}

		fmt.Printf("[%d/%d] ALB: %s (%s, %s, records: %s)\n",
			seq+1, len(unknowns), e.ALBName, e.Environment, e.Status, recordStr)

		// Ask for solution
		fmt.Println("  Select solution:")
		for i, s := range solutions {
			fmt.Printf("  %2d. %s", i+1, s)
			if (i+1)%5 == 0 {
				fmt.Println()
			}
		}
		if len(solutions)%5 != 0 {
			fmt.Println()
		}

		solIdx := askChoice(reader, "  Selection: ", len(solutions))
		e.Solution = solutions[solIdx]

		// Ask for action
		fmt.Println("  Select action:")
		for i, a := range actions {
			fmt.Printf("  %d. %s  ", i+1, a)
		}
		fmt.Println()

		actIdx := askChoice(reader, "  Selection: ", len(actions))
		e.Action = actions[actIdx]

		fmt.Printf("  -> %s: solution=%s, action=%s\n\n", e.ALBName, e.Solution, e.Action)
	}

	fmt.Println("=== Tikitaka complete ===")
}

// askChoice prompts until the user enters a valid number in [1, max].
func askChoice(reader *bufio.Reader, prompt string, max int) int {
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		n, err := strconv.Atoi(line)
		if err != nil || n < 1 || n > max {
			fmt.Printf("  Please enter a number between 1 and %d.\n", max)
			continue
		}
		return n - 1
	}
}
