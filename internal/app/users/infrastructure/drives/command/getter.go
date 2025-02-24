package command

import (
	"context"
	"errors"

	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/command"
	"github.com/spf13/cobra"
)

func SetCommandGetAllUsers(commandGetAllUsers *CommandGetAllUsers) (*cobra.Command, error) {
	if commandGetAllUsers == nil {
		return nil, errors.New("missing 'command get all users' dependency")
	}

	cmd := &cobra.Command{
		Use:       "get-all",
		Short:     "get all users based on command",
		ValidArgs: []string{userIDFlag, usernameFlag},
		Run:       UserErrorCommand(commandGetAllUsers.GetAll()),
	}

	return cmd, nil
}

type CommandGetAllUsers struct {
	getter business.UserBusiness
}

func NewCommandGetAllUsers(getter business.UserBusiness) (*CommandGetAllUsers, error) {
	if getter == nil {
		return nil, errors.New("missing 'user business' dependency")
	}

	return &CommandGetAllUsers{
		getter: getter,
	}, nil
}

func (c CommandGetAllUsers) GetAll() command.CommandFunc {
	return func(cmd *cobra.Command, args []string) (err error) {
		ctx := context.Background()

		users, err := c.getter.SearchAll(ctx)
		if err != nil {
			return
		}

		usersCommand := NewUsersCommand(users)

		cmd.Println(usersCommand)

		return
	}
}
