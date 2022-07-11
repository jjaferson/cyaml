# CYaml

CYaml is a CLI that generates [AWS CloudFormation](https://aws.amazon.com/cloudformation/) template to run containers on AWS from a simplified yaml file.

> **Note:** CYaml is still under development and is only capable of generating CloudFormation templates for ECS using Fargate in a pre-defined VPC.

## AWS Resources generated

* Application Load Balancer
  * Target Group
* IAM Role to run containers
* EC2 Security Group
* ECS Cluster
* ECS Services (Type: Fargate)
* ECS Task Definition

## How it works

CYaml abstracts the generation of AWS CloudFormation templates to get containers up and running on AWS. 

The CLI receives a YAML file with all bare minimum information required for setting up the containers and outputs the CF template.

## How to use

To use the CLI, you need to clone the repo and use go to build the binaries.

1) Clone the repo
2) Build the CLI using Go

```bash
$ go install github.com/jjaferson/cyaml
```
Once the CLI is installed you can run:

```bash
$ cyaml ecs --yaml <path-to-yaml-file> > <path-to-cloud-formation-file>
```

You can find an example of the YAML file [here](./templates/ecs.yaml)

### YAML File

TO DO

## Resources 

* https://github.com/awslabs/goformation
