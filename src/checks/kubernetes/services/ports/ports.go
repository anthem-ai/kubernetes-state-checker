package ports

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
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

	// Run kube stuff
	i.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	for {

		pods, err := i.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "example-xxxxx"
		_, err = i.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}

	testResult := Results{
		DidPass: true,
		Message: "Port found",
	}
	return testResult
}
