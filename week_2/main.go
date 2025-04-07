package main

import (
	"strings"
	"time"

	"github.com/Kirby980/study/week_2/internal/repository"
	"github.com/Kirby980/study/week_2/internal/repository/dao"
	"github.com/Kirby980/study/week_2/internal/service"
	"github.com/Kirby980/study/week_2/internal/web"
	"github.com/Kirby980/study/week_2/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()

	server := initWebServer()
	initUserHdl(db, server)
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowHeaders:     []string{"Content-Type", "authorization"},
		//AllowHeaders: []string{"content-type"},
		//AllowMethods: []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://192.168.3.97") || strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的 Middleware")
	})

	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	//store := cookie.NewStore([]byte("secret"))
	//基于内存的实现
	//store := memstore.NewStore([]byte("secret"))
	//存到redis
	store, err := redis.NewStore(16, "tcp", "192.168.3.97:6379", "", "", []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"), []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
	if err != nil {
		panic(err)
	}
	//未使用jwt
	//server.Use(sessions.Sessions("ssid", store), login.IgnorePaths("/users/login").IgnorePaths("/users/signup").Build())
	// 使用jwt
	server.Use(sessions.Sessions("ssid", store), middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signuo").IgnorePaths("/users/login").Build())
	return server
}
