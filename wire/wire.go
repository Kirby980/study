//go:build wireinject

package wire

import (
	"github.com/Kirby980/study/wire/repository"
	"github.com/Kirby980/study/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	// 我只在这里声明我要用的各种东西，但是具体怎么构造，怎么编排顺序
	// 这个方法里面传入各个组件的初始化方法
	wire.Build(InitDB, repository.NewUserRepository,
		dao.NewUserDAO)
	return new(repository.UserRepository)
}
