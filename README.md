[![logo](https://github.com/iboware/matriarch/raw/master/assets/matriarch.png]

# Matriarch
Matriarch is a CLI utility and a Kubernetes Operator to Deploy High Available and Scalable PostgreSQL Clusters

The PostgreSQL logo elephant is named "Slonik". The elephant herd is led by the oldest and largest female cow known as the **matriarch**. This was the inspiration for our projects name.

***Important: This is a work in progress***

Matriarch currently can create and manage **simple** High Available and Scalable Clusters of PostgreSQL based on [bitnami/bitnami-docker-postgresql-repmgr](http://github.com//bitnami/bitnami-docker-postgresql-repmgr) image via Matriarch CLI utility to deploy and manage the operator. It is also possible to create and manage clusters via YAML/JSON definitions without CLI:

**How to Install**
1. Download [**matriarch**](https://github.com/iboware/matriarch/releases/download/v0.3.7/matriarch)
2. Put it under any binary path.
4. Deploy operator via `matriarch init` to the active cluster in your `kubeconfig` file.
3. Start creating and managing PostgreSQL clusters.

**Matriarch Demonstration** (Init/Create/List/Delete/Scale operations)

[![asciicast](https://asciinema.org/a/351880.svg)](https://asciinema.org/a/351880)

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
