- [Introduction](#introduction)
- [How to Install](#how-to-install)
  - [Prerequisites](#prerequisites)
  - [Steps](#steps)
- [How to Use](#how-to-use)
  - [Matriarch CLI](#matriarch-cli)
  - [Manualy via kubectl](#manualy-via-kubectl)
- [How to Build & Develop the Operator](#how-to-build--develop-the-operator)
  - [Prerequisites](#prerequisites-1)
  - [Steps for Operator](#steps-for-operator)
  - [Steps for Matriarch CLI](#steps-for-matriarch-cli)
- [How to Uninstall](#how-to-uninstall)

# Introduction
![logo](https://github.com/iboware/matriarch/raw/master/assets/matriarch128.png "")

Matriarch is a CLI utility and a Kubernetes Operator to Deploy High Available and Scalable PostgreSQL Clusters

The PostgreSQL logo elephant is named "Slonik". The elephant herd is led by the oldest and the largest female cow known as the **matriarch**. This was the inspiration for our projects name.

***Important: This is a work in progress, not suitable for Production use!***

Matriarch currently can create and manage **simple** High Available and Scalable Clusters of PostgreSQL based on [bitnami/bitnami-docker-postgresql-repmgr](http://github.com//bitnami/bitnami-docker-postgresql-repmgr) image via Matriarch CLI utility. 

# How to Install
## Prerequisites 
* Access to a Kubernetes v1.11.3+ cluster
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.11.3+
* Linux (Currently only tested under Linux)
## Steps
1. Download [**matriarch**](https://github.com/iboware/matriarch/releases/download/v0.4.0/matriarch)
2. Put it under any binary path. (eg. /usr/local/bin)
4. Deploy operator via `matriarch init` to the active cluster in your `kubeconfig` file.
3. Start creating and managing PostgreSQL clusters.

**Note**: Matriarch currently uses kubectl config to identify the Kubernetes cluster(s). Kubernetes cluster(s) used should be defined in kubectl config file. To learn how to configure kubectl check [here](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/).

# How to Use
## Matriarch CLI

Init and Create operations will be demonstrated in the video.

[![asciicast](https://asciinema.org/a/351880.svg)](https://asciinema.org/a/351880)

## Manualy via kubectl
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
  pgpool: false
  postgrespassword: verysecurepassword
  repmgrpassword: verysecurepassword
  pgpoolpassword: verysecurepassword
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
    "pgpool": false,
    "postgrespassword": "verysecurepassword",
    "repmgrpassword": "verysecurepassword",
    "pgpoolpassword": "verysecurepassword",
    "namespace": "mynamespace"
  }
}
```

# How to Build & Develop the Operator
## Prerequisites 
* git
* go version v1.13+.
* docker version 17.03+.
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.11.3+
* Access to a Kubernetes v1.11.3+ cluster
* Install [operator-sdk](https://sdk.operatorframework.io/docs/installation/install-operator-sdk/#install-from-github-release) and its prequisites.
## Steps for Operator
1. Clone the repository to your local path `$ git clone https://github.com/iboware/matriarch.git`
2. Install CRD to Kubernetes Cluster. `$ make install`
3. Run operator locally without deploying to cluster. `$ make run ENABLE_WEBHOOKS=false`
   
For more advanced scenarios check [Operator SDK](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/#build-and-run-the-operator) documentations.

## Steps for Matriarch CLI
1. Change directory to matriarch `$ cd matriarch`
2. Build the source code with go. `$ go build -o ./bin`
3. Run the binary `$ ./bin/matriarch`


We recommend using [Visual Studio Code](https://code.visualstudio.com/) to test and debug both projects. We also provide launch configuration files under **vscode** folder. Copy the launch.json file under **.vscode** folder in the project, before opening the folder with Visual Studio Code.


# How to Uninstall
Simply remove the CRD from the kubernetes cluster via Kubectl.

```bash
$ kubectl delete -f https://github.com/iboware/matriarch/releases/download/v0.4.0/postgresql-operator.crd.yaml
```
**Important**: It will also **remove** all deployed PostgreSQL clusters! This behavior will be changed in the future releases.