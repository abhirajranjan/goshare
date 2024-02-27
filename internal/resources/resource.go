package resources

import "time"

type File struct {
	Size     int    `json:"size"`
	FileName string `json:"filename"`
}

type BroadcastRequest struct {
	Files []File `json:"files"`
}

type BroadcastResponse struct {
	OTP            string        `json:"otp"`
	GeneratingTime time.Time     `json:"generating_time"`
	Interval       time.Duration `json:"interval"`
}

type LookupRequest struct {
	Otp string `json:"otp"`
}
