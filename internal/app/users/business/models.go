package business

import (
	"fmt"
	"strings"

	"github.com/hexagonal/pkg/identifier"
)

var _ fmt.Stringer = UserID("")

type UserID string

func (u UserID) Validate() (UserID, error) {
	if _, err := identifier.Identifier(u).Validate(); err != nil {
		return "", WrapErrorf(
			err,
			ErrorInvalidArgument,
			"invalid identifier %[1]v",
			u,
		)
	}
	return u, nil
}

func (u UserID) String() string {
	return string(u)
}

var _ fmt.Stringer = UserName("")

type UserName string

func (u UserName) Validate() (UserName, error) {
	if strings.TrimSpace(string(u)) == "" {
		return "", NewErrorf(
			ErrorInvalidArgument,
			"missing username field",
		)
	}
	return u, nil
}

func (u UserName) String() string {
	return string(u)
}

type (
	User struct {
		userID   UserID
		userName UserName
	}
	Users []User
)

func NewUser(
	userID,
	userName string,
) (*User, error) {
	id, err := UserID(userID).Validate()
	if err != nil {
		return nil, err
	}

	name, err := UserName(userName).Validate()
	if err != nil {
		return nil, err
	}

	return &User{
		userID:   id,
		userName: name,
	}, nil
}

func (u *User) UserID() *UserID {
	return &u.userID
}

func (u *User) UserName() *UserName {
	return &u.userName
}
