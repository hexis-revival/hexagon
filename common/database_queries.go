package common

import (
	"strings"
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

func FetchUserRelationships(userId int, status RelationshipStatus, state *State) ([]*Relationship, error) {
	relationships := []*Relationship{}
	result := state.Database.Where("user_id = ? AND status = ?", userId, status).Find(&relationships)

	if result.Error != nil {
		return nil, result.Error
	}

	return relationships, nil
}
