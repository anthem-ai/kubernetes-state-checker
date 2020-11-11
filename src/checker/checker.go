package checker

import (
	"fmt"
	"checks/kubernetes/services"
)

type check struct {
	ttype string
	name string
	description string
	namespace string
}

type results struct {
	DidPass bool
	Message string
}

// New New
func New(ttype string, name string, description string, namespace string) check {
	c := check {ttype, name, description, namespace}
	return c
}

// Run - runner
func (c check) Run() results {
	fmt.Println("Starting runner ")

	var returnResults results

	switch c.ttype {
	case "doesServicePortExist":
		checkResults := services.DoesPortExist()

		returnResults = results{
			DidPass: checkResults.DidPass,
			Message: checkResults.Message,
		}
	}

	return returnResults
}
