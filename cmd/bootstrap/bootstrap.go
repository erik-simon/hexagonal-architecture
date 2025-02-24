package bootstrap

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/internal/app/users/infrastructure/driven/memory"
	cmdUser "github.com/hexagonal/internal/app/users/infrastructure/drives/command"
	"github.com/hexagonal/internal/app/users/infrastructure/drives/handlers"
	"github.com/spf13/cobra"
)

func Inject(ctx context.Context, logger *slog.Logger, a any) error {
	userCases, err := buildUserCases(ctx, logger)
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case *chi.Mux:
		return injectHandler(ctx, userCases, a)
	case *cobra.Command:
		return injectCommand(ctx, userCases, a)
	}

	return errors.New("unprocessable container")
}

func buildUserCases(_ context.Context, logger *slog.Logger) (_ *business.UserCases, err error) {
	userMap := make(map[business.UserID]*business.User)

	userStore, err := memory.NewUserStore(
		userMap,
		logger,
	)
	if err != nil {
		return
	}

	userCases, err := business.NewUserCases(userStore)
	if err != nil {
		return
	}

	return userCases, nil
}

func injectHandler(_ context.Context, userCases *business.UserCases, routes *chi.Mux) (err error) {
	userPostHandler, err := handlers.NewPostUserHandler(userCases)
	if err != nil {
		return
	}

	userDeleteHandler, err := handlers.NewDeleteUserHandler(userCases)
	if err != nil {
		return
	}

	userSearchAllHandler, err := handlers.NewGetAllUserHandler(userCases)
	if err != nil {
		return
	}

	return handlers.SetUserRoutes(
		routes,
		userPostHandler,
		userDeleteHandler,
		userSearchAllHandler,
	)
}

func injectCommand(_ context.Context, userCases *business.UserCases, cmd *cobra.Command) (err error) {
	commandCreateUser, err := cmdUser.NewCommandCreateUser(userCases)
	if err != nil {
		return
	}

	commandGetAllUsersUser, err := cmdUser.NewCommandGetAllUsers(userCases)
	if err != nil {
		return
	}

	commandRemoveUser, err := cmdUser.NewCommandRemoveUser(userCases)
	if err != nil {
		return
	}

	setCommandCreateUser, err := cmdUser.SetCommandCreateUser(commandCreateUser)
	if err != nil {
		return
	}

	setCommandGetAllUsers, err := cmdUser.SetCommandGetAllUsers(commandGetAllUsersUser)
	if err != nil {
		return
	}

	setCommandRemoveUser, err := cmdUser.SetCommandRemoveUser(commandRemoveUser)
	if err != nil {
		return
	}

	*cmd = cobra.Command{
		Use:   "hexagonal",
		Short: "cli to create users using hexagonal architecture",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(setCommandCreateUser)
	cmd.AddCommand(setCommandGetAllUsers)
	cmd.AddCommand(setCommandRemoveUser)
	return
}
