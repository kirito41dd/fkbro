package data

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zshorz/ezlog"
	"os"
)

var DB *sql.DB
var Log *ezlog.EzLogger

func init() {
	Log = ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll)
}

func Setup(dbaddr, dbname, dbuser, dbpasswd string) {
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbuser, dbpasswd, dbaddr, dbname)
	var err error
	DB, err = sql.Open("mysql", str)
	if err != nil {
		Log.Fatal(err)
	}

	_, err = DB.Query("select * from alert where id = 0")
	if err != nil {
		Log.Error(err)
		Log.Fatal("please check database")
	}

	Log.Info("sql set up success")
}

