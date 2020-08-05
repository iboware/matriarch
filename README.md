# Matriarch
Matriarch is a CLI tool and a Kubernetes Operator to Deploy High Available and Scalable PostgreSQL Clusters

The elephant herd is led by the oldest and largest female cow known as the **matriarch**. This was the inspiration for our projects name.

***Important: This is a work in progress***

Matriarch currently can create and manage **simple** High Available and Scalable Clusters of PostgreSQL based on [bitnami/bitnami-docker-postgresql-repmgr](http://github.com//bitnami/bitnami-docker-postgresql-repmgr) image via Matriarch CLI utility to deploy and manage the operator. It is also possible to create and manage clusters via YAML/JSON definitions without CLI:

**Matriarch Demonstration** (Init/Create/List/Delete/Scale operations)

[![asciicast](https://asciinema.org/a/tULY7wnMRyyTHojc79eKamDS2.svg)](https://asciinema.org/a/tULY7wnMRyyTHojc79eKamDS2)
ls
**Example Custom Resource to create a cluster**

*YAML:*
```yaml
apiVersion: database.iboware.com/v1alpha1
kind: PostgreSQL
metadata:
  name: mycluster
spec:
  disksize: 8Gi
  replicas: 3
  postgrespassword: verysecurepassword
  repmgrpassword: verysecurepassword2
  namespace: mynamespace
```
*JSON:*
```json
{
  "apiVersion": "database.iboware.com/v1alpha1",
  "kind": "PostgreSQL",
  "metadata": {
    "name": "mycluster"
  },
  "spec": {
    "disksize": "8Gi",
    "replicas": 3,
    "postgrespassword": "verysecurepassword",
    "repmgrpassword": "verysecurepassword2",
    "namespace": "mynamespace"
  }
}
```
