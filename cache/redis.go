package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCfg struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
	LockPrefix string `yaml:"lock_prefix"`
	LockTTLSec int    `yaml:"lock_ttl_sec"`
}

func NewRedis(cfg *RedisCfg) (*RedisClient, error) {
	opts := &redis.Options{}
	opts.Network = "tcp"
	opts.Addr = cfg.Host + ":" + cfg.Port
	opts.Password = cfg.Password
	opts.DB = cfg.DB
	client := redis.NewClient(opts)
	locker := NewLocker(client, cfg.LockPrefix, time.Duration(cfg.LockTTLSec)*time.Second)
	return &RedisClient{
		Client: client,
		Locker: locker,
	}, nil
}

type RedisClient struct {
	*redis.Client
	Locker Locker
}
