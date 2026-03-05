package infra

// DB simulates a database connection handle.
type DB struct {
	DSN string
}

// Connect opens a database connection.
func Connect(dsn string) *DB {
	return &DB{DSN: dsn}
}
