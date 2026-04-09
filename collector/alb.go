package collector

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"

	"new-lb/model"
)

// collectALBs fetches all ALBs from ELBv2 API, filtering out k8s-managed ones.
func collectALBs(ctx context.Context, elbClient *elasticloadbalancingv2.Client) ([]model.Entry, error) {
	log.Println("[collector] fetching ALB list...")

	var entries []model.Entry
	var marker *string

	for {
		input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
			Marker: marker,
		}
		out, err := elbClient.DescribeLoadBalancers(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, lb := range out.LoadBalancers {
			name := ""
			if lb.LoadBalancerName != nil {
				name = *lb.LoadBalancerName
			}
			// k8s ALB controller managed ALBs are excluded
			if strings.HasPrefix(name, "k8s") {
				log.Printf("[collector] skipping k8s ALB: %s", name)
				continue
			}

			arn := ""
			if lb.LoadBalancerArn != nil {
				arn = *lb.LoadBalancerArn
			}
			dns := ""
			if lb.DNSName != nil {
				dns = *lb.DNSName
			}

			entries = append(entries, model.Entry{
				ALBName:     name,
				ALBArn:      arn,
				ALBDNS:      dns,
				Solution:    "unknown",
				Environment: "unknown",
				Status:      "unknown",
				Action:      "미정",
			})
		}

		if out.NextMarker == nil {
			break
		}
		marker = out.NextMarker
	}

	log.Printf("[collector] found %d ALBs (k8s excluded)", len(entries))
	return entries, nil
}
