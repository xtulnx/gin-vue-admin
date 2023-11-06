package request

import (
	"github.com/gofrs/uuid/v5"
	jwt "github.com/golang-jwt/jwt/v4"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
}

func (R *CustomClaims) GetUserID() uint {
	return R.BaseClaims.ID
}
func (R *CustomClaims) GetUserUuid() uuid.UUID {
	return R.BaseClaims.UUID
}
func (R *CustomClaims) GetUserAuthorityId() uint {
	return R.BaseClaims.AuthorityId
}
func (R *CustomClaims) GetUserName() string {
	return R.BaseClaims.Username
}
func (R *CustomClaims) GetNickName() string {
	return R.BaseClaims.NickName
}
