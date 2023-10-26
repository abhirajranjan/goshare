package pubsub

const DirectProtoTag = "direct"

func directProto(Id uint64, data []byte) (resid uint64, resData []byte) {
	return Id, data
}
