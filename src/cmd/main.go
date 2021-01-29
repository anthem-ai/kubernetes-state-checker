package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/anthem-ai/kubernetes-state-checker/src/checker"

	"github.com/olekukonko/tablewriter"
	yaml "gopkg.in/yaml.v2"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeNamespace string
var kscConfig string

type conf struct {
	Some                   string          `yaml:"some"`
	Random                 int64           `yaml:"random"`
	KubernetesStateChecker []checker.Check `yaml:"kubernetes-state-checker"`
}

func (c *conf) getConf() *conf {

	configFileLocation := ""

	// Input config from CLI flags
	if kscConfig != "" {
		configFileLocation = kscConfig
	} else {
		// Input config file from environment variable
		configFileLocation = os.Getenv("KSC_CONFIG")

		if configFileLocation == "" {
			fmt.Println("ERROR: You must set the environment variable KSC_CONFIG which points to the input config file.")
			os.Exit(1)
		}
	}

	yamlFile, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func init() {
	// kubernetes namespace override
	flag.StringVar(&kubeNamespace, "namespace", "", "(optional) to override the namespace value in the input config file")
	flag.StringVar(&kubeNamespace, "n", "", "(optional) to override the namespace value in the input config file")

	// kubernetes-state-checker input config via cli flag
	flag.StringVar(&kscConfig, "config", "", "(optional) a flag to set the config file.  Will override the environment variable KSC_CONFIG")
	flag.StringVar(&kscConfig, "c", "", "(optional) a flag to set the config file.  Will override the environment variable KSC_CONFIG")
}

func main() {

	// Get kubeconfig
	// kubeconfig setup example: https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// Parse all flags
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get input yaml with checks
	var c conf
	c.getConf()

	// Output data
	outpuData := [][]string{}

	for _, aCheck := range c.KubernetesStateChecker {

		// convert to yaml
		valuesYaml, err := yaml.Marshal(aCheck)
		if err != nil {
			panic(err.Error())
		}

		// Use namespace override or not
		namespace := aCheck.Namespace
		if kubeNamespace != "" {
			namespace = kubeNamespace
			fmt.Println("Setting from input flag")
		}


		// Execute the check runner
		chk := checker.New(
			BytesToString(valuesYaml),
			clientset,
			aCheck.Ttype,
			aCheck.Name,
			aCheck.Description,
			namespace,
			aCheck.Values,
		)
		results := chk.Run()

		outpuData = append(outpuData, []string{aCheck.Ttype, c.KubernetesStateChecker[0].Name, strconv.FormatBool(results.DidPass), results.Message})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Test Type", "Name", "Status", "Message"})
	table.SetRowLine(true)
	table.SetReflowDuringAutoWrap(false)
	table.SetColWidth(100)

	for _, v := range outpuData {
		table.Append(v)
	}
	table.Render() // Send output

}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
