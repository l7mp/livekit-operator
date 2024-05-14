# LiveKit-Operator
The LiveKit-Operator is a Kubernetes operator which aims to ease the deployment of the [LiveKit](https://github.com/livekit/livekit)
media server into Kubernetes.

## Notice
Work in progress, not ready for deployment! Can be used to play with.
This operator only supports the deployment of the LiveKit server and no other media servers. This might change in
the future if there is interest in it by the people.

## Table of contents
* [Description](#description)
* [Features](#current-features)
* [Running the operator](#running-the-operator)
  * [Running locally](#running-locally)
  * [Deploying into the cluster via kubectl](#deploying-into-the-cluster-via-kubectl)
  * [Deploying into the cluster using Helm](#deploying-into-the-cluster-using-helm)

## Description
Since WebRTC and Kubernetes doesn't walk holding hands it can be a real struggle to deploy your application into your cluster.
Although STUNner provides many thorough examples of deploying some media servers, everything needs to be done manually.
There is no automation involved in the process, you must configure resources and deploy them one-by-one just to wait nervously
for a LoadBalancer IP to show up which can be taken and put into another resource which then to be deployed and so on.
It is hard for beginners and not easy for experienced developers.
The purpose of the operator is to provide a simple solution to all struggles, and create and deploy all resources into your
cluster in a few minutes. Users just need to configure a single custom resource which will be picked up by the operator
and used to configure the required resources.
In the upcoming sections all the necessary things will be explained in order to fire up the Operator and configure it
accordingly to your needs.
NOTE: This operator only supports the deployment of the LiveKit server and no other media servers.

## Current capabilities

As was mentioned above the operator is still in progress, lots of issues to be fixed, tested out and the API is changing
a lot. Although, it is not ready, yet it can fire up a LiveKit server along with LiveKit's Ingress and Egress if needed.
All the resources are created to have a running and usable setup. This includes three additional Helm charts
([STUNner-Gateway-Operator](https://github.com/l7mp/stunner/blob/main/docs/INSTALL.md),
[Envoy-Gateway](https://github.com/envoyproxy/gateway/tree/main/charts/gateway-helm),
and [Cert-Manager](https://cert-manager.io/docs/installation/helm/)), Redis resources, 
[External DNS](https://github.com/kubernetes-sigs/external-dns) resources, 
[Gateway API](https://gateway-api.sigs.k8s.io/) resources and so on.  

> [!Warning]
>
> The following constraints apply to the current version:
> - External DNS: only the CloudFlare configuration is possible and supported via the LiveKitMesh custom resource.
> - Only tested in GKE on a standard cluster. No experience running with different cloud providers.


## LiveKitMesh custom resource
The Operator uses a CR (Custom Resource) to trigger the deployment and control the configuration of a new setup.
This custom resource is called the `LiveKitMesh`, and it specifies separate components in the `spec` field in order to handle
all resources and configuration in a different context.
The Custom Resource Definition file for the CR can be found [here](config/crd/bases/livekit.stunner.l7mp.io_livekitmeshes.yaml).

The below yaml shows a potential configuration for the LiveKitMesh custom resource.
The things must be explained are as follows:
* The `metadata.name` and `metadata.namespace` pair MUST be unique otherwise things will break.
* `liveKit`, `stunner`, and `applicationExpose` components are REQUIRED, `ingress` and `egress` are OPTIONAL
* 

```yaml
apiVersion: livekit.stunner.l7mp.io/v1alpha1
kind: LiveKitMesh
metadata:
  labels:
    app.kubernetes.io/name: livekitmesh
    app.kubernetes.io/instance: livekitmesh-sample
    app.kubernetes.io/part-of: livekit-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: livekit-operator
  name: livekitmesh-sample
  namespace: default
spec:
  components:
    liveKit:
      deployment:
        replicas: 1
        config:
          keys:
            access_token: secretsecretsecretsecretsecretsecret
          log_level: debug
          port: 7880
# If you have your own redis deployment uncomment the lines under and specify them.          
#          redis:
#            address: redis.default.svc:6379
#            username: user
#            password: pw
#            db: db
          rtc:
            port_range_end: 60000
            port_range_start: 50000
            tcp_port: 7801
        container:
          image: livekit/livekit-server:v1.4.2
          imagePullPolicy: Always
          args: ["--disable-strict-config"]
          terminationGracePeriodSeconds: 3600
          resources:
            limits:
              cpu: 2
              memory: 512Mi
            requests:
              cpu: 500m
              memory: 128Mi

    applicationExpose:
      hostName: livekit.<your-domain>
      externalDNS:
        cloudFlare:
          token: <your-token>
          email: <your-email>
      certManager:
        issuer:
          apiToken: <your-token>
          challengeSolver: cloudflare
          email: <your-email>
    stunner:
      gatewayConfig:
        realm: stunner.l7mp.io
        authType: static
        userName: "username"
        password: "password"
      gatewayListeners:
        - name: udp-listener
          port: 3478
          protocol: TURN-UDP

    ingress:
      config:
        rtmp_port: 1935
        whip_port: 8080
        cpu_cost:
          rtmp_cpu_cost: 2
        http_relay_port: 9090
        logging:
          level: debug
        prometheus_port: 7889

    egress:
      config:
        log_level: debug
      container:
        resources:
          requests:
            memory: "256Mi"
            cpu: "1"
          limits:
            memory: "512Mi"
            cpu: "2"
```

### Component liveKit
This component is REQUIRED.  
In the `liveKit` component users can configure the startup configuration for the LiveKit server, and the container spec.
Note that the `spec.components.liveKit.deployment.container` is not the full container spec 

| Field                            | Type              | Description                                                                                                                                                          | Required |
|----------------------------------|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| deployment.replicas              | int               | Number of desired pods. This is a pointer to distinguish between  explicit zero and not specified. Defaults to `7`.                                                  | False    |
| deployment.config.keys           | map[string]string | API key / secret pairs. Keys are used for JWT authentication,  server APIs would require a keypair in order to generate access tokens  and make calls to the server. | True     |
| deployment.config.log_level      | string            | LogLevel is the level used in the LiveKit server. Defaults to `info` Valid values: `debug`,`info`,`warn`,`error`                                                     | False    |
| deployment.config.port           | int               | Port is main TCP port for RoomService and RTC endpoint. Defaults to `7880`.                                                                                          | False    |
| deployment.config.redis          | {}                | Redis configuration. Redis in case `redis` is configured no Redis  resources will be created by the operator.                                                        | False    |
| deployment.config.redis.address  | string            | Address of the Redis service.                                                                                                                                        | False    |
| deployment.config.redis.username | string            | Username for the Redis connection.                                                                                                                                   | False    |
| deployment.config.redis.password | string            | Password for the Redis connection.                                                                                                                                   | False    |
| deployment.config.redis.db       | string            | Database in Redis to use.                                                                                                                                            | False    |
| deployment.config.rtc            |                   |                                                                                                                                                                      |          |

### Component stunner
This component is REQUIRED.

### Component applicationExpose
This component is REQUIRED.

### Component ingress
This component is OPTIONAL.

### Component egress
This component is OPTIONAL.


## Running the operator

### Running locally
 
### Deploying into the cluster via kubectl

### Deploying into the cluster using Helm




### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/livekit-operator:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/livekit-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2021-2024 by its authors. Some rights reserved. See [AUTHORS](AUTHORS).

MIT License - see [LICENSE](LICENSE) for full text.

