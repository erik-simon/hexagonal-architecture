package command

import (
	"context"
	"errors"

	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/command"
	"github.com/spf13/cobra"
)

func SetCommandRemoveUser(commandRemoveUser *CommandRemoveUser) (*cobra.Command, error) {
	if commandRemoveUser == nil {
		return nil, errors.New("missing 'command remove user' dependency")
	}
	cmd := &cobra.Command{
		Use:       "remove",
		Short:     "remove an user based on command",
		ValidArgs: []string{userIDFlag, usernameFlag},
		Run:       UserErrorCommand(commandRemoveUser.Remove()),
	}

	flags := cmd.PersistentFlags()

	_ = flags.String(userIDFlag, "", "indicates the user id")

	err := cmd.MarkPersistentFlagRequired(userIDFlag)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

type CommandRemoveUser struct {
	remover business.UserBusiness
}

func NewCommandRemoveUser(remover business.UserBusiness) (*CommandRemoveUser, error) {
	if remover == nil {
		return nil, errors.New("missing 'user business' dependency")
	}

	return &CommandRemoveUser{
		remover: remover,
	}, nil
}

func (c CommandRemoveUser) Remove() command.CommandFunc {
	return func(cmd *cobra.Command, args []string) (err error) {
		ctx := context.Background()

		if err = cmd.ParseFlags(args); err != nil {
			return
		}

		id, err := cmd.Flags().GetString(userIDFlag)
		if err != nil {
			return
		}

		userId, err := business.UserID(id).Validate()
		if err != nil {
			return
		}

		if err = c.remover.Delete(ctx, &userId); err != nil {
			return
		}

		cmd.Printf(
			"Deleted user with id %[1]v\n", userId,
		)
		return
	}
}
