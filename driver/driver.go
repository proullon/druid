package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// init register driver into database/sql package
// allowing other package to sql.Open with it
func init() {
	sql.Register("druid", NewDriver())
}

// Driver is the driver entrypoint,
// implementing database/sql/driver interface
type Driver struct {
}

// NewDriver creates a driver object
func NewDriver() *Driver {
	d := &Driver{}
	return d
}

func (d *Driver) Open(dsn string) (conn driver.Conn, err error) {
	if dsn == "" {
		return nil, fmt.Errorf("invalid dsn")
	}

	c := NewConn(dsn+"/v2/druid/sql", "proullon/druid/1.0")
	return c, nil
}
