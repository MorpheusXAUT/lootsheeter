// player
package models

type Player struct {
	Id       int64
	PlayerId int64
	Name     string
	Corporation
	AccessMask
}

func NewPlayer(id int64, playerId int64, name string, corp Corporation, access AccessMask) Player {
	player := Player{
		Id:          id,
		PlayerId:    playerId,
		Name:        name,
		Corporation: corp,
		AccessMask:  access,
	}

	return player
}
