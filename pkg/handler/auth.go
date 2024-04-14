package handler

import (
	"banner"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type registerInput struct {
	NickName        string `json:"nickname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Role            string `json:"role" binding:"required"`
}

func (h *Handler) register(c *gin.Context) {
	var input registerInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !validateEmail(input.Email) {
		newErrorResponse(c, http.StatusBadRequest, "enter a different email")
		return
	}

	if input.Password != input.PasswordConfirm {
		newErrorResponse(c, http.StatusBadRequest, "passwords does not match")
		return
	}

	_, err := h.services.Authorization.CheckNickNameAndEmail(input.NickName, input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	user := banner.User{
		NickName: input.NickName,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	}
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type loginInput struct {
	NickName string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) login(c *gin.Context) {
	var input loginInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := h.services.GetPasswordHash(input.NickName)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = comparePasswordHash(passwordHash, input.Password); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.NickName, passwordHash)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
