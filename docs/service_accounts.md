# Tiller and Service Accounts

In Kubernetes, granting a role to an application-specific service account is a best practice to ensure that your application is operating in the scope that you have specified. Read more about service account permissions in Kubernetes [here](https://kubernetes.io/docs/admin/authorization/rbac/#service-account-permissions).

You can add a service account to Tiller using the `--service-account <NAME>` flag while you're configuring helm. As a prerequisite, you'll have to create a role binding which specifies a [role](https://kubernetes.io/docs/admin/authorization/rbac/#role-and-clusterrole) and a [service account](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/) name that have been set up in advance.

Once you have satisfied the pre-requisite and have a service account with the correct permissions, you'll run a command like this: `helm init --service-account <NAME>`

## Example: Service account with cluster-admin role

```console
$ kubectl create serviceaccount tiller --namespace kube-system
```

In `rbac-config.yaml`:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
```

_Note: The cluster-admin role is created by default in a Kubernetes cluster, so you don't have to define it explicitly._

```console
$ kubectl create -f rbac-config.yaml
$ helm init --service-account tiller
```

## Example: Service account restricted to a namespace
In the example above, we gave Tiller admin access to the entire cluster. You are not at all required to give Tiller cluster-admin access for it to work. Instead of specifying a ClusterRole or a ClusterRoleBinding, you can specify a Role and RoleBinding to limit Tiller's scope to a particular namespace.

```console
$ kubectl create namespace application-world
$ kubectl create serviceaccount tiller --namespace kube-system
```

In `role-deployment-manager.yaml`,
```yaml
kind: Role
  apiVersion: rbac.authorization.k8s.io/v1beta1
  metadata:
    namespace: application-world
    name: deployment-manager
  rules:
  - apiGroups: ["", "extensions", "apps"]
    resources: ["deployments", "replicasets", "pods", "configmaps", "secrets"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"] # You can also use ["*"]
```

```console
$ kubectl create -f role-deployment-manager.yaml
```

In `role-binding-deployment-manager.yaml`,
```yaml
kind: RoleBinding
  apiVersion: rbac.authorization.k8s.io/v1beta1
  metadata:
    name: deployment-manager-binding
    namespace: office
  subjects:
  - kind: User
    name: employee
    apiGroup: ""
  roleRef:
    kind: Role
    name: deployment-manager
    apiGroup: ""
```

```console
$ kubectl create -f rolebinding-deployment-manager.yaml
```

```console
$ helm init --service-account tiller
```

