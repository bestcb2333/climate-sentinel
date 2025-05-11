package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
)

type ListResourceDTO struct {
	Type   string `form:"type"`
	Region uint   `form:"region"`
	ListDTO
}

func AddResourceRoutes(r *gin.Engine, pbc *PreloaderBaseConfig) {

	p.A

	r.GET("/resources", CreateListHandler[Resource](
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
			Bind:                BindConfig{Query: true},
		},
		&ListResourceDTO{"", 0, ListDTO{1, 10}},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListResourceDTO) *gorm.DB {
			query = query.Preload("Region", Select("id", "name"))
			if r.Type != "" {
				query = query.Where("type = ?", r.Type)
			}
			if r.Region != 0 {
				query = query.Where("region_id = ?", r.Region)
			}
			return query
		},
	))

	r.POST("/resources", CreateAddHandler(
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
			Bind:                BindConfig{JSON: true},
		},
		new(ResourceDTO),
		func(data *Resource, u *User, dto *ResourceDTO) *Resource {
			data.ResourceDTO = *dto
			return data
		},
	))

	r.DELETE("/resources", CreateDeleteHandler[Resource](&PreloaderConfig{
		PreloaderBaseConfig: pbc,
		Bind:                BindConfig{Query: true},
	}))
}

type RadarResp struct {
	Value []uint `json:"value"`
	Name  string `json:"name"`
}

func AddResourceRadarRoutes(r *gin.Engine, pbc *PreloaderBaseConfig) {
	r.GET("/radar", Preload(
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
			Bind:                BindConfig{Query: true},
		},
		&struct{}{},
		func(c *gin.Context, u *User, r *struct{}) {

			var regions []Region

			if err := pbc.DB.Select("id", "name").Find(&regions).Error; err != nil {
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
					if err := pbc.DB.Model(new(Resource)).Where(
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
}
