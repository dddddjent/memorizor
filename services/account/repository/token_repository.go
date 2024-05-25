package repository

import (
	"context"
	"fmt"
	"log"
	"memorizor/services/account/util"
	"time"

	"github.com/redis/go-redis/v9"
)

type sTokenRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewSTokenRepository(rdb *redis.Client) ITokenRepository {
	return &sTokenRepository{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *sTokenRepository) SetRefreshToken(userID, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.rdb.Set(r.ctx, key, 0, expiresIn).Err(); err != nil {
		log.Println(err.Error())
		return &util.Error{
			Type:    util.Internal,
			Message: "Could not set the refresh token",
		}
	}
	return nil
}
func (r *sTokenRepository) DeleteRefreshToken(userID, previousTokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, previousTokenID)
	if err := r.rdb.Del(r.ctx, key).Err(); err != nil {
		log.Println(err.Error())
		return &util.Error{
			Type:    util.Internal,
			Message: "Could not delete the previous refresh token",
		}
	}
	return nil
}
