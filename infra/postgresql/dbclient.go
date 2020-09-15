package postgresql

import (
	"fmt"
	"home-broker/infra"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBClient handles database connections.
type DBClient struct {
	infra.DBClient
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	db       *gorm.DB
}

// NewDBClient creates a new database client.
func NewDBClient(host string, port int, user string, password string, dbName string, sslMode string) *DBClient {
	return &DBClient{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}

// SetDB sets a pointer to a database connection.
func (dbClient *DBClient) SetDB(db *gorm.DB) {
	dbClient.db = db
}

// Open connects to the database.
func (dbClient *DBClient) Open() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbClient.Host, dbClient.Port, dbClient.User, dbClient.Password, dbClient.DBName, dbClient.SSLMode)
	/*
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // Disable color
			},
		)
	*/
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	dbClient.SetDB(db)
	return err
}

// GetDB returns a pointer to a GORM DB connection.
func (dbClient *DBClient) GetDB() *gorm.DB {
	return dbClient.db
}
