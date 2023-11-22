在kubernetes上安装binlogo
============

- k8s相关的配置在[地址](docs/kubernetes)中
- 修改配置为自己需要的，然后手动或者使用命令
    > kubectl apply -f ./docs/kubernetes/
  
- 说明 
  > 因为本身binlogo的部署并不复杂，没有节点的证书。解析的日志也是直接发送到输出端，例如消息队列中，不需要配置持久化的硬盘。
  > 所以目前只是提供了写好的kubernetes配置文件，一个service和一个stateful的配置。
  > 后期会补充helm的配置文件，helm安装的方法。

