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

func (r *sTokenRepository) SetRefreshToken(userID, tokenID uuid.UUID, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID.String(), tokenID.String())
	if err := r.rdb.Set(r.ctx, key, 0, expiresIn).Err(); err != nil {
		log.Println(err.Error())
		return util.NewInternal("Could not set the refresh token")
	}
	return nil
}
func (r *sTokenRepository) DeleteRefreshToken(userID, previousTokenID uuid.UUID) error {
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
