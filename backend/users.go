package main

import (
	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListUserReq struct {
	RegionID uint `form:"region_id"`
	p.PageConfig
}

func AddUserRoutes(r *gin.Engine, bc *p.BaseConfig) {

	r.GET("/users", p.CreateListHandler[User](
		&p.Config[ListUserReq]{
			Base: bc,
			Bind: &p.BindConfig{Query: true},
		},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListUserReq) *gorm.DB {
			query = query.Preload("Region", Select("id", "name"))
			if r.RegionID != 0 {
				query = query.Where("region_id = ?", r.RegionID)
			}
			return query
		},
	))

	r.GET("/users/stats", p.Preload(
		&p.Config[struct{}]{
			Base: bc,
		},
		func(c *gin.Context, u *User, r *struct{}) {
			var items []BoardItem
			query := bc.DB.Model(new(User)).Session(new(gorm.Session))
			total := BoardItem{Label: "total"}
			if err := query.Count(&total.Value).Error; err != nil {
				c.JSON(500, bc.Resp("统计总数失败", err, nil))
				return
			}
			items = append(items, total)
			volunteers := BoardItem{Label: "volunteers"}
			if err := query.Where("region_id IS NOT null").Count(
				&volunteers.Value,
			).Error; err != nil {
				c.JSON(500, bc.Resp("统计志愿者数量失败", err, nil))
				return
			}
			items = append(items, volunteers)
			admins := BoardItem{Label: "admins"}
			if err := query.Where("admin = true").Count(
				&admins.Value,
			).Error; err != nil {
				c.JSON(500, bc.Resp("查询管理员数量失败", err, nil))
				return
			}
			items = append(items, admins)
			c.JSON(200, bc.Resp("统计成功", nil, items))
		},
	))
}
