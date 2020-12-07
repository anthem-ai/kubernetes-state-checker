How to add a check
===================
This document walks you through on how to add a check into this application.

This will use the `deploymentChecks` check as an example.

## Create the input data structure
The input data structure is the yaml that the user of this check will input into the application and the application will use this information to perform the check.

This is the entire check information


```yaml
- ttype: deploymentChecks
  name: Checks for various aspects of a deployment
  description: Allows you to check for various aspects of a deployment
  namespace: hos-m2
  # Input values for this specific check
  values:
    # The service name to act on
    deploymentName: healthapp-caregaps
    checksEnabled:
      # Check if a set of environment variables are present
      environmentVars:
        # Specify the container in the pod that should contain this set of envars
        containers:
          # A specific named container
          - name: container1
            envars:
            - foo: bar
            - foo2: bar2
          # A wildcard to find this envar in any container in the deployment
          - name: "*"
            envars:
            - foo: bar
            - foo2: bar2            
      # Check if a configmap mount is set
      configMapMounts:
      - foo
      - bar
```

This section is mandatory and the structure of it cannot be changed:
```yaml
- ttype: deploymentChecks
  name: Checks for various aspects of a deployment
  description: Allows you to check for various aspects of a deployment
  namespace: hos-m2
  # Input values for this specific check
  values: {}
```

The section under the `values: {}` is unique for all checks.  The `values` section is all unique to the check itself depending on what type of input it requires from the user to perform the check.

## Add new check to `checker.go`
In the file `checker.go`, there is a function `Run()`.  This function holds a list of the checks and executes it.  You will have to add your new check function into the `switch` statement.

## Adding Check files
You can mostly do TDD here with a little caveat.  I find the k8s API data structure to be pretty complex and confusing.  There are v1, betas, extentions and I don't know it well enough or really want to remember the entire stucture on where things live.  The way I am doing this right now is to scaffold the check enough so that I can run the app and get back the response for the k8s resource I want.  Then from there I can put vscode into debug mode and just look at the data structure and construct my unit tests k8s fake data structure so that I can make a unit test.

Here is what I do (which is just easy for me).  Feel free to do or suggest another method.

Ill add in the deployment check file: `src/checks/kubernetes/deployments/deployment.go`

I will scaffold this to a point where I can get back the k8s deployment object.  Basically get it to run and query the k8s api for info:

```go
deployment, err := kubeClientSet.AppsV1().Deployments(i.namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
```
 
You can also take a look at this commit for an example: https://github.com/anthem-ai/kubernetes-state-checker/commit/a80bbb1bcfebd0140fc05c0e1b52125550481cba#diff-54c8b8aaa53e128cb37fffc9f8dac55169dd1aef27d8800b0ab9ebe2ffc5eaf0R79

If you run this in the debugger and put the break point to after this section runs, you can see the entire data structure of the k8s resource and it's type: 

![k8s api spec](./img/k8s-api-spec.png "k8s-api-spec")

This made it a lot easier for me when constructing the unit test.

## Adding the unit test
Vscode can help you add in the unit test and get you going.  Open the `deployment.go` file or the file you want to generate a unit test file for.  Then press `F1` then select or type in: `Go: Generate Unit Test For File`.  A new file for the file you were just on will be created with the unit tests in it.  It will create a unit test for every funtion in the source file.

For the deployment, we are mostly interested in the `Test_inputs_GeneralCheck` unit test.

The simpliest thing you can do is to just add in the input variables and the expected output for this function.  There is a comment in each unit test `// TODO: Add test cases.`.  This is where you will add in the input and expected outputs.

Refer to what is currently in this project for further expamples on how to add in the data structure.

Now you can fill out the k8s fake resource and what you expect the output to be.  Then you can do a TDD workflow when developing this new check.
