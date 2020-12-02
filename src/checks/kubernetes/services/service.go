package services

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var kubeClientSet *kubernetes.Clientset

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
func New(valuesYaml string, clientSet *kubernetes.Clientset, checkName string, namespace string) inputs {
	s := inputs{valuesYaml, checkName, namespace}
	kubeClientSet = clientSet

	return s
}

type serviceStruct struct {
	Values struct {
		ServiceName   string `yaml:"serviceName"`
		Port          int32  `yaml:"port"`
		ChecksEnabled struct {
			Ports     bool `yaml:"ports"`
			Endpoints bool `yaml:"endpoints"`
			ClusterIP bool `yaml:"clusterIP"`
			HostPort  bool `yaml:"hostPort"`
		} `yaml:"checksEnabled"`
	} `yaml:"values"`
}

func serviceParse(valuesYaml string, v *serviceStruct) error {

	err := yaml.Unmarshal([]byte(valuesYaml), &v)
	if err != nil {
		return errors.New(fmt.Sprintf("YAML Parse Error: %v", err))
	}

	if v.Values.ServiceName == "" {
		return errors.New("Check values: no `ServiceName` set")
	}

	if v.Values.Port < 1 || v.Values.Port > 65353 {
		return errors.New("Check values: invalid `Port` specified, allowed range (1 - 65353)")
	}

	return nil
}

// GeneralCheck GeneralCheck
func (i inputs) GeneralCheck() Results {

	var values serviceStruct

	// Set initial check results
	checkResult := Results{
		DidPass: false,
		Message: "Port not found: " + fmt.Sprint(values.Values.Port),
	}

	didValuesParse := false

	err := serviceParse(i.valuesYaml, &values)
	if err != nil {
		didValuesParse = true
		checkResult.Message = fmt.Sprintf("%v", err)
	}

	if !didValuesParse {
		// Run kube stuff
		services, err := kubeClientSet.CoreV1().Services(i.namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		if values.Values.ChecksEnabled.ClusterIP {
			fmt.Println("fooooooo")
		}

		if values.Values.ChecksEnabled.Endpoints {
			fmt.Println("fooooooo")
		}

		if values.Values.ChecksEnabled.HostPort {
			fmt.Println("fooooooo")
		}

		if values.Values.ChecksEnabled.Ports {
			for _, aService := range services.Items {

				if aService.ObjectMeta.Name == values.Values.ServiceName {

					for _, port := range aService.Spec.Ports {
						if port.Port == values.Values.Port {

							checkResult.DidPass = true
							checkResult.Message = "Port found: " + fmt.Sprint(values.Values.Port)
						}
					}
				}
			}
		}

	}

	return checkResult
}
