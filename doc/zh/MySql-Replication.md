MySql Replication （MySql的主从复制）
=========

[English](../en/Mysql-Replication.md)

> MySql的主从节点同步数据的机制叫做replication。简单来说，主节点会记录二进制日志（数据的增量变化），通过某种方式将这些数据同步给从节点，达到数据同步的目的。
> <br>
> 大量的Mysql数据同步到其他异构数据的应用场景都会利用MySql的replication实现，因为该方式可以达到准实时、稳定性也很好。
> 


