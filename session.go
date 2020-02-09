package telegrabotlib

import (
    "github.com/go-redis/redis/v7"
    "github.com/pkg/errors"
    "time"
)

type Session interface {
    Set(user, key string, value string) error
    SetForever(user, key string, value string) error
    Get(user,key string) (string, error)
    GetAllLike(key string) ([]string, error)
    Delete(user string, key string) error
    DeleteAll(user string) error
}

var (
    ErrKeyNotFound = errors.New("Not found key")
)

type defaultRedisSession struct {
    client *redis.Client
}

func NewRedisSession(redisURL string) *defaultRedisSession {
    options, _ := redis.ParseURL(redisURL)
    return &defaultRedisSession{
        client: redis.NewClient(options),
    }
}

func (d *defaultRedisSession) SetForever(user, key string, value string) error {
    d.client.Set(user+key, value, 0)
    return nil
}

func (d *defaultRedisSession) Set(user, key string, value string) error {
    d.client.Set(user+key, value, time.Hour * 100)
    return nil
}

func (d *defaultRedisSession) Get(user, key string) (string, error) {
    r := d.client.Get(user+key)
    return r.Result()
}

func (d *defaultRedisSession) GetAllLike(key string) ([]string, error) {
    iter := d.client.Scan(0, key, 0).Iterator()

    keys := make([]string, 0)
    for iter.Next() {
        keys = append(keys, iter.Val())
    }

    if len(keys) == 0 {
        return nil, ErrKeyNotFound
    }

    r := d.client.MGet(keys...)

    vals, err := r.Result()
    if err != nil {
        return nil, err
    }

    res := make([]string, len(vals))
    for _, v := range vals {
        res = append(res, v.(string))
    }

    return res, nil
}

func (d *defaultRedisSession) DeleteAll(user string) error {
    iter := d.client.Scan(0, user+"*", 0).Iterator()
    for iter.Next() {
        err := d.client.Del(iter.Val()).Err()
        if err != nil {
            return err
        }
    }
    if err := iter.Err(); err != nil {
        return err
    }

    return nil
}

func (d *defaultRedisSession) Delete(user, key string) error {
    err := d.client.Del(user+key).Err()
    if err != nil {
        return err
    }

    return nil
}






