# aws-ssm-operator

A Kubernetes operator that automatically maps what are stored in AWS SSM Parameter Store into Kubernetes Secrets.

`aws-ssm-operator` Custom Resources defines desired state of Kubernetes Secret fetched from SSM Parameter Store. Otherwise, `parameterstore-controller` Controller monitors user's request and cached parameter values or credentials as plaintext into Kubernetes Secret.

## Before you begin

You have to store configuration parameter or credentials into SSM Parameter Store.

Let's say, your application inside a Pod wish to connect to Aurora instances by using database user and password.

```bash
# Store database user with simple string
aws ssm put-parameter \
    --name "/aurora/mysql/stg/dbuser" \
    --type "String" \
    --value "dbuser" \
    --overwrite
# Store database password encrypted with default KMS key
aws ssm put-parameter \
    --name "/aurora/mysql/stg/dbpassword" \
    --type "SecureString" \
    --value "dbpassword" \
    --overwrite
```

In case of using EKS Cluster as your Kubernetes platform, attach `AmazonSSMReadOnlyAccess` policy or equivalent of that to Instance Profiles of EKS worker nodes.

## Installation

```bash
# Setup Service Account
kubectl apply -f deploy/

# Verify that a Pod is running
kubectl get pod -l app=aws-ssm-operator --watch -n kube-system
```

## Usage

Create an sample Paramter Store resource::

```bash
# Create an example Parameter Store resource
$ kubectl create -f example/database.yaml
```

## Verifying

You can find your credentials in separate Secret resources as follows. The key to secret data is 'name' which is hardcoded and cannot be changed.

```bash
$ kubectl describe secret dbpassword
Name:         dbpassword
Namespace:    default
Labels:       app=dbpassword
Annotations:  <none>

Type:  Opaque

Data
====
name:  9 bytes

$ kubectl describe secret dbuser
Name:         dbuser
Namespace:    default
Labels:       app=dbuser
Annotations:  <none>

Type:  Opaque

Data
====
name:  7 bytes
```

You can reference secret data in your Pod declaration as follows:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: sample-app-pod
spec:
  containers:
  - name: sample-app
    image: toversus/sample-app
    env:
      - name: DB_USER
        valueFrom:
          secretKeyRef:
            name: dbuser
            key: name
      - name: DB_PASSWORD
        valueFrom:
          secretKeyRef:
            name: dbpassword
            key: name
  restartPolicy: Never
```

In Knative Service resource, you can define them in containers spec:

```yaml
   apiVersion: serving.knative.dev/v1alpha1
   kind: Service
   metadata:
     name: sample-app
   spec:
     template:
       spec:
         containers:
           - image: toversus/sample-app
             env:
               - name: DB_USER
                 valueFrom:
                   secretKeyRef:
                     name: dbuser
                     key: name
               - name: DB_PASSWORD
                 valueFrom:
                   secretKeyRef:
                     name: dbpassword
                     key: name
```

## Clean up

To clean up all the components:

```bash
kubectl delete -f deploy/crds/ssm_v1alpha1_parameterstore_cr.yaml \
    -f deploy/operator.yaml \
    -f deploy/crds/ssm_v1alpha1_parameterstore_crd.yaml \
    -f deploy/role.yaml -f deploy/role_binding.yaml \
    -f deploy/service_account.yaml
```

## Acknowledgements

The idea behind this project is fully based on [mumoshu/aws-secret-operator](https://github.com/mumoshu/aws-secret-operator). Thank you for your awesome work!
