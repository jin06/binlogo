Attention Points
====
### ETCD
> Due to the utilization of ETCD for storing MySQL binlog position information and global event logs, and as ETCD updates data based on logs, failure to clear historical log files regularly can result in > continuous data expansion. Therefore, it is necessary to periodically compress ETCD. You can perform the following settings during startup for regular cleaning of historical data: <br>
> etcd --auto-compaction-retention=1 --max-request-bytes=$((32*1024*1024)) --quota-backend-bytes=$((8*1024*1024*1024)) <br>

> Where,

- --auto-compaction-retention=1 indicates cleaning every 1 hour.
- --max-request-bytes=$((32*1024*1024)) sets the limit of key-value size to 32M for ETCD.
- --quota-backend-bytes=$((8*1024*1024*1024)) sets the limit of storage space that ETCD can use to 8G.
- For more restrictions and settings, refer to the official documentation:
  - https://etcd.io/docs/v3.4/op-guide/maintenance/
  - https://etcd.io/docs/v3.4/dev-guide/limit/
