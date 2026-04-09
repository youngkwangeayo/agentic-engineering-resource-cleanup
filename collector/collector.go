package collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"golang.org/x/sync/errgroup"

	"new-lb/model"
)

const (
	maxTGWorkers     = 10
	maxHealthWorkers = 5
	httpTimeout      = 10 * time.Second
)

// Run executes the full collection pipeline:
// Phase A: ALB list + Route53 records (parallel)
// Phase B: ALB DNS <-> Route53 matching
// Phase C: TG health (10 workers) + HTTPS healthcheck (5 workers)
func Run(ctx context.Context) ([]model.Entry, error) {
	log.Println("[collector] starting collection pipeline...")

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-2"))
	if err != nil {
		return nil, err
	}

	elbClient := elasticloadbalancingv2.NewFromConfig(cfg)
	r53Client := route53.NewFromConfig(cfg)

	// Phase A: Parallel ALB + Route53 collection
	var entries []model.Entry
	var recordMap map[string][]model.Record

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		entries, err = collectALBs(gctx, elbClient)
		return err
	})

	g.Go(func() error {
		var err error
		recordMap, err = collectRoute53Records(gctx, r53Client)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Phase B: Match ALB DNS to Route53 records
	log.Println("[collector] matching ALB DNS to Route53 records...")
	for i := range entries {
		albDNS := normalizeDNS(entries[i].ALBDNS)
		if records, ok := recordMap[albDNS]; ok {
			entries[i].Records = records
		}
	}

	// Phase C: TG health + HTTPS healthcheck (worker pools)
	log.Println("[collector] collecting Target Group health (workers=10)...")
	collectTGsConcurrently(ctx, elbClient, entries)

	log.Println("[collector] running HTTPS healthchecks (workers=5)...")
	runHealthChecks(entries)

	// Determine final status for each entry
	for i := range entries {
		entries[i].Status = determineStatus(&entries[i])
	}

	log.Printf("[collector] collection complete: %d entries", len(entries))
	return entries, nil
}

// collectTGsConcurrently fetches TG info for all entries using a worker pool.
func collectTGsConcurrently(ctx context.Context, elbClient *elasticloadbalancingv2.Client, entries []model.Entry) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxTGWorkers)

	for i := range entries {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			tgs, err := collectTargetGroups(ctx, elbClient, entries[idx].ALBArn)
			if err != nil {
				log.Printf("[collector] warning: TG collection failed for %s: %v", entries[idx].ALBName, err)
				entries[idx].Status = "unknown"
				return
			}
			entries[idx].TGs = tgs
		}(i)
	}
	wg.Wait()
}

// runHealthChecks performs HTTPS healthchecks on all records of all entries.
func runHealthChecks(entries []model.Entry) {
	type healthJob struct {
		entryIdx  int
		recordIdx int
		domain    string
	}

	var jobs []healthJob
	for i := range entries {
		for j := range entries[i].Records {
			domain := entries[i].Records[j].Name
			if domain != "" {
				jobs = append(jobs, healthJob{entryIdx: i, recordIdx: j, domain: domain})
			}
		}
	}

	if len(jobs) == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxHealthWorkers)

	for _, job := range jobs {
		wg.Add(1)
		sem <- struct{}{}
		go func(j healthJob) {
			defer wg.Done()
			defer func() { <-sem }()

			code := checkHealth(j.domain)
			entries[j.entryIdx].Records[j.recordIdx].HealthCode = code
			log.Printf("[healthcheck] %s -> %d", j.domain, code)
		}(job)
	}
	wg.Wait()
}

// determineStatus determines the final status of an entry based on TG and healthcheck data.
// Priority: no_target > no_record > error > unhealthy > healthy
func determineStatus(e *model.Entry) string {
	// If status was already set to "unknown" due to TG collection failure, keep it
	if e.Status == "unknown" && len(e.TGs) == 0 {
		// Check if it was explicitly set to unknown (TG failure case)
		// vs just never collected
	}

	// If status was already set to "unknown" due to TG collection failure, keep it
	if e.Status == "unknown" && len(e.TGs) == 0 {
		return "unknown"
	}

	// Check no_target: no TGs at all, or TGs exist but all have targetCount == 0
	if len(e.TGs) == 0 {
		return "no_target"
	}
	allNoTarget := true
	for _, tg := range e.TGs {
		if tg.TargetCount > 0 {
			allNoTarget = false
			break
		}
	}
	if allNoTarget {
		return "no_target"
	}

	// Check no_record
	if len(e.Records) == 0 {
		return "no_record"
	}

	// Check healthcheck results for error
	hasError := false
	hasUnhealthy := false
	for _, rec := range e.Records {
		code := rec.HealthCode
		if code == -1 || code == -2 || code == 0 || code >= 500 {
			hasError = true
		}
	}

	// Check TG health for unhealthy
	for _, tg := range e.TGs {
		if tg.UnhealthyCount > 0 {
			hasUnhealthy = true
		}
	}

	if hasError {
		return "error"
	}
	if hasUnhealthy {
		return "unhealthy"
	}

	// Check if all TGs have only unused targets
	allUnused := true
	for _, tg := range e.TGs {
		if tg.HealthyCount > 0 || tg.UnhealthyCount > 0 {
			allUnused = false
			break
		}
	}
	if allUnused && len(e.TGs) > 0 {
		// All targets are in unused state
		return "unhealthy"
	}

	return "healthy"
}
