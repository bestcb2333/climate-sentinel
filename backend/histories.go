package main

import (
	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListHistoryReq struct {
	RegionID uint   `form:"region_id"`
	Future   bool   `form:"future"`
	Type     string `form:"type"`
	p.PageConfig
}

func AddHistoryRoutes(r *gin.Engine, bc *p.BaseConfig) {

	r.GET("/histories", p.CreateListHandler[History](
		&p.Config[ListHistoryReq]{
			Base: bc,
			Bind: &p.BindConfig{Query: true},
		},
		func(query *gorm.DB, c *gin.Context, u *User, r *ListHistoryReq) *gorm.DB {
			query = query.Preload("Region", Select("id", "name"))
			if r.RegionID != 0 {
				query = query.Where("region_id = ?", r.RegionID)
			}
			if r.Future {
				query = query.Where("time > NOW()")
			}
			if r.Type != "all" {
				query = query.Where("type = ?", r.Type)
			}
			return query
		},
	))

	r.GET("/stats/histories", p.Preload(
		&p.Config[struct{}]{},
		func(c *gin.Context, u *User, r *struct{}) {

			var items []BoardItem
			query := bc.DB.Model(new(History)).Session(new(gorm.Session))

			weathers := []struct {
				ID    string
				Label string
			}{
				{ID: "sunny", Label: "晴天"},
				{ID: "rainy", Label: "雨天"},
				{ID: "cloudy", Label: "多云"},
				{ID: "foggy", Label: "雾天"},
				{ID: "snowy", Label: "下雪"},
				{ID: "windy", Label: "大风"},
				{ID: "overcast", Label: "阴天"},
			}

			for _, weather := range weathers {
				item := BoardItem{Label: weather.Label}
				if err := query.Where("type = ?", weather.ID).Count(
					&item.Value,
				).Error; err != nil {
					c.JSON(500, bc.Resp("统计失败", err, nil))
					return
				}
				items = append(items, item)
			}

			c.JSON(200, bc.Resp("统计完成", nil, items))

		},
	))

	r.POST("/histories", p.CreateAddHandler[History](
		&p.Config[HistoryDTO]{
			Base: bc,
		},
		func(c *gin.Context, u *User, dto *HistoryDTO) *History {
			data := new(History)
			data.HistoryDTO = dto
			return data
		},
	))
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
