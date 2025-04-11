package service

import (
	"context"
	"errors"

	"github.com/Kirby980/study/week_2/internal/domain"
	"github.com/Kirby980/study/week_2/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	err := svc.repo.Update(ctx, u)
	return err
}

func (svc *UserService) Select(ctx context.Context, u int64) (domain.User, error) {
	user, err := svc.repo.FindByID(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (svc *UserService) FindORCreate(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, phone)
	if err == nil {
		return user, nil
	}
	if err != repository.ErrUserNotFound {
		return domain.User{}, err
	}
	u := domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	return u, nil
}
