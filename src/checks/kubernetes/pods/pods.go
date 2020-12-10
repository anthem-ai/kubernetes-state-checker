package pods

import (
	"context"
	"errors"
	"fmt"
	"regexp"

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

type podStruct struct {
	Values struct {
		ChecksEnabled struct {
			State []struct {
				PodName      string `yaml:"podName"`
				DesiredState string `yaml:"desiredState"`
			} `yaml:"state,omitempty"`
		} `yaml:"checksEnabled"`
	} `yaml:"values"`
}

func podParse(valuesYaml string, v *podStruct) error {

	err := yaml.Unmarshal([]byte(valuesYaml), &v)
	if err != nil {
		return errors.New(fmt.Sprintf("YAML Parse Error: %v", err))
	}

	// if v.Values.PodName == "" {
	// 	return errors.New("Check values: no `PodName` set")
	// }

	return nil
}

// GeneralCheck GeneralCheck
func (i inputs) GeneralCheck(kubeClientSet kubernetes.Interface) Results {

	var values podStruct

	// Set initial check results
	checkResult := Results{
		DidPass: false,
		Message: "",
	}

	didValuesParse := false

	err := podParse(i.valuesYaml, &values)
	if err != nil {
		didValuesParse = true
		checkResult.Message = fmt.Sprintf("%v", err)
	}

	if !didValuesParse {

		// Get data from kubernetes
		pods, err := kubeClientSet.CoreV1().Pods(i.namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		//
		// Check pod state
		//
		didFindError := false
		for _, inputPod := range values.Values.ChecksEnabled.State {

			didFindContainer := false

			for _, aPod := range pods.Items {

				// Match string for the input pod
				match, _ := regexp.MatchString(inputPod.PodName, aPod.ObjectMeta.Name)

				if match {

					didFindContainer = true

					if inputPod.DesiredState == string(aPod.Status.Phase) {
						checkResult.DidPass = true
						checkResult.Message += "* Pod " + aPod.ObjectMeta.Name + " is in " + inputPod.DesiredState + " state\n"
					} else {
						checkResult.Message += "* Pod " + aPod.ObjectMeta.Name + " is NOT in " + inputPod.DesiredState + " state\n"
						didFindError = true
					}
				}
			}

			if !didFindContainer {
				checkResult.Message += "* Did not find pod: " + inputPod.PodName + "\n"
				didFindError = true
			}
		}

		// Found one or more errors
		if didFindError {
			checkResult.DidPass = false
		}

	}

	return checkResult
}
