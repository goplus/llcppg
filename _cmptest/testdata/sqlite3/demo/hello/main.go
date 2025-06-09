package main

import (
	"sqlite3"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/os"
)

func main() {
	os.Remove(c.Str("test.db"))

	var db *sqlite3.Sqlite3
	err := sqlite3.DoOpen(c.Str("test.db"), &db)
	check(err, db, "sqlite: Open")

	err = db.Exec(c.Str("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)"), nil, nil, nil)
	check(err, db, "sqlite: Exec CREATE TABLE")

	var stmt *sqlite3.Stmt
	sql := "INSERT INTO users (id, name) VALUES (?, ?)"
	err = db.DoPrepareV3(c.GoStringData(sql), c.Int(len(sql)), 0, &stmt, nil)
	check(err, db, "sqlite: PrepareV3 INSERT")

	stmt.BindInt(1, 100)
	stmt.BindText(2, c.Str("Hello World"), -1, nil)

	err = stmt.Step()
	checkDone(err, db, "sqlite: Step INSERT 1")

	stmt.Reset()
	stmt.BindInt(1, 200)
	stmt.BindText(2, c.Str("This is llgo"), -1, nil)

	err = stmt.Step()
	checkDone(err, db, "sqlite: Step INSERT 2")

	stmt.Close()

	sql = "SELECT * FROM users"
	err = db.DoPrepareV3(c.GoStringData(sql), c.Int(len(sql)), 0, &stmt, nil)
	check(err, db, "sqlite: PrepareV3 SELECT")

	for {
		if err = stmt.Step(); err != sqlite3.ROW {
			break
		}
		c.Printf(c.Str("==> id=%d, name=%s\n"), stmt.ColumnInt(0), stmt.ColumnText(1))
	}
	checkDone(err, db, "sqlite: Step done")

	stmt.Close()
	db.Close()
}

func check(err c.Int, db *sqlite3.Sqlite3, at string) {
	if err != sqlite3.OK {
		c.Printf(c.Str("==> %s Error: (%d) %s\n"), c.AllocaCStr(at), err, db.Errmsg())
		c.Exit(1)
	}
}

func checkDone(err c.Int, db *sqlite3.Sqlite3, at string) {
	if err != sqlite3.DONE {
		check(err, db, at)
	}
}
