package resources

import (
	"fmt"

	"github.com/awslabs/goformation/cloudformation"
)

func NewLoadBalancer(name, securityGroupResourceName string, subnets []string) *cloudformation.AWSElasticLoadBalancingV2LoadBalancer {
	deploymentName := name
	subnetes := subnets

	return &cloudformation.AWSElasticLoadBalancingV2LoadBalancer{
		Name:    fmt.Sprintf("%s-loadbalancer", deploymentName),
		Scheme:  "internet-facing",
		Type:    "application",
		Subnets: subnetes,
		SecurityGroups: []string{
			cloudformation.GetAtt(securityGroupResourceName, "GroupId"),
		},
	}
}

func NewTargetGroup(name, vpcId string) *cloudformation.AWSElasticLoadBalancingV2TargetGroup {
	name = fmt.Sprintf("%s-lb-targetgroup", name)
	return &cloudformation.AWSElasticLoadBalancingV2TargetGroup{
		Name:       name,
		VpcId:      vpcId,
		TargetType: "ip",
	}
}

func NewLoadBalancerListener(loadbalancerResourceName string, ports []int) []*cloudformation.AWSElasticLoadBalancingV2Listener {
	var listeners []*cloudformation.AWSElasticLoadBalancingV2Listener
	for _, port := range ports {
		protocal := "HTTPS"
		if port == 80 {
			protocal = "HTTP"
		}

		listeners = append(listeners, &cloudformation.AWSElasticLoadBalancingV2Listener{
			LoadBalancerArn: cloudformation.Ref(loadbalancerResourceName),
			Port:            port,
			Protocol:        protocal,
		})
	}
	return listeners
}
