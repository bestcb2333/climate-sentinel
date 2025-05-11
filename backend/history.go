package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListHistoryDTO struct {
	RegionID uint `form:"regionId"`
	ListDTO
}

func AddHistoryRoutes(r *gin.Engine, pbc *PreloaderBaseConfig) {

	r.GET("/histories", CreateListHandler[History](
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
			Bind:                BindConfig{Query: true},
		},
		&ListHistoryDTO{0, ListDTO{1, 10}},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListHistoryDTO) *gorm.DB {
			query = query.Preload("Region", Select("id", "name"))
			if r.RegionID != 0 {
				query = query.Where("region_id = ?", r.RegionID)
			}
			return query
		},
	))

	r.POST("/histories", CreateAddHandler[History](
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
			Bind:                BindConfig{JSON: true},
		},
		new(HistoryDTO),
		func(data *History, u *User, dto *HistoryDTO) *History {
			return data
		},
	))

	r.DELETE("/histories", CreateDeleteHandler[History](&PreloaderConfig{
		PreloaderBaseConfig: pbc,
		Bind:                BindConfig{Query: true},
	}))
}

type Trend struct {
	Name   string           `json:"name"`
	Series []LineChartSerie `json:"series"`
}

type LineChartSerie struct {
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Stack string    `json:"stack"`
	Data  []float64 `json:"data"`
}

func AddHistoryTrendRoutes(r *gin.Engine, pbc *PreloaderBaseConfig) {

	r.GET("/trends", func(c *gin.Context) {

		var trends []Trend

		var weatherIndicators = []struct {
			Column string
			Label  string
		}{
			{"max_temperature", "最高气温"},
			{"min_temperature", "最低气温"},
			{"avg_temperature", "平均气温"},
			{"wind_speed", "风速"},
			{"visibility", "能见度"},
			{"rain_fall", "降水量"},
		}

		var regions []Region

		if err := pbc.DB.Select("id", "name").Find(&regions).Error; err != nil {
			c.JSON(500, Resp("获取区域列表失败", err, nil))
			return
		}

		for _, indicator := range weatherIndicators {
			trend := Trend{
				Name:   indicator.Label,
				Series: make([]LineChartSerie, 0),
			}
			for _, region := range regions {
				serie := LineChartSerie{
					Name: region.Name,
					Type: "line",
				}
				if err := pbc.DB.Model(new(History)).Where(
					"region_id = ?", region.ID,
				).Pluck(
					indicator.Column, &serie.Data,
				).Order("time desc").Limit(7).Error; err != nil {
					c.JSON(500, Resp("查询特定数据失败", err, nil))
					return
				}
				trend.Series = append(trend.Series, serie)
			}
			trends = append(trends, trend)
		}

		c.JSON(200, Resp("查询成功", nil, trends))
	})
}
