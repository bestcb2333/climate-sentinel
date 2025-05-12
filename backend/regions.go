package main

import (
	"time"

	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListRegionReq struct {
	History bool `form:"history"`
	p.PageConfig
}

func AddRegionRoutes(r *gin.Engine, bc *p.BaseConfig) {

	r.GET("/regions", p.CreateListHandler[Region](
		&p.Config[ListRegionReq]{
			Base: bc,
			DefReq: ListRegionReq{PageConfig: p.PageConfig{
				Page:     1,
				PageSize: 100,
			}},
		},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListRegionReq) *gorm.DB {
			if r.History == true {
				query = query.Preload("Histories", "time > ?", time.Now()).Select("id", "name")
			}
			return query
		},
	))
}
