package handlers

import "github.com/hexagonal/internal/app/users/business"

type (
	UserResponse struct {
		Id       string `json:"id"`
		UserName string `json:"username"`
	}
	UsersResponse []*UserResponse
)

func NewUserResponse(user *business.User) *UserResponse {
	return &UserResponse{
		Id:       user.UserID().String(),
		UserName: user.UserName().String(),
	}
}

func NewUsersResponse(users *business.Users) UsersResponse {
	usersResponse := make(UsersResponse, 0, len(*users))

	for _, user := range *users {
		usersResponse = append(usersResponse,
			NewUserResponse(&user),
		)
	}

	return usersResponse
}

func (u *UserResponse) ToBusiness() (*business.User, error) {
	user, err := business.NewUser(u.Id, u.UserName)
	if err != nil {
		return nil, err
	}

	return user, nil
}
