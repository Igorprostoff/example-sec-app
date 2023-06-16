package v1

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"net/http"

	"securewebapp/internal/entity"
	"securewebapp/internal/usecase"
	"securewebapp/pkg/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	t usecase.UserRepo
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserRepo, l logger.Interface) {
	r := &userRoutes{t, l}
	handler.POST("/login", r.login)

}

type loginRequest struct {
	Email    string `json:"email"       binding:"required"  example:"igor@prostov.ru"`
	Password string `json:"password"  binding:"required"  example:"owgjng2u38^^2hbf"`
}

type loginResponse struct {
	IsAuthenticated bool `json:"isAuthenticated"       binding:"required"`
}

// @Summary     Login
// @Description Login a user
// @ID          login
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body loginReques
// @Success     200 {object} isAuthenticated
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /login [post]
func (r *userRoutes) login(c *gin.Context) {
	session := sessions.Default(c)

	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - login")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	hash := sha256.New()
	hash_bytes := hash.Sum([]byte(request.Password))
	encoded_hash := b64.StdEncoding.EncodeToString([]byte(hash_bytes))
	isLogin, err := r.t.Login(
		c.Request.Context(),
		entity.User{
			Email:    request.Email,
			Password: encoded_hash,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - login")
		errorResponse(c, http.StatusInternalServerError, "login service problems")
		return
	}

	response := loginResponse{IsAuthenticated: isLogin}
	session.Set("auth", isLogin)
	session.Save()
	c.JSON(http.StatusOK, response)

}
