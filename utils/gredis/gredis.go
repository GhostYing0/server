package gredis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

// 定义redis链接池
var RedisClient *redis.Client

//var RedisClient *redis.ClusterClient

func Setup() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:6379",
		DB:          0,
		IdleTimeout: 30 * time.Minute, // 空闲链接超时时间
	})

	//RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs:    []string{setting.RedisSetting.Host},
	//	Password: setting.RedisSetting.Password,
	//	//DB:          0,
	//	IdleTimeout: setting.RedisSetting.IdleTimeout, // 空闲链接超时时间
	//})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
		panic(err)
	}
}

func HashSet(key string, field string, value interface{}) error {
	_, err := RedisClient.HSet(key, field, value).Result()
	if err != nil {
		return err
	}
	return nil
}

func HashMSet(key string, fields map[string]interface{}) error {
	_, err := RedisClient.HMSet(key, fields).Result()
	if err != nil {
		return err
	}
	return nil
}

func BatchHashSet(key string, fields map[string]interface{}) error {
	_, err := RedisClient.HMSet(key, fields).Result()
	if err != nil {
		return err
	}
	return nil
}

func HashIncrBy(key string, field string, incr int64) error {
	_, err := RedisClient.HIncrBy(key, field, incr).Result()
	if err == redis.Nil {
		return errors.New("Key Doesn't Exists: " + field)
	} else if err != nil {
		return err
	}
	return nil
}

func HashSetnx(key string, field string, value interface{}) error {
	_, err := RedisClient.HSetNX(key, field, value).Result()
	if err != nil {
		return err
	}
	return nil
}

func HExists(key, field string) (bool, error) {
	val, err := RedisClient.HExists(key, field).Result()
	if err != nil {
		return val, err
	}
	return val, nil
}

func HGetAll(key string) (interface{}, error) {
	arr, err := RedisClient.Do("hgetall", key).Result()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func HGet(key, field string) (string, error) {
	str, err := RedisClient.HGet(key, field).Result()
	if err == redis.Nil {
		return "", errors.New("Key Doesn't Exists: " + field)
	} else if err != nil {
		return "", err
	}
	return str, nil
}

func HashDel(key string, fields []string) error {

	_, err := RedisClient.HDel(key, fields...).Result()
	if err != nil {
		return err
	}
	return nil
}

func Del(key string) error {
	val, err := RedisClient.Del(key).Result()
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	val, err := RedisClient.SetNX(key, value, expiration).Result()
	if err != nil {
		return val, err
	}
	return val, err
}

func Set(key string, value interface{}, expiration time.Duration) error {
	err := RedisClient.Set(key, value, expiration).Err()
	return err
}

func Get(key string) (str string, err error) {
	str, err = RedisClient.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("Key Doesn't Exists Redis")
		err = nil
	}
	return
}

func HExpire(key string, expiration time.Duration) (bool, error) {
	val, err := RedisClient.Expire(key, expiration).Result()
	if err != nil {
		return val, err
	}
	return val, err
}

func HyperloglogAdd(key, val string) error {
	_, err := RedisClient.PFAdd(key, val).Result()
	RedisClient.PFAdd(key, val)
	return err
}

// HyperloglogCount flag标志位，true需要回库，false不需要
func HyperloglogCount(key string) int64 {
	count, err := RedisClient.PFCount(key).Result()
	if err != nil {
		// 获取缓存失败，需要查询数据库
		return 0
	}
	return count
}

func HashSetTimeout(key string, field string, value interface{}, expire time.Duration) error {
	_, err := RedisClient.HSet(key, field, value).Result()
	if err != nil {
		return err
	}

	err = RedisClient.Expire(key, expire).Err()
	return err
}
