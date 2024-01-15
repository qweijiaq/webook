package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/qweijiaq/webook/internal/repository"
	"github.com/qweijiaq/webook/internal/repository/dao"
	"github.com/qweijiaq/webook/internal/service"
	"github.com/qweijiaq/webook/internal/web"
	"github.com/qweijiaq/webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	server.Run(":8080")
}

// initDB 初始化数据库
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/webook"))
	if err != nil {
		// 只在初始化过程中 panic
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

// initWebServer 初始化 Web 服务器
func initWebServer() *gin.Engine {
	server := gin.Default()

	// 解决跨域问题
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// ExposeHeaders:    []string{},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境
				return true
			}
			// 生产环境
			return strings.Contains(origin, "your domain name")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 实现 session 登录校验
	// 步骤1
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	// 步骤3
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())

	return server
}

// initUser 初始化 UserHandler
func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
