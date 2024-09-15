package main

import (
	"Go_Learn/conf"
	"Go_Learn/pkg/route"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Writer để log ra console
		logger.Config{
			SlowThreshold: time.Second, // Thời gian để log các truy vấn chậm hơn 1 giây
			LogLevel:      logger.Info, // Log toàn bộ query (Info, Warn, Error, Silent)
			Colorful:      true,        // Màu sắc của log (tuỳ vào môi trường console)
		},
	)
	config, err := conf.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config")
	}

	db, err := gorm.Open(postgres.Open(config.DBSource), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	runGinServer(config, db)
}
func runGinServer(config conf.Config, pg *gorm.DB) {
	server, err := route.NewServer(config, pg)
	if err != nil {
		log.Fatal("can't create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can't start server")
	}
}
