package main

import (
	"gopkg.in/redis.v4"
	"database/sql"
	"time"
	"log"
	"github.com/cocotyty/summer"
)

const key = "/key"

func init() {
	summer.Put(&SyncDataWorker{})
}

type SyncDataWorker struct {
	RedisProvider    *RedisProvider `sm:"*"`
	DatabaseProvider *DatabaseProvider `sm:"*"`
	redisClient      *redis.Client
	db               *sql.DB
}

func (worker *SyncDataWorker) Run() {
	for {
		if result := worker.redisClient.Get(key); result.Err() == nil {
			worker.db.Exec("update `test_table` set `text` = ? where `key` = ? ", result.String(), key)
		} else {
			log.Println(result.Err())
		}
		time.Sleep(time.Minute)
	}
}
func (worker *SyncDataWorker) Ready() {
	worker.redisClient = worker.RedisProvider.Provide()
	worker.db = worker.DatabaseProvider.Provide()
	go worker.Run()
}
