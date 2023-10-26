package models

type BroadcastModel struct {
	Id    uint64
	Proto string
	Data  []byte
}
