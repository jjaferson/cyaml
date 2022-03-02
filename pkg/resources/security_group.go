package resources

import (
	"fmt"

	"github.com/awslabs/goformation/cloudformation"
)

const (
	accessFromAnyIP = "0.0.0.0/0"
)

func NewSecurityGroup(name, vpcId string, ports []int) *cloudformation.AWSEC2SecurityGroup {

	sgName := fmt.Sprintf("%s-sg", name)

	var ingress []cloudformation.AWSEC2SecurityGroup_Ingress
	for _, port := range ports {
		ingress = append(ingress, cloudformation.AWSEC2SecurityGroup_Ingress{
			CidrIp:   accessFromAnyIP,
			FromPort: port,
			ToPort:   port,
		})
	}

	// there are only support for creating security group at ingress level (incoming traffic)
	// any the engress traffic is allowed
	securityGroup := &cloudformation.AWSEC2SecurityGroup{
		VpcId:                vpcId,
		GroupName:            sgName,
		SecurityGroupIngress: ingress,
		SecurityGroupEgress: []cloudformation.AWSEC2SecurityGroup_Egress{
			{
				CidrIp:     accessFromAnyIP,
				IpProtocol: "-1",
			},
		},
	}
	return securityGroup

}
