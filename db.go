package main

import (
	"log"
)

// save key value to Redis
func SaveToDB(key string, value []byte) error {
	err := DB.Set(key, value, 0).Err()
	if err != nil {
		log.Println("Error occurred with save to Redis:", err)
		return err
	}
	return nil
}

// load data from Redis by filtered key
func LoadFromDB(filter string) ([]string, error) {
	keys, err := DB.Keys(filter).Result()
	if err != nil {
		log.Println("Error occurred with loading data from Redis:", err)
		return nil, err
	}

	values := make([]string, len(keys))

	for i, key := range keys {
		value, err := DB.Get(key).Result()
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
func DeleteFromDB(key string) error {
	err := DB.Del(key).Err()
	if err != nil {
		log.Printf("Error occurred with deleting value by key %s: %s", key, err)
		return err
	}
	return nil
}
