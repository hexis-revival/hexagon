package common

import (
	"strings"

	"gorm.io/gorm"
)

func CreateUser(user *User, state *State) error {
	result := state.Database.Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchUserById(id int, state *State) (*User, error) {
	user := &User{}
	result := state.Database.Preload("Stats").First(user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByName(name string, state *State) (*User, error) {
	user := &User{}
	query := state.Database.Where("name = ?", name)
	result := query.Preload("Stats").First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByNameCaseInsensitive(name string, state *State) (*User, error) {
	user := &User{}
	result := state.Database.Where("lower(name) = ?", strings.ToLower(name))
	result = result.Preload("Stats").First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func CreateStats(stats *Stats, state *State) error {
	result := state.Database.Create(stats)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchStatsByUserId(userId int, state *State) (*Stats, error) {
	stats := &Stats{}
	result := state.Database.First(stats, "user_id = ?", userId)

	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func CreateUserRelationship(relationship *Relationship, state *State) error {
	result := state.Database.Create(relationship)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RemoveUserRelationship(relationship *Relationship, state *State) error {
	result := state.Database.Delete(relationship)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchUserRelationships(userId int, status RelationshipStatus, state *State) ([]*Relationship, error) {
	relationships := []*Relationship{}
	result := state.Database.Where("user_id = ? AND status = ?", userId, status).Find(&relationships)

	if result.Error != nil {
		return nil, result.Error
	}

	return relationships, nil
}

func FetchUserRelationship(userId int, targetId int, state *State) (*Relationship, error) {
	relationship := &Relationship{}
	result := state.Database.First(relationship, "user_id = ? AND target_id = ?", userId, targetId)

	if result.Error != nil {
		return nil, result.Error
	}

	return relationship, nil
}

func FetchBeatmapsetById(id int, state *State, preload ...string) (*Beatmapset, error) {
	beatmapset := &Beatmapset{}
	result := preloadQuery(state, preload).First(beatmapset, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmapset, nil
}

func FetchBeatmapById(id int, state *State, preload ...string) (*Beatmap, error) {
	beatmap := &Beatmap{}
	result := preloadQuery(state, preload).First(beatmap, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmap, nil
}

func preloadQuery(state *State, preload []string) *gorm.DB {
	result := state.Database

	for _, p := range preload {
		result = result.Preload(p)
	}

	return result
}
