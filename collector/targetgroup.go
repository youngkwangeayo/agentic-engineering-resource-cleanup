package collector

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"new-lb/model"
)

// collectTargetGroups fetches TGs for a given ALB ARN and returns TGInfo slice.
// Returns nil slice (not error) if the ALB has no TGs.
func collectTargetGroups(ctx context.Context, elbClient *elasticloadbalancingv2.Client, albArn string) ([]model.TGInfo, error) {
	out, err := elbClient.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{
		LoadBalancerArn: &albArn,
	})
	if err != nil {
		return nil, err
	}

	var tgs []model.TGInfo
	for _, tg := range out.TargetGroups {
		tgName := ""
		if tg.TargetGroupName != nil {
			tgName = *tg.TargetGroupName
		}
		tgArn := ""
		if tg.TargetGroupArn != nil {
			tgArn = *tg.TargetGroupArn
		}

		info := model.TGInfo{
			Name: tgName,
			ARN:  tgArn,
		}

		// Get target health for this TG
		healthOut, err := elbClient.DescribeTargetHealth(ctx, &elasticloadbalancingv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err != nil {
			log.Printf("[collector] warning: failed to get target health for TG %s: %v", tgName, err)
			tgs = append(tgs, info)
			continue
		}

		info.TargetCount = len(healthOut.TargetHealthDescriptions)
		for _, thd := range healthOut.TargetHealthDescriptions {
			if thd.TargetHealth == nil {
				continue
			}
			switch thd.TargetHealth.State {
			case elbtypes.TargetHealthStateEnumHealthy:
				info.HealthyCount++
			case elbtypes.TargetHealthStateEnumUnhealthy, elbtypes.TargetHealthStateEnumDraining:
				info.UnhealthyCount++
			case elbtypes.TargetHealthStateEnumUnused:
				info.UnusedCount++
			default:
				info.UnusedCount++
			}
		}

		tgs = append(tgs, info)
	}

	return tgs, nil
}
