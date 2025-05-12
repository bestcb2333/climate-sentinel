package main

import (
	"time"

	"github.com/paulmach/orb"
	"gorm.io/gorm"
)

var Tables = []any{
	new(User),
	new(Region),
	new(Event),
	new(History),
	new(Resource),
	new(Notice),
	new(Route),
}

type ListDTO struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func (ld ListDTO) GetListDTO() *ListDTO {
	return &ld
}

type DeleteDTO struct {
	IDs []uint `form:"id"`
}

type IDField struct {
	ID uint `json:"id" gorm:"primarykey;comment:ID"`
}

type CreatedAtField struct {
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"comment:创建时间"`
}

type UpdatedAtField struct {
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"comment:更新时间"`
}

type DeletedAt struct {
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index;comment:删除时间"`
}

// 用户表
type User struct {
	IDField
	CreatedAtField
	UpdatedAtField
	Name     string   `json:"name" gorm:"type:VARCHAR(20);not null;unique;comment:用户名"`
	Password string   `json:"-" gorm:"type:CHAR(64);not null;comment:密码"`
	Email    string   `json:"email" gorm:"type:VARCHAR(50);not null;unique;comment:邮箱"`
	Admin    bool     `json:"admin" gorm:"not null;comment:是管理员"`
	RegionID *uint    `json:"-" gorm:"index;comment:志愿服务的区域"`
	Region   *Region  `json:"region"`
	Notices  []Notice `json:"notices" gorm:"constraint:OnDelete:SET NULL"`
	Events   []Event  `json:"events" gorm:"constraint:OnDelete:SET NULL"`
}

// 区域表
type RegionDTO struct {
	Name        string           `json:"name" gorm:"type:VARCHAR(20);not null;unique;comment:区域名"`
	Description string           `json:"description,omitempty" gorm:"type:VARCHAR(200);not null;comment:区域描述"`
	Coordinate  orb.MultiPolygon `json:"coordinate,omitempty" gorm:"type:JSON;serializer:json;not null;comment:坐标范围"`
}

type Region struct {
	IDField
	UpdatedAtField
	*RegionDTO
	Users     []User     `json:"users,omitempty" gorm:"constraint:OnDelete:SET NULL"`
	Events    []Event    `json:"events,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Histories []History  `json:"histories,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Resources []Resource `json:"resources,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Routes    []Route    `json:"routes,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

// 灾害事件
type EventDTO struct {
	Name        string     `json:"name" gorm:"type:VARCHAR(20);not null;unique;comment:事件名称"`
	RegionID    uint       `json:"-" gorm:"not null;index;comment:所在区域ID"`
	StartTime   time.Time  `json:"startTime" gorm:"not null;comment:开始时间"`
	EndTime     *time.Time `json:"endTime" gorm:"comment:结束时间"`
	Type        string     `json:"type" gorm:"type:VARCHAR(20);not null;comment:类型"`
	Severity    string     `json:"severity" gorm:"VARCHAR(20);not null;comment:严重性"`
	Coordinate  orb.Point  `json:"coordinate" gorm:"type:JSON;serializer:json;not null;comment:坐标"`
	Description string     `json:"description" gorm:"type:TEXT;not null;comment:描述"`
}

type Event struct {
	IDField
	CreatedAtField
	*EventDTO
	Region *Region `json:"region" gorm:"constraint:OnDelete:CASCADE"`
	UserID *uint   `json:"-" gorm:"index;comment:上传的用户ID"`
	User   *User   `json:"user"`
}

// 历史数据
type HistoryDTO struct {
	RegionID       uint      `json:"-" gorm:"not null;index;comment:相关的区域ID"`
	Type           string    `json:"type" gorm:"type:VARCHAR(20);not null;comment:天气类型"`
	Time           time.Time `json:"time" gorm:"not null;comment:对应的时间"`
	MaxTemperature *float64  `json:"maxTemperature" gorm:"comment:最高气温"`
	MinTemperature *float64  `json:"minTemperature" gorm:"comment:最低气温"`
	AvgTemperature *float64  `json:"avgTemperature" gorm:"comment:平均气温"`
	WindSpeed      *float64  `json:"windSpeed" gorm:"comment:风速"`
	Visibility     *float64  `json:"visibility" gorm:"comment:能见度"`
	RainFall       *float64  `json:"rainFall" gorm:"comment:降水量"`
	Severity       *float64  `json:"severity" gorm:"comment:严重性"`
	Source         string    `json:"source,omitempty" gorm:"type:VARCHAR(20);not null;comment:数据源"`
}

type History struct {
	IDField
	CreatedAtField
	*HistoryDTO
	Region *Region `json:"region,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

// 资源
type ResourceDTO struct {
	Type       string    `json:"type" gorm:"type:VARCHAR(50);not null;comment:资源类型"`
	Name       string    `json:"name" gorm:"type:VARCHAR(50);not null;comment:资源名称"`
	Quantity   uint      `json:"quantity" gorm:"not null;comment:数量"`
	RegionID   uint      `json:"regionId" gorm:"not null;index;comment:所在的区域ID"`
	Coordinate orb.Point `json:"coordinate" gorm:"type:JSON;serializer:json;not null;comment:坐标"`
	Available  bool      `json:"available" gorm:"not null;comment:是否可用"`
}

type Resource struct {
	IDField
	UpdatedAtField
	*ResourceDTO
	Region *Region `json:"region" gorm:"constraint:OnDelete:CASCADE"`
}

// 通知公告
type NoticeDTO struct {
	Title   string `json:"title" gorm:"type:VARCHAR(20);not null;comment:标题"`
	Content string `json:"content" gorm:"type:VARCHAR(200);not null;comment:内容"`
}

type Notice struct {
	IDField
	CreatedAtField
	UpdatedAtField
	*NoticeDTO
	UserID *uint `json:"-" gorm:"index;comment:编写者ID"`
	User   *User `json:"user"`
}

// 救援路线
type RouteDTO struct {
	Type        string              `json:"type" gorm:"type:VARCHAR(20);not null;comment:道路类型"`
	Name        string              `json:"name" gorm:"type:VARCHAR(20);not null;unique;comment:名称"`
	Coordinate  orb.MultiLineString `json:"coordinate" gorm:"type:JSON;serializer:json;not null;comment:坐标"`
	Description string              `json:"description" gorm:"type:VARCHAR(200);not null;comment:描述"`
	Available   bool                `json:"available" gorm:"not null;comment:是否可用"`
	Rate        *float64            `json:"rate" gorm:"comment:道路限速"`
	RegionID    uint                `json:"-" gorm:"not null;index;comment:所在的区域ID"`
}

type Route struct {
	IDField
	CreatedAtField
	*RouteDTO
	Region Region `json:"region"`
}
