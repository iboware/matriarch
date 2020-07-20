# postgresql-operator
A Kubernetes Operator to Deploy High Available and Scalable PostgreSQL Clusters
***Important: This is a work in progress***

PostgreSQL Operator currently can create and manage **simple** High Available and Scalable Clusters of PostgreSQL based on bitnami/postgresql-repmgr image via PGCTL, which is a CLI tool to manage deployed operator. It is also possible to create and manage clusters via YAML/JSON definitions:

**PGCTL Demonstration**

![alt text](https://iboware.com/assets/img/pgctl-demo.gif "PGCTL Demonstration")

**Example Custom Resource to create a cluster**

```yaml
apiVersion: database.iboware.com/v1alpha1
kind: PostgreSQL
metadata:
  name: mycluster
  namespace: mynamespace
spec:
  disksize: 8Gi
  replicas: 3
  postgrespassword: verysecurepassword
  repmgrpassword: verysecurepassword2
```
