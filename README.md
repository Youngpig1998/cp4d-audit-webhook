# CP4D-Audit-Webhook

This webhook looks for pod that has `cp4d-audit: "yes"` label and injects the sidecar container to your pod.

## Preqreqs

#### Environmental checklist🧾

| Configuration item | Specific configuration |
| :----------------: | :--------------------: |
|         OS         |        centos7         |
|     Kubernetes     |         v1.22+         |

Have Installed the cert-manager operator（Be care of the version of the operator）  ,it will create the secret we need.

**PS:This cert-manager operator is the newest version.**

```shell
kubectl apply -f $(pwd)/deploy/cert-manager.yaml
```



## Steps to install cp4d-audit-webhook

1. Create an issuer

   ```shell
   kubectl apply -f $(pwd)/deploy/issuer.yaml
   ```

   

2. Create a certificate

   ```shell
   kubectl apply -f $(pwd)/deploy/certificate.yaml
   ```

PS: Watch out   the default is the namespace!!!!
spec:
  dnsNames:
  - audit-webhook-service.default.svc


3. Create ConfigMap、Service、deployment and Mutatingwebhookconfiguration

   ```shell
   kubectl apply -f $(pwd)/deploy/audit-webhook-configmap.yaml
   kubectl apply -f $(pwd)/deploy/audit-webhook-server-service.yaml
   kubectl apply -f $(pwd)/deploy/audit-webhook-server-deployment.yaml
   kubectl apply -f $(pwd)/deploy/audit-mutating-webhook-configuration.yaml
   ```

PS: Watch out   the default is the namespace!!!!
In audit-mutating-webhook-configuration.yaml:
   1)、cert-manager.io/inject-ca-from: "default/serving-cert"
   2）、clientConfig:
         service:
            name: audit-webhook-service
            namespace: default


To test this, use the logwriter example

```shell
kubectl apply -f $(pwd)/deploy/example/logwriter.yaml
```

## To utilize cp4d-audit-webhook

1. Include the label `cp4d-audit: "yes"` in the pod/deployment/statefulset/job type of resource where the audit logs are generated. Inlcude the label in the below mentioned location.
   
   i) For Pod
   
   `.metadata.labels`
   
   iii) For Deployment, StatefulSet, Job
   
   `.spec.template.metadata.labels`
2. Push the audit logs to the volume named `varlog` [like here](deploy/example/logwriter.yaml)
   
   ```yaml
   volumes:
   - name: varlog
       emptyDir: {}
   ```
3. Currently the sidecar image is available.

## Log Augmentation and Auditing

The log data generated by the IBM Watson team may not be CloudPak compliant. We compared the public log and private log data and we found the attachments field is missing. So we augment this log field by using the webhook we created.

1. Push your necessary environment variables to the volume named`varlog`[like here](deploy/example/logwriter.yaml)

```yaml
volumes:
  - name: varlog
    emptyDir: {}
```

The environment variables include NAMESPACE, CONTAINERNAME, NODENAME, PODIPADDRESS and CONTAINERID. 
The first four environment variables can be specified in your deployment YAML file. An example is shown in [here](deploy/example/logwriter.yaml). 
CONTAINERID, however, needs to be shared via a directory path. An example is shown in [here](logwriter/writer.sh).
`More details` [Link here](https://github.ibm.com/PrivateCloud-analytics/zen-dev-test-utils/blob/gh-pages/docs/audit-logging.md#adding-system-environment-variables) in here.
They will be written to a EmptyDir that will be picked up by the webhook-injected sidecar container later on.


2. The log data is read by a fluentd in_tail plugin in the sidecar container. It also[augments](fluentd/example.rb) the log with the required audit information, and POST it to the stdout.
   More details can be refered[here](fluentd/fluent.conf).

## AuditWebhook

It will receive the message which sent by the mutatingwebhookconfiguration,and set the sidecar container to the loogwriter pod.
