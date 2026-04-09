package collector

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"

	"new-lb/model"
)

// r53Record is an intermediate struct for Route53 record data before ALB matching.
type r53Record struct {
	RecordName string
	ZoneID     string
	ZoneName   string
	Type       string // "ALIAS" or "CNAME"
	TargetDNS  string // the ALB DNS this record points to (normalized, lowercase, no trailing dot)
}

// collectRoute53Records fetches all hosted zones and their records,
// returning a map keyed by normalized ALB DNS to []model.Record.
func collectRoute53Records(ctx context.Context, r53Client *route53.Client) (map[string][]model.Record, error) {
	log.Println("[collector] fetching Route53 hosted zones...")

	// Step 1: List all hosted zones
	var zones []r53types.HostedZone
	var zoneMarker *string
	for {
		out, err := r53Client.ListHostedZones(ctx, &route53.ListHostedZonesInput{
			Marker: zoneMarker,
		})
		if err != nil {
			return nil, err
		}
		zones = append(zones, out.HostedZones...)
		if !out.IsTruncated {
			break
		}
		zoneMarker = out.NextMarker
	}
	log.Printf("[collector] found %d hosted zones", len(zones))

	// Step 2: For each zone, list all records and extract ALIAS/CNAME pointing to ELB
	var allRecords []r53Record
	for _, zone := range zones {
		zoneID := ""
		if zone.Id != nil {
			zoneID = strings.TrimPrefix(*zone.Id, "/hostedzone/")
		}
		zoneName := ""
		if zone.Name != nil {
			zoneName = strings.TrimSuffix(*zone.Name, ".")
		}

		var nextName *string
		var nextType r53types.RRType
		hasMore := true

		for hasMore {
			input := &route53.ListResourceRecordSetsInput{
				HostedZoneId: zone.Id,
			}
			if nextName != nil {
				input.StartRecordName = nextName
				input.StartRecordType = nextType
			}

			out, err := r53Client.ListResourceRecordSets(ctx, input)
			if err != nil {
				log.Printf("[collector] warning: failed to list records for zone %s: %v", zoneName, err)
				break
			}

			for _, rr := range out.ResourceRecordSets {
				recName := ""
				if rr.Name != nil {
					recName = strings.TrimSuffix(*rr.Name, ".")
				}

				// Check ALIAS records pointing to ELB
				if rr.AliasTarget != nil && rr.AliasTarget.DNSName != nil {
					target := normalizeDNS(*rr.AliasTarget.DNSName)
					if isELBDNS(target) {
						allRecords = append(allRecords, r53Record{
							RecordName: recName,
							ZoneID:     zoneID,
							ZoneName:   zoneName,
							Type:       "ALIAS",
							TargetDNS:  target,
						})
					}
				}

				// Check CNAME records pointing to ELB
				if rr.Type == r53types.RRTypeCname && len(rr.ResourceRecords) > 0 {
					for _, val := range rr.ResourceRecords {
						if val.Value != nil {
							target := normalizeDNS(*val.Value)
							if isELBDNS(target) {
								allRecords = append(allRecords, r53Record{
									RecordName: recName,
									ZoneID:     zoneID,
									ZoneName:   zoneName,
									Type:       "CNAME",
									TargetDNS:  target,
								})
							}
						}
					}
				}
			}

			if !out.IsTruncated {
				hasMore = false
			} else {
				nextName = out.NextRecordName
				nextType = out.NextRecordType
			}
		}
	}

	log.Printf("[collector] found %d Route53 records pointing to ELB", len(allRecords))

	// Step 3: Build map[albDNS][]Record
	result := make(map[string][]model.Record)
	for _, r := range allRecords {
		rec := model.Record{
			Name:     r.RecordName,
			ZoneID:   r.ZoneID,
			ZoneName: r.ZoneName,
			Type:     r.Type,
		}
		result[r.TargetDNS] = append(result[r.TargetDNS], rec)
	}

	return result, nil
}

// normalizeDNS removes trailing dot and converts to lowercase.
func normalizeDNS(dns string) string {
	return strings.ToLower(strings.TrimSuffix(dns, "."))
}

// isELBDNS checks if a DNS name looks like an ELB DNS.
func isELBDNS(dns string) bool {
	return strings.Contains(dns, ".elb.amazonaws.com")
}
