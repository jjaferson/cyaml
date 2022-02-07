package ecs

import (
	"fmt"

	"github.com/jjaferson/cyaml/api"
	"github.com/jjaferson/cyaml/pkg/template"
	"github.com/jjaferson/cyaml/utils"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type ECSOptions struct {
	path string
}

//eg: https://github.com/cli/cli/blob/e17964cf0377110526106d4c3f4654b4d8e9c971/pkg/cmd/pr/checkout/checkout.go#L39

func NewECSTemplate() *cobra.Command {

	options := &ECSOptions{}

	cmd := &cobra.Command{
		Use:   "ecs",
		Short: "Creates ECS CloudFormation template",
		Long:  `Outputs CloudFormation template for ECS based on the defined yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ecsTemplateGenerate(options)
		},
	}

	cmd.PersistentFlags().StringVarP(&options.path, "yaml", "y", "", "Path to the yaml with the definition of the cluster")
	if err := cmd.MarkPersistentFlagRequired("yaml"); err != nil {
		panic(err)
	}

	return cmd
}

func ecsTemplateGenerate(options *ECSOptions) (err error) {
	ecsYamlDefinition, err := utils.ReadExternalResource(options.path)
	if err != nil {
		return
	}

	ecsDeployment := &api.ECSDeployment{}
	err = yaml.Unmarshal(ecsYamlDefinition, ecsDeployment)
	if err != nil {
		return
	}

	ecsTemplate := template.NewECSTemplateGenerator(ecsDeployment)
	cloudFormationtemplate, err := ecsTemplate.Generate()
	if err != nil {
		return
	}

	fmt.Print(cloudFormationtemplate)

	return nil
}
