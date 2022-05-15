package driver

import (
	"database/sql/driver"
	"net/http"
)

// Conn implements sql/driver Conn interface
type Conn struct {
	// Druid router endpoint
	endpoint string
	// User-Agent used in request
	userAgent string
	// Conn http client
	client *http.Client
}

// NewConn initialise basically an http.Client wrapper
func NewConn(endpoint string, userAgent string) driver.Conn {
	return &Conn{
		endpoint:  endpoint,
		userAgent: userAgent,
		client:    &http.Client{},
	}
}

// Begin starts and returns a new transaction.
func (c *Conn) Begin() (driver.Tx, error) {

	tx := Tx{
		conn: c,
	}

	return &tx, nil
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (c *Conn) Close() error {

	return nil
}

// Prepare returns a prepared statement, bound to this connection.
func (c *Conn) Prepare(query string) (driver.Stmt, error) {

	stmt := prepareStatement(c.endpoint, c.userAgent, query, c.client)

	return stmt, nil
}
