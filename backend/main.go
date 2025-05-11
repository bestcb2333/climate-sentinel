package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

func main() {

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("配置读取失败：%w\n", err)
	}

	db, err := InitDB(&config)
	if err != nil {
		log.Fatalf("项目初始化失败: %v\n", err.Error())
	}

	r := GetRouter(db, &config)

	if err := r.Run(":" + config.Port); err != nil {
		log.Fatalf("Failed to run app: %w\n", err)
	}
}
