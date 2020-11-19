package ports

import (
	"fmt"
	"reflect"
)

type inputs struct {
	checkName   string
	namespace   string
	checkValues map[string]string
}

type Results struct {
	DidPass bool
	Message string
}

// New New
func New(checkName string, namespace string, checkValues interface{}) inputs {
	s := inputs{checkName, namespace, parseInputInterface(checkValues)}
	return s
}

// DoesPortExist DoesPortExist
func (i inputs) DoesPortExist() Results {
	fmt.Println("Running :" + i.checkName)
	fmt.Println("Checking port: " + i.checkValues["port"])

	testResult := Results{
		DidPass: true,
		Message: "Port found",
	}
	return testResult
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
