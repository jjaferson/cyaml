package resources

import "github.com/awslabs/goformation/cloudformation"

func NewECSCluster(name string) *cloudformation.AWSECSCluster {
	return &cloudformation.AWSECSCluster{
		ClusterName: name,
	}
}
