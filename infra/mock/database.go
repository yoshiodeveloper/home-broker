package mock

// Database is a mocked database connection.
type Database struct {
}

// NewDatabase returns a new mocked database.
func NewDatabase() Database {
	return Database{}
}
