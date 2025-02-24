package business

import "context"

// right ports
type UserBusiness interface {
	Create(context.Context, *User) error
	Delete(context.Context, *UserID) error
	SearchAll(context.Context) (*Users, error)
}

// left ports
type UserStore interface {
	Create(context.Context, *User) error
	Delete(context.Context, *UserID) error
	SearchAll(context.Context) (*Users, error)
}
