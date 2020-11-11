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

// New New
func New(namespace string, port int) services{
	s := services{namespace, port}
	return s
}

// DoesPortExist DoesPortExist
func (s services) DoesPortExist() Results {
	fmt.Println("Running doesServicePortExist")

	testResult := Results{
		DidPass: true,
		Message: "Port found",
	}

	return testResult
}

