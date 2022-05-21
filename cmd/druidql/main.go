package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/proullon/druid/cli"
	_ "github.com/proullon/druid/driver"
)

func main() {

	dsn := "http://127.0.0.1:8888"
	if ge := os.Getenv("DRUID_DSN"); ge != "" {
		dsn = ge
	}

	db, err := sql.Open("druid", dsn)
	if err != nil {
		fmt.Printf("Error : cannot open connection : %s\n", err)
		return
	}
	cli.Run(db)
}
