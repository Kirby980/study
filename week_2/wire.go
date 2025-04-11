//go:build wireinject

package main

import (
	"github.com/Kirby980/study/week_2/internal/repository"
	"github.com/Kirby980/study/week_2/internal/repository/cache"
	"github.com/Kirby980/study/week_2/internal/repository/dao"
	"github.com/Kirby980/study/week_2/internal/service"
	"github.com/Kirby980/study/week_2/internal/web"
	"github.com/Kirby980/study/week_2/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	// 这里是 wire 的入口
	// 你可以在这里声明你需要的各种组件
	// 然后 wire 会根据你声明的组件，自动生成代码
	// 你只需要在这里调用 wire.Build() 方法
	// wire.Build(
	// 	InitRedis,
	// 	InitSMSService,
	// )
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// 初始化 DAO
		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		// 直接基于内存实现
		ioc.InitSMSService,
		web.NewUserHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)

}
