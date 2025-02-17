package pasetos

import (
	"fmt"
	"time"

	"github.com/myrachanto/grpcgateway/src/support"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker() (Maker, error) {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file in routes")
	// }
	// symmetricKey := os.Getenv("EncryptionKey")
	symmetricKey := "Myrachanto"
	symmetricKey = support.Hasher2(symmetricKey)
	symmetricKey = symmetricKey[0:32]
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("something is wrong with the symmetric key length")
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(data *Data, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(data, duration)
	if err != nil {
		return "", nil, err
	}
	str, errs := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if errs != nil {
		return "", nil, fmt.Errorf("something went paseto encryption, %s", errs)
	}
	return str, payload, nil
}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("something went Paseto Encryption, %d", err)
	}
	errs := payload.Valid()
	if errs != nil {
		return nil, errs
	}
	return payload, nil
}
