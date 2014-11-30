// player
package models

import (
	"time"
)

type Player struct {
	Id   int64
	Name string
	Corporation
	AccessMask
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPlayer(name string, corp Corporation, access AccessMask) Player {
	player := Player{
		Id:          -1,
		Name:        name,
		Corporation: corp,
		AccessMask:  access,
	}

	return player
}
