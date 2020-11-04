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


