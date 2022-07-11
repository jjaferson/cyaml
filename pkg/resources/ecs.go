package resources

import (
	"fmt"

	"github.com/awslabs/goformation/cloudformation"
	"github.com/jjaferson/cyaml/api"
)

const (
	// network mode default to awsvpc as we are using fargate
	networkMode         = "awsvpc"
	instanceTypeFargate = "FARGATE"
)

var (
	compatibilities = []string{instanceTypeFargate}
)

func NewECSCluster(name string) *cloudformation.AWSECSCluster {
	return &cloudformation.AWSECSCluster{
		ClusterName: name,
	}
}

func NewECSTaskDefinition(clusterName string, service api.Service, ecsExecutionRoleResourceName string) *cloudformation.AWSECSTaskDefinition {
	//TODO: implement envs
	//TODO: add log
	var envs map[string]string

	return &cloudformation.AWSECSTaskDefinition{
		Family:                  fmt.Sprintf("%s-%s", service.Name, clusterName),
		Cpu:                     service.CPU,
		Memory:                  service.Memory,
		NetworkMode:             networkMode,
		ExecutionRoleArn:        cloudformation.Ref(ecsExecutionRoleResourceName),
		ContainerDefinitions:    getContainerDefinitions(service.Name, service.Image, service.Port.To, envs),
		RequiresCompatibilities: compatibilities,
	}
}

func NewECSService(serviceName, ecsClusterResourceName, ecsTaskDefinitionResourceName, targetGroupResourceName, lbListenerResourceName, serviceSecurityGroupResourceName string, subnets []string, port int) *cloudformation.AWSECSService {
	svc := &cloudformation.AWSECSService{
		ServiceName:    serviceName,
		LaunchType:     instanceTypeFargate,
		Cluster:        cloudformation.Ref(ecsClusterResourceName),
		TaskDefinition: cloudformation.Ref(ecsTaskDefinitionResourceName),
		NetworkConfiguration: &cloudformation.AWSECSService_NetworkConfiguration{
			AwsvpcConfiguration: &cloudformation.AWSECSService_AwsVpcConfiguration{
				AssignPublicIp: "true",
				SecurityGroups: []string{
					cloudformation.Ref(serviceSecurityGroupResourceName),
				},
				Subnets: subnets,
			},
		},
		LoadBalancers: []cloudformation.AWSECSService_LoadBalancer{
			{
				ContainerName:  serviceName,
				ContainerPort:  port,
				TargetGroupArn: cloudformation.Ref(targetGroupResourceName),
			},
		},
	}
	svc.SetDependsOn([]string{lbListenerResourceName})
	return svc
}

func getContainerDefinitions(serviceName, image string, port int, envs map[string]string) []cloudformation.AWSECSTaskDefinition_ContainerDefinition {
	return []cloudformation.AWSECSTaskDefinition_ContainerDefinition{
		{
			Name:  serviceName,
			Image: image,
			PortMappings: []cloudformation.AWSECSTaskDefinition_PortMapping{
				{
					ContainerPort: port,
				},
			},
			Environment: getContainerEnvs(envs),
		},
	}
}

func getContainerEnvs(envs map[string]string) []cloudformation.AWSECSTaskDefinition_KeyValuePair {
	var containerEnvs []cloudformation.AWSECSTaskDefinition_KeyValuePair
	for key, value := range envs {
		containerEnvs = append(containerEnvs, cloudformation.AWSECSTaskDefinition_KeyValuePair{
			Name:  key,
			Value: value,
		})
	}
	return containerEnvs
}
