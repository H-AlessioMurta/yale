# Yet Another Library Example
As a student developer, I happened to practice a lot using libraries almost always as representations.
Y.A.L.E.
The aim is to structure three microservices, put them in synchronous communication and exploit the Githubactions to build a Kubernetes cluster able to keep them in complete execution.

## Goals reached 
-Learned an entry-level knowledge about Go!
-Learned an entry-level knowledge about Grapqhl
-Learned an entry-level knowledg about Kubernetes

-Microservices are in CI/CD
## Fails&Bugs at the moment (to do list)
-No enviroment's key automation, the code is pretty forced because no configmap was settled on k8's cluster
-No testing in CI/CD's phase because the tests were not written with "unit test pattern", so they try to connect to cluster's db and sadly
-Cluster k8's not accessible from external sources.
-Notification svc does not build during github actions

### Design a solution of [microservices's assignment](https://github.com/sunnyvale-academy/ITS-ICT_Microservices/tree/master/assignments/01-Library_application)
