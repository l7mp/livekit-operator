# LiveKit-Operator
The LiveKit-Operator is a Kubernetes operator which aims to ease the deployment of the [LiveKit](https://github.com/livekit/livekit)
media server into Kubernetes.

## Notice
Work in progress, not ready for deployment! Can be used to play with.
This operator only supports the deployment of the LiveKit server and no other media servers. This might change in
the future if there is interest in it by the people.

## Table of contents
* [Description](#description)
* [Features](#current-capabilities)
* [LiveKitMesh custom resource](#livekitmesh-custom-resource)
* [Component liveKit](#component-livekit)
* [Component stunner](#component-stunner)
* [Component applicationExpose](#component-applicationexpose)
* [Component ingress](#component-ingress)
* [Component egress](#component-egress)
* [Running the operator](#running-the-operator)
  * [Running locally](#running-locally)
  * [Deploying into the cluster via kubectl](#deploying-into-the-cluster-via-kubectl)
  * [Deploying into the cluster using Helm](#deploying-into-the-cluster-using-helm)
* [Uninstall CRDs](#uninstall-crds)
* [Licence](#license)

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
Note that the `spec.components.liveKit.deployment.container` is not the full container spec.
ALL fields in the below table are in the `spec.components.liveKit` field object.

| Field                                  	 | Type              	 | Description                                                                                                                                                          	  | Required 	 |
|------------------------------------------|---------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| deployment.replicas                    	 | int               	 | Number of desired pods. This is a pointer to distinguish between  explicit zero and not specified. Defaults to `1`.                                                  	  | False    	 |
| deployment.config.keys                 	 | map[string]string 	 | API key / secret pairs. Keys are used for JWT authentication,  server APIs would require a keypair in order to generate access tokens  and make calls to the server. 	  | True     	 |
| deployment.config.log_level            	 | string            	 | LogLevel is the level used in the LiveKit server. Defaults to `info` Valid values: `debug`,`info`,`warn`,`error`                                                     	  | False    	 |
| deployment.config.port                 	 | int               	 | Port is main TCP port for RoomService and RTC endpoint. Defaults to `7880`.                                                                                          	  | False    	 |
| deployment.config.redis                	 | {}                	 | Redis configuration. Redis in case `redis` is configured no Redis  resources will be created by the operator.                                                        	  | False    	 |
| deployment.config.redis.address        	 | string            	 | Address of the Redis service.                                                                                                                                        	  | False    	 |
| deployment.config.redis.username       	 | string            	 | Username for the Redis connection.                                                                                                                                   	  | False    	 |
| deployment.config.redis.password       	 | string            	 | Password for the Redis connection.                                                                                                                                   	  | False    	 |
| deployment.config.redis.db             	 | string            	 | Database in Redis to use.                                                                                                                                            	  | False    	 |
| deployment.config.rtc                  	 | {}                	 | WebRTC configuration for LiveKit                                                                                                                                     	  | False    	 |
| deployment.config.rtc.port_range_start 	 | int               	 | UDP ports to use for client traffic. Defaults to `50000`.                                                                                                            	  | False    	 |
| deployment.config.rtc.port_range_end   	 | int               	 | UDP ports to use for client traffic. Defaults to `60000`.                                                                                                            	  | False    	 |
| deployment.config.rtc.tcp_port         	 | int               	 | NOT WORKING as of now! When set, LiveKit enable WebRTC ICE over TCP when UDP isn't available. Defaults to `7801`.                                                     	 | False    	 |                                                                                                                                                               |          |

### Component stunner
This component is REQUIRED.  
In the `stunner` component users can configure and personalize the STUNner related resources, 
such as the authentication method, protocol and port.
ALL fields in the below table are in the `spec.components.stunner` field object.

| Field            	 | Type 	 | Description                                                                                                                                                                                                                                                	 | Required 	 |
|--------------------|--------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| gatewayConfig    	 | {}   	 | GatewayConfig is the configuration for the STUNner deployment's GatewayConfig object. This object is copied from the STUNner API. See how to configure it here: https://pkg.go.dev/github.com/l7mp/stunner-gateway-operator/api/v1alpha1#GatewayConfigSpec 	 | True     	 |
| gatewayListeners 	 | []   	 | GatewayListeners list is the configuration for the STUNner deployment's Gateway object. The list takes the Gateway API V1's Listener object. See how to configure each element: https://pkg.go.dev/sigs.k8s.io/gateway-api@v1.1.0/apis/v1#Listener         	 | True     	 |

### Component applicationExpose
This component is REQUIRED.
In the `applicationExpose` component users can configure all the resources related to expose TCP and HTTP based
applications (LiveKit server and Ingress).  
ALL fields in the below table are in the `spec.components.applicationExpose` field object.

| Field                                                	 | Type   	 | Description                                                                                                                                                                                                                                                                                                                          	 | Required 	 |
|--------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| hostName                           	                   | string 	 | HostName is the DNS hostname that will be used by both Cert-Manager, External DNS and Envoy GW. DnsZone Certificate requests will be issued against this  HostName. This ChallengeSolver will use this to solve the challenge.                                                                                                       	 | True     	 |
| certManager                        	                   | {}     	 | CertManager will obtain certificates from a variety of Issuers, both popular public Issuers and private Issuers, and ensure the certificates are valid and up-to-date, and will attempt to renew certificates at a configured time before expiry.                                                                                    	 | False    	 |
| certManager.issuer                 	                   | {}     	 | Issuer holds the necessary configuration for the used Issuer                                                                                                                                                                                                                                                                         	 | False    	 |
| certManager.issuer.email           	                   | string 	 | Email is the email address to be associated with the ACME account. This field is optional, but it is strongly recommended to be set.  It will be used to contact you in case of issues with your account or certificates, including expiry notification emails. This field may be updated after the account is initially registered. 	 | False    	 |
| certManager.issuer.challangeSolver 	                   | string 	 | ChallengeSolver is used to configure a DNS01 challenge provider to be used when solving DNS01 challenges. Valid values:  `cloudflare`, `route53`, `clouddns`, `digitalocean`, `azuredns`. NOTE: currently only `cloudflare` is supported.                                                                                            	 | True     	 |
| certManager.issuer.apiToken        	                   | string 	 | ApiToken is the API token for the CloudFlare account that owns the challenged DnsZone.                                                                                                                                                                                                                                               	 | True     	 |
| externalDNS                                            | TODO     | TODO                                                                                                                                                                                                                                                                                                                                   |            |

### Component ingress
This component is OPTIONAL.
In the `ingress` component users can configure the LiveKit Ingress related resources.
ALL fields in the below table are in the `spec.components.ingress` field object.

| Field                                            	 | Type   	 | Description                                                                                                                                                                                             	 | Required 	 |
|----------------------------------------------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| config                                           	 | {}     	 | Config holds configuration for the LiveKit Ingress                                                                                                                                                      	 | False    	 |
| config.cpu_cost                                  	 | {}     	 | CPU resources to reserve when accepting sessions, in fraction of core count                                                                                                                             	 | False    	 |
| config.cpu_cost.rtmp_cpu_cost                    	 | int    	 | CPU resources to reserve when accepting RTMP sessions, in fraction of core count                                                                                                                        	 | False    	 |
| config.cpu_cost.whip_cpu_cost                    	 | int    	 | CPU resources to reserve when accepting WHIP sessions, in fraction of core count                                                                                                                        	 | False    	 |
| config.cpu_cost.whip_bypass_transcoding_cpu_cost 	 | int    	 | CPU resources to reserve when accepting WHIP sessions with transcoding disabled, in fraction of core count                                                                                              	 | False    	 |
| config.health_port                               	 | int    	 | If used, will open an http port for health checks.                                                                                                                                                      	 | False    	 |
| config.prometheus_port                           	 | int    	 | Port used to collect Prometheus metrics.                                                                                                                                                                	 | False    	 |
| config.rtmp_port                                 	 | int    	 | Port for the RMTP service. If specified the all the necessary rescources (e.g.: gatewayclass, gateway listener, tcp route) will be created. If omitted completely no resources will be created for it.  	 | False    	 |
| config.whip_port                                 	 | int    	 | Port for the WHIP service. If specified the all the necessary rescources (e.g.: gatewayclass, gateway listener, http route) will be created. If omitted completely no resources will be created for it. 	 | False    	 |
| config.http_relay_port                           	 | int    	 | TCP port for communication between the main service process and session handler processes                                                                                                               	 | False    	 |
| config.logging.level                             	 | string 	 | Sets the log level of the deployed Ingress resource. debug, info, warn, or error. Defaults to `info`.                                                                                                   	 | False    	 |

### Component egress
This component is OPTIONAL.
In the `egress` component users can configure the LiveKit Ingress related resources.
ALL fields in the below table are in the `spec.components.egress` field object.

| Field                  	 | Type   	 | Description                                                                                          	 | Required 	 |
|--------------------------|----------|--------------------------------------------------------------------------------------------------------|------------|
| config                 	 | {}     	 | Config holds configuration for the LiveKit Egress                                                    	 | False    	 |
| config.health_port     	 | int    	 | If used, will open an http port for health checks.                                                   	 | False    	 |
| config.template_port   	 | int    	 | Port used to host default templates.                                                                 	 | False    	 |
| config.prometheus_port 	 | int    	 | Port used to collect Prometheus metrics.                                                             	 | False    	 |
| config.log_level       	 | string 	 | Sets the log level of the deployed Egress resource. debug, info, warn, or error. Defaults to `info`. 	 | False    	 |
| config.s3              	 | {}     	 | S3 configuration. See  https://docs.livekit.io/home/self-hosting/egress/#Config                      	 | False    	 |
| config.azure           	 | {}       | Azure configuration. See  https://docs.livekit.io/home/self-hosting/egress/#Config                   	 | False    	 |
| config.gcp             	 | {}       | GCP configuration. See  https://docs.livekit.io/home/self-hosting/egress/#Config                     	 | False    	 |

## Running the operator

There are multiple ways to run the operator. In the following subsections there are three different methods.
But first the `LiveKitMesh` CRD must be installed into the cluster.
```sh
make manifests generate
make install
```

### Running locally

```sh
make run
```

The above command will start the operator locally. Basically, it runs  `go run ./main.go`.
This method will not create any resources in the operator. It uses the current `kubeconfig` to interact with your cluster.

### Deploying into the cluster via kubectl

```sh
make deploy IMG=<some-registry>/livekit-operator:<tag>
```
The above comamnd will deploy the controller to the cluster with the image specified by `IMG`. `deploy` uses `kubectl`
to apply all the generated resources into the cluster.

##### Undeploy controller
Undeploy the controller from the cluster:

```sh
make undeploy
```


### Deploying into the cluster using Helm

TODO

###### Uninstall the Helm chart from your cluster

```sh
helm uninstall livekit-operator -n <your-namespace>
```

### Command-line flags

It is possible to modify the log level of the operator and which Helm charts to deploy.
The list of the flags is as follows:
 * The default log level is `info`; however, if you would like to see everything that happens in the container feed the
`--zap-log-level=10` flag to the startup argument list. 
   * Usage: `go run ./main.go --zap-log-level=10`
 * By default, three Helm charts are installed in the cluster, these are the STUNner-Gateway-Operator, Envoy-Gateway 
and Cert-Manager. You can turn them off one-by-one in case you have your own in your cluster.
   * Usage (this enables the stunner-gw-operator chart, while disabling the envoy-gateway and cert-manager charts:
    ```
       go run ./main.go \
          --install-stunner-gateway-chart=true \
          --install-cert-manager-chart=false \
          --install-envoy-gateway-chart=false
    ```

### Uninstall CRDs
In case you would like to clean up everything after the Operator
you need to delete the CRDs from the cluster.

If you deployed the crd with `make install`:

```sh
make uninstall
```

If you deployed using the Helm chart you need to delete the CRD manually since Helm does not delete CRDs.

```sh 
kubectl delete crd livekitmeshes.livekit.stunner.l7mp.io
```

## License

Copyright 2021-2024 by its authors. Some rights reserved. See [AUTHORS](AUTHORS).

MIT License - see [LICENSE](LICENSE) for full text.

## Acknowledgments

Huge thanks to [LiveKit](https://github.com/livekit). Without their great open-source tools
and infrastructure this project would not be possible.