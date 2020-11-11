package main

import (
	"checker"
	"fmt"
	"strconv"
)

func main() {

	// Get kubeconfig

	// Get input yaml with checks

	// Execute the check runner
	ttype := "doesServicePortExist"
	name := "Does microservice 1 have a kubernetes service with port 5000 exposed"
	description := "This checks if microservice 1 has a Kubernetes service with port 5000 exposed"
	namespace := "app"
	c := checker.New(ttype, name, description, namespace)
	results := c.Run()

	fmt.Println(fmt.Sprintf(`+----------------------------------------------------------------------+------+-------------+	
| Name                                                                 | Pass | Message     |
+----------------------------------------------------------------------+------+-------------+
| %s | %s | %s  |
+----------------------------------------------------------------------+------+-------------+`,
name, strconv.FormatBool(results.DidPass), results.Message))

}
