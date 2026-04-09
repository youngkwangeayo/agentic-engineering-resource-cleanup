package model

// Entry represents a single ALB with all collected/classified data.
type Entry struct {
	ALBName     string   `json:"albName"`
	ALBArn      string   `json:"albArn"`
	ALBDNS      string   `json:"albDns"`
	Solution    string   `json:"solution"`    // signage|cms|nserise|bss|ncount|srt|aiagent|wine|ws2025|dooh|unknown
	Environment string   `json:"environment"` // dev|stg|prd|unknown
	Status      string   `json:"status"`      // healthy|unhealthy|no_target|no_record|error|unknown
	Action      string   `json:"action"`      // 유지|합치기|삭제|미정
	Records     []Record `json:"records"`
	TGs         []TGInfo `json:"targetGroups"`
	MergeTarget string   `json:"mergeTarget,omitempty"`
	MergedName  string   `json:"mergedName,omitempty"`
	Note        string   `json:"note,omitempty"`
}

// Record represents a Route53 record pointing to an ALB.
type Record struct {
	Name       string `json:"name"`       // Route53 record name (e.g. dev-scms.nextpay.co.kr)
	ZoneID     string `json:"zoneId"`     // Hosted Zone ID
	ZoneName   string `json:"zoneName"`   // Hosted Zone name (e.g. nextpay.co.kr)
	Type       string `json:"type"`       // ALIAS or CNAME
	HealthCode int    `json:"healthCode"` // HTTP response code, -1=cert error, 0=timeout, -2=DNS failure
}

// TGInfo represents a Target Group and its health summary.
type TGInfo struct {
	Name           string `json:"name"`
	ARN            string `json:"arn"`
	HealthyCount   int    `json:"healthyCount"`
	UnhealthyCount int    `json:"unhealthyCount"`
	UnusedCount    int    `json:"unusedCount"`
	TargetCount    int    `json:"targetCount"`
}

// MergePlan describes a merge operation between ALBs.
type MergePlan struct {
	TargetALB  string   `json:"targetAlb"`
	SourceALBs []string `json:"sourceAlbs"`
	Records    []Record `json:"records"`
	TGs        []TGInfo `json:"targetGroups"`
	Note       string   `json:"note,omitempty"`
}

// Report is the summary response for GET /api/report.
type Report struct {
	Total            int            `json:"total"`
	BySolution       map[string]int `json:"bySolution"`
	ByStatus         map[string]int `json:"byStatus"`
	ByAction         map[string]int `json:"byAction"`
	ByEnvironment    map[string]int `json:"byEnvironment"`
	DeleteCandidates []string       `json:"deleteCandidates"`
	MergeCandidates  []MergePlan    `json:"mergeCandidates"`
}
