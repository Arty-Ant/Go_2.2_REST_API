package api

import (
	db "Bankstore/db/sqlc"
	"Bankstore/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
type createUserResponse struct {
	Username          string             `json:"username"`
	FullName          string             `json:"full_name"`
	Email             string             `json:"email"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	// получаем JSON и выполняем десериализацию
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// получаем hash паролья по строке пароля
	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// создаем структуру для запроса на создание нового пользователя в БД
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	// создаем пользователя в БД и обрабатывем ошибки(уникальность и БД)
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// сформировать структуру для ответа
	rsp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
