package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
)

type ListResourceReq struct {
	Type      string `form:"type"`
	RegionID  uint   `form:"region_id"`
	Available bool   `form:"available"`
	p.PageConfig
}

type RadarResp struct {
	Value []uint `json:"value"`
	Name  string `json:"name"`
}

func AddResourceRoutes(r *gin.Engine, bc *p.BaseConfig) {

	group := r.Group("/resources")

	group.GET("", p.CreateListHandler[Resource](
		&p.Config[ListResourceReq]{Base: bc},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListResourceReq) *gorm.DB {
			query = query.Preload("Region", Select("id", "name"))
			if r.Type != "all" {
				query = query.Where("type = ?", r.Type)
			}
			if r.RegionID != 0 {
				query = query.Where("region_id = ?", r.RegionID)
			}
			if r.Available {
				query = query.Where("available = true")
			}
			return query
		},
	))

	group.GET("/stats", p.Preload(
		&p.Config[struct{}]{Base: bc},
		func(c *gin.Context, u *User, r *struct{}) {

			var items []BoardItem
			query := bc.DB.Model(new(Resource)).Session(new(gorm.Session))

			total := BoardItem{Label: "资源总数"}
			if err := query.Count(&total.Value).Error; err != nil {
				c.JSON(500, bc.Resp("统计总数失败", err, nil))
				return
			}
			items = append(items, total)

			available := BoardItem{Label: "可用资源数"}
			if err := query.Where("available = true").Count(
				&available.Value,
			).Error; err != nil {
				c.JSON(500, bc.Resp("统计可用数量失败", err, nil))
				return
			}
			items = append(items, available)

			items = append(items, BoardItem{
				Label: "不可用资源数",
				Value: total.Value - available.Value,
			})

			items = append(items, BoardItem{
				Label: "资源种类数",
				Value: 6,
			})

			c.JSON(200, bc.Resp("统计完成", nil, items))

		},
	))

	group.POST("", p.CreateAddHandler(
		&p.Config[ResourceDTO]{Base: bc},
		func(c *gin.Context, u *User, dto *ResourceDTO) *Resource {
			data := new(Resource)
			data.ResourceDTO = dto
			return data
		},
	))

	group.GET("/radar", p.Preload(
		&p.Config[struct{}]{
			Base: bc,
			Bind: &p.BindConfig{Query: true},
		},
		func(c *gin.Context, u *User, r *struct{}) {

			var regions []Region

			if err := bc.DB.Select("id", "name").Find(&regions).Error; err != nil {
				c.JSON(500, Resp("获取区域列表失败", err, nil))
				return
			}

			indicators := []struct {
				Label  string
				Column string
			}{
				{"救援车辆", "vehicle"},
				{"救援人员", "personnel"},
				{"通讯工具", "comm"},
				{"天气监测设备", "weather"},
				{"道路维护设备", "maintain"},
				{"其他资源", "other"},
			}

			var datas []RadarResp
			for _, region := range regions {
				data := RadarResp{
					Value: make([]uint, 0),
					Name:  region.Name,
				}
				for _, indicator := range indicators {
					var value uint
					if err := bc.DB.Model(new(Resource)).Where(
						"region_id = ? AND type = ?", region.ID, indicator.Column,
					).Select("SUM(quantity)").Scan(&value).Error; err != nil {
						c.JSON(500, Resp("统计数量失败", err, nil))
						return
					}
					data.Value = append(data.Value, value)
				}
				datas = append(datas, data)
			}

			c.JSON(200, Resp("查询成功", nil, datas))
		},
	))

	group.DELETE("", p.CreateDeleteHandler[Resource, User](
		&p.Config[p.DelReq]{Base: bc},
	))

}
