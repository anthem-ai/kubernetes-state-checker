package deployments

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type inputs struct {
	valuesYaml string
	checkName  string
	namespace  string
}

type Results struct {
	DidPass bool
	Message string
}

// New New
func New(valuesYaml string, checkName string, namespace string) inputs {
	s := inputs{valuesYaml, checkName, namespace}

	return s
}

type deploymentStruct struct {
	Values struct {
		DeploymentName string `yaml:"deploymentName"`
		ChecksEnabled  struct {
			EnvironmentVars environmentVars `yaml:"environmentVars"`
		} `yaml:"checksEnabled"`
	} `yaml:"values"`
}

// type environmentVars map[string]string
type environmentVars []map[string]interface{}

func deploymentParse(valuesYaml string, v *deploymentStruct) error {

	err := yaml.Unmarshal([]byte(valuesYaml), &v)
	if err != nil {
		return errors.New(fmt.Sprintf("YAML Parse Error: %v", err))
	}

	if v.Values.DeploymentName == "" {
		return errors.New("Check values: no `DeploymentName` set")
	}

	return nil
}

// GeneralCheck GeneralCheck
func (i inputs) GeneralCheck(kubeClientSet kubernetes.Interface) Results {

	var values deploymentStruct

	// Set initial check results
	checkResult := Results{
		DidPass: false,
		Message: "",
	}

	didValuesParse := false

	err := deploymentParse(i.valuesYaml, &values)
	if err != nil {
		didValuesParse = true
		checkResult.Message = fmt.Sprintf("%v", err)
	}

	if !didValuesParse {

		// Get data from kubernetes
		deployment, err := kubeClientSet.AppsV1().Deployments(i.namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// Loop through all of the services found
		for _, aDeployment := range deployment.Items {
			if aDeployment.ObjectMeta.Name == values.Values.DeploymentName {
				fmt.Println("woot")
			}
		}

	}

	return checkResult
}
