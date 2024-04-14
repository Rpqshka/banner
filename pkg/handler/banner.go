package handler

import (
	"banner"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type createBannerInput struct {
	TagsIds   []int          `json:"tag_ids" binding:"required"`
	FeatureId int            `json:"feature_id" binding:"required"`
	Content   banner.Content `json:"content" binding:"required"`
	IsActive  bool           `json:"is_active"`
}

func (h *Handler) createBanner(c *gin.Context) {
	var input createBannerInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if role != "admin" {
		newErrorResponse(c, http.StatusForbidden, "only admin can create banner")
		return
	}

	exist, err := h.services.Banner.CheckBanner(input.TagsIds, input.FeatureId)

	if exist {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	currentTime := getTime()

	banner := banner.Banner{
		TagIds:    input.TagsIds,
		FeatureId: input.FeatureId,
		Content:   input.Content,
		IsActive:  input.IsActive,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	id, err := h.services.Banner.CreateBanner(banner)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"banner_id": id,
	})
}

type updateContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type updateBannerInput struct {
	TagsIds   []int         `json:"tag_ids"`
	FeatureId int           `json:"feature_id"`
	Content   updateContent `json:"content"`
	IsActive  bool          `json:"is_active"`
}

func (h *Handler) updateBanner(c *gin.Context) {
	var input updateBannerInput

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if role != "admin" {
		newErrorResponse(c, http.StatusForbidden, "only admin can update banner")
		return
	}

	oldBanner, err := h.services.GetBannerById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	currentTime := getTime()

	content := banner.Content{
		Title: input.Content.Title,
		Text:  input.Content.Text,
		Url:   input.Content.Url,
	}

	banner := banner.Banner{
		TagIds:    input.TagsIds,
		FeatureId: input.FeatureId,
		Content:   content,
		IsActive:  input.IsActive,
		UpdatedAt: currentTime,
	}

	updatedBanner := getUpdatedBanner(oldBanner, banner)

	if err = h.services.Banner.UpdateBannerById(id, updatedBanner); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Handler) deleteBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	role, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if role != "admin" {
		newErrorResponse(c, http.StatusForbidden, "only admin can delete banner")
		return
	}

	if err = h.services.DeleteBannerById(id); err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, map[string]interface{}{})
}

func (h *Handler) getUserBanner(c *gin.Context) {
	var input banner.UserBannerInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	content, err := h.services.GetUserBanner(input, role)
	if err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *Handler) getAllBanners(c *gin.Context) {
	var input banner.FilterInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if role != "admin" {
		newErrorResponse(c, http.StatusForbidden, "only admin can get all banners")
		return
	}

	type getAllBannersResponse struct {
		Data []banner.Banner `json:" "`
	}

	banners, err := h.services.GetAllBanners(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllBannersResponse{
		Data: banners,
	})
}
