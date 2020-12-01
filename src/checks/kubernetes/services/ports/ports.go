package ports

import (
	"context"
	"errors"
	"fmt"
	"log"

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

type doesPortExistStruct struct {
	Values struct {
		ServiceName string `yaml:"serviceName"`
		Port        int32  `yaml:"port"`
	} `yaml:"values"`
}

func doesPortExistParse(valuesYaml string, v *doesPortExistStruct) error {

	err := yaml.Unmarshal([]byte(valuesYaml), &v)
	if err != nil {
		return errors.New("YAML Parse Error: "+err)
	}

	if v.Values.ServiceName == "" {
		return errors.New("Check values: no `ServiceName` set")
	}

	if v.Values.Port < 1 || v.Values.Port > 65353 {
		return errors.New("Check values: invalid `Port` specified, allowed range (1 - 65353)")
	}

	return nil
}

// DoesPortExist DoesPortExist
func (i inputs) DoesPortExist() Results {

	var values doesPortExistStruct

	// Set initial check results
	checkResult := Results{
		DidPass: false,
		Message: "Port not found: " + fmt.Sprint(values.Values.Port),
	}

	didValuesParse := false

	err := doesPortExistParse(i.valuesYaml, &values)
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

	return checkResult
}
