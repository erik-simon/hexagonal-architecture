package command

import (
	"context"
	"errors"

	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/command"
	"github.com/spf13/cobra"
)

func SetCommandCreateUser(commandCreateUser *CommandCreateUser) (*cobra.Command, error) {
	if commandCreateUser == nil {
		return nil, errors.New("missing 'command create user' dependency")
	}
	cmd := &cobra.Command{
		Use:       "create",
		Short:     "create a new user based on command",
		ValidArgs: []string{userIDFlag, usernameFlag},
		Run:       UserErrorCommand(commandCreateUser.Create()),
	}

	flags := cmd.PersistentFlags()

	_ = flags.String(userIDFlag, "", "indicates the user id")
	_ = flags.String(usernameFlag, "", "indicates the user name")

	err := cmd.MarkPersistentFlagRequired(userIDFlag)
	if err != nil {
		return nil, err
	}

	err = cmd.MarkPersistentFlagRequired(usernameFlag)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

type CommandCreateUser struct {
	creator business.UserBusiness
}

func NewCommandCreateUser(creator business.UserBusiness) (*CommandCreateUser, error) {
	if creator == nil {
		return nil, errors.New("missing 'user business' dependency")
	}

	return &CommandCreateUser{
		creator: creator,
	}, nil
}

func (c CommandCreateUser) Create() command.CommandFunc {
	return func(cmd *cobra.Command, args []string) (err error) {
		ctx := context.Background()

		if err = cmd.ParseFlags(args); err != nil {
			return
		}

		userID, err := cmd.Flags().GetString(userIDFlag)
		if err != nil {
			return
		}

		username, err := cmd.Flags().GetString(usernameFlag)
		if err != nil {
			return
		}

		userCmd := &UserCommand{
			Id:       userID,
			UserName: username,
		}

		user, err := userCmd.ToBusiness()
		if err != nil {
			return
		}

		if err = c.creator.Create(ctx, user); err != nil {
			return
		}

		cmd.Printf(
			"Created user: %[1]v\n", userCmd,
		)
		return
	}
}
