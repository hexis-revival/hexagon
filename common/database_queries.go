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

func CreateBeatmapset(beatmapset *Beatmapset, state *State) error {
	result := state.Database.Create(beatmapset)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchBeatmapsetById(id int, state *State, preload ...string) (*Beatmapset, error) {
	beatmapset := &Beatmapset{}
	result := preloadQuery(state, preload).First(beatmapset, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmapset, nil
}

func RemoveBeatmapset(beatmapset *Beatmapset, state *State) error {
	result := state.Database.Delete(beatmapset)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchBeatmapsetsByCreatorId(userId int, state *State, preload ...string) ([]Beatmapset, error) {
	beatmapsets := []Beatmapset{}
	result := preloadQuery(state, preload).Find(&beatmapsets, "creator_id = ?", userId)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmapsets, nil
}

func FetchBeatmapsetCountByCreatorId(userId int, state *State) (int, error) {
	var count int64
	result := state.Database.Model(&Beatmapset{}).Where("creator_id = ?", userId).Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func FetchBeatmapsetRankedCountByCreatorId(userId int, state *State) (int, error) {
	var count int64
	result := state.Database.Model(&Beatmapset{}).Where("creator_id = ? AND status >= ?", userId, BeatmapStatusRanked).Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func FetchBeatmapsetUnrankedCountByCreatorId(userId int, state *State) (int, error) {
	var count int64
	result := state.Database.Model(&Beatmapset{}).Where("creator_id = ? AND status < ?", userId, BeatmapStatusRanked).Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func FetchBeatmapsetsByStatus(userId int, status BeatmapStatus, state *State, preload ...string) ([]Beatmapset, error) {
	beatmapsets := []Beatmapset{}
	result := preloadQuery(state, preload).Find(&beatmapsets, "creator_id = ? AND status = ?", userId, status)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmapsets, nil
}

func UpdateBeatmapset(beatmapset *Beatmapset, state *State) error {
	result := state.Database.Save(beatmapset)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateBeatmap(beatmap *Beatmap, state *State) error {
	result := state.Database.Create(beatmap)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateBeatmaps(beatmaps []Beatmap, state *State) error {
	for _, beatmap := range beatmaps {
		err := CreateBeatmap(&beatmap, state)

		if err != nil {
			return err
		}
	}
	return nil
}

func FetchBeatmapById(id int, state *State, preload ...string) (*Beatmap, error) {
	beatmap := &Beatmap{}
	result := preloadQuery(state, preload).First(beatmap, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmap, nil
}

func FetchBeatmapByChecksum(checksum string, state *State, preload ...string) (*Beatmap, error) {
	beatmap := &Beatmap{}
	result := preloadQuery(state, preload).First(beatmap, "checksum = ?", checksum)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmap, nil
}

func FetchBeatmapsBySetId(setId int, state *State, preload ...string) ([]Beatmap, error) {
	beatmaps := []Beatmap{}
	result := preloadQuery(state, preload).Find(&beatmaps, "set_id = ?", setId)

	if result.Error != nil {
		return nil, result.Error
	}

	return beatmaps, nil
}

func UpdateBeatmap(beatmap *Beatmap, state *State) error {
	result := state.Database.Save(beatmap)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RemoveBeatmap(beatmap *Beatmap, state *State) error {
	result := state.Database.Delete(beatmap)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RemoveBeatmapsBySetId(setId int, state *State) error {
	result := state.Database.Delete(&Beatmap{}, "set_id = ?", setId)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateForum(forum *Forum, state *State) error {
	result := state.Database.Create(forum)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchForumById(id int, state *State, preload ...string) (*Forum, error) {
	forum := &Forum{}
	result := preloadQuery(state, preload).First(forum, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return forum, nil
}

func FetchForumByName(name string, state *State, preload ...string) (*Forum, error) {
	forum := &Forum{}
	result := preloadQuery(state, preload).First(forum, "name = ?", name)

	if result.Error != nil {
		return nil, result.Error
	}

	return forum, nil
}

func UpdateForum(forum *Forum, state *State) error {
	result := state.Database.Save(forum)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteForum(forum *Forum, state *State) error {
	result := state.Database.Delete(forum)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateTopic(topic *ForumTopic, state *State) error {
	result := state.Database.Create(topic)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchTopicById(id int, state *State, preload ...string) (*ForumTopic, error) {
	topic := &ForumTopic{}
	result := preloadQuery(state, preload).First(topic, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return topic, nil
}

func UpdateTopic(topic *ForumTopic, state *State) error {
	result := state.Database.Save(topic)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteTopic(topic *ForumTopic, state *State) error {
	result := state.Database.Delete(topic)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreatePost(post *ForumPost, state *State) error {
	result := state.Database.Create(post)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchPostById(id int, state *State, preload ...string) (*ForumPost, error) {
	post := &ForumPost{}
	result := preloadQuery(state, preload).First(post, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

func FetchInitialPost(topicId int, state *State, preload ...string) (*ForumPost, error) {
	post := &ForumPost{}
	result := preloadQuery(state, preload).Order("id asc").First(post, "topic_id = ?", topicId)

	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

func UpdatePost(post *ForumPost, state *State) error {
	result := state.Database.Save(post)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeletePost(post *ForumPost, state *State) error {
	result := state.Database.Delete(post)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateScore(score *Score, state *State) error {
	result := state.Database.Create(score)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchScoreById(id int, state *State, preload ...string) (*Score, error) {
	score := &Score{}
	result := preloadQuery(state, preload).First(score, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return score, nil
}

func UpdateScore(score *Score, state *State) error {
	result := state.Database.Save(score)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteScore(score *Score, state *State) error {
	result := state.Database.Delete(score)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func preloadQuery(state *State, preload []string) *gorm.DB {
	result := state.Database

	for _, p := range preload {
		result = result.Preload(p)
	}

	return result
}
