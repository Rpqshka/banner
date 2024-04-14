package repository

import (
	"banner"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

type BannerPostgres struct {
	db *sqlx.DB
}

func NewBannerPostgres(db *sqlx.DB) *BannerPostgres {
	return &BannerPostgres{db: db}
}

func (r *BannerPostgres) CheckBanner(tagIds []int, featureId int) (bool, error) {
	checkQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE tag_ids = $1 AND feature_id = $2)", bannersTable)

	var exists bool
	err := r.db.QueryRow(checkQuery, pq.Array(tagIds), featureId).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return exists, errors.New("banner with specified tag_ids and feature_id already exists")
	}

	return exists, nil
}

func (r *BannerPostgres) CreateBanner(banner banner.Banner) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (tag_ids, feature_id, title, text, url, is_active, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, bannersTable)
	row := r.db.QueryRow(query,
		pq.Array(banner.TagIds), banner.FeatureId, banner.Content.Title, banner.Content.Text, banner.Content.Url, banner.IsActive,
		banner.CreatedAt, banner.UpdatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *BannerPostgres) GetBannerById(id int) (banner.Banner, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return banner.Banner{}, err
	}
	defer tx.Rollback()

	var content banner.Content
	var banner banner.Banner
	var tagIDs []uint8

	queryBanner := fmt.Sprintf("SELECT tag_ids, feature_id, is_active FROM %s WHERE id = $1", bannersTable)
	if err = r.db.QueryRow(queryBanner, id).Scan(&tagIDs, &banner.FeatureId, &banner.IsActive); err != nil {
		return banner, err
	}

	tagID := uint(0)
	for _, b := range tagIDs {
		tagID = tagID*10 + uint(b-'0')
	}

	queryContent := fmt.Sprintf("SELECT title, text, url FROM %s WHERE id = $1", bannersTable)
	if err = r.db.QueryRowx(queryContent, id).StructScan(&content); err != nil {
		return banner, err
	}

	banner.Content = content

	if err = tx.Commit(); err != nil {
		return banner, err
	}

	return banner, nil
}

func (r *BannerPostgres) UpdateBannerById(id int, banner banner.Banner) error {

	query := fmt.Sprintf(`
		UPDATE %s 
		SET 
			tag_ids = $1, 
			feature_id = $2, 
			title = $3, 
			text = $4, 
			url = $5, 
			is_active = $6, 
			updated_at = $7
		WHERE 
			id = $8
	`, bannersTable)
	_, err := r.db.Exec(query, pq.Array(banner.TagIds), banner.FeatureId, banner.Content.Title, banner.Content.Text,
		banner.Content.Url, banner.IsActive, banner.UpdatedAt, id)

	if err != nil {
		return err
	}
	return nil
}

func (r *BannerPostgres) DeleteBannerById(id int) error {
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", bannersTable)
	result, err := r.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *BannerPostgres) GetUserBanner(input banner.UserBannerInput, role string) (banner.Content, error) {
	var content banner.Content

	if role == "admin" {
		adminQuery := fmt.Sprintf("SELECT title, text, url FROM %s WHERE $1 = ANY(tag_ids) AND feature_id = $2", bannersTable)
		err := r.db.QueryRowx(adminQuery, input.TagId, input.FeatureId).StructScan(&content)
		if err != nil {
			return content, err
		}
		return content, nil
	} else {
		userQuery := fmt.Sprintf(`SELECT title, text, url FROM %s WHERE ($1 = ANY(tag_ids) AND feature_id = $2)
                                AND is_active = true`, bannersTable)
		err := r.db.QueryRowx(userQuery, input.TagId, input.FeatureId).StructScan(&content)
		if err != nil {
			return content, err
		}
		return content, nil
	}

}

func (r *BannerPostgres) GetAllBanners(input banner.FilterInput) ([]banner.Banner, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var banners []banner.Banner

	query := fmt.Sprintf(`
        SELECT tag_ids, feature_id, is_active, created_at, updated_at, title, text, url
        FROM %s 
        WHERE $1 = ANY(tag_ids) OR feature_id = $2
        LIMIT $3 OFFSET $4`,
		bannersTable)

	rows, err := tx.Query(query, input.TagId, input.FeatureId, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b banner.Banner
		var tagIDs []byte
		err = rows.Scan(&tagIDs, &b.FeatureId, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
			&b.Content.Title, &b.Content.Text, &b.Content.Url)
		if err != nil {
			return nil, err
		}

		tagIDsStr := string(tagIDs)
		tagIDsStr = strings.Trim(tagIDsStr, "{}")
		tagIDsSplit := strings.Split(tagIDsStr, ",")

		tagIDsInt := make([]int, len(tagIDsSplit))
		for i, idStr := range tagIDsSplit {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return nil, err
			}
			tagIDsInt[i] = id
		}

		b.TagIds = tagIDsInt

		banners = append(banners, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return banners, nil
}
