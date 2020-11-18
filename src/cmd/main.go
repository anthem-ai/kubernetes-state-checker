package main

import (
	"checker"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	Some                   string          `yaml:"some"`
	Random                 int64           `yaml:"random"`
	KubernetesStateChecker []checker.Check `yaml:"kubernetes-state-checker"`
}

// type KubernetesStateChecker struct {
// 	Ttype       string      `yaml:"ttype"`
// 	Name        string      `yaml:"name"`
// 	Description string      `yaml:"description"`
// 	Namespace   string      `yaml:"namespace"`
// 	Values      interface{} `yaml:"values"`
// }

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {

	// Get kubeconfig

	// Get input yaml with checks
	var c conf
	c.getConf()

	for _, aCheck := range c.KubernetesStateChecker {

		// Execute the check runner
		chk := checker.New(
			aCheck.Ttype,
			aCheck.Name,
			aCheck.Description,
			aCheck.Namespace,
			aCheck.Values,
		)
		results := chk.Run()

		fmt.Println(fmt.Sprintf(`+----------------------+----------------------------------------------------------------------+------+-------------+
	| Test Type          | Name                                                                 | Pass | Message     |
	+----------------------+----------------------------------------------------------------------+------+-------------+
	| %s | %s | %s | %s  |
	+----------------------+----------------------------------------------------------------------+------+-------------+`,
			aCheck.Ttype, c.KubernetesStateChecker[0].Name, strconv.FormatBool(results.DidPass), results.Message))

	}

}
