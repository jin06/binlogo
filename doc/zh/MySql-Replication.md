MySql Replication （MySql的主从复制）
=========

[English](../en/Mysql-Replication.md)

> MySql的主从节点同步数据的机制叫做replication。简单来说，主节点会记录二进制日志（数据的增量变化），通过某种方式将这些数据同步给从节点，达到数据同步的目的。
> 
> 大量的Mysql数据同步到其他异构数据的应用场景都会利用MySql的replication实现，因为该方式可以达到准实时、稳定性也很好。
> 其基本原理是，MySql服务器会把数据的变更、表结构的变更记录到日志文件中，也就是binlog，从库只要拿到这些日志文件就可以将主库的数据同步过来。
> 同样的，基于binlog日志的一些围绕MySql的中间件通过获取binlog就可以得到MySql的数据变化的信息。
> 
> 做一个获取MySql增量信息的工具就要稍微了解一下MySql的binlog以及同步的方式。
>> 例如主库和从库同步数据的时候通信过程是怎么样的，怎么建立的连接，是主库主动发送binlog还是从库主动拉取。
>
>> binlog是二进制文件，不能像普通的文本文件一样直接读取。其实也不是十分必要了解它的物理组织方式，因为已经有很多开源代码可以做解析。
> 
>> 要了解BinlogEvent，这是一个重要的逻辑概念，BinlogEvent就是binlog记录的内容，要了解它的基本概念和有哪些版本。不同的MySql版本之间的差异。
> 
>> Mysql 配置相关的东西，如开启binlog记录，记录的方式。
> 
>> 伪装成从库的工具对于MySql主库的影响，比如MySql主从复制中会有半同步复制

> 总结一下，写这个工具要考虑的问题就是：1. binlog的基本原理 2. 不同MySql版本之间的差异（binlogEvent、GTID）3. Mysql主库要做的配置  4. 工具对于Mysql的影响


### Binlog 基本原理
#### 1.1 Mysql建立连接的过程，同步binlog日志的过程

### 不同Mysql之间的差异
#### 2.1 市面上的Mysql版本
##### MariaDB 与 Mysql 

> 

##### Mysql 各个版本
> 主流云厂商(全球前4)对于mysql的支持：
> -  亚马逊 AWS:
> -  微软 Azure:
> -  谷歌 : 
> -  阿里云: 的Mysql实例支持5.5 到 8之间的版本。
> 
> 没有查询到mysql各个版本用户存量数据，查询发行时间，Mysql5.5 发行时间在2009年，结合云厂商的数据，考虑这个数据同步的项目支持的Mysql最低版本为5.5 。 
> 

#### 2.2 GTID
#### 2.3 binglog event 版本

### Mysql的配置
#### 3.1 binlog日志格式
> Mysql 的三种binlog日志格式statement、row、mixed。
> 
> statement是记录sql语句；row模式记录每条数据具体的变化；mixed介于两者之间，会根据不同的sql在statement和row之间选择。
> 
> 所以需要将binlog的日志格式配置成row，因为这种模式的binlog会记录每一行数据的变化。但是也有一些情况是记录不到的，例如在变更数据表结构时添加新的字段并赋予一个默认值的情况就只会记录这个sql，而不会将以往数据的变化都记录一份。
> 

### 工具对于Mysql的影响、风险
#### 4.1 一个伪装成slave的中间件对于原有的Mysql主从复制的影响

