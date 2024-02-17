package resources

type File struct {
	Size     int    `json:"size"`
	FileName string `json:"filename"`
}

type BroadcastRequest struct {
	CurrentDeviceId []byte `json:"id"`
	Files           []File `json:"files"`
}
