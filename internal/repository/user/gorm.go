package user

import (
	"errors"
	"strings"

	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrDuplicateTGID = errors.New("duplicate tg_id")
)

type userGormRepo struct {
	orm *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userGormRepo{orm: db}
}

func (r *userGormRepo) Create(u *db.User) error {
	// Генерируем UUID если он не установлен
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	err := r.orm.Create(u).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// по желанию можно ещё дополнительно проверить имя constraint:
			if strings.Contains(pgErr.ConstraintName, "users_tg_id_key") {
				return ErrDuplicateTGID
			}
			return ErrDuplicateTGID
		}
	}
	return err
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
		First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, err
}

func (r *userGormRepo) FindByTGID(tgID int64) (*db.User, error) {
	var u db.User
	err := r.orm.
		Preload("Subscriptions").
		Preload("Payments").
		First(&u, "tg_id = ?", tgID).
		Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userGormRepo) Update(u *db.User) error {
	return r.orm.Save(u).Error
}

func (r *userGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.User{}, "id = ?", id).Error
}
