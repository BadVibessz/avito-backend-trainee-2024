package banner

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	sliceutils "avito-backend-trainee-2024/pkg/utils/slice"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"strconv"
	"strings"
)

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

// general method for fetching banners with where condition
func (r *Repo) getBannersWhere(ctx context.Context, condition string, offset, limit int) ([]*entity.Banner, error) {
	var query string

	// TODO: maybe without join but two sql queries?
	if limit == math.MaxInt64 {
		query = fmt.Sprintf(`SELECT banner.id,  feature_id, is_active, created_at, updated_at, title, text, url, c.content_id 
FROM banner JOIN public.content c ON c.content_id = banner.content_id %v
ORDER BY banner.feature_id OFFSET %v`, condition, offset)
	} else {
		query = fmt.Sprintf(`SELECT banner.id,  feature_id, is_active, created_at, updated_at, title, text, url, c.content_id 
FROM banner JOIN public.content c ON c.content_id = banner.content_id %v
ORDER BY banner.feature_id LIMIT %v OFFSET %v`, condition, limit, offset)

	}

	// execute in transaction
	tx, err := r.DB.BeginTxx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	rows, err := tx.Queryx(query)
	if err != nil {
		return nil, err
	}

	var banners []*entity.Banner

	for rows.Next() {
		var banner entity.Banner

		if err = rows.StructScan(&banner); err != nil {
			return nil, err
		}

		banners = append(banners, &banner)
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// find tagIds associated with each banner
	rows, err = tx.Queryx("SELECT * FROM banner_tag") // todo: performance bottleneck (understand how to use offset and limit here for optimization)
	if err != nil {
		return nil, err
	}

	var bannerTagsMap = make(map[int][]entity.BannerTag)

	for rows.Next() {
		var bannerTag entity.BannerTag

		if err = rows.StructScan(&bannerTag); err != nil {
			return nil, err
		}

		_, exists := bannerTagsMap[bannerTag.BannerID]
		if !exists {
			bannerTagsMap[bannerTag.BannerID] = []entity.BannerTag{bannerTag}
		} else {
			bannerTagsMap[bannerTag.BannerID] = append(bannerTagsMap[bannerTag.BannerID], bannerTag)
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	for _, banner := range banners {
		_, exists := bannerTagsMap[banner.ID]
		if exists {
			for _, tag := range bannerTagsMap[banner.ID] {
				banner.TagIDs = append(banner.TagIDs, tag.TagID)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *Repo) GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error) {
	return r.getBannersWhere(ctx, "", offset, limit)
}

func (r *Repo) GetBannerByID(ctx context.Context, id int) (*entity.Banner, error) {
	query := fmt.Sprintf(`SELECT banner.id,
       is_active,
       title,
       text,
       url,
       array_agg(bt.tag_id ORDER BY bt.tag_id) AS tag_ids
FROM banner
         JOIN public.content c ON c.content_id = banner.content_id
         JOIN public.banner_tag bt ON banner.id = bt.banner_id
WHERE banner.id = %v
GROUP BY banner.id, c.content_id`,
		id,
	)

	rows, err := r.DB.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	type Row struct {
		ID        int    `db:"id"`
		IsActive  bool   `db:"is_active"`
		Title     string `db:"title"`
		Text      string `db:"text"`
		Url       string `db:"url"`
		TagIDsStr string `db:"tag_ids"`
		TagIDsInt []int
	}

	var row Row

	if rows.Next() {
		if err = rows.StructScan(&row); err != nil {
			return nil, err
		}

		// row.TagIDs have structure {1,2,...}
		row.TagIDsInt = make([]int, 0, len(row.TagIDsStr)-2)

		for _, tagStr := range strings.Split(row.TagIDsStr[1:len(row.TagIDsStr)-1], ",") {
			tag, convErr := strconv.Atoi(tagStr)
			if convErr != nil {
				return nil, convErr
			}

			row.TagIDsInt = append(row.TagIDsInt, tag)
		}
	}

	content := entity.Content{
		Title: row.Title,
		Text:  row.Text,
		Url:   row.Url,
	}

	return &entity.Banner{
			Content:  content,
			IsActive: row.IsActive,
		},
		nil
}

/*
GetBannerByFeatureAndTags

	TODO: this method can take long time,
		openapi.yaml spec says that we need to return only banner.Content to user,
		so there's no need to return from db banner model with all initialized fields
*/
func (r *Repo) GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error) {
	query := fmt.Sprintf(`SELECT banner.id,
       is_active,
       title,
       text,
       url,
       array_agg(bt.tag_id ORDER BY bt.tag_id) AS tag_ids
FROM banner
         JOIN public.content c ON c.content_id = banner.content_id
         JOIN public.banner_tag bt ON banner.id = bt.banner_id
WHERE feature_id = %v
GROUP BY banner.id, c.content_id
`, featureID)

	dbRows, err := r.DB.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	type Row struct {
		ID        int    `db:"id"`
		IsActive  bool   `db:"is_active"`
		Title     string `db:"title"`
		Text      string `db:"text"`
		Url       string `db:"url"`
		TagIDsStr string `db:"tag_ids"` // todo: pgx cannot convert sql.array to golang slice, maybe use pq instead?
		TagIDsInt []int
	}

	var rows []*Row

	for dbRows.Next() {
		var row Row

		if err = dbRows.StructScan(&row); err != nil {
			return nil, err
		}

		// row.TagIDs have structure {1,2,...}
		row.TagIDsInt = make([]int, 0, len(row.TagIDsStr)-2)

		for _, tagStr := range strings.Split(row.TagIDsStr[1:len(row.TagIDsStr)-1], ",") {
			tag, convErr := strconv.Atoi(tagStr)
			if convErr != nil {
				return nil, convErr
			}

			row.TagIDsInt = append(row.TagIDsInt, tag)
		}

		rows = append(rows, &row)
	}

	// each row represents banner with banner.feature_id = featureID => find banner with banner.tag_ids = tagIDs
	for _, row := range rows {
		if sliceutils.Equals(row.TagIDsInt, tagIDs) { // todo: here tagIDs gotta be sorted by asc, row.TagIDs already sorted
			content := entity.Content{
				Title: row.Title,
				Text:  row.Text,
				Url:   row.Url,
			}

			return &entity.Banner{
					Content:  content,
					IsActive: row.IsActive,
				},
				nil
		}
	}

	return nil, nil // todo: maybe return error?
}

func (r *Repo) CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error) {
	// execute in transaction
	tx, err := r.DB.BeginTxx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	// firstly add content to Content table
	rows, err := tx.NamedQuery(`INSERT INTO content (title, text, url) VALUES (:title, :text, :url) RETURNING *`, &banner.Content)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if rows.Next() {
		if err = rows.StructScan(&content); err != nil {
			return nil, err
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// then insert new banner into banner table
	query := fmt.Sprintf(`INSERT INTO banner (feature_id,is_active, content_id) 
VALUES (:feature_id, :is_active, %v) 
RETURNING id, feature_id, is_active, created_at, updated_at`,
		content.ID)

	rows, err = tx.NamedQuery(query, &banner)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		if err = rows.StructScan(&banner); err != nil {
			return nil, err
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// for each tag id in entity.banner create new row (banner.ID, tag.ID) in BannerTag table
	for _, tag := range banner.TagIDs {
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO banner_tag (banner_id, tag_id) VALUES (%v, %v)", banner.ID, tag))
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	banner.Content = content

	return &banner, nil
}

func (r *Repo) UpdateBanner(ctx context.Context, id int, updateModel entity.Banner) error {
	tx, err := r.DB.BeginTxx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	// update some fields in banner table
	setQuery := fmt.Sprintf("is_active = %v, updated_at = now()", updateModel.IsActive)

	if updateModel.FeatureID != 0 {
		setQuery += fmt.Sprintf(", feature_id = %v", updateModel.FeatureID)
	}

	rows, err := tx.QueryxContext(
		ctx,
		fmt.Sprintf("UPDATE banner SET %v WHERE id = %v RETURNING content_id", setQuery, id),
	)
	if err != nil {
		return err
	}

	contentIdStruct := struct {
		ContentID int `db:"content_id"`
	}{}

	// TODO: change this method for optimization

	// fetch content id
	if rows.Next() {
		if err = rows.StructScan(&contentIdStruct); err != nil {
			return err
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return err
	}

	// update content associated with this banner
	_, err = tx.QueryxContext(
		ctx,
		`UPDATE content
SET title = COALESCE($1, title),
    text = COALESCE($2, text),
    url = COALESCE($3, text)
WHERE content_id = $4;`,
		updateModel.Content.Title, updateModel.Content.Text, updateModel.Content.Url, contentIdStruct.ContentID,
	)
	if err != nil {
		return err
	}

	/* update tag ids in banner_tag table:
	to do this we need firstly delete all rows from banner_tag where banner_id = id,
	then add new rows in this table of form (banner_id = id, tag_id = updateModel.tagIds[i])
	*/
	if updateModel.TagIDs != nil {
		_, err = tx.QueryxContext(
			ctx,
			"DELETE FROM banner_tag WHERE banner_id = $1",
			id,
		)
		if err != nil {
			return err
		}

		for _, tag := range updateModel.TagIDs {
			_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO banner_tag (banner_id, tag_id) VALUES (%v, %v)", id, tag))
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *Repo) DeleteBanner(ctx context.Context, id int) (*entity.Banner, error) {
	row, err := r.DB.QueryxContext(ctx, "DELETE FROM banner WHERE id = $1 RETURNING *", id)
	if err != nil {
		return nil, err
	}

	var banner entity.Banner

	if row.Next() {
		if err = row.StructScan(&banner); err != nil {
			return nil, err
		}
	}

	return &banner, nil
}
