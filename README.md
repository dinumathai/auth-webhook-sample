# Kubernetes authentication/authorization webhook using golang in minikube

Here we are discussing mainly on how webhooks can be used to delegate authentication and authorization in Kubernetes. 

We will discuss about 
1. [Authentication webhook](#authentication-webhook)
1. [RBAC Authorization](#rbac-authorization)
1. [Authorization webhook](#authorization-webhook)
1. [Build and deploy in minikube](#build-and-deploy-in-minikube)

[auth-webhook-sample](https://github.com/dinumathai/auth-webhook-sample) is a sample Kubernetes authentication and authorization webhook application. The code is structured to extend for further use cases like authentication against AD or some other open id provider like Azure AD.

# Authentication webhook
Kubernetes authentication webhook can be used to delegate authentication outside of the Kubernetes.

## Why authentication webhook
Kubernetes has below way of managing authentication.
1. Using valid certificate signed by the cluster's certificate authority (CA).
1. Using static token file.
1. OpenID Connect Tokens.
1. Kubernetes service account.

Read about them at [Kubernetes Authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/). If none of them serve your purpose, Kubernetes authentication webhook is your best option(preferred over authenticating proxy). Usually webhook is used for integration with authentication system like LDAP, SAML etc.

## How authentication webhook works
![authentication webhook flow](./doc/webhook-flow.png)

1. User generates token. In our case we will be using `auth-webhook-sample` application to generate token using the username/password. The token can be generated from different source for example if we use Azure AD, the token will be generated using Azure AD api.
2. User uses the token to call Kubernetes api by setting the token in the api header. In our case we will set the token in `kubectl` config and execute `kubectl get pods` command to call the Kubernetes api.
3. The token received by Kubernetes api will be passes to authentication webhook in [predefined format](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#tokenreview-request-0)
4. The webhook validates the token a returns the `status` and `groups` for the user in [required format](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#tokenreview-response-success-0)
5. Kubernetes will return response after validating the user permission to access the requested resource using the `groups` from webhook. The access can be configured using Kubernetes roles / clusterroles and Kubernetes rolebindings /clusterrolebindings. Read more at [Kubernetes RBAC Authorization
](https://kubernetes.io/docs/reference/access-authn-authz/rbac/).


## What is authentication webhook?
Authentication webhook is a HTTPS service that receives a request in defined format. The request is validated and must return back a response in defined format.

### Request Details
__Method__ : POST
__Request Param__: `timeout=30s` - Configured timeout will be received as request parameter.
__Request body__ : 
```
{
  "kind": "TokenReview",
  "apiVersion": "authentication.k8s.io/v1beta1",
  "metadata": {
    "creationTimestamp": null
  },
  "spec": {
    "token": "### AUTHENTICATION BEARER TOKEN ###",
    "audiences": [
      "https://kubernetes.default.svc.cluster.local"
    ]
  },
  "status": {
    "user": {}
  }
}
```
### Response Details
__Response body__ : 
```
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": true,
    "user": {
      "username": "some-username",
      "uid": "some-uid",
      "groups": [
        "group_one",
        "group_two",
        "group_n"
      ]
    }
  }
}
```

# RBAC Authorization
We have see that each user will be having a set of groups - [config/user_details.yaml](config/user_details.yaml). Now we will discuss on how to give permission(authorize) the groups for kubernetes objects(deployment.secret etc).

##  Type of Kubernetes User Authorization
1. [ABAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/abac/) helps to provide the access to users using a file(`--authorization-policy-file=SOME_FILENAME`) which contain all the policies.
1. [RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) granting permission to the groups of user using `Roles`/`ClusterRoles` and `RoleBindings`/`ClusterRoleBindings`
1. [Authorization Webhook](https://kubernetes.io/docs/reference/access-authn-authz/webhook/) a way to delegate authorization in Kubernetes to an external Application/HTTP-API.

## At what level the permission can be granted?
1. __apiGroups__: indicates the core API group
1. __resources__: Kubernetes resources like "pods", "secrets", "deployments" etc
1. __verbs__: Action like  "list", "get", "update" etc. Refer [here](https://kubernetes.io/docs/reference/access-authn-authz/authorization/#determine-the-request-verb)
1. __resourceNames__: List of resource name to which access must be given.
1. __namespaces__: To grand access to cluster level resources and resources on all the namespace `ClusterRoles` and `ClusterRoleBindings` are used. To grant permission resources on a single namespace `Roles` and `RoleBindings` are used.

## Kubernetes RBAC Authorization objects.
1. __Roles__: Defines the Kubernetes access at a namespace level setting the `apiGroups`, `resources` and `verbs`.
1. __ClusterRoles__: Same as that of `Roles` except used to define access of cluster level Kubernetes objects and for access across the namespace.
1. __RoleBindings__: Links a `Role` to a `User`, `Group` or `ServiceAccount`
1. __ClusterRoleBindings__: Links a `ClusterRoles` to a `User`, `Group` or `ServiceAccount`

## Default user-facing roles
Kubernetes allows you to use default user-facing roles, including, but not limited to:

1. __cluster-admin__: This “superuser” can perform any action on any resource in a cluster. You can use this in a ClusterRoleBinding to grant full control over every resource in the cluster (and in all namespaces) or in a RoleBinding to grant full control over every resource in the respective namespace.
1. __admin__: This role permits unlimited read/write access to resources within a namespace. This role can create roles and role bindings within a particular namespace. It does not permit write access to the namespace itself.
1. __edit__: This role grants read/write access within a given Kubernetes namespace. It cannot view or modify roles or role bindings. 
1. __view__: This role allows read-only access within a given namespace. It does not allow viewing or modifying of roles or role bindings. 
You can find even more information about these user-facing roles and others in the [Kubernetes documentation](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles).

## Examples 
1. `admin-crb` is a `ClusterRoleBinding` that links the default `ClusterRole` - `cluster-admin` to a user group `g_admin`. So that any users with group `g_admin` will have cluster admin privilege.
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-crb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: g_admin
```

2. `app-watcher-role` is a `Role` which defines read access to `pods` and `deployments` in the `my-app-namespace` namespace. And `RoleBinding` `app_watcher-rb` links the role  `app-watcher-role` to group `g_app_watcher`. So that any users with group `g_app_watcher` will have read permission to `pods` and `deployments` in the `my-app-namespace` namespace.
```
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: my-app-namespace
  name: app-watcher-role
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods","deployments"]
  verbs: ["get", "watch", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: my-app-namespace
  name: app_watcher-rb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: app-watcher-role
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: g_app_watcher
```
Read more at [Kubernetes Authorization documentation](https://kubernetes.io/docs/reference/access-authn-authz/authorization/).


# Authorization webhook
Authorization webhook is used to delegate authorization in Kubernetes. Also please note that this is an option rarely used in Kubernetes clusters. It is always advice to use built-in RBAC Authorization.

## What is Authorization webhook?

Authorization webhook is a HTTPS service that receives a request in defined format. The request is validated and must return back a response in defined format.

__Request body__
```
{
  "kind": "SubjectAccessReview",
  "apiVersion": "authorization.k8s.io/v1beta1",
  "metadata": {
    "creationTimestamp": null
  },
  "spec": {
    "resourceAttributes": {
      "verb": "get",
      "group": "storage.k8s.io",
      "version": "v1",
      "resource": "csinodes",
      "name": "minikube"
    },
    "user": "system:node:minikube",
    "group": [
      "system:nodes",
      "system:authenticated"
    ]
  },
  "status": {
    "allowed": false
  }
}
```
__Access Granted Response__
```
{
  "apiVersion": "authorization.k8s.io/v1",
  "kind": "SubjectAccessReview",
  "status": {
    "allowed": true
  }
}
```
__Access Denied Response__
```
{
  "apiVersion": "authorization.k8s.io/v1",
  "kind": "SubjectAccessReview",
  "status": {
    "allowed": false,
    "denied": true,
    "reason": "User do not have access to resource"
  }
}
```

# Build and deploy in minikube
To get the webhooks up and running in minikube. First we have have generate certificates for webhooks, bring up the webhook and then configure the minikuke to use it. And finally test it :-).

## Prerequisites
1. Basic understanding of Kubernetes.
1. Minikube running in local machine.
1. openssl
1. Kubectl must be installed locally and must have a basic understanding of config for `kubectl`.
1. Docker Or Golang must be installed locally.

## Create the certificate

The [deploy/ca/server.conf](deploy/ca/server.conf) must be modified and add your local system ip instead of `192.168.1.35`(Update on two places). And [deploy/ca/server.crt](deploy/ca/server.crt) must be regenerated using below commands.

```
git clone git@github.com/dinumathai/auth-webhook-sample.git
cd auth-webhook-sample
openssl req -new -key deploy/ca/server.key -out deploy/ca/server.csr -config deploy/ca/server.conf

# Sign the server certificate with the above CA
openssl x509 -req -in deploy/ca/server.csr -CA deploy/ca/ca.crt -CAkey deploy/ca/ca.key -CAcreateserial -out deploy/ca/server.crt -days 100000 -extensions v3_req -extfile deploy/ca/server.conf
```
Commands to generate the all certificate files are available at [deploy/ca/README.md](deploy/ca/README.md).

## Building & Run webhook using docker
The below docker is already uploaded to docker hub. So you can directly run the docker run command to bring up the authentication webhook.
```
docker build -t dmathai/auth-webhook-sample:latest -f Dockerfile .

# GENERATE the server.crt and server.key
export AUTH_CERT_TLS_CRT=$(cat deploy/ca/server.crt)
export AUTH_CERT_TLS_KEY=$(cat deploy/ca/server.key)
docker run --env AUTH_CERT_TLS_KEY=$AUTH_CERT_TLS_KEY --env AUTH_CERT_TLS_CRT=$AUTH_CERT_TLS_CRT -p 8443:8443 dmathai/auth-webhook-sample:latest
```
The webhook application will at https://localhost:8443/.

## Building & Run webhook locally
If you want to run the application locally with our docker. Please follow below commands.
```
go build github.com/dinumathai/auth-webhook-sample

# GENERATE the server.crt and server.key
export AUTH_CERT_TLS_CRT=$(cat deploy/ca/server.crt)
export AUTH_CERT_TLS_KEY=$(cat deploy/ca/server.key)
./auth-webhook-sample
```
The webhook application will at https://localhost:8443/.

## API Details

### Generate Auth JWT token
In this api the user credentials/details are managed by the auth service. Refer [config/user_details.yaml](config/user_details.yaml) to see the list of user and the groups configured for the users. The filepath of user details is configured in `v0.userDetailFilePath` of [config/auth_config.yaml](config/auth_config.yaml). 

[Configuration file is explained here](doc/configuration.md)
```
curl -X POST --insecure https://localhost:8443/v0/login  -u __YOUR_USERNAME__:__YOUR_PASSWORD__
```

### Validate the Token
This URL will be used by Kubernetes to validate the token.
```
curl -X POST --insecure https://localhost:8443/v0/authenticate  -H 'Authorization: Bearer XXXXXXXXX'
```

```
curl -X POST --insecure https://localhost:8443/v0/authenticate -d '{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "spec": {
    "token": "XXXXXXX"
  }
}'
```

## Deploy in minikube
Assuming that the authentication webhook is running in https://192.168.1.35:8443/. If not you have to make sure [deploy/auth-webhook-conf.yaml](deploy/auth-webhook-conf.yaml) is updated with proper url. Also [deploy/ca/server.conf](deploy/ca/server.conf) is modified and [deploy/ca/server.crt](deploy/ca/server.crt) is regenerated.

1. Start minikube
1. Create `ClusterRoleBinding` using - `kubectl apply deploy/create-cluster-role-binding.yaml`. We are creating the cluster-role-binding for a groups `g_admin`, `g_write` and `g_read`. The user `admin` is configured to have groups `g_admin`, refer [config/user_details.yaml](config/user_details.yaml). Read more at [Kubernetes RBAC Authorization
](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
1. Stop minikube.
1. Create a folder with name `var/lib/minikube/certs/auth`(`mkdir -p var/lib/minikube/certs/auth`) inside `$HOME/.minikube/files`.
1. Copy [deploy/ca/ca.crt](deploy/ca/ca.crt), [deploy/authorize-webhook-conf.yaml](deploy/authorize-webhook-conf.yaml) and [deploy/auth-webhook-conf.yaml](deploy/auth-webhook-conf.yaml) to `$HOME/.minikube/files/var/lib/minikube/certs/auth` folder. All these files will be available inside `/var/lib/minikube/certs/auth` folder of minikube container. You can confirm this by restarting minikube and doing `minikube ssh`.
1. Restart minikube with below command.

Only with authentication webhook
```
minikube start --driver=docker --extra-config apiserver.authorization-mode=RBAC --extra-config apiserver.authentication-token-webhook-config-file=/var/lib/minikube/certs/auth/auth-webhook-conf.yaml
```

With authentication and authorization webhook
```
minikube start --driver=docker --extra-config apiserver.authorization-mode=RBAC,Webhook --extra-config apiserver.authentication-token-webhook-config-file=/var/lib/minikube/certs/auth/auth-webhook-conf.yaml --extra-config apiserver.authorization-webhook-config-file=/var/lib/minikube/certs/auth/authorize-webhook-conf.yaml
```

## Test authentication webhook
1. Generate the token using the curl `curl -X POST --insecure https://localhost:8443/v0/login  -u admin:admin`
2. Add new user to the `kubectl` config with token and change the context to point to the new user. The below commands will help to do this.
```
# Add new user "admin" with generated token
kubectl config set-credentials admin --token=XXXXXXXXXXX
# Change the context "minikube" to point to user "admin"
kubectl config set-context minikube --user=admin
```
3. `kubectl get pods --all-namespaces` must return some pods and you must get some logs in webhook application. Done !!!

## Debugging tips
1. If the `minikube` is not starting with webhook config. Do `minikube ssh` to get into the minikube docker container. Run command `docker ps | grep apiserver` to get the api-server container. `docker logs <container_id>` to get the logs.
1. If you are getting error `error: You must be logged in to the server (Unauthorized)`. View the logs of webhook application to see whether any request is reaching the webhook application. Also refer the apiserver logs to make sure that cluster is  able to communicate with authentication webhook.
1. If request if not reaching webhook application but the `kubectl` commands are working. Each time the `minikube` is restarted the `kubectl` config will be reset. Please make sure that the context is pointing to the user with token. Also there is a default cache time of 30sec for which the cluster will cache the response from webhook.
1. If you are getting some error like `Error: pods is forbidden: User "admin" cannot list resource "pods" in API group ""`. Open [https://jwt.io/](https://jwt.io/) make sure that the token is having expected `groups` in the jwt token. If expected `groups` are there it has to do something with with Kubernetes `Roles/ClusterRoles` or `RoleBinding/ClusterRoleBinding`. Continue reading to learn more.

## References 

1. [Kubernetes Webhook Token Authentication documentation](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication)
1. [Kubernetes Authorization documentation](https://kubernetes.io/docs/reference/access-authn-authz/authorization/)
