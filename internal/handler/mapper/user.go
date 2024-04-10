package mapper

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"avito-backend-trainee-2024/internal/handler/request"
	"avito-backend-trainee-2024/internal/handler/response"
)

func MapRegisterRequestToEntity(req *request.RegisterRequest, isAdmin bool) entity.User {
	return entity.User{
		Username:       req.Username,
		IsAdmin:        isAdmin,
		HashedPassword: req.Password,
	}
}

func MapUserToRegisterUserResponse(user *entity.User) response.RegisterUserResponse {
	return response.RegisterUserResponse{
		Username: user.Username,
	}
}

func MapUserToRegisterAdminResponse(user *entity.User) response.RegisterAdminResponse {
	return response.RegisterAdminResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
