package main

import (
	"database/sql"
	"fmt"

	"github.com/proullon/druid/cli"
	_ "github.com/proullon/druid/driver"
)

func main() {
	db, err := sql.Open("druid", "")
	if err != nil {
		fmt.Printf("Error : cannot open connection : %s\n", err)
		return
	}
	cli.Run(db)
}
