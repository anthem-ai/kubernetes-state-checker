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
			Containers []struct {
				Name string `yaml:"name"`
				Env  []struct {
					Name  string `yaml:"name,omitempty"`
					Value string `yaml:"value,omitempty"`
				} `yaml:"env,omitempty"`
			} `yaml:"containers,omitempty"`
		} `yaml:"checksEnabled"`
	} `yaml:"values"`
}

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

			// Find the deployment we want to look at
			if aDeployment.ObjectMeta.Name == values.Values.DeploymentName {

				//
				// Check for envars
				//
				// Number of containers to check
				numberOfContainers := len(values.Values.ChecksEnabled.Containers)
				numberOfContainersEnvarsFound := 0

				// Loop through the containers in the input values
				for _, inputContainer := range values.Values.ChecksEnabled.Containers {

					// Find the container in the Deployment
					for _, container := range aDeployment.Spec.Template.Spec.Containers {
						if inputContainer.Name == container.Name {

							// The number of envars that should exist
							numberOfEnvars := len(inputContainer.Env)
							numberOfEnvarsFound := 0

							// Find the envars in the k8s pod's containers
							for _, inputContainerEnv := range inputContainer.Env {
								for _, k8sDeploymentEnv := range container.Env {
									if inputContainerEnv.Name == k8sDeploymentEnv.Name &&
										inputContainerEnv.Value == k8sDeploymentEnv.Value {
										// Found the envar
										numberOfEnvarsFound++
									}
								}
							}

							if numberOfEnvars > 0 {
								if numberOfEnvars == numberOfEnvarsFound {
									// Found the correct amount of envars
									numberOfContainersEnvarsFound++
									checkResult.Message += "* Found all envars in Deployment: " + values.Values.DeploymentName + " | container: " + container.Name + "\n"
								}
							}
						}
					}
				}

				if numberOfContainers == numberOfContainersEnvarsFound {
					// Found the envars in all of the input check's envar(s)
					checkResult.DidPass = true
				} else {
					checkResult.DidPass = false
				}

				//
				// Check for the the containers that has the `containerMustBePresent` flag set to true
				//
				if len(values.Values.ChecksEnabled.Containers) > 0 {

					didFindAllContainers := true

					// Find each container in the deployment based on the user input
					for _, inputContainer := range values.Values.ChecksEnabled.Containers {

						didFindContainer := false

						// Search for the user inputted container in the deployment
						for _, container := range aDeployment.Spec.Template.Spec.Containers {
							if inputContainer.Name == container.Name {
								didFindContainer = true
							}
						}

						if !didFindContainer {
							didFindAllContainers = false
						}
					}

					if didFindAllContainers {
						checkResult.DidPass = true
						checkResult.Message += "* Found the correct number of containers in this deployment\n"
					} else {
						checkResult.DidPass = false
					}
				}
			}
		}

	}

	return checkResult
}
