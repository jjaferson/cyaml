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
	deploymentName := ecsTemplate.ecsDeployment.Name
	subnetes := ecsTemplate.ecsDeployment.Network.Subnets
	cfTemplate := cloudformation.NewTemplate()

	cfTemplate.Resources["LoadBalancer"] = &cloudformation.AWSElasticLoadBalancingV2LoadBalancer{
		Name:    fmt.Sprintf("%s-loadbalancer", deploymentName),
		Scheme:  "internet-facing",
		Type:    "application",
		Subnets: subnetes,
	}

	y, err := cfTemplate.YAML()
	if err != nil {
		return
	}
	template = string(y)

	return template, nil
}

func (ecsTemplate *ecsTemplateGenerator) getLoadBalancerSecurityGroup() {

}
