# Yet Another Library Example
As a developer student, I practice a lot using libraries examples, almost always. This is the better name i can assign to my work!
The aim is to structure three microservices, put them in synchronous communication in a Kubernetes cluster able to keep them in execution.

## Personal Goals reached 
- [x] Learned an entry-level knowledge about Go!

- [x] Learned an entry-level knowledge about Grapqhl

- [x] Learned an entry-level knowledg aboute Kubernetes

- [x] Applying a design pattern: Saga Choreography 

- [x] Microservices are in CI/CD

## Fails&Bugs at the moment ( AKA Improvements to do list)
- [ ] No enviroment's keys automation, the code is pretty forced because no configmap was settled on k8's cluster

- [ ] No testing in CI/CD's phase because, the Go's tests weren't written with "unit test pattern", so they try to connect to cluster's db and sadly they fails during a Github actions phase.

- [ ] Cluster k8's not accessible from external sources like internet, Cluster's network was not properly planned.

- [ ] Notification svc does not build during github actions

- [ ] Refactoring of the code using scaffolding as Go's best pratices for avoid repetition

---

### Design a solution of [microservices's assignment](https://github.com/sunnyvale-academy/ITS-ICT_Microservices/tree/master/assignments/01-Library_application)
![Yale Cluster](https://github.com/H-AlessioMurta/yale/blob/main/K8s%20YALE.jpg)

I chose the **saga pattern**, putting a lot of responsibility on the borrowingsvc microservice as center of choreograpy of the logic.
It manages all crud operation of the entity defined by the borrowed model, and also requests and manages all CRUD operations on books and customers model.
The main benefit of the Saga Pattern is that it helps maintain data consistency across multiple services without tight coupling, spreading errors and messages between all the microservices involved.
Each microservices will properly log on stdout what it is doing, and who asked for.
If a fatal error is detected, an API REST call to notificationsvc with error description will be send in a kafka producer.
ELK stack is also implemented, Elastisearch, Logstash, Filebeat, Kibana pods are running in backgroud whatching over K8's cluster
### Booksvc & Customesvc
They are twins microservices, same logic and same implementation. Both of them handles http methods for operating CRUD on a Postgresql database.
They were written on GO language, I chose this language because of its high efficiency and its cosistency.

![Gorilla Logo](https://cloud-cdn.questionable.services/gorilla-icon-64.png)

Http handler was written with the open library [Gorilla Mux](https://github.com/gorilla/mux) instead of go's vanilla "http", because it seemed to me easier to read the code and manage multiple http methods.


<img src="https://external-preview.redd.it/SmsJqB8DdKq1FhsuBSAMN2rpZVEumG2wcBsHqKJEVK4.jpg?auto=webp&s=c2b78c143fe2f6e9e2c228db02c96ad88314e052" width="85">
[Postgresql's Driver](https://github.com/lib/pq)

### Borrowingsvc


<img src="https://avatars.githubusercontent.com/u/36954732?v=4" width="100">

Borrowingsvc is based on [gqlgen](https://github.com/99designs/gqlgen) a powerful, smart way to implement a Grapql's server starting from a [schema](https://github.com/H-AlessioMurta/yale/blob/main/borrowing/graph/schema.graphqls).
It generated all the infrastructure, leaving me to just implement logic resolvers in schema_resolvers.go.

<img src="https://github.com/mongodb/mongo-go-driver/raw/v1.8.2/etc/assets/mongo-gopher.png" width="85">

Mongo was a complicated choice to pair with Go, the [drivers](https://github.com/mongodb/mongo-go-driver) are quite recent, but they have some instability and don't always provide clarifying errors, golang is not compatible with versions higher than fourth to interface with docker images.
After many attempts, it worked.
In borrowingsvc's router module, i chose to use the vanilla http library instead of Gorilla Mux mainly to get experience with different forms and libraries. Here we transform grapql queries in http methods requests.

### Notificationsvc
It doesn't work at the moment, while the code seems to propose a sensible and coherent logic its implementation caused me a lot of problems.
I would have preferred that it was an independent microservice, and that autonomously read the logs on the stdout of mine microservices, so I tried to use a k8 client that asks with kube tail all logs, looking for the syntax  ** [fatal] ** to generate the message on kafka.
Failing this, I preferred to reduce the complexity by using an api rest communication again.
For the moment.

---

## Deploy on Kubernetes
![Road](https://miro.medium.com/max/873/1*NII9Htj87LjmNIa1PJzgCA.png)

For each service i wrote a proper Dockerfile, in github/workflow i used a githubaction for build from that Dockerfile's and push to my Docker hub repository.
In a directory called K8s you can find .yaml for uploading mine microservices and ELK, kafka, in a cluster using helm.
You need it install on your cluster.

<img src="https://dashboard.snapcraft.io/site_media/appmedia/2017/06/helm.png" width="85">

From here i will post the sequence of shell's commands, you can use it with shell terminal insde K8s directory:

```console
$ minikube start --memory 8192 --cpus 4
$ helm install --values dbBorrow.yaml borrowing-mongodb bitnami/mongodb
$ helm install --values dbCustomer.yaml customer-postgres bitnami/postgresql
$ helm install --values dbBook.yaml book-postgres bitnami/postgresql
$ kubectl apply -f customersvc.yaml
$ kubectl apply -f booksvc.yaml
$ kubectl apply -f borrowingsvc.yaml
$ helm install --values values.yaml elasticsearch elastic/elasticsearh
$ helm install --values kibana-values.yaml kibana elastic/kibana
$ helm install --values filebeat-values.yaml filebeat elastic/filebeat
$ helm install --values logstash-values.yaml logstash elastic/logstash
```
---
