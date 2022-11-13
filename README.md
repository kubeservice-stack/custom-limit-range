# Custom Limit Range CRD 使用

## Pod网络限流

### 一、设计文档

设计文档地址: [https://kubeservice.cn/2022/08/10/k8s-pod-bandwidth-limit/](https://kubeservice.cn/2022/08/10/k8s-pod-bandwidth-limit/)

### 二、部署使用步骤

0. Check 底层 `kubelet` 是否配置 `cni` 模式. 目前 `1.21.5是默认开启`的

```bash
$ /usr/bin/kubelet 
    --cni-bin-dir=/usr/lib/cni/ 
    --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf 
    --kubeconfig=/etc/kubernetes/kubelet.conf 
    --config=/var/lib/kubelet/config.yaml 
    --cgroup-driver=cgroupfs 
    --cni-bin-dir=/opt/cni/bin
    --cni-conf-dir=/etc/cni/net.d   ## cni-conf-dir设置为cni 配
    --network-plugin=cni      ## network-plugin设置为cni
    ... 

```

1. Check cni 配置. 目前 `1.21.5是默认开启`

查看`/etc/cni/net.d/`下的 calico 配置

```json
{
  "name": "k8s-pod-network",
  "cniVersion": "0.4.0",
  "plugins": [
    {
      "type": "calico",
      "log_level": "info",
      "datastore_type": "kubernetes",
      "nodename": "127.0.0.1",
      "ipam": {
        "type": "host-local",
        "subnet": "usePodCidr"
      },
      "policy": {
        "type": "k8s"
      },
      "kubernetes": {
        "kubeconfig": "/etc/cni/net.d/calico-kubeconfig"
      }
    },
    # 设置包含为以下部分
    {
      "type": "bandwidth",
      "capabilities": {"bandwidth": true}
    }
  ]
}
```

2. 部署基础 cert-manager 管理 webhook ca证书。 
如果集群中有这部分了， 可以跳过这步骤

```bash
$ kubectl apply -f hack/deployment/certs/cert-manager.yaml
```

确保启动成功

```bash
$ kubectl get pod -A
NAMESPACE      NAME                                               READY   STATUS    RESTARTS          AGE
cert-manager   cert-manager-6789f4858b-6975h                      1/1     Running   5 (158m ago)      4d
cert-manager   cert-manager-cainjector-659cdf7887-ftqqw           1/1     Running   44 (157m ago)     4d
cert-manager   cert-manager-webhook-69b9966fb-r2hj7               1/1     Running   4 (158m ago)      4d
```

3. 创建 custom limit range ca文件、https pem文件

```bash
$ kubectl apply -f hack/deployment/certs/issuer.yaml
$ kubectl apply -f hack/deployment/certs/certificate.yaml
```

4. 创建CRD文件

```bash
$ kubectl apply -f hack/deployment/crds/custom.cmss.com_customlimitranges.yaml
```

验证CRD创建成功

```bash
$ kubectl get crd   
NAME                                  CREATED AT
certificaterequests.cert-manager.io   2022-08-18T03:01:29Z
certificates.cert-manager.io          2022-08-18T03:01:29Z
challenges.acme.cert-manager.io       2022-08-18T03:01:29Z
clusterissuers.cert-manager.io        2022-08-18T03:01:30Z
customlimitranges.custom.cmss.com     2022-08-10T16:04:00Z
issuers.cert-manager.io               2022-08-18T03:01:31Z
orders.acme.cert-manager.io           2022-08-18T03:01:32Z
$ kubectl get crd | grep "customlimitranges"
customlimitranges.custom.cmss.com     2022-08-10T16:04:00Z
```

5. 部署webhook

```bash
$ kubectl apply -f hack/deployment/webhook/
```

验证webhook部署成功：

```bash
$ kubectl get service -n kube-system
NAME                               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                  AGE
customlimitrange-webhook-service   ClusterIP   10.106.25.21     <none>        443/TCP                  4d4h

$ kubectl get pod -n kube-system
NAME                                               READY   STATUS    RESTARTS           AGE
customlimitrange-webhook-ffddbd9f9-74kgz           1/1     Running   1 (4h40m ago)      2d19h
```

6. 验证`customlimitrange`生效
在`hack/deployment/example/`中，包括各种case

### 三、 验证注意点
分为`2`部分:

> 基于`Pod`和`Deployment`添加`annotations` 方式验证
> 
> 通过namespace添加`CustomLimitRange`, 设置上下限Max/Min、添加默认Default 方式

验证细节：

> 1）ingress 生效；
> 
> 2）egress 生效；

1. ingress: client -> server. server端入口流量限流成功
2. egress: client -> server. client端出口流量限流成功

验证工具，设计方案中的验证部分：
[https://kubeservice.cn/2022/08/10/k8s-pod-bandwidth-limit/](https://kubeservice.cn/2022/08/10/k8s-pod-bandwidth-limit/)


