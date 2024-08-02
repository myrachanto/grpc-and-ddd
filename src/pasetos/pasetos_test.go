package pasetos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var data = &Data{
	Username: "myrachanto",
	Email:    "myrachanto@gmail.com",
}

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker()
	require.Equal(t, nil, err)

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, payload, err := maker.CreateToken(data, duration)
	require.Equal(t, nil, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Equal(t, nil, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.IDs)
	require.Equal(t, data.Username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker()
	require.Equal(t, nil, err)
	duration := -time.Minute
	token, payload, err := maker.CreateToken(data, duration)
	require.Equal(t, nil, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Equal(t, ErrExpiredToken, err.Error())
	require.Nil(t, payload)
}
