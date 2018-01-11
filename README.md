# kaudit
Auditing tool for resources in Kubernetes.


# App Def Working Group

The App Def working group has develop a guide line for
labels and annotations [here](https://docs.google.com/document/d/1EVy0wRJRm5nogkHl38fNKbFrhERmSL_CLNE4cxcsc_M/edit#).

This project attempts to do two things:


## JSON Spec

The `app-def.json` file in this repo defines the
[JSON Schema](http://json-schema.org/) for labels 
and annotations.

## Audit Tool

The `kaudit` tool accepts a JSON Schema config file
and validates all objects in the workload API adhere
adhere to the schema.



# Usage

## Installation
```bash
$ go install github.com/runyontr/kaudit
```

## Deploy Samples
Execute the following from the command line to deploy two different deployments.  The deployment
`foo` are configured with the appropriate labels and annotations, where `bar` is missing all of the 
labels and annotations 

### Kubernetes 1.9.0+
```bash
$ kubectl apply -f ./deployments/1.9.0/
```

### Kubernetes <1.9.0
```bash
$ kubectl apply -f ./deployments/1.8.0/
```

## Validate
Validate there are two deployments
```bash
$ kubectl get deployments
NAME             DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
bar-deployment   3         3         3            3           35m
foo-deployment   1         1         1            1           36m
```

running the `kaudit` command should demonstrate which labels and annotations are missing from which 
applications:

```bash
$ kaudit --spec app-def.json 
Deployment: bar-deploymentThe document is not valid. see errors :
- app.kubernetes.io/name: app.kubernetes.io/name is required
- app.kubernetes.io/version: app.kubernetes.io/version is required
- app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
- app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
- app.kubernetes.io/component: app.kubernetes.io/component is required
- app.kubernetes.io/manager: app.kubernetes.io/manager is required
- app.kubernetes.io/usage: app.kubernetes.io/usage is required
- app.kubernetes.io/url: app.kubernetes.io/url is required
Deployment foo-deployment is valid
```



# Future Work

Include all objects in Workloads API:  Currently
only Deployments are searched

Dynamic resource definitions:
* Query against [ServerResourcesForGroupVersion](https://github.com/kubernetes/client-go/blob/master/discovery/discovery_client.go#L75)
to obtain all resource types inside of `apps/v1` ( or other group).  Use the APIResource object
to query for all instances of this resource and validate each instance
