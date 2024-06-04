package otp

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
)

type Totp struct {
	ValidDuration time.Duration
	salt          []byte
}

type TotpOpts func(*Totp)

func NewTotp(opts ...TotpOpts) (otp Totp) {
	otp.setDefault()
	for _, opt := range opts {
		opt(&otp)
	}

	return
}

func (otp *Totp) setDefault() {
	otp.ValidDuration = time.Second * 30
}

func (otp *Totp) GenerateOTP(secret string) (token string, origin time.Time, interval time.Duration) {
	return otp.GenerateOTPAt(secret, time.Now())
}

func (otp *Totp) GenerateOTPAt(secret string, At time.Time) (token string, origin time.Time, interval time.Duration) {
	currentTime := At
	currentTimeUnix := currentTime.Unix()
	timeInterval := currentTimeUnix / int64(otp.ValidDuration.Seconds())

	// Convert the time interval to a byte slice
	timeBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		timeBytes[i] = byte(timeInterval & 0xff)
		timeInterval >>= 8
	}

	// Generate HMAC-SHA256 hash
	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write(timeBytes)
	if otp.salt != nil {
		hasher.Write(otp.salt)
	}
	hash := hasher.Sum(nil)

	// Dynamic truncation to extract 4 bytes
	offset := hash[len(hash)-1] & 0x0f
	binary := (int(hash[offset]) & 0x7f) << 24
	binary |= (int(hash[offset+1]) & 0xff) << 16
	binary |= (int(hash[offset+2]) & 0xff) << 8
	binary |= int(hash[offset+3]) & 0xff

	// Generate 6-digit Totp
	Totp := binary % 1000000
	return fmt.Sprintf("%06d", Totp), currentTime, otp.ValidDuration
}

func WithDuration(d time.Duration) TotpOpts {
	return func(t *Totp) {
		t.ValidDuration = d
	}
}

func WithSalt(salt []byte) TotpOpts {
	return func(t *Totp) {
		t.salt = salt
	}
}
