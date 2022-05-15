package driver

// Tx implements SQL transaction method
type Tx struct {
	conn *Conn
}

// Commit the transaction on server
func (t *Tx) Commit() error {
	// Nothing to do
	return nil
}

// Rollback all changes
func (t *Tx) Rollback() error {
	// Nothing to do
	return nil
}
