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
	c := checker.New(
		 "doesServicePortExist",
		  "Does microservice 1 have a kubernetes service with port 5000 exposed",
		  "This checks if microservice 1 has a Kubernetes service with port 5000 exposed",
		  "app",
	)
	results := c.Run()

	fmt.Println(fmt.Sprintf(`+----------------------------------------------------------------------+------+-------------+	
| Name                                                                 | Pass | Message     |
+----------------------------------------------------------------------+------+-------------+
| %s | %s | %s  |
+----------------------------------------------------------------------+------+-------------+`,
		name, strconv.FormatBool(results.DidPass), results.Message))

}
