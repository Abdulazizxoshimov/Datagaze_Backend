package repo

import (
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/logger"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/interfaces"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/postgres"
	"github.com/jmoiron/sqlx"
)

type StorageI interface{
	User() interfaces.User
}

type storagePG struct {
	user interfaces.User
}

func NewStoragePG (db *sqlx.DB,log logger.Logger)StorageI{
	return &storagePG{
		user: postgres.NewUserRepo(db, log),
	}
}

func (s *storagePG)User()interfaces.User{
	return s.user
}

