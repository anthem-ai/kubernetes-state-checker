package checker

import (
	"checks/kubernetes/services"
	"fmt"
)

type Check struct {
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
func New(ttype string, name string, description string, namespace string, values interface{}) Check {
	c := Check{ttype, name, description, namespace, values}
	return c
}

// Run - runner
func (c Check) Run() results {
	fmt.Println("Starting runner ")

	var returnResults results

	switch c.Ttype {
	case "doesServicePortExist":
		t := services.New("app", 5000)
		r := t.DoesPortExist()

		returnResults = results{
			DidPass: r.DidPass,
			Message: r.Message,
		}
	}

	return returnResults
}
