package memory

import (
	"context"
	"errors"
	"log/slog"

	"github.com/hexagonal/internal/app/users/business"
)

var _ business.UserStore = &UserStore{}

type UserStore struct {
	storage map[business.UserID]*business.User
	logger  *slog.Logger
}

func NewUserStore(
	storage map[business.UserID]*business.User,
	logger *slog.Logger,
) (*UserStore, error) {
	if storage == nil {
		return nil, errors.New("missing 'storage' dependency")
	}

	return &UserStore{
		storage: storage,
		logger:  logger,
	}, nil
}

func (u *UserStore) Create(ctx context.Context, user *business.User) error {
	if _, exists := u.storage[*user.UserID()]; exists {
		err := errors.New("23505 duplicate key value violates unique constraint 'user_id'")

		u.logger.ErrorContext(ctx, err.Error())

		return business.WrapErrorf(
			err,
			business.ErrorDuplicateRecord,
			"duplicated user with id '%[1]s'",
			user.UserID().String(),
		)
	}

	u.storage[*user.UserID()] = user
	return nil
}

func (u *UserStore) Delete(ctx context.Context, userID *business.UserID) error {
	if _, exists := u.storage[*userID]; !exists {
		err := errors.New("no rows")

		u.logger.ErrorContext(ctx, err.Error())

		return business.WrapErrorf(
			err,
			business.ErrorRecordNotFound,
			"user with id '%[1]s' not found",
			userID.String(),
		)
	}

	delete(u.storage, *userID)
	return nil
}

func (u *UserStore) SearchAll(ctx context.Context) (*business.Users, error) {
	users := make(business.Users, 0, len(u.storage))

	for _, user := range u.storage {
		users = append(users, *user)
	}

	if len(users) == 0 {
		u.logger.InfoContext(ctx, "not found users")
	}

	return &users, nil
}
