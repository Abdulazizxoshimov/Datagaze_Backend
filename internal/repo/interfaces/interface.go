package interfaces

import (
	"context"
	"time"

	"github.com/Abdulazizxoshimov/Datagaze_Backend/entity"
)

type Redis interface{
	Set(ctx context.Context, key string, value any, expiration time.Duration)error
	Get(ctx context.Context, key string)([]byte, error)
	Del(ctx context.Context, key string)error
}

type User interface {
	CreateUser(*entity.User)(*entity.User, error)
	UpdateUser(*entity.User)(*entity.User, error)
	GetUser(map[string]interface{})(*entity.User, error)
	GetAllUsers(int, int)(*entity.UserListResponse,  error)
	DeleteUser(*entity.DeleteRequest)error
	IsUnique(field string, value string) (*entity.Response, error) 
	UpdateRefresh(*entity.UpdateRefresh) (*entity.Response, error) 
	UpdatePassword(*entity.UpdatePassword)(*entity.Response, error) 
	
}	