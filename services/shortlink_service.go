package services

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jiaqi-yin/go-url-shortener/utils"
	"github.com/pilu/go-base62"
)

const (
	URLIDKEY     = "next.url.id"
	ShortlinkKey = "shortlink:%s"
	URLHashKey   = "urlhash:%s"
)

type ShortlinkServiceInterface interface {
	Shorten(string) (string, utils.RestErr)
	Unshorten(string) (string, utils.RestErr)
}

type RedisCli struct {
	Cli *redis.Client
}

func NewRedisCli(addr string, password string, db int) *RedisCli {
	ctx := context.Background()
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err := c.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: c}
}

func toSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", string(h.Sum(nil)))
}

func (rc *RedisCli) Shorten(url string) (string, utils.RestErr) {
	h := toSha1(url)

	ctx := context.Background()
	d, err := rc.Cli.Get(ctx, fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		// Do nothing
	} else if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	} else {
		if d == "" {
			// Do nothing
		} else {
			return d, nil
		}
	}

	err = rc.Cli.Incr(ctx, URLIDKEY).Err()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}

	id, err := rc.Cli.Get(ctx, URLIDKEY).Int64()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}
	eid := base62.Encode(int(id))

	err = rc.Cli.Set(ctx, fmt.Sprintf(ShortlinkKey, eid), url, 0).Err()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}

	err = rc.Cli.Set(ctx, fmt.Sprintf(URLHashKey, h), eid, 0).Err()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}

	return eid, nil
}

func (rc *RedisCli) Unshorten(eid string) (string, utils.RestErr) {
	log.Println("eid", eid)
	ctx := context.Background()
	url, err := rc.Cli.Get(ctx, fmt.Sprintf(ShortlinkKey, eid)).Result()
	log.Println("url", url)
	if err == redis.Nil {
		return "", utils.NewNotFoundError("shortlink key not exist")
	} else if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	} else {
		return url, nil
	}
}
