package repository

import (
	"context"
	"fmt"
	"log"
	"memorizor/services/account/util"
	"time"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

type sTokenRepositoryRedis struct {
	rdb *redis.Client
	ctx context.Context
}

func NewSTokenRepository(rdb *redis.Client) ITokenRepository {
	return &sTokenRepositoryRedis{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *sTokenRepositoryRedis) SetRefreshToken(userID, tokenID uuid.UUID, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID.String(), tokenID.String())
	if err := r.rdb.Set(r.ctx, key, 0, expiresIn).Err(); err != nil {
		log.Println(err.Error())
		return util.NewInternal("Could not set the refresh token")
	}
	return nil
}
func (r *sTokenRepositoryRedis) DeleteRefreshToken(userID, previousTokenID uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", userID.String(), previousTokenID.String())
	delResult := r.rdb.Del(r.ctx, key)
	if err := delResult.Err(); err != nil {
		log.Println(err.Error())
		return util.NewInternal("Could not delete the previous refresh token")
	}
	if delResult.Val() < 1 {
		log.Println("Previous refresh token has expired")
		return util.NewAuthorization("Invalid refresh token: token has expired")
	}
	return nil
}

func (r *sTokenRepositoryRedis) DeleteUserRefreshTokens(userID uuid.UUID) error {
	pattern := fmt.Sprintf("%s*", userID.String())
	iter := r.rdb.Scan(r.ctx, 0, pattern, 5).Iterator()
	failedCnt := 0
	for iter.Next(r.ctx) {
		if err := r.rdb.Del(r.ctx, iter.Val()).Err(); err != nil {
			log.Println("Failed to delete a refresh token")
			failedCnt++
		}
	}
	if failedCnt > 0 {
		return util.NewInternal("Failed to delete a refresh token")
	}
	if iter.Err() != nil {
		return util.NewInternal("Failed to find the next user:refresh token pair")
	}
	return nil
}
