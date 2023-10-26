package handler

type RequestBroadcastModel struct {
	Proto string `json:"proto"`
	Data  string `json:"data"`
	Id    uint64 `json:"id"`
}

type RequestRoomModel struct {
	RoomID uint64 `json:"room_id"`
}
