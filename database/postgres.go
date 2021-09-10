package database

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	database *gorm.DB
	once     sync.Once
)

func GetTransaction() *gorm.DB {
	once.Do(func() {
		database = connectToDb()
	})
	return database.Begin()
}

func GetTransactionWithContext(ctx context.Context) *gorm.DB {
	tx := GetTransaction()
	return tx.WithContext(ctx)
}

func connectToDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.pass"),
		viper.GetString("db.name"),
		viper.GetString("db.port"))
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error when connecting to db: %s", err)
	}
	return d
}
