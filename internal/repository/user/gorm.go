package user

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"gorm.io/gorm"
)

type userGormRepo struct {
	orm *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userGormRepo{orm: db}
}

func (r *userGormRepo) Create(u *db.User) error {
	return r.orm.Create(u).Error
}

func (r *userGormRepo) List(limit, offset int) ([]db.User, error) {
	var users []db.User
	err := r.orm.
		Preload("Subscriptions").
		Preload("Payments").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

func (r *userGormRepo) FindByID(id string) (*db.User, error) {
	var u db.User
	err := r.orm.
		Preload("Subscriptions").
		Preload("Payments").
		First(&u, "id = ?").Error
	if err != nil {
		return nil, err
	}
	return &u, err
}

func (r *userGormRepo) Update(u *db.User) error {
	return r.orm.Save(u).Error
}

func (r *userGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.User{}, "id = ?", id).Error
}
