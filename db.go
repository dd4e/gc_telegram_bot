package main

import (
	"github.com/go-redis/redis"
	"log"
)

type redisDB struct {
	client   *redis.Client
	Addr     string
	Password string
	numDB    int
}

// initial connection to Redis
func (db *redisDB) Connect() error {
	db.client = redis.NewClient(
		&redis.Options{
			DB:       db.numDB,
			Addr:     db.Addr,
			Password: db.Password,
		},
	)

	// try ping redis
	if err := db.client.Ping().Err(); err != nil {
		return err
	}
	return nil
}

// save key value to Redis
func (db redisDB) SaveToDB(key string, value []byte) error {
	err := db.client.Set(key, value, 0).Err()
	if err != nil {
		log.Println("Error occurred with save to Redis:", err)
		return err
	}
	return nil
}

// load data from Redis by filtered key
func (db redisDB) LoadFromDB(filter string) ([]string, error) {
	keys, err := db.client.Keys(filter).Result()
	if err != nil {
		log.Println("Error occurred with loading data from Redis:", err)
		return nil, err
	}

	values := make([]string, len(keys))

	for i, key := range keys {
		value, err := db.client.Get(key).Result()
		if err != nil {
			log.Printf(
				"Error occurred with getting value by key %s: %s. Skip...", key, err)
			continue
		}
		values[i] = value
	}
	return values, nil
}

// delete value by key from Redis
func (db redisDB) DeleteFromDB(key string) error {
	err := db.client.Del(key).Err()
	if err != nil {
		log.Printf("Error occurred with deleting value by key %s: %s", key, err)
		return err
	}
	return nil
}
