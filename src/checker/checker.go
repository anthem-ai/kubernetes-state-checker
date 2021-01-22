package checker

import (
	"kubernetes-state-checker/src/checks/kubernetes/deployments"
	"kubernetes-state-checker/src/checks/kubernetes/pods"
	"kubernetes-state-checker/src/checks/kubernetes/services"

	"k8s.io/client-go/kubernetes"
)

var kubeClientSet *kubernetes.Clientset

type Check struct {
	valuesYaml  string      `yaml:"checkYaml"`
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
func New(valuesYaml string, clientSet *kubernetes.Clientset, ttype string, name string, description string, namespace string, values interface{}) Check {
	c := Check{valuesYaml, ttype, name, description, namespace, values}
	kubeClientSet = clientSet
	return c
}

// Run - runner
func (c Check) Run() results {

	var returnResults results

	switch c.Ttype {

	case "serviceChecks":
		check := services.New(c.valuesYaml, c.Name, c.Namespace)
		r := check.GeneralCheck(kubeClientSet)

		returnResults = results{
			DidPass: r.DidPass,
			Message: r.Message,
		}
	case "deploymentChecks":
		check := deployments.New(c.valuesYaml, c.Name, c.Namespace)
		r := check.GeneralCheck(kubeClientSet)

		returnResults = results{
			DidPass: r.DidPass,
			Message: r.Message,
		}
	case "podChecks":
		check := pods.New(c.valuesYaml, c.Name, c.Namespace)
		r := check.GeneralCheck(kubeClientSet)

		returnResults = results{
			DidPass: r.DidPass,
			Message: r.Message,
		}
	}

	return returnResults
}
