注意事项
=====
 - ETCD需要配置定时清理历史版本和ETCD数据库大小，可以使用下列配置
> etcd --auto-compaction-mode=revision \ <br>
--auto-compaction-retention=1000 \ <br>
--max-request-bytes=$((32*1024*1024)) \ <br>
--quota-backend-bytes=$((8*1024*1024*1024)) <br>
