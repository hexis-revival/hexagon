package hnet

type PlayerCollection struct {
	idMap   map[uint32]*Player
	nameMap map[string]*Player
}

func (collection *PlayerCollection) Add(player *Player) {
	collection.idMap[player.Presence.Id] = player
	collection.nameMap[player.Presence.Name] = player
}

func (collection *PlayerCollection) Remove(player *Player) {
	delete(collection.idMap, player.Presence.Id)
	delete(collection.nameMap, player.Presence.Name)
}

func (collection *PlayerCollection) Count() int {
	return len(collection.idMap)
}

func (collection *PlayerCollection) ByID(id uint32) *Player {
	if val, ok := collection.idMap[id]; ok {
		return val
	}

	return nil
}

func (collection *PlayerCollection) ByName(name string) *Player {
	if val, ok := collection.nameMap[name]; ok {
		return val
	}

	return nil
}

func (collection *PlayerCollection) All() []Player {
	players := make([]Player, 0, len(collection.idMap))

	for _, player := range collection.idMap {
		players = append(players, *player)
	}

	return players
}

func NewPlayerCollection() PlayerCollection {
	return PlayerCollection{
		idMap:   make(map[uint32]*Player),
		nameMap: make(map[string]*Player),
	}
}
