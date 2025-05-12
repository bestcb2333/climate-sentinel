package main

import (
	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListEventReq struct {
	RegionID uint   `form:"region_id"`
	Type     string `form:"type"`
	p.PageConfig
}

func AddEventRoutes(r *gin.Engine, bc *p.BaseConfig) {

	r.GET("/events", p.CreateListHandler[Event](
		&p.Config[ListEventReq]{
			Base:   bc,
			DefReq: ListEventReq{Type: "all"},
		},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListEventReq) *gorm.DB {
			query = query.Preload("User", Select("id", "name")).Preload("Region", Select("id", "name"))
			if r.RegionID != 0 {
				query = query.Where("region_id = ?", r.RegionID)
			}
			if r.Type != "all" {
				query = query.Where("type = ?", r.Type)
			}
			return query
		},
	))

	r.GET("/events/stats", p.Preload(
		&p.Config[struct{}]{Base: bc},
		func(c *gin.Context, u *User, r *struct{}) {
			query := bc.DB.Model(new(Event)).Session(new(gorm.Session))
			var items []BoardItem
			for value, label := range []string{"safe", "low", "medium", "high"} {
				item := BoardItem{Label: label}
				if err := query.Where("severity = ?", value).Count(&item.Value).Error; err != nil {
					c.JSON(500, bc.Resp("统计失败", err, nil))
					return
				}
				items = append(items, item)
			}
			for _, typ := range []string{
				"blizzard", "typhoon", "hail", "fog", "thunder", "others",
			} {
				item := BoardItem{Label: typ}
				if err := query.Where("type = ?", typ).Count(&item.Value).Error; err != nil {
					c.JSON(500, bc.Resp("统计失败", err, nil))
					return
				}
				items = append(items, item)
			}
			c.JSON(200, bc.Resp("统计成功", nil, items))
		},
	))

	r.POST("/events", p.CreateAddHandler[Event](
		&p.Config[EventDTO]{
			Base: bc,
		},
		func(c *gin.Context, u *User, dto *EventDTO) *Event {
			data := new(Event)
			data.EventDTO = dto
			return data
		},
	))
}
