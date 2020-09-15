package postgresql

import (
	"errors"
	"fmt"
	"home-broker/domain"
	"home-broker/infra"
	"strings"
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
	CreatedAt time.Time      `gorm:"not null;index:,sort:desc"`
	UpdatedAt time.Time      `gorm:"not null;index:,sort:desc"`
	DeletedAt gorm.DeletedAt `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of User.
// It is used by GORM to perfom operations on user table (queries, migrations, etc.).
func (UserGORM) TableName() string {
	return "user"
}

// UserRepo handles database commands for user table.
type UserRepo struct {
	infra.UserRepoI
	dbClient *DBClient
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(dbClient *DBClient) *UserRepo {
	return &UserRepo{
		dbClient: dbClient,
	}
}

// ToEntity returns an User entity from the ORM model.
func (repo *UserRepo) ToEntity(modelGORM *UserGORM) (entity domain.User) {
	// "userGORM.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if modelGORM.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = modelGORM.DeletedAt.Time
	}
	entity = domain.User{
		ID:        domain.UserID(modelGORM.ID),
		CreatedAt: modelGORM.CreatedAt,
		UpdatedAt: modelGORM.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return
}

// GetByID returns an user from the database by ID.
func (repo *UserRepo) GetByID(id domain.UserID) (*domain.User, error) {
	modelGORM := UserGORM{}
	res := repo.dbClient.db.Take(&modelGORM, int64(id))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := repo.ToEntity(&modelGORM)
	return &entity, nil
}

// Insert inserts a new user.
func (repo *UserRepo) Insert(entity domain.User) (*domain.User, error) {
	modelGORM := UserGORM{ID: int64(entity.ID)}
	res := repo.dbClient.db.Create(&modelGORM)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "duplicate key value") {
			// Original error: "duplicate key value violates unique constraint"
			return nil, fmt.Errorf("%w: ID %d", infra.ErrUserAlreadyExists, entity.ID)
		}
		return nil, res.Error
	}
	newEntity := repo.ToEntity(&modelGORM)
	return &newEntity, nil
}
