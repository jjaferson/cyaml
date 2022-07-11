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

func NewTargetGroup(name, vpcId string, port int) *cloudformation.AWSElasticLoadBalancingV2TargetGroup {
	name = fmt.Sprintf("%s-lb-targetgroup", name)
	return &cloudformation.AWSElasticLoadBalancingV2TargetGroup{
		Name:       name,
		VpcId:      vpcId,
		TargetType: "ip",
		Port:       port,
	}
}

func NewLoadBalancerListener(loadbalancerResourceName string, port int) *cloudformation.AWSElasticLoadBalancingV2Listener {
	protocol := "HTTPS"
	if port == 80 {
		protocol = "HTTP"
	}

	return &cloudformation.AWSElasticLoadBalancingV2Listener{
		LoadBalancerArn: cloudformation.Ref(loadbalancerResourceName),
		Port:            port,
		Protocol:        protocol,
	}
}

func AddNewLoadBalancerListenerRules(loadbalancerListenerResourceName, targetGroupResourceName, serviceName string) *cloudformation.AWSElasticLoadBalancingV2ListenerRule {
	// at the moment we only support forward type via the path-pattern
	lbListenerRules := cloudformation.AWSElasticLoadBalancingV2ListenerRule{
		Actions: []cloudformation.AWSElasticLoadBalancingV2ListenerRule_Action{
			{
				TargetGroupArn: cloudformation.Ref(targetGroupResourceName),
				Type:           "forward",
			},
		},
		ListenerArn: cloudformation.Ref(loadbalancerListenerResourceName),
		Conditions: []cloudformation.AWSElasticLoadBalancingV2ListenerRule_RuleCondition{
			{
				Field:  "path-pattern",
				Values: []string{serviceName},
			},
		},
	}

	return &lbListenerRules
}
