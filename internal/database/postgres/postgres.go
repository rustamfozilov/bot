package postgres

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/ssharifzoda/bot/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func NewPostgresGorm() (*gorm.DB, error) {
	log := logging.GetLogger()

	Host := viper.GetString("db.host")
	Port := viper.GetUint16("db.port")
	Username := viper.GetString("db.username")
	Password := os.Getenv("DB_PASSWORD")
	DBName := viper.GetString("db.dbname")
	connString := fmt.Sprintf("host=%s user=%s password=%v dbname=%s port=%d sslmode=disable TimeZone=Asia/Dushanbe",
		Host, Username, Password, DBName, Port)
	conn, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Infof("%s GetPostgresConnection -> Open error: ", err.Error())
		return nil, err
	}
	log.Info("Postgres Connection success: ", Host)
	return conn, nil
}
