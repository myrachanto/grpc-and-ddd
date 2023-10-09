package pasetos

import (
	"fmt"
	"time"

	"github.com/myrachanto/grpcgateway/src/support"
	httperrors "github.com/myrachanto/erroring"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker() (Maker, httperrors.HttpErr) {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file in routes")
	// }
	// symmetricKey := os.Getenv("EncryptionKey")
	symmetricKey := "Myrachanto"
	symmetricKey = support.Hasher2(symmetricKey)
	symmetricKey = symmetricKey[0:32]
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, httperrors.NewBadRequestError("Something is wrong with the symmetric key length")
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(data *Data, duration time.Duration) (string, *Payload, httperrors.HttpErr) {
	payload, err := NewPayload(data, duration)
	if err != nil {
		return "", nil, err
	}
	str, errs := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, httperrors.NewBadRequestError(fmt.Sprintf("something went Paseto Encryption, %d", errs))
	}
	return str, payload, nil
}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, httperrors.HttpErr) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("something went Paseto Encryption, %d", err))
	}
	errs := payload.Valid()
	if errs != nil {
		return nil, errs
	}
	return payload, nil
}
