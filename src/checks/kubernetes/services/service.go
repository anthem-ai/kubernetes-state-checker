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
		DidPass: true,
		Message: "",
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

		// Loop through all of the services found
		for _, aService := range services.Items {

			// Find the service with the name we are interested in
			if aService.ObjectMeta.Name == values.Values.ServiceName {
				//
				// Run only enabled checks
				//

				if values.Values.ChecksEnabled.ClusterIP {
					if aService.Spec.ClusterIP == "" {
						checkResult.DidPass = false
						checkResult.Message += "* No ClusterIP Found\n"
					} else {
						checkResult.Message += "* ClusterIP Found\n"
					}
				}

				if values.Values.ChecksEnabled.Endpoints {
					// Look to see if the endpoints for this service exists or not
					endpoints, err := kubeClientSet.CoreV1().Endpoints(i.namespace).List(context.TODO(), metav1.ListOptions{})
					if err != nil {
						panic(err.Error())
					}

					for _, anEndpoint := range endpoints.Items {
						if anEndpoint.ObjectMeta.Name == values.Values.ServiceName {
							if len(anEndpoint.Subsets) == 1 {
								if len(anEndpoint.Subsets[0].Addresses) > 0 {
									for _, anAddress := range anEndpoint.Subsets[0].Addresses {
										if anAddress.IP != "" {
											checkResult.Message += "* Endpoint found: " + anAddress.IP + "\n"
										} else {
											checkResult.DidPass = false
											checkResult.Message += "* No Endpoint found in the Subsets[0].Addresses[x].IP field\n"
										}
									}
								} else {
									checkResult.DidPass = false
									checkResult.Message += "* No Endpoint found in the Subsets[0].Addresses list\n"
								}
							} else {
								checkResult.DidPass = false
								checkResult.Message += "* No Endpoint found in the subsets list\n"
							}
						}
					}
				}

				if values.Values.ChecksEnabled.HostPort {
					// TBD
				}

				if values.Values.ChecksEnabled.Ports {
					for _, port := range aService.Spec.Ports {
						if port.Port != values.Values.Port {

							checkResult.DidPass = false
							checkResult.Message += "* Port NOT found: " + fmt.Sprint(values.Values.Port) + "\n"
						} else {
							checkResult.Message += "* Port found: " + fmt.Sprint(values.Values.Port) + "\n"
						}
					}
				}
			}
		}

	}

	return checkResult
}
