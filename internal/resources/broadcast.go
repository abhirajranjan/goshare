package resources

import "time"

type file struct {
	Size     int    `json:"size"`
	FileName string `json:"filename"`
}

type BroadcastRequest struct {
	Files []file `json:"files"`
}

type BroadcastResponse struct {
	OTP            string
	GeneratingTime time.Time
	Interval       time.Duration
}
