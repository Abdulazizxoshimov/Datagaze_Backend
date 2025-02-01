package postgres

import (
	"fmt"
	"strings"

	"github.com/Abdulazizxoshimov/Datagaze_Backend/entity"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/logger"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

type userRepo struct {
	db *sqlx.DB
	log logger.Logger
}

func NewUserRepo(db *sqlx.DB, log logger.Logger) interfaces.User {
	return &userRepo{
		db: db,
		log: log,
	}
}

func (u *userRepo) CreateUser(user *entity.User) (*entity.User, error) {
	SqlStr := `INSERT INTO users (
		id,
		name,
		surname,
		username,
		email,
		password,
		role,
		refresh_t,
		created_at
	) VALUES (
		:id,
		:name,
		:surname,
		:username,
		:email,
		:password,
		:role,
		:refresh_t,
		:created_at
	)
`
	_, err := u.db.NamedExec(SqlStr, map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"surname":    user.SurName,
		"username":   user.UserName,
		"email":      user.Email,
		"password":   user.Password,
		"role":       user.Role,
		"refresh_t": user.RefreshT,
		"created_at": user.CreatedAt,
	})
	if err != nil {
		u.log.Error("error while adding user to database", zap.Error(err))
		return nil, err
	}
	filter := map[string]interface{}{
		"id": user.ID,
	}

	return u.GetUser(filter)
}

func (u *userRepo) UpdateUser(user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users
		SET 
			name = :name,
			surname = :surname,
			username = :username,
			email = :email,
			password = :password
		WHERE id = :id
	`

	_, err := u.db.NamedExec(query, map[string]interface{}{
		"id":       user.ID,
		"name":     user.Name,
		"surname":  user.SurName,
		"username": user.UserName,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		u.log.Error("error while updating user in database", zap.Error(err))
		return nil, err
	}
	filter := map[string]interface{}{
		"id"  :  user.ID,
	}
	responseUser, err := u.GetUser(filter)
	if err != nil {
		u.log.Error("error while getting user", zap.Error(err))
		return nil, err
	}

	return responseUser, nil
}

func (u *userRepo) GetUser(filter map[string]interface{}) (*entity.User, error) {
	query := `
		SELECT 
			id,
			name,
			surname,
			username,
			email,
			password,
			refresh_t,
			role,
			created_at
		FROM users
	`

	var args []interface{}
	var conditions []string
	i := 1

	for key, value := range filter {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var user entity.User
	pp.Println(query)
	pp.Println(args...)
	pp.Println(&user)

	err := u.db.Get(&user, query, args...)
	pp.Println(user)
	if err != nil {
		u.log.Error("error while retrieving user from database", zap.Error(err))
		return nil, err
	}

	return &user, nil
}


func (r *userRepo) GetAllUsers(page, limit int) (*entity.UserListResponse, error) {
	offset := (page - 1) * limit
	var users entity.UserListResponse

	var totalCount int
	err := r.db.Get(&totalCount, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, err
	}
	users.UserCount = totalCount
	query := `
		SELECT id, 
				name, 
				surname, 
				username, 
				email, 
				role, 
				created_at 
		FROM users
		ORDER BY created_at
		LIMIT $1 
		OFFSET $2
	`
	err = r.db.Select(&users.User, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (u *userRepo) DeleteUser(req *entity.DeleteRequest) error {
	query := `
		DELETE FROM users		
		WHERE id = $1
	`

	_, err := u.db.Exec(query, req.ID)
	if err != nil {
		u.log.Error("error while deleting user from database", zap.Error(err))
		return err
	}

	return nil
}

func (u *userRepo) IsUnique(field string, value string) (*entity.Response, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM users WHERE %s = $1", field)
	var count int
	err := u.db.QueryRow(query, value).Scan(&count)
	if err != nil {
		u.log.Error("error while checking uniqueness for field", zap.Error(err))
		return &entity.Response{Status: false}, err
	}
	if count != 0 {
		return &entity.Response{Status: true}, nil
	}
	return &entity.Response{Status: false}, nil
}

func (p *userRepo) UpdateRefresh(request *entity.UpdateRefresh) (*entity.Response, error) {
	sqlStr := `UPDATE users SET 
							refresh_token = :refresh_token 
				WHERE id = :user_id 
				AND deleted_at IS NULL`
	_, err := p.db.NamedExec(sqlStr, map[string]interface{}{
		"refresh_token": request.RefreshToken,
		"user_id":       request.UserID,
	})
	if err != nil {
		p.log.Error("error while updating refresh token for user_id ", zap.Error(err))
		return &entity.Response{Status: false}, err
	}
	return &entity.Response{Status: true}, nil
}

func (p *userRepo) UpdatePassword(request *entity.UpdatePassword) (*entity.Response, error) {
	sqlStr := `UPDATE users 
					SET 
						password = :password 
					WHERE id = :user_id 
					AND deleted_at IS NULL`
	_, err := p.db.NamedExec(sqlStr, map[string]interface{}{
		"password": request.NewPassword,
		"user_id":  request.UserID,
	})
	if err != nil {
		p.log.Error("error while updating password for user_id ", zap.Error(err))
		return &entity.Response{Status: false}, err
	}
	return &entity.Response{Status: true}, nil
}
