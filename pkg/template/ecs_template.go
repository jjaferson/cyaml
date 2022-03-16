package template

import (
	"fmt"

	"github.com/awslabs/goformation/cloudformation"

	"github.com/jjaferson/cyaml/api"
	awsResources "github.com/jjaferson/cyaml/pkg/resources"
)

const (
	loadBalancerResourceName              = "LoadBalancer"
	loadBalancerSecurityGroup             = "LoadBalancerSecurityGroup"
	targetGroupResourceName               = "TargetGroup"
	loadBalancerListener                  = "LBListener"
	loadBalancerListenerRulesResourceName = "LBListenerRules"
	clusterResourceName                   = "ECSClusterName"
	ecsIAMExecutionRoleResourceName       = "ECSIAMExecutionRole"
	ecsTaskDefinitionResourceName         = "ECSTaskDefinition"
	ecsServiceResourceName                = "ECSService"
	ecsServiceSecurityGroup               = "ECSServiceSecurityGroup"
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
	ecsTemplate.addResource(clusterResourceName, awsResources.NewECSCluster(clusterName))

	// Creates cloudformation loadbalance
	ecsTemplate.addResource(loadBalancerResourceName,
		awsResources.NewLoadBalancer(
			clusterName, loadBalancerSecurityGroup, network.Subnets))

	// Creates security group for loadbalance (might become an dependency of NewLoadBalancer)
	ecsTemplate.addResource(loadBalancerSecurityGroup,
		awsResources.NewSecurityGroup(clusterName, network.ID, getIngressPorts()))

	// add service resources
	ecsTemplate.addDeploymentServiceResource()

	y, err := ecsTemplate.cfTemplate.YAML()
	if err != nil {
		return
	}
	template = string(y)

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
	lbListenerResourceName := map[int]string{}
	clusterName := ecsTemplate.ecsDeployment.Name

	// ECS Execution role for the task definitions
	ecsTemplate.addResource(ecsIAMExecutionRoleResourceName, awsResources.NewECSExecutionRole())

	// Creates load balancer listener for the ingress ports
	for i, port := range getIngressPorts() {
		lbListenerResourceName[port] = fmt.Sprintf("%s%d", loadBalancerListener, i)
		ecsTemplate.addResource(lbListenerResourceName[port],
			awsResources.NewLoadBalancerListener(loadBalancerResourceName, port))
	}

	for i, service := range services {
		targetGroupSvcResourceName := fmt.Sprintf("%s%d", targetGroupResourceName, i)
		lbListenerRulesSvcResourceName := fmt.Sprintf("%s%d", loadBalancerListenerRulesResourceName, i)
		ecsTaskDefinitionSvcResourceName := fmt.Sprintf("%s%d", ecsTaskDefinitionResourceName, i)
		ecsSvcResourceName := fmt.Sprintf("%s%d", ecsServiceResourceName, i)
		ecsSvcSecurityGroup := fmt.Sprintf("%s%d", ecsServiceSecurityGroup, i)

		//TODO: validate services
		// 1) port.from container in the ingress port

		// Creates target group
		ecsTemplate.addResource(
			targetGroupSvcResourceName,
			awsResources.NewTargetGroup(service.Name, network.ID, service.Port.To))

		//create ListenerRules
		ecsTemplate.addResource(lbListenerRulesSvcResourceName,
			awsResources.AddNewLoadBalancerListenerRules(
				lbListenerResourceName[service.Port.From], targetGroupSvcResourceName, service.Name))

		//Create task definition
		ecsTemplate.addResource(ecsTaskDefinitionSvcResourceName,
			awsResources.NewECSTaskDefinition(clusterName, service, ecsIAMExecutionRoleResourceName))

		// Creates security group to access the service
		ecsTemplate.addResource(ecsSvcSecurityGroup,
			awsResources.NewSecurityGroup(service.Name, network.ID, []int{service.Port.To}))

		// Creates ecs cluster service
		ecsTemplate.addResource(ecsSvcResourceName,
			awsResources.NewECSService(
				service.Name,
				clusterResourceName,
				ecsTaskDefinitionSvcResourceName,
				targetGroupSvcResourceName,
				lbListenerResourceName[service.Port.From],
				ecsSvcSecurityGroup,
				network.Subnets,
				service.Port.To,
			),
		)
	}
}
