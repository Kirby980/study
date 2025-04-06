package repository

import (
	"context"
	"time"

	"github.com/Kirby980/study/week_2/internal/domain"
	"github.com/Kirby980/study/week_2/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
func (repo *UserRepository) Update(ctx context.Context, user domain.User) error {
	err := repo.dao.Edit(ctx, dao.User{
		Id:       user.Id,
		Email:    user.Email,
		Birthday: user.Birthday,
		Profile:  user.Profile,
		Nickname: user.Nickname,
		Utime:    time.Now().UnixMilli(),
	})
	return err
}
func (repo *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday,
		Profile:  u.Profile,
		Nickname: u.Nickname,
	}
}
