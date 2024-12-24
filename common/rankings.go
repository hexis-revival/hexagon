package common

import (
	"strconv"
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

func RemoveRankingsEntry(stats *Stats, country string, state *State) error {
	country = strings.ToLower(country)
	errors := NewErrorCollection()

	entries := []string{
		"rankings:rscore",
		"rankings:tscore",
		"rankings:clears",
	}

	for _, key := range entries {
		result := state.Redis.ZRem(
			*state.RedisContext, key, stats.UserId,
		)
		errors.Add(result.Err())

		result = state.Redis.ZRem(
			*state.RedisContext, key+":"+country, stats.UserId,
		)
		errors.Add(result.Err())
	}

	return errors.Next()
}

func GetScoreRank(userId int, state *State) (int, error) {
	result := state.Redis.ZRevRank(
		*state.RedisContext,
		"rankings:rscore", strconv.Itoa(userId),
	)
	return int(result.Val()) + 1, result.Err()
}

func GetCountryScoreRank(userId int, country string, state *State) (int, error) {
	result := state.Redis.ZRevRank(
		*state.RedisContext,
		"rankings:rscore:"+strings.ToLower(country), strconv.Itoa(userId),
	)
	return int(result.Val()) + 1, result.Err()
}

func GetTotalScoreRank(userId int, state *State) (int, error) {
	result := state.Redis.ZRevRank(
		*state.RedisContext,
		"rankings:tscore", strconv.Itoa(userId),
	)
	return int(result.Val()) + 1, result.Err()
}

func GetCountryTotalScoreRank(userId int, country string, state *State) (int, error) {
	result := state.Redis.ZRevRank(
		*state.RedisContext,
		"rankings:tscore:"+strings.ToLower(country), strconv.Itoa(userId),
	)
	return int(result.Val()) + 1, result.Err()
}

func GetClearsRank(userId int, state *State) (int, error) {
	result := state.Redis.ZRevRank(
		*state.RedisContext,
		"rankings:clears", strconv.Itoa(userId),
	)
	return int(result.Val()) + 1, result.Err()
}

func GetCountryClearsRank(userId int, countryCode string, state *State) (int, error) {
	result := state.Redis.ZRevRank(
		*state.RedisContext,
		"rankings:clears:"+strings.ToLower(countryCode), strconv.Itoa(userId),
	)
	return int(result.Val()) + 1, result.Err()
}
