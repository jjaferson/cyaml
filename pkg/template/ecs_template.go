package template

import (
	"fmt"

	"github.com/awslabs/goformation/cloudformation"

	"github.com/jjaferson/cyaml/api"
	awsResources "github.com/jjaferson/cyaml/pkg/resources"
)

const (
	loadBalancerResourceName  = "LoadBalancer"
	loadBalancerSecurityGroup = "LoadBalancerSecurityGroup"
	targetGroupResourceName   = "TargetGroup"
	loadBalancerListener      = "LBListener"
	loadBalancerListenerRules = "LBListenerRules"
	clusterName               = "ECSClusterName"
)

type ECSTemplate interface {
	Validate()
	Generate() (template string, err error)
}

type ecsTemplateGenerator struct {
	ecsDeployment *api.ECSDeployment
	cfTemplate    *cloudformation.Template
}

var _ ECSTemplate = &ecsTemplateGenerator{}

func NewECSTemplateGenerator(ecsDeployment *api.ECSDeployment) *ecsTemplateGenerator {
	cfTemplate := cloudformation.NewTemplate()

	return &ecsTemplateGenerator{
		ecsDeployment: ecsDeployment,
		cfTemplate:    cfTemplate,
	}
}

func (ecsTemplate *ecsTemplateGenerator) Validate() {

}

func (ecsTemplate *ecsTemplateGenerator) Generate() (template string, err error) {

	clusterName := ecsTemplate.ecsDeployment.Name
	network := ecsTemplate.ecsDeployment.Network

	// Creates cloudformation ECS cluster
	ecsTemplate.addResource(clusterName, awsResources.NewECSCluster(clusterName))

	// Creates cloudformation loadbalance
	ecsTemplate.addResource(loadBalancerResourceName,
		awsResources.NewLoadBalancer(
			clusterName, loadBalancerSecurityGroup, network.Subnets))

	// Creates security group for loadbalance (might become an dependency of NewLoadBalancer)
	ecsTemplate.addResource(loadBalancerSecurityGroup,
		awsResources.NewSecurityGroup(clusterName, network.ID, getIngressPorts()))

	// y, err := ecsTemplate.cfTemplate.YAML()
	// if err != nil {
	// 	return
	// }
	// template = string(y)

	return template, nil
}

func (ecsTemplate *ecsTemplateGenerator) addResource(resourceName string, resource cloudformation.Resource) {
	ecsTemplate.cfTemplate.Resources[resourceName] = resource
}

func getIngressPorts() []int {
	return []int{80, 433}
}

func (ecsTemplate *ecsTemplateGenerator) addDeploymentServiceResource() {
	services := ecsTemplate.ecsDeployment.Services
	network := ecsTemplate.ecsDeployment.Network
	lbListenerResourceNames := map[int]string{}

	// Creates load balancer listener for the ingress ports
	for i, port := range getIngressPorts() {
		lbListenerResourceNames[port] = fmt.Sprintf("%s%d", loadBalancerListener, i)
		ecsTemplate.addResource(lbListenerResourceNames[port],
			awsResources.NewLoadBalancerListener(loadBalancerResourceName, port))
	}

	for _, service := range services {
		targetGroupResourceName := fmt.Sprintf("%s%s", service.Name, targetGroupResourceName)
		lbListenerRulesResourceName := fmt.Sprintf("%s%s", service.Name, loadBalancerListenerRules)

		//TODO: validate services
		// 1) port.from container in the ingress port

		// Creates target group
		ecsTemplate.addResource(
			targetGroupResourceName,
			awsResources.NewTargetGroup(service.Name, network.ID, service.Port.To))

		//create ListenerRules
		ecsTemplate.addResource(lbListenerRulesResourceName,
			awsResources.AddNewLoadBalancerListenerRules(
				lbListenerResourceNames[service.Port.From], targetGroupResourceName, service.Name))

		//TODO Create task definition

	}
}
