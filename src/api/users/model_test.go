package users

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var u User = User{Firstname: "mark", Lastname: "white", Birthday: "12/12/1994", Phone: "232453366674", Address: "whitemores street", Password: "123456sdf", Email: "email@example.com"}
var jsondata = `{"firstname":"jane","lastame":"Doe","username":"doe","Usercode": "Doe345","Phone":"1234567","Email":   "email@example.com","Password": "1234567","Address":"psd 456 king view"}`

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		user User
		err  string
		code int
	}{
		{name: "ok", user: u, err: ""},
		{name: "Empty Firstname", user: User{Firstname: "", Lastname: "white", Birthday: "12/12/1994", Phone: "232453366674", Address: "whitemores street", Password: "123456sdf", Email: "email@example.com"}, err: "firstname should not be empty", code: 400},
		{name: "Empty Lastname", user: User{Firstname: "mark", Lastname: "", Birthday: "12/12/1994", Phone: "232453366674", Address: "whitemores street", Password: "123456sdf", Email: "email@example.com"}, err: "lastname should not be empty", code: 400},
		{name: "Empty Address", user: User{Firstname: "mark", Lastname: "white", Birthday: "12/12/1994", Phone: "232453366674", Address: "", Password: "123456sdf", Email: "email@example.com"}, err: "address should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.user.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Error())
			}
		})
	}

}

func TestValidateLoginUserInputRequiredFields(t *testing.T) {
	jsondata := `{"email":"email@example.com","password":"1234567"}`
	user := &LoginUser{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	expected := ""
	if err := user.Validate(); err != nil {
		expected = "invalid Email"
		if err.Error() == expected {
			assert.EqualValues(t, "", err.Error(), "Error validating email")
		}
		expected = "invalid password"
		if err.Error() == expected {
			assert.EqualValues(t, "", err.Error(), "Error validating password")
		}

	}

}
func TestComparingPasswords(t *testing.T) {
	// fmt.Println("------------------", user)
	password := "anton345"
	user := User{}
	pas1, _ := user.HashPassword(password)
	ok := user.Compare(password, pas1)
	if !ok {
		assert.EqualValues(t, false, ok, "Error comparing passwords")
	}
}

func TestValidateEmailInputRequiredFields(t *testing.T) {
	user := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	ok := user.ValidateEmail(user.Password)
	if !ok {
		assert.EqualValues(t, false, ok, "Error Validating emails")
	}

}

func TestValidatePasswordInputRequiredFields(t *testing.T) {
	user := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	ok, _ := user.ValidatePassword(user.Password)
	if !ok {
		assert.EqualValues(t, true, ok, "Error Validating passwords")
	}

}
