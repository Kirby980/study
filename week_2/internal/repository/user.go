package repository

import (
	"context"
	"database/sql"

	"github.com/Kirby980/study/week_2/internal/domain"
	"github.com/Kirby980/study/week_2/internal/repository/cache"
	"github.com/Kirby980/study/week_2/internal/repository/dao"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
func (repo *UserRepository) Update(ctx context.Context, user domain.User) error {
	err := repo.dao.Edit(ctx, repo.domainToEntity(user))
	return err
}
func (repo *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	ue, err := repo.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = repo.toDomain(ue)
	go func() {
		repo.cache.Set(ctx, u)
		if err != nil {
			log.Debug("设置缓存失败")
		}
	}()
	return u, err
}

func (repo *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		Birthday: u.Birthday,
		Profile:  u.Profile,
		Nickname: u.Nickname,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Birthday: u.Birthday,
		Profile:  u.Profile,
		Nickname: u.Nickname,
	}
}
