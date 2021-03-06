# kubernetes-state-checker

Kubernetes State Checker as the name suggests helps you check the state of your Kubernetes cluster.  

You might ask, with Kubernetes, don’t you declare the state you want and Kubernetes makes that happen, why do I have to check the state?  You would be correct but there are usually interactions from one Kubernetes resource with another or it might depend on an external or third party resource.  Anyone of these items might not be able to get into the state you wanted which can have downstream effects or effects on the environment.  Even if Kubernetes put the state of the cluster/application in the state you declared, sometimes what you declared could be wrong because someone changed the setting looking to make a fix which was not communicated to downstream dependencies and that could break the environment.

With a microservice architecture, multiple teams can contribute to one application environment.  Ensuring everything coming together and working correctly is often a challenge.  There are often the application expert(s) that know how all of these applications integrate with each other and are supposed to work for your environment and this person is relied on to debug integration problems.  This is a tedious task and relying on certain people to debug these types of issues makes them a bottleneck.

## How Does Kubernetes State Checker help?

Here is a real example that has been sanitized and made generic.  We have an environment with a bunch of microservices.  In this scenario, we will talk about 2 of them.  Microservice 1 was set to listen on port 5000 and Microservice 2 was set to connect to Microservice 1 on port 30001.  Based on the configuration each state of the deployments made it to it’s desired state but when Microservice 2 tried to connect to Microservice 1, that connection failed.

We just told you the exact problem but when this occurred, it was not apparent that the ports were set incorrectly.  The initial thing that alerted us that this environment was now broken was an e2e test that failed.  However, the e2e tests only tests from certain entry points into the system and can’t tell us why the system is broken.  This led to developers looking at application logs on the various services since they know how the call flows through the system.  From the logs, the developers were able to localize the problem to a few microservices but was not able to say exactly why the call is broken.  The next step was to loop in the infrastructure people to take a look.  The infrastructure people had to catch up on what was happening and with that information started to check various Kubernetes things.  After some tedious task of tracing out how the application was configured to what it was trying to communicate with, it was found that the port settings were off and one side of the port numbers would have to change.  

With Kubernetes State Checker, we will be able to declare these states on how the port numbers should be configured and then run a check to make sure it is in that state.  If it is not, like in this case, it would tell you that this particular port is not in the state that you said you wanted it in.  
This was essentially an integration problem between Microservice 1 and Microservice 2.  Microservice 1 was listening on one port but Microservice 2 thought it should connect to Microservice 1 on another port.  Who is correct?  This is an understanding between the two microservices on how they will connect to each other but nothing is enforcing or checking that.  Kubernetes State Checker can be that “check” or “enforcement”.


From a developers point of view, this gives them a tool to check the layers underneath the application to ensure that everything in those layers is set to what is expected.


From a DevOps/Infrastructure person’s point of view, this allows them to set up an expected state and enable other groups to check for that.  It also gives this group a tool where they can run to check the state which can help them eliminate what could possibly be wrong and look at other areas that this tool did not cover.

## Use cases

### Is an environment setup correctly?
If you had a series of checks that represents a correctly setup environment, then you can use those checks and run it against another environment to make sure everything is in place.  If something is not correctly set, then kubernetes-state-checker will output what is not correct.

### Feedback loop from dev to production
As the developer is developing the application, this person usually knows what the application needs in order for it to function properly and probably have found some hiccups while debugging this application.  The developer can create check(s) very similiar to how they can create unit tests on areas that are known to failures.  Then other teams that are managing other environments can uses these check(s) to make sure their environments are correctly setup.

The check(s) can flow from the other way as well.  Usually in a larger organization, the developers might not be the ones that are taking care of the production systems.  There might be another team for that.  As you run the application(s) in production, you will usually see other operational issues that are not an issue in dev.  This team can also write check(s) for these items so that these issue(s) are detected before it becomes a problem.  With these checks, the other teams including the developers can run it to make sure everything is good and that the assumptions from development to production are adhere to.

## Example usage

