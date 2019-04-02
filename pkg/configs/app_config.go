package configs

import (
	"fmt"
	"github.com/buicongtan1997/manabie/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type appConfig struct {
	Env string
	mysqlConfig
}

type mysqlConfig struct {
	Host             string
	Port             string
	Database         string
	UserName         string
	Password         string
	ConnectionString string
}

func init(){
	// initialize logger
	if err := logger.Init(-1, "2006-01-02T15:04:05Z07:00"); err != nil {
		log.Fatal("failed to initialize logger: %v", err)
	}
	cfg.loadConfig()
}



var cfg appConfig

func GetAppConfig() appConfig {
	return cfg
}

func (rcv appConfig) loadConfig() {
	err := godotenv.Load("./resources/.env")
	if err != nil {
		logger.Log.Fatal("Error loading .env file")
	}

	cfg.Env = os.Getenv("MANABIE_ENV")
	if "" == cfg.Env {
		cfg.Env = "local"
	}

	err = godotenv.Overload("./resources/.env." + cfg.Env)
	if err != nil {
		logger.Log.Fatal("Error loading ./resources/.env." + cfg.Env  + " file")
	}
	cfg.New(
		os.Getenv("MYSQL_DB_HOST"),
		os.Getenv("MYSQL_DB_PORT"),
		os.Getenv("MYSQL_DB_DATABASE"),
		os.Getenv("MYSQL_DB_USERNAME"),
		os.Getenv("MYSQL_DB_PASSWORD"),
	)
}

func (rcv *mysqlConfig) New(host, port, database, userName, password string) {
	rcv.Host = host
	rcv.Port = port
	rcv.Database = database
	rcv.UserName = userName
	rcv.Password = password
	rcv.ConnectionString = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4,utf8",
		userName,
		password,
		host,
		port,
		database,
	)
}
