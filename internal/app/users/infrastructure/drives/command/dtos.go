package command

import (
	"encoding/json"
	"strings"

	"github.com/hexagonal/internal/app/users/business"
)

type (
	UserCommand struct {
		Id       string `json:"id"`
		UserName string `json:"username"`
	}
	UsersCommand []*UserCommand
)

func (users UsersCommand) String() string {
	strJsonMessage := &strings.Builder{}

	for _, user := range users {
		jsonMessage, _ := json.MarshalIndent(user, "", "  ")
		strJsonMessage.WriteString(string(jsonMessage))
	}

	return strJsonMessage.String()
}

func (s *UserCommand) String() string {
	jsonMessage, _ := json.MarshalIndent(s, "", "  ")
	return string(jsonMessage)
}

func (u *UserCommand) ToBusiness() (*business.User, error) {
	user, err := business.NewUser(u.Id, u.UserName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserCommand(user *business.User) *UserCommand {
	return &UserCommand{
		Id:       user.UserID().String(),
		UserName: user.UserName().String(),
	}
}

func NewUsersCommand(users *business.Users) UsersCommand {
	usersCommand := make(UsersCommand, 0, len(*users))

	for _, user := range *users {
		usersCommand = append(usersCommand,
			NewUserCommand(&user),
		)
	}

	return usersCommand
}
