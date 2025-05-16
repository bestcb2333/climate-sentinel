package main

import (
	"time"

	"github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRouter(db *gorm.DB, config *Config) *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	pbc := &PreloaderBaseConfig{
		DB:     db,
		JWTKey: config.JWTKey,
	}

	bc := &preloader.BaseConfig{
		DB:            db,
		UserTableName: "users",
		AdminColName:  "admin",
		JWTKey:        config.JWTKey,
		Resp: func(message string, err error, data any) gin.H {
			var errStr *string
			if err != nil {
				str := err.Error()
				errStr = &str
			}
			return gin.H{
				"message": message,
				"error":   errStr,
				"data":    data,
			}
		},
	}

	r.GET("/ping", Ping)
	r.GET("/captcha", GetCaptcha)

	RegGetMyinfoHandler(r, bc)
	AddSignupRoutes(r, bc)
	AddLoginRoutes(r, bc)
	AddHistoryTrendRoutes(r, pbc)

	AddUserRoutes(r, bc)
	AddRegionRoutes(r, bc)
	AddEventRoutes(r, bc)
	AddHistoryRoutes(r, bc)
	AddNoticeRoutes(r, pbc)
	AddResourceRoutes(r, bc)
	AddRouteRoutes(r, pbc)
	AddSendEmailRoutes(r, bc, &config.SMTP)

	return r
}
