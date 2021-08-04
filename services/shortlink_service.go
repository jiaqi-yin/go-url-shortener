package services

import (
	"context"
	"fmt"

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

type ShortlinkService struct {
	Cli *redis.Client
}

func NewShortlinkService(addr string, password string, db int) *ShortlinkService {
	ctx := context.Background()
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err := c.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return &ShortlinkService{Cli: c}
}

func (s *ShortlinkService) Shorten(url string) (string, utils.RestErr) {
	h := utils.ToSha1(url)

	ctx := context.Background()
	d, err := s.Cli.Get(ctx, fmt.Sprintf(URLHashKey, h)).Result()
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

	err = s.Cli.Incr(ctx, URLIDKEY).Err()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}

	id, err := s.Cli.Get(ctx, URLIDKEY).Int64()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}
	eid := base62.Encode(int(id))

	err = s.Cli.Set(ctx, fmt.Sprintf(ShortlinkKey, eid), url, 0).Err()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}

	err = s.Cli.Set(ctx, fmt.Sprintf(URLHashKey, h), eid, 0).Err()
	if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	}

	return eid, nil
}

func (s *ShortlinkService) Unshorten(eid string) (string, utils.RestErr) {
	ctx := context.Background()
	url, err := s.Cli.Get(ctx, fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", utils.NewNotFoundError("shortlink key not exist")
	} else if err != nil {
		return "", utils.NewInternalServerError(err.Error())
	} else {
		return url, nil
	}
}
