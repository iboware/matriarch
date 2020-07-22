# postgresql-operator
A Kubernetes Operator to Deploy High Available and Scalable PostgreSQL Clusters

***Important: This is a work in progress***

PostgreSQL Operator currently can create and manage **simple** High Available and Scalable Clusters of PostgreSQL based on bitnami/postgresql-repmgr image via **pgctl**, which is a CLI tool to manage deployed PostgreSQL Operators. It is also possible to create and manage clusters via YAML/JSON definitions:

**pgctl Demonstration** (Create/List/Delete operations)

[![asciicast](https://asciinema.org/a/tULY7wnMRyyTHojc79eKamDS2.svg)](https://asciinema.org/a/tULY7wnMRyyTHojc79eKamDS2)

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
