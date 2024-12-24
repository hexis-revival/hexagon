package common

import (
	"strings"

	"github.com/redis/go-redis/v9"
)

func UpdateRankingsEntry(stats *Stats, country string, state *State) error {
	country = strings.ToLower(country)
	errors := NewErrorCollection()

	entries := map[string]float64{
		"rankings:rscore": float64(stats.RankedScore),
		"rankings:tscore": float64(stats.TotalScore),
		"rankings:clears": float64(stats.Clears()),
	}

	for key, score := range entries {
		result := state.Redis.ZAdd(
			*state.RedisContext, key,
			redis.Z{
				Score:  score,
				Member: stats.UserId,
			},
		)
		errors.Add(result.Err())

		result = state.Redis.ZAdd(
			*state.RedisContext, key+":"+country,
			redis.Z{
				Score:  score,
				Member: stats.UserId,
			},
		)
		errors.Add(result.Err())
	}

	return errors.Next()
}
