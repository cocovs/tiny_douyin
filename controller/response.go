package controller

import "github.com/cocovs/tiny-douyin/models"

type UserLoginResponse struct {
	models.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	models.Response
	User models.User `json:"user"`
}
