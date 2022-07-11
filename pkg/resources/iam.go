package resources

import "github.com/awslabs/goformation/cloudformation"

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
