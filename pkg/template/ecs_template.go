package template

import (
	"fmt"

	"github.com/awslabs/goformation/cloudformation"

	"github.com/jjaferson/cyaml/api"
)

type ECSTemplate interface {
	Validate()
	Generate() (template string, err error)
}

type ecsTemplateGenerator struct {
	ecsDeployment *api.ECSDeployment
}

var _ ECSTemplate = &ecsTemplateGenerator{}

func NewECSTemplateGenerator(ecsDeployment *api.ECSDeployment) *ecsTemplateGenerator {
	return &ecsTemplateGenerator{
		ecsDeployment: ecsDeployment,
	}
}

func (ecsTemplate *ecsTemplateGenerator) Validate() {

}

func (ecsTemplate *ecsTemplateGenerator) Generate() (template string, err error) {

	cfTemplate := cloudformation.NewTemplate()
	loadBalancerSecGroupResourceName := "LoadBalancerSecurityGroup"

	cfTemplate.Resources[loadBalancerSecGroupResourceName] = ecsTemplate.getLoadBalancerSecurityGroup()
	cfTemplate.Resources["LoadBalancer"] = ecsTemplate.getLoadBalancer(loadBalancerSecGroupResourceName)

	y, err := cfTemplate.YAML()
	if err != nil {
		return
	}
	template = string(y)

	return template, nil
}

func (ecsTemplate *ecsTemplateGenerator) getLoadBalancer(securityGroupResourceName string) *cloudformation.AWSElasticLoadBalancingV2LoadBalancer {
	deploymentName := ecsTemplate.ecsDeployment.Name
	subnetes := ecsTemplate.ecsDeployment.Network.Subnets

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

func (ecsTemplate *ecsTemplateGenerator) getLoadBalancerSecurityGroup() *cloudformation.AWSEC2SecurityGroup {

	sgName := fmt.Sprintf("%s-security-group", ecsTemplate.ecsDeployment.Name)
	//TODO validate name size
	//TODO add loadbalance external port to deployment and if not specified use 80 and 433

	securityGroup := &cloudformation.AWSEC2SecurityGroup{
		VpcId:     ecsTemplate.ecsDeployment.Network.ID,
		GroupName: sgName,
		SecurityGroupIngress: []cloudformation.AWSEC2SecurityGroup_Ingress{
			{
				CidrIp:     "0.0.0.0/0",
				FromPort:   80,
				ToPort:     80,
				IpProtocol: "tcp",
			},
			{
				CidrIp:     "0.0.0.0/0",
				FromPort:   433,
				ToPort:     433,
				IpProtocol: "tcp",
			},
		},
		SecurityGroupEgress: []cloudformation.AWSEC2SecurityGroup_Egress{
			{
				CidrIp:     "0.0.0.0/0",
				IpProtocol: "-1",
			},
		},
	}
	return securityGroup
}
