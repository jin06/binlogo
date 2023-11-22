Install binlogo on kubernetes
============

- K8s related configuration [Address](../kubernetes)
- Modify the configuration to your own needs, and then manually or use the command
  > kubectl apply -f ./docs/kubernetes/

- Explain
  > Because the deployment of binlogo itself is not complex, there is no node certificate. The parsed logs are also sent directly to the output side, for example, in the message queue. There is no need to configure a persistent hard disk.
  > Therefore, only the written kubernetes configuration file, a service configuration and a stateful configuration are provided.
  > The configuration file of helm and the method of installing by helm will be supplemented later.

