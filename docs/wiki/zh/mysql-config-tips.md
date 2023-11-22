MySQL Configuration （MySQL配置）
=======


- MySQL配置文件地址
> MySQL配置文件一般存放在 /etc/my.cnf /etc/mysql/my.cnf ~/.my.cnf
> </br> 
> 找不到使用命令 find / -name my.cnf
> 
> 

- 开启binlog记录
> 在配置文件my.cnf中添加配置项 log_bin=mysql-bin
> <br>
> 上述配置中的mysql-bin为binlog的日志名称，可以修改为其他的值
> 

- 修改binlog的记录类型为ROW
> 在配置文件中修改  binlog_format=ROW

- 创建一个有同步权限的MySQL用户
```code
mysql> CREATE USER 'repl'@'%.example.com' IDENTIFIED BY 'password';
mysql> GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%.example.com'; 
```

- 重启MySQL后生效

- 检查是否生效, 一般使用 show master status 能显示出日志记录的状态就是成功了
```code
-- 查看配置项目是否成功
show variables like "%log_bin%";
-- 查看是否已经开始记录binlog日志
show master status;
-- 查看binlog日志
show binlog events;
```

