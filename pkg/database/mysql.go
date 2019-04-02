package database

import (
	"github.com/buicongtan1997/manabie/pkg/configs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var db *gorm.DB

func init() {
	db = connect(configs.GetAppConfig().Database, configs.GetAppConfig().ConnectionString)

	//migrate table
	migrate()
	if configs.GetAppConfig().Env != "test" {
		initData()
	}
}

func GetMainDatabase() *gorm.DB {
	return db
}

func connect(databaseName, connectionString string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Fail! ", err)
	}

	db.DB().SetConnMaxLifetime(30 * time.Minute)
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(100)
	if gin.Mode() == gin.DebugMode {
		db = db.Debug()
	}
	gorm.DefaultCallback.Create().Remove("gorm:force_reload_after_create")
	log.Println("Connected to database", databaseName)
	return db
}
