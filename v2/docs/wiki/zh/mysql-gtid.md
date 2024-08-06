MySQL GTID
=======

> 说明 <br>
> 1. GTID是MySQL5.6以后推出的功能。 
> 2. 云厂商的MySQL服务可能和社区版的有差异，这里主要讲MySQL官方版本，不包括云厂商和MariaDB
> 3. MySQL官网对GTID的介绍有很多，我只是看了其中涉及到主从复制和故障恢复的部分，然后记录一下这部分内容。
>  

- 为什么建议使用GTID
> 简单来说GTID就是另一种binlog日志记录的机制, 格式：
> <br>
> GTID = source_id:transaction_id
> <br>
> 其中的source_id对应的是当前MySQL实例的ID，这个ID是全局唯一的，transaction_id代表的是事务提交ID。
> <br>
> 相较于文件名加上文件内位置的方式，GTID可以保证这个记录是全局唯一的。当你有两个双主互备的MySQL，有三个从库，如果主库某一天发生了主备切换，
> <br>
> 传统方式下，你需要人工定位文件名和位置，这个过程漫长且容易出错。GTID方式则是直接使用GTID值，MySQL内部会给你处理好对应的位置。
> <br>
> Binlogo也是一个伪从库，如果你使用文件名加位置这种传统方式，需要人工判断新的主库位置。。
> <br>
> 另外从MySQL8开始，GTID的启用会提升记录binlog的效率（具体可参考官网）[MySQL GTID](https://dev.mysql.com/doc/refman/8.0/en/replication-gtids-concepts.html)

- 开启GTID
>  在配置文件中修改  gtid_mode=on
> <br>开启gtid_mode还有另一个值ON_PERMISSIVE，和ON的区别在于该数据库作为从库时对于主库binlog的是否接受的设置。具体可以参考这个文档末尾的介绍
> [replication-mode-change-online-concepts](https://dev.mysql.com/doc/refman/8.0/en/replication-mode-change-online-concepts.html)
> <br>另外ON_PERMISSIVE这个值是5.7.6以后才有的
