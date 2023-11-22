注意事项
=====
## ETCD
> 由于使用了ETCD存储mysql binlog的位置信息和全局的事件日志，ETCD是基于日志更新数据，历史版本日志不做清理会导致数据一直膨胀，所以ETCD需要定期进行压缩。可以在启动时进行如下设置，定期清理历史数据  <br>
> etcd --auto-compaction-retention=1 --max-request-bytes=$((32*1024*1024)) --quota-backend-bytes=$((8*1024*1024*1024)) <br>

> 其中，
 - --auto-compaction-retention=1，表示每隔1小时进行一次清理
 - --max-request-bytes=$((32*1024*1024)) 表示将etcd的键值大小的限制更改为32M
 - --quota-backend-bytes=$((8*1024*1024*1024)) 表示将etcd可以使用的存储空间限制更改为 8G
 - 更多限制以及设置可以参考官方文档
   - https://etcd.io/docs/v3.4/op-guide/maintenance/
   - https://etcd.io/docs/v3.4/dev-guide/limit/
