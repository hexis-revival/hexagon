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

func FetchUserById(id int, state *State, preload ...string) (*User, error) {
	user := &User{}
	result := preloadQuery(state, preload).First(user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByName(name string, state *State, preload ...string) (*User, error) {
	user := &User{}
	query := preloadQuery(state, preload).Where("name = ?", name)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByNameCaseInsensitive(name string, state *State, preload ...string) (*User, error) {
	user := &User{}
	query := preloadQuery(state, preload).Where("lower(name) = ?", strings.ToLower(name))
	result := query.First(user)

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

func FetchUserRelationships(userId int, status RelationshipStatus, state *State, preload ...string) ([]*Relationship, error) {
	relationships := []*Relationship{}
	query := preloadQuery(state, preload).Where("user_id = ? AND status = ?", userId, status)
	result := query.Find(&relationships)

	if result.Error != nil {
		return nil, result.Error
	}

	return relationships, nil
}

func FetchUserRelationship(userId int, targetId int, state *State, preload ...string) (*Relationship, error) {
	relationship := &Relationship{}
	result := preloadQuery(state, preload).First(relationship, "user_id = ? AND target_id = ?", userId, targetId)

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
