package db

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gorp/gorp"
	_redis "github.com/go-redis/redis/v7"
)

type DB struct {
	*sql.DB
}

var db *gorp.DbMap

func Init() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	var err error
	db, err = ConnectDB(dbInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB(dataSourceName string) (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	// dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests
	return dbmap, nil
}

// GetDB returns the DB
func GetDB() *gorp.DbMap {
	return db
}

// RedisClient ...
var RedisClient *_redis.Client

// InitRedis
func InitRedis(selectDB ...int) {
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	RedisClient = _redis.NewClient(&_redis.Options{
		Addr:               redisHost,
		Password:           redisPassword,
		DB:                 selectDB[0],
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
}

// GetRedis ...
func GetRedis() *_redis.Client {
	return RedisClient
}
