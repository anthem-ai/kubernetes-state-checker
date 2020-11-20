package checker

import (
	"checks/kubernetes/services/ports"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

type Check struct {
	ClientSet   *kubernetes.Clientset
	Ttype       string      `yaml:"ttype"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Namespace   string      `yaml:"namespace"`
	Values      interface{} `yaml:"values"`
}

type results struct {
	DidPass bool
	Message string
}

// New New
func New(clientSet *kubernetes.Clientset, ttype string, name string, description string, namespace string, values interface{}) Check {
	c := Check{clientSet, ttype, name, description, namespace, values}
	return c
}

// Run - runner
func (c Check) Run() results {
	fmt.Println("Starting runner ")

	var returnResults results

	switch c.Ttype {
	case "doesServicePortExist":
		check := ports.New(c.ClientSet, c.Name, c.Namespace, c.Values)
		r := check.DoesPortExist()

		returnResults = results{
			DidPass: r.DidPass,
			Message: r.Message,
		}
	}

	return returnResults
}
