package handler

import (
	"banner"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
	userCtxId           = "userId"
	userCtxRole         = "role"
)

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func comparePasswordHash(hash, inputPassword string) error {
	hashedPassword := []byte(hash)
	inputPasswordBytes := []byte(inputPassword)

	return bcrypt.CompareHashAndPassword(hashedPassword, inputPasswordBytes)
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	//parse token
	userId, role, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtxId, userId)
	c.Set(userCtxRole, role)
}

func getUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get(userCtxRole)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user role not found")
		return "", errors.New("user role not found")
	}

	roleStr, ok := role.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user role is of invalid type")
		return "", errors.New("user role is of invalid type")
	}

	return roleStr, nil
}

func getTime() string {
	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02T15:04:05.999Z")

	return formattedTime
}

func getUpdatedBanner(oldBanner banner.Banner, inputBanner banner.Banner) banner.Banner {
	updatedBanner := oldBanner

	if len(inputBanner.TagIds) > 0 {
		updatedBanner.TagIds = inputBanner.TagIds
	}

	if inputBanner.FeatureId != 0 {
		updatedBanner.FeatureId = inputBanner.FeatureId
	}

	if inputBanner.IsActive != oldBanner.IsActive {
		updatedBanner.IsActive = inputBanner.IsActive
	}

	//Content
	if inputBanner.Content.Title != "" {
		updatedBanner.Content.Title = inputBanner.Content.Title
	}

	if inputBanner.Content.Text != "" {
		updatedBanner.Content.Text = inputBanner.Content.Text
	}

	if inputBanner.Content.Url != "" {
		updatedBanner.Content.Url = inputBanner.Content.Url
	}

	updatedBanner.UpdatedAt = inputBanner.UpdatedAt

	return updatedBanner
}
