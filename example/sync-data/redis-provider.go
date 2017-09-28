package main

import (
	"github.com/cocotyty/summer"
	"gopkg.in/redis.v4"
)

func init() {
	summer.Put(&RedisProvider{})
}

type RedisProvider struct {
	Client *redis.Client
}

func (provider *RedisProvider) Init() {
	provider.Client = redis.NewClient(&redis.Options{
		Addr: Conf.RedisAddr,
	})
	err := provider.Client.Ping().Err()
	if err != nil {
		panic(err)
	}
}

func (provider *RedisProvider) Provide() (client *redis.Client) {
	return provider.Client
}
