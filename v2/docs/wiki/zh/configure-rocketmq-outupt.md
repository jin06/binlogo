> - 一般用户都是直接使用阿里云的RocketMQ云服务，所以在做RocketMQ的输出端时直接采用了云服务的形式。
> - Binlogo 会把消息发送到相应的RocketMQ的主题中
> - 用户需要自行登录阿里云去配置相关的主题和接入点等，然后在binlogo后台去创建流水线
> - 参考 [RocketMQ Topic Pattern](https://help.aliyun.com/document_detail/255810.html?spm=5176.rocketmq.help.dexternal.248c7d10NtDnwh)

- 配置流水线到RocketMQ

![avatar](/docs/wiki/assets/pic/config_output_rocketmq.png)










