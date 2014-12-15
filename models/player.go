// player
package models

type Player struct {
	ID       int64
	PlayerID int64
	Name     string
	Corp     *Corporation
	AccessMask
}

func NewPlayer(id int64, playerID int64, name string, corp *Corporation, access AccessMask) *Player {
	player := &Player{
		ID:         id,
		PlayerID:   playerID,
		Name:       name,
		Corp:       corp,
		AccessMask: access,
	}

	return player
}
