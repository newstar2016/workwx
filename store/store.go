package store

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenCache interface {
	GetSuiteTicket(suiteId string) (string, error)
	SetSuiteTicket(suiteId string, ticket string) error
	GetSuiteAccessToken(suiteId string) (string, error)
	SetSuiteAccessToken(suiteId string, accessToken string, expiresIn int) error
	GetSuiteCorpAccessToken(suiteId string, corpId string) (string, error)
	SetSuiteCorpAccessToken(suiteId string, corpId string, accessToken string, expiresIn int) error
}

type RedisTokenStore struct {
	prefix string
	redis  *redis.Client
}

func NewRedisTokenStore(prefix string, opt *redis.Options) TokenCache {
	return &RedisTokenStore{
		prefix: prefix,
		redis:  redis.NewClient(opt),
	}
}

const cacheKeySuiteTicket = "%s:%s:ticket"
const cacheKeySuiteAccessToken = "%s:%s:access_token"
const cacheKeyCorpAccessToken = "%s:%s:auth:access_token:%s"

func (s *RedisTokenStore) SetSuiteTicket(suiteId string, ticket string) error {
	r := s.redis.Set(s.redis.Context(), fmt.Sprintf(cacheKeySuiteTicket, s.prefix, suiteId), ticket, time.Duration(0))
	return r.Err()
}

func (s *RedisTokenStore) GetSuiteTicket(suiteId string) (string, error) {
	r := s.redis.Get(s.redis.Context(), fmt.Sprintf(cacheKeySuiteTicket, s.prefix, suiteId))
	if r.Err() != nil {
		return "", r.Err()
	}
	return r.Val(), nil
}

func (s *RedisTokenStore) SetSuiteAccessToken(suiteId string, accessToken string, expiresIn int) error {
	r := s.redis.Set(s.redis.Context(), fmt.Sprintf(cacheKeySuiteAccessToken, s.prefix, suiteId), accessToken, time.Second*time.Duration(expiresIn))
	return r.Err()
}

func (s *RedisTokenStore) GetSuiteAccessToken(suiteId string) (string, error) {
	r := s.redis.Get(s.redis.Context(), fmt.Sprintf(cacheKeySuiteAccessToken, s.prefix, suiteId))
	err := r.Err()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", nil
	}

	return r.Val(), nil
}

func (s *RedisTokenStore) SetSuiteCorpAccessToken(suiteId string, corpId string, accessToken string, expiresIn int) error {
	r := s.redis.Set(s.redis.Context(), fmt.Sprintf(cacheKeyCorpAccessToken, s.prefix, suiteId, corpId), accessToken, time.Second*time.Duration(expiresIn))
	return r.Err()
}

func (s *RedisTokenStore) GetSuiteCorpAccessToken(suiteId string, corpId string) (string, error) {
	r := s.redis.Get(s.redis.Context(), fmt.Sprintf(cacheKeyCorpAccessToken, s.prefix, suiteId, corpId))
	err := r.Err()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", nil
	}
	return r.Val(), nil
}
