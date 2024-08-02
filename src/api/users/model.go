package users

import (
	"fmt"
	"regexp"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname      string             `json:"firstname,omitempty"`
	Lastname       string             `json:"lastname,omitempty"`
	Username       string             `json:"username,omitempty"`
	Birthday       string             `json:"birthday,omitempty"`
	Address        string             `json:"address,omitempty"`
	Phone          string             `json:"phone,omitempty"`
	Email          string             `json:"email,omitempty"`
	Password       string             `json:"password,omitempty"`
	HashedPassword string             `json:"hashed_password,omitempty"`
	Usercode       string             `json:"usercode,omitempty"`
	Role           string             `json:"role,omitempty"`
	Picture        string             `json:"picture,omitempty"`
	Base           `json:"base,omitempty"`
}
type UserDto struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty"`
	Username  string             `json:"username,omitempty"`
	Birthday  string             `json:"birthday,omitempty"`
	Address   string             `json:"address,omitempty"`
	Phone     string             `json:"phone,omitempty"`
	Email     string             `json:"email,omitempty"`
	Usercode  string             `json:"usercode,omitempty"`
	Role      string             `json:"role,omitempty"`
	Picture   string             `json:"picture,omitempty"`
	Base      `json:"base,omitempty"`
}
type LoginUser struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Auth struct {
	Usercode            string    `json:"usercode,omitempty"`
	UserName            string    `json:"username,omitempty"`
	Picture             string    `json:"picture,omitempty"`
	Token               string    `bson:"token" json:"token,omitempty"`
	TokenExpires        time.Time `json:"token_expires,omitempty"`
	RefleshToken        string    `json:"reflesh_token,omitempty"`
	RefleshTokenExpires time.Time `json:"reflesh_token_expires,omitempty"`
	SessionCode         string    `json:"session_code,omitempty"`
	Role                string    `json:"role,omitempty"`
}
type Base struct {
	Created_At time.Time  `bson:"created_at"`
	Updated_At time.Time  `bson:"updated_at"`
	Delete_At  *time.Time `bson:"deleted_at"`
}

func (user User) UserConvter() *UserDto {
	return &UserDto{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Birthday:  user.Birthday,
		Email:     user.Email,
		Address:   user.Address,
		Phone:     user.Phone,
		Usercode:  user.Usercode,
		Role:      user.Role,
		Picture:   user.Picture,
		Base:      user.Base,
	}
}

func (user User) ValidateEmail(email string) (matchedString bool) {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return false
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func (user User) ValidatePassword(password string) (bool, error) {
	if len(password) == 0 {
		return false, fmt.Errorf("your password is empty")
	}
	if len(password) < 5 {
		return false, fmt.Errorf("your password need more characters")
	} else if len(password) > 32 {
		return false, fmt.Errorf("your password is way too long")
	}
	return true, nil
}
func (user User) HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", fmt.Errorf("your password Must not be empty")
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("something went wrong")
	}
	return string(pass), nil

}

func (user User) Compare(p1, p2 string) bool {
	if len(p1) == 0 {
		return false
	}
	if len(p2) == 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}
func (user LoginUser) Compare(p1, p2 string) bool {
	stringresults := httperrors.ValidStringNotEmpty(p1)
	if stringresults.Noerror() {
		return false
	}
	stringresults2 := httperrors.ValidStringNotEmpty(p2)
	if stringresults2.Noerror() {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}
func (u User) Validate() error {
	if u.Firstname == "" {
		return fmt.Errorf("firstname should not be empty")
	}
	if u.Lastname == "" {
		return fmt.Errorf("lastname should not be empty")
	}
	if u.Address == "" {
		return fmt.Errorf("address should not be empty")
	}
	if u.Email == "" {
		return fmt.Errorf("email should not be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password should not be empty")
	}
	return nil
}

func (u LoginUser) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if u.Password == "" {
		return fmt.Errorf("invalid password")
	}
	return nil
}
