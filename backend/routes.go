package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListRouteDTO struct {
	RegionID uint `form:"regionId"`
	ListDTO
}

func AddRouteRoutes(r *gin.Engine, pbc *PreloaderBaseConfig) {

	r.GET("/routes", CreateListHandler[Route](
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
			Bind:                BindConfig{Query: true},
		},

		&ListRouteDTO{0, ListDTO{1, 10}},
		func(query *gorm.DB, c *gin.Context, u *User, dto *ListRouteDTO) *gorm.DB {
			query = query.Preload("Region", Select("id", "name"))
			if dto.RegionID != 0 {
				query = query.Where("region_id = ?", dto.RegionID)
			}

			return query
		},
	))

	r.POST("/routes", CreateAddHandler[Route](
		&PreloaderConfig{
			PreloaderBaseConfig: pbc,
		},
		new(RouteDTO),
		func(data *Route, u *User, dto *RouteDTO) *Route {
			return data
		},
	))

	r.DELETE("/routes", CreateDeleteHandler[Route](&PreloaderConfig{
		PreloaderBaseConfig: pbc,
		Bind:                BindConfig{Query: true},
	}))
}
