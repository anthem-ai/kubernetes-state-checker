package ports

import (
	"context"
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

// DoesPortExist DoesPortExist
func (i inputs) DoesPortExist() Results {

	var values doesPortExistStruct

	err := yaml.Unmarshal([]byte(i.valuesYaml), &values)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	// Set initial check results
	checkResult := Results{
		DidPass: false,
		Message: "Port not found: " + fmt.Sprint(values.Values.Port),
	}

	// Run kube stuff
	services, err := kubeClientSet.CoreV1().Services(i.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, aService := range services.Items {

		if aService.ObjectMeta.Name == values.Values.ServiceName {
			fmt.Println("Found service: " + aService.ObjectMeta.Name)

			for _, port := range aService.Spec.Ports {
				if port.Port == values.Values.Port {
					fmt.Println("Found port: " + fmt.Sprint(values.Values.Port))

					checkResult.DidPass = true
					checkResult.Message = "Port found: " + fmt.Sprint(values.Values.Port)
				}
			}
		}
	}

	return checkResult
}
