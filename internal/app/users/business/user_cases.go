package business

import (
	"context"
	"errors"
)

type UserCases struct {
	UserStore UserStore
}

func NewUserCases(UserStore UserStore) (*UserCases, error) {
	if UserStore == nil {
		return nil, errors.New("missing 'user repository' dependency")
	}

	return &UserCases{
		UserStore: UserStore,
	}, nil
}

func (u *UserCases) Create(ctx context.Context, user *User) error {
	return u.UserStore.Create(ctx, user)
}

func (u *UserCases) Delete(ctx context.Context, userID *UserID) error {
	return u.UserStore.Delete(ctx, userID)
}

func (u *UserCases) SearchAll(ctx context.Context) (*Users, error) {
	return u.UserStore.SearchAll(ctx)
}
