package postgresql

import (
	"errors"
	"fmt"

	"home-broker/core/implem/postgresql"
	"home-broker/users"
	"strings"
	"time"

	"gorm.io/gorm"
)

// UserModel is the ORM version of User entity.
// Notice that the PK is not an auto-increment.
//   For this project we are assuming that users will come from an external service and the IDs are big integers.
type UserModel struct {
	gorm.Model
	ID        users.UserID   `gorm:"primaryKey;type:bigint;autoIncrement:false"`
	CreatedAt time.Time      `gorm:"not null;index:,sort:desc"`
	UpdatedAt time.Time      `gorm:"not null;index:,sort:desc"`
	DeletedAt gorm.DeletedAt `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of User.
// It is used by GORM to perfom operations on user table (queries, migrations, etc.).
func (UserModel) TableName() string {
	return "user"
}

// UserDB handles database commands for user table.
type UserDB struct {
	users.UserDBInterface
	db postgresql.DB
}

// NewUserDB creates a new UserDB.
func NewUserDB(db postgresql.DB) UserDB {
	return UserDB{db: db}
}

// ToEntity returns an User entity from the ORM model.
func (UserDB) ToEntity(model UserModel) (entity users.User) {
	// "model.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	// Notice that a "time.Time" with zero value represents a "null".
	deletedAt := time.Time{}
	if model.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = model.DeletedAt.Time
	}
	entity = users.User{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return
}

// GetByID returns an user from the database by ID.
// A nil entity will be returned if it does not exist.
func (userDB UserDB) GetByID(id users.UserID) (*users.User, error) {
	model := UserModel{}
	res := userDB.db.GetDB().Take(&model, int64(id))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := userDB.ToEntity(model)
	return &entity, nil
}

// Insert inserts a new user.
// A nil entity will be returned if an error occurs.
// The following expected error can happen: ErrUserAlreadyExists.
func (userDB UserDB) Insert(entity users.User) (*users.User, error) {
	model := UserModel{ID: entity.ID}
	res := userDB.db.GetDB().Create(&model)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "duplicate key value") {
			// Original error: "duplicate key value violates unique constraint"
			return nil, fmt.Errorf("%w: ID %d", users.ErrUserAlreadyExists, entity.ID)
		}
		return nil, res.Error
	}
	newEntity := userDB.ToEntity(model)
	return &newEntity, nil
}
