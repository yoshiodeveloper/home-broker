package postgresql

import (
	"home-broker/domain"
	"home-broker/infra"
	"time"

	"gorm.io/gorm"
)

// GORMUser is the ORM version of User entity.
// Notice that PK is not an auto-increment.
//   For this project we are assuming that users will come from an external service and the IDs are big integers.
// In the future this table can be expanded to include specific fields for this service.
type GORMUser struct {
	gorm.Model
	ID        int64          `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time      `gorm:"index:,sort:desc"`
	UpdatedAt time.Time      `gorm:"index:,sort:desc"`
	DeletedAt gorm.DeletedAt `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of User.
// It is used by GORM to perfom operations on user table (queries, migrations, etc.).
func (GORMUser) TableName() string {
	return "user"
}

// UserRepo handles database commands for user table.
type UserRepo struct {
	infra.UserRepo
	dbClient *DBClient
}

// NewUserRepo creates a new NewUserRepo.
func NewUserRepo(dbClient *DBClient) *UserRepo {
	return &UserRepo{
		dbClient: dbClient,
	}
}

// ToEntity returns an User entity from the ORM model.
func (repo *UserRepo) ToEntity(model *GORMUser) *domain.User {
	// "model.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if model.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = model.DeletedAt.Time
	}
	user := domain.User{
		ID:        domain.UserID(model.ID),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return &user
}

// GetByID gets an user from the database by an ID.
func (repo *UserRepo) GetByID(id domain.UserID) (user domain.User, found bool) {
	gormUser := GORMUser{}
	res := repo.dbClient.db.Take(&gormUser, id)
	if res.Error == gorm.ErrRecordNotFound {
		return user, false
	}
	return user, true
}
