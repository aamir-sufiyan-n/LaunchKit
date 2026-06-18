package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/utils"
	"github.com/redis/go-redis/v9"
)

// the store
type Store interface {
	StoreRefreshToken(ctx context.Context, tokenID, userID string, expiry time.Duration) error
	BlackListRefreshtoken(ctx context.Context, tokenID string, ttl time.Time) error
	UpgradeTokenVersion(ctx context.Context, tokenID string) error
	DeleteRefreshToken(ctx context.Context, tokenID string) error
	GetTokenVersion(ctx context.Context, userID string) (int64, error)
	IsRefreshBlacklisted(ctx context.Context, tokenID string) (time.Time, error)
}

//redis store
type redisStore struct {
	client *redis.Client
}

//returns a new store
func NewStore(c *redis.Client) Store {
	return &redisStore{
		client:c,
	}
}

//store refresh tokens
func (r *redisStore)StoreRefreshToken(ctx context.Context,tokenID,userID string, expiry time.Duration)error{
	key:=fmt.Sprintf("refreshToken:%s",tokenID)
	err:=r.client.Set(ctx,key,userID,expiry).Err()
	if err!=nil{
		return fmt.Errorf("error storing refresh token:%w",err)
	}
	return nil
}

//blacklist a refresh token
func (r *redisStore)BlackListRefreshtoken(ctx context.Context,tokenID string,ttl time.Time)error{
	key:=fmt.Sprintf("blacklist:%s",tokenID)
	expiresAt:=time.Until(ttl)
	timestamp:=utils.ToUnixTimestamp(time.Now())
	return r.client.Set(ctx,key,timestamp,expiresAt).Err()
}

//upgrade token version 
func( r *redisStore)UpgradeTokenVersion(ctx context.Context,userID string)error{
	key:=fmt.Sprintf("auth:user:%s:version",userID)
	err:=r.client.Incr(ctx,key).Err()
	if err!=nil{
		return fmt.Errorf("redis increment version error: %w",err)
	}
	return nil
}

//delete a refresh token
func(r *redisStore)DeleteRefreshToken(ctx context.Context, tokenID string)error{
	key:=fmt.Sprintf("refreshToken:%s",tokenID)
	if err:= r.client.Del(ctx,key).Err(); err !=nil{
		return fmt.Errorf("deleting refresh token:%w",err)
	}
	return  nil
}

//get token version
func (r *redisStore) GetTokenVersion(ctx context.Context, userID string) (int64, error) {
	key := fmt.Sprintf("auth:user:%s:version", userID)
	version, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("getting token version: %w", err)
	}
	return version, nil
}

//check if token is blacklisted
func (r *redisStore) IsRefreshBlacklisted(ctx context.Context, tokenID string) (time.Time, error) {
	key := fmt.Sprintf("blacklist:%s", tokenID)
	issuedTime, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("redis check failed: %w", err)
	}
	return utils.FromUnixTimestamp(issuedTime)
}




