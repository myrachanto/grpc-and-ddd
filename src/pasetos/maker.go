package pasetos

import (
	"time"
)

type Maker interface {
	//creates a new token
	CreateToken(data *Data, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
