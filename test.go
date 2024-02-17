// package main

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"fmt"
// )

// func mmain() {
// 	// Generate ECDSA key pair
// 	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Message to be signed
// 	message := []byte("hello")

// 	// Compute the hash of the message
// 	hash := sha256.Sum256(message)

// 	// Sign the hash with the private key
// 	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Verify the signature using the public key
// 	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
// 	if valid {
// 		fmt.Println("Signature is valid")
// 	} else {
// 		fmt.Println("Signature is not valid")
// 	}
// }

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
)

func generateTOTP(secret string) (string, error) {
	// Decode the secret from base32
	// Calculate the number of 30-second intervals since Unix epoch
	currentTime := time.Now().Unix()
	timeInterval := currentTime / 5

	// Convert the time interval to a byte slice
	timeBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		timeBytes[i] = byte(timeInterval & 0xff)
		timeInterval >>= 8
	}

	// Generate HMAC-SHA256 hash
	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write(timeBytes)
	hasher.Write(timeBytes)
	hash := hasher.Sum(nil)

	// Dynamic truncation to extract 4 bytes
	offset := hash[len(hash)-1] & 0x0f
	binary := (int(hash[offset]) & 0x7f) << 24
	binary |= (int(hash[offset+1]) & 0xff) << 16
	binary |= (int(hash[offset+2]) & 0xff) << 8
	binary |= int(hash[offset+3]) & 0xff

	fmt.Printf("%06d", binary)
	// Generate 6-digit TOTP
	totp := binary % 1000000
	return fmt.Sprintf("%06d", totp), nil
}

func main() {
	// Replace this with your secret key
	secret := "NBSWY3DPEB3W64TMMQ"

	totp, err := generateTOTP(secret)
	if err != nil {
		fmt.Println("Error generating TOTP:", err)
		return
	}

	fmt.Println("Generated TOTP:", totp)
}
