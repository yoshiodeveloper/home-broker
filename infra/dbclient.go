package infra

// DBClient is an interface that handles database connections.
type DBClient interface {

	// Open must connect to the database.
	Open() error

	// SetDB must set a pointer to a database connection. It is ORM specific.
	SetDB(db *interface{})

	// GetDB must returns a pointer to a database connection. It is ORM specific.
	GetDB() *interface{}
}
