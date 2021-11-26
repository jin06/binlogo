
#### Insert

> INSERT INTO `test_database`.`users` (`id`, `name`, `address`, `age`) VALUES ('1', 'roy', 'abcdefgh', '10');

```json
{
    "head":{
        "type":"insert",
        "time":1637551072,
        "database":"test_database",
        "table":"users",
        "position":{
            "binlog_file":"mysql-bin.000004",
            "binlog_position":13547,
            "gtid_set":"045c649a-408d-11ec-ae21-0242ac110006:1-50",
            "pipeline_name":"gtid-mode"
        }
    },
    "data":{
        "new":{
            "address":"abcdefgh",
            "age":10,
            "id":1,
            "name":"roy"
        }
    }
}
```

#### Update

> UPDATE `test_database`.`users` SET `address` = 'china' WHERE (`id` = '1');

```json
{
    "head":{
        "type":"update",
        "time":1637551289,
        "database":"test_database",
        "table":"users",
        "position":{
            "binlog_file":"mysql-bin.000004",
            "binlog_position":14147,
            "gtid_set":"045c649a-408d-11ec-ae21-0242ac110006:1-52",
            "pipeline_name":"gtid3"
        }
    },
    "data":{
        "old":{
            "address":"abcdefgh",
            "age":10,
            "id":1,
            "name":"roy"
        },
        "new":{
            "address":"china",
            "age":10,
            "id":1,
            "name":"roy"
        }
    }
}
```

#### Delete

> DELETE FROM `test_database`.`users` WHERE (`id` = '1');

```json
{
    "head":{
        "type":"delete",
        "time":1637551359,
        "database":"test_database",
        "table":"users",
        "position":{
            "binlog_file":"mysql-bin.000004",
            "binlog_position":14434,
            "gtid_set":"045c649a-408d-11ec-ae21-0242ac110006:1-53",
            "pipeline_name":"gtid-mode"
        }
    },
    "data":{
        "old":{
            "address":"china",
            "age":10,
            "id":1,
            "name":"roy"
        }
    }
}
```

