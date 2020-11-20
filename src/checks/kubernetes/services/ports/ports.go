package ports

import (
	"context"
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type inputs struct {
	clientSet   *kubernetes.Clientset
	checkName   string
	namespace   string
	checkValues map[string]string
}

type Results struct {
	DidPass bool
	Message string
}

// New New
func New(clientSet *kubernetes.Clientset, checkName string, namespace string, checkValues interface{}) inputs {
	s := inputs{clientSet, checkName, namespace, parseInputInterface(checkValues)}
	return s
}

func parseInputInterface(m interface{}) map[string]string {
	var extractedValues = make(map[string]string)

	v := reflect.ValueOf(m)

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			fmt.Println(key.Interface(), "::", strct.Interface())
			extractedValues[fmt.Sprintf(key.Interface().(string))] = fmt.Sprintf(string(strct.Interface().(string)))
		}
	}

	return extractedValues
}

// DoesPortExist DoesPortExist
func (i inputs) DoesPortExist() Results {
	fmt.Println("Running :" + i.checkName)
	fmt.Println("Checking port: " + i.checkValues["port"])

	// Set initial check results
	checkResult := Results{
		DidPass: false,
		Message: "Port not found",
	}

	// Run kube stuff
	services, err := i.clientSet.CoreV1().Services(i.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, aService := range services.Items {

		if aService.ObjectMeta.Name == i.checkValues["serviceName"] {
			fmt.Println("Found service: " + aService.ObjectMeta.Name)

			for _, port := range aService.Spec.Ports {
				fmt.Println(fmt.Sprint(port.Port) + " : " + string(i.checkValues["port"]))
				if fmt.Sprint(port.Port) == string(i.checkValues["port"]) {
					fmt.Println("Found port: " + i.checkValues["port"])

					checkResult.DidPass = true
					checkResult.Message = "Port found"
				}
			}
		}
	}

	return checkResult
}
