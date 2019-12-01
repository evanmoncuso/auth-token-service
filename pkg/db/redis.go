package datastore

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var Pool *redis.Pool

func InitializeRedis(url string) error {
	Pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
		},
	}

	fmt.Println("Successful Redis connection")

	return nil
}

// SetInvalidToken adds a token to the pool if it doesn't exist
func SetInvalidToken(t string, expireTime int64) error {
	conn := Pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", t))
	if err != nil {
		return err
	} else if exists == 1 {
		log.Printf("Token %s has already been invalidated\n", t)

		// I'm not sure the end user needs to know it's already been invalidated
		return nil
	}

	conn.Send("MULTI")
	conn.Send("SET", t, expireTime)
	conn.Send("EXPIREAT", t, expireTime)
	_, err = conn.Do("EXEC")

	return err
}

// CheckForInvalidToken checks the store to see if this token is invalid
//  it returns an error if it is invalid, returns nil if okay
func CheckForInvalidToken(t string, u string) error {
	conn := Pool.Get()
	defer conn.Close()

	if exists, err := redis.Int(conn.Do("EXISTS", t)); err != nil {
		return err
	} else if exists == 1 {
		return errors.New("Token present in invalidation store")
	}

	return nil
}