### Check Kubernetes service port

```
kubernetes-state-checker:
- type: doesServicePortExist
  name: Does microservice 1 have a kubernetes service with port 5000 exposed
  description: This checks if microservice 1 has a Kubernetes service with port 5000 exposed
  namespace: app
  # Input values for this specific check
  values:
    serviceName: microservice-1
    port: 5000
```

### Check environment in a pod

```
kubernetes-state-checker:
- type: doesEnvarExistInDeployment
  name: Check that the microservice 2 deployments has the correct envar for microservice 1
  description: The microservice 2 uses the "MICROSERVICE_1_HOST_PORT" envar to find microservice 1.  This checks to make sure that this envar is there and set to the correct value.
  Namespace: app
  # Input values for this specific check
  values:
    deploymentName: microservice-2
    envarKey: MICROSERVICE_1_HOST_PORT
    envarValue: microservice-1:5000
```

### Check if a port is open
Maybe even kube exec telnet/nc to test the connection


## Open discussions

### A more dynamic way to read configurations
During a peer review of this document, there was an idea put out to see if we can read in configurations in a more dynamic way so that these configurations don't have to be in more than one place.  Taking our our microservice 1 and microservice 2 example from above.  The actual ports and envars are defined in each services Helm values files.  Then with this test, we once again have to define what ports maps to what.  This means that the same information is in two places now.  When someone wants to update the port for microservice 1, they would have to update it in microservice 1's Helm values and then go into the kubernetes-state-checker's check config yaml and change the value in there as well.  This make repetitive and tedious amount of work.

We would like a way where we can tell this kubernetes-state-checker's check config yaml that here is the port that microservice-1 is listening on and here is the file and here is the envar that microservice-2 is using to reach that port.  These values should be the same.

#### Proposed solutions:

Be able to read from any yaml file with a `kubernetes-state-checker` section
We can point it to any yaml file and it will find and only the `kubernetes-state-checker` section.

```
fullnameOverride: &name "hos-core-authentication"

image:
  repository: 1234.dkr.ecr.us-west-2.amazonaws.com/hos/hos-core-authentication
  pullPolicy: Always
  tag: &tag dev

service:
  port: 20004
  targetPort: 20004

some-other-yaml:
  foo: bar

kubernetes-state-checker:
- type: doesServicePortExist
  name: Does hos-core-authentication have a kubernetes service with port 20004 exposed
  description: This checks if hos-core-authentication has a Kubernetes service with port 20004 exposed
  namespace: app
  # Input values for this specific check
  values:
    serviceName: hos-core-authentication
    port: 20004
```

It will ignore all sections and just use the information in the `kubernetes-state-checker` section.

# Running this in VScode
This has the configuration files that helps you bring up the enviroment you need to develop against this locally via a Docker container.

This is based on this example: https://github.com/microsoft/vscode-remote-try-go

You will have to install the `Visual Studio Code Remote - Containers extension`: https://code.visualstudio.com/docs/remote/containers#_getting-started

Once you install this, and restart VScode in this repository, it will ask if you want to open it in a container.  The first time you do this, it will take some time since it is building the container based on the Dockerfile in the `.devcontainer` directory.  After it builds the container, it will open this project inside of this container and you can code away as normal.

Install the go dependencies:
```
go get gopkg.in/yaml.v2
go get k8s.io/apimachinery/pkg/api/errors
go get k8s.io/apimachinery/pkg/apis/meta/v1
go get k8s.io/client-go/kubernetes
go get k8s.io/client-go/tools/clientcmd
go get k8s.io/client-go/util/homedir
go get github.com/evanphx/json-patch
go get k8s.io/kube-openapi/pkg/util/proto
go get github.com/olekukonko/tablewriter
```
## Set go path
```
export GOPATH=$GOPATH:$PWD
```

## Run all unit tests
```
go test ./...
```

## Kubeconfig
The `.devcontainer/devcontainer.json` file specifies a local mount from your `$HOME/.kube/config` into the `$HOME/.kube/config` inside the container.
