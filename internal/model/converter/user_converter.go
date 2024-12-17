package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		UUID:       user.UUID,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		VerifiedAt: user.VerifiedAt,
	}
}

func UserToTokenResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		UUID:         user.UUID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
}

func DTOUserToResponse(user entity.User) model.UserResponse {
	return model.UserResponse{
		UUID:       user.UUID,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		VerifiedAt: user.VerifiedAt,
	}
}
