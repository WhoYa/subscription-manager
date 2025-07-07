package user

import "github.com/WhoYa/subscription-manager/pkg/db"

type UserRepository interface {
	Create(u *db.User) error
	List(limit, offset int) ([]db.User, error)
	FindByID(id string) (*db.User, error)
	Update(u *db.User) error
	Delete(id string) error
}
