package benchmark

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func BenchmarkInsertMySQL(b *testing.B) {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:13308)/test_database?charset=utf8")
	if err != nil {
		b.Error(err)
	}
	// c, err := db.Conn(context.Background())
	// if err != nil {
	// 	b.Error(err)
	// }
	one := "INSERT INTO `test_database`.`users` (`name`, `address`, `age`) VALUES ('tester', 'Hangzhou', '101');"
	for i := 0; i < b.N; i++ {
		stmt, e := db.Prepare(one)
		if e != nil {
			b.Error(e)
		}
		// stmt.Exec(name, address, age)
		stmt.Exec()
	}
}
