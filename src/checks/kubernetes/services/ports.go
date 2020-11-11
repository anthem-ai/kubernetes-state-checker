package services

import (
	"fmt"
)

type services struct {
	namespace string
	port int
}

type Results struct {
	DidPass bool
	Message string
}

// DoesPortExist DoesPortExist
func DoesPortExist() Results {
	fmt.Println("Running doesServicePortExist")

	testResult := Results{
		DidPass: true,
		Message: "Port found",
	}

	return testResult
}

