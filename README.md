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
$ go get github.com/runyontr/kaudit
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


deployments: 
bar-deployment:	Errors:
	 - app.kubernetes.io/name: app.kubernetes.io/name is required
	 - app.kubernetes.io/version: app.kubernetes.io/version is required
	 - app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
	 - app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
	 - app.kubernetes.io/component: app.kubernetes.io/component is required
	 - app.kubernetes.io/manager: app.kubernetes.io/manager is required
	 - app.kubernetes.io/usage: app.kubernetes.io/usage is required
	 - app.kubernetes.io/url: app.kubernetes.io/url is required
foo-deployment:	Ok!


replicasets: 
bar-deployment-589f55cb9d:	Errors:
	 - app.kubernetes.io/name: app.kubernetes.io/name is required
	 - app.kubernetes.io/version: app.kubernetes.io/version is required
	 - app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
	 - app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
	 - app.kubernetes.io/component: app.kubernetes.io/component is required
	 - app.kubernetes.io/manager: app.kubernetes.io/manager is required
	 - app.kubernetes.io/usage: app.kubernetes.io/usage is required
	 - app.kubernetes.io/url: app.kubernetes.io/url is required
foo-deployment-57fc95945b:	Ok!

```

To compare against `v1` resources (e.g. services, pods) use the following:

```bash
$ kaudit --spec app-def.json --version v1


pods: 
bar-deployment-589f55cb9d-2zbdz:	Errors:
	 - app.kubernetes.io/name: app.kubernetes.io/name is required
	 - app.kubernetes.io/version: app.kubernetes.io/version is required
	 - app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
	 - app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
	 - app.kubernetes.io/component: app.kubernetes.io/component is required
	 - app.kubernetes.io/manager: app.kubernetes.io/manager is required
bar-deployment-589f55cb9d-6msgt:	Errors:
	 - app.kubernetes.io/name: app.kubernetes.io/name is required
	 - app.kubernetes.io/version: app.kubernetes.io/version is required
	 - app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
	 - app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
	 - app.kubernetes.io/component: app.kubernetes.io/component is required
	 - app.kubernetes.io/manager: app.kubernetes.io/manager is required
bar-deployment-589f55cb9d-8jw56:	Errors:
	 - app.kubernetes.io/name: app.kubernetes.io/name is required
	 - app.kubernetes.io/version: app.kubernetes.io/version is required
	 - app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
	 - app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
	 - app.kubernetes.io/component: app.kubernetes.io/component is required
	 - app.kubernetes.io/manager: app.kubernetes.io/manager is required
foo-deployment-57fc95945b-hf2p6:	Ok!


services: 
kubernetes:	Errors:
	 - app.kubernetes.io/name: app.kubernetes.io/name is required
	 - app.kubernetes.io/version: app.kubernetes.io/version is required
	 - app.kubernetes.io/component-name: app.kubernetes.io/component-name is required
	 - app.kubernetes.io/component-version: app.kubernetes.io/component-version is required
	 - app.kubernetes.io/component: app.kubernetes.io/component is required
	 - app.kubernetes.io/manager: app.kubernetes.io/manager is required
exit status 24

```
