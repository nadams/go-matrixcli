package matrixext

import (
	"net/http"

	"github.com/matrix-org/gomatrix"
)

type Room struct {
	RoomID  string   `json:"room_id"`
	Servers []string `json:"servers"`
}

func GetRoomByAlias(cl *gomatrix.Client, alias string) (*Room, error) {
	var room Room

	if err := cl.MakeRequest(http.MethodGet, cl.BuildURL("directory", "room", alias), nil, &room); err != nil {
		return nil, err
	}

	return &room, nil
}
