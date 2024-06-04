package resources

import "time"

type Config struct {
	Crypto Crypto
	Server ServerConfig
}

func (c *Config) LoadDefault() {
	c.Crypto = Crypto{
		P: []byte{1, 2, 3, 4, 5, 6, 7},
		G: []byte{1, 2, 3},
	}

	c.Server = ServerConfig{
		HostPort: ":8000",
	}
}

type Crypto struct {
	P []byte
	G []byte
}

type ServerConfig struct {
	HostPort string
}

type Media struct {
	Name string `json:"name"`
}

type EventType struct {
	Type string

	Otp       string    `json:",omitempty"`
	ValidTill time.Time `json:",omitempty"`

	Code        string `json:",omitempty"`
	ResponseURL string `json:",omitempty"`
}

func OTPEvent(otp string, validTill time.Time) EventType {
	return EventType{
		Type:      "otp",
		Otp:       otp,
		ValidTill: validTill,
	}
}

func RecieverEvent(recvCode string, responseURL string) EventType {
	return EventType{
		Type:        "receiver",
		Code:        recvCode,
		ResponseURL: responseURL,
	}
}
