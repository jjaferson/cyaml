package resources

import (
	"github.com/awslabs/goformation/cloudformation"
	"github.com/jjaferson/cyaml/api"
)

const (
	// network mode default to awsvpc as we are using fargate
	networkMode = "awsvpc"
)

func NewECSTaskDefinition(name string, service api.Service) *cloudformation.AWSECSTaskDefinition {

	return &cloudformation.AWSECSTaskDefinition{
		Family:      name,
		Cpu:         service.CPU,
		Memory:      service.Memory,
		NetworkMode: networkMode,
	}
}

func NewECSExecutionRole() *cloudformation.AWSIAMRole {
	return &cloudformation.AWSIAMRole{
		AssumeRolePolicyDocument: map[string]interface{}{
			"Statement": map[string]interface{}{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "ecs-tasks.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
			},
		},
		ManagedPolicyArns: []string{
			"arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
		},
	}
}
