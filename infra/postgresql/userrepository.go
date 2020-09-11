package postgresql

import (
	"broker-dealer/domain"
	"broker-dealer/infra"
	"time"

	"gorm.io/gorm"
)

// UserGORM is the ORM version of User entity.
// Notice that PK is not an auto-increment.
//   For this project we are assuming that users will come from an external service and the IDs are big integers.
// In the future this table can be expanded to include specific fields for this service.
type UserGORM struct {
	gorm.Model
	ID        int64          `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time      `gorm:"index:,sort:desc"`
	UpdatedAt time.Time      `gorm:"index:,sort:desc"`
	DeletedAt gorm.DeletedAt `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of User.
// It is used by GORM to perfom operations on user table (queries, migrations, etc.).
func (UserGORM) TableName() string {
	return "user"
}

// UserRepository handles database commands for user table.
type UserRepository struct {
	infra.UserRepository
	dbClient *DBClient
}

// NewUserRepository creates a new NewUserRepository.
func NewUserRepository(dbClient *DBClient) UserRepository {
	return UserRepository{
		dbClient: dbClient,
	}
}

// GetDBClient returns a pointer to the DBClient.
func (repo *UserRepository) GetDBClient() *DBClient {
	return repo.dbClient
}

// GetByID gets an user from the database by an ID.
func (repo *UserRepository) GetByID(id int64) (user domain.User, found bool, err error) {
	userGORM := UserGORM{}
	res := repo.dbClient.db.Take(&userGORM, id)
	if res.Error == gorm.ErrRecordNotFound {
		return user, false, nil
	}
	if res.Error != nil {
		return user, false, res.Error
	}
	user = domain.User{
		ID:        userGORM.ID,
		CreatedAt: userGORM.CreatedAt,
		UpdatedAt: userGORM.UpdatedAt,
	}
	return user, true, nil
}
