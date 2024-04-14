package banner

type Content struct {
	Title string `json:"title" db:"title" binding:"required"`
	Text  string `json:"text" db:"text" binding:"required"`
	Url   string `json:"url" db:"url" binding:"required"`
}

type Banner struct {
	TagIds    []int   `json:"tag_ids" db:"tag_ids" binding:"required"`
	FeatureId int     `json:"feature_id" db:"feature_id" binding:"required"`
	Content   Content `json:"content" binding:"required"`
	IsActive  bool    `json:"is_active" db:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type UserBannerInput struct {
	TagId           int  `json:"tag_id"`
	FeatureId       int  `json:"feature_id"`
	UseLastRevision bool `json:"use_last_revision"`
}

type FilterInput struct {
	TagId     int `json:"tag_id"`
	FeatureId int `json:"feature_id"`
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}
