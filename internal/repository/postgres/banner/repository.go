package banner

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	sliceutils "avito-backend-trainee-2024/pkg/utils/slice"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"time"

	stringutils "avito-backend-trainee-2024/pkg/utils/string"
)

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func setQueryForUpdateModel(updateModel entity.Banner) string {
	setQuery := ""
	isCommaNeeded := false

	if updateModel.Content.Title != "" {
		setQuery += fmt.Sprintf("title = '%v'", updateModel.Content.Title)
		isCommaNeeded = true
	}

	if updateModel.Content.Text != "" {
		if isCommaNeeded {
			setQuery += ", "
		}

		setQuery += fmt.Sprintf("text = '%v'", updateModel.Content.Text)
		isCommaNeeded = true
	}

	if updateModel.Content.Url != "" {
		if isCommaNeeded {
			setQuery += ", "
		}

		setQuery += fmt.Sprintf("url = '%v'", updateModel.Content.Url)
		isCommaNeeded = true
	}

	return setQuery
}

func (r *Repo) GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error) {
	query := `SELECT banner.id,
       feature_id,
       is_active,
       created_at,
       updated_at,
       title,
       text,
       url,
       array_agg(bt.tag_id ORDER BY bt.tag_id) AS tag_ids
FROM banner
         JOIN public.content c ON c.content_id = banner.content_id
         JOIN public.banner_tag bt ON banner.id = bt.banner_id
GROUP BY c.content_id, banner.id, feature_id
ORDER BY feature_id`

	if limit == math.MaxInt64 {
		query = fmt.Sprintf(`%v OFFSET %v`, query, offset)
	} else {
		query = fmt.Sprintf(`%v LIMIT %v OFFSET %v`, query, limit, offset)
	}

	type Row struct {
		ID        int       `db:"id"`
		FeatureID int       `db:"feature_id"`
		TagIDsStr string    `db:"tag_ids"`
		IsActive  bool      `db:"is_active"`
		Title     string    `db:"title"`
		Text      string    `db:"text"`
		Url       string    `db:"url"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`

		TagIDsInt []int
	}

	rows, err := r.DB.Queryx(query)
	if err != nil {
		return nil, err
	}

	var banners []*entity.Banner

	for rows.Next() {
		var row Row

		if err = rows.StructScan(&row); err != nil {
			return nil, err
		}

		// row.TagIDsStr have structure {1,2,...}
		row.TagIDsInt, err = stringutils.FillIntSliceFromString(row.TagIDsStr[1 : len(row.TagIDsStr)-1])
		if err != nil {
			return nil, err
		}

		content := entity.Content{
			Title: row.Title,
			Text:  row.Text,
			Url:   row.Url,
		}

		banner := entity.Banner{
			ID:        row.ID,
			TagIDs:    row.TagIDsInt,
			FeatureID: row.FeatureID,
			Content:   content,
			IsActive:  row.IsActive,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}

		banners = append(banners, &banner)
	}

	return banners, nil
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

		// row.TagIDsStr have structure {1,2,...}
		row.TagIDsInt, err = stringutils.FillIntSliceFromString(row.TagIDsStr[1 : len(row.TagIDsStr)-1])
		if err != nil {
			return nil, err
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
		TagIDsStr string `db:"tag_ids"`
		TagIDsInt []int
	}

	var rows []*Row

	for dbRows.Next() {
		var row Row

		if err = dbRows.StructScan(&row); err != nil {
			return nil, err
		}

		// row.TagIDs have structure {1,2,...}
		row.TagIDsInt, err = stringutils.FillIntSliceFromString(row.TagIDsStr[1 : len(row.TagIDsStr)-1])
		if err != nil {
			return nil, err
		}

		rows = append(rows, &row)
	}

	// each row represents banner with banner.feature_id = featureID => find banner with banner.tag_ids = tagIDs
	for _, row := range rows {
		if sliceutils.Equals(row.TagIDsInt, tagIDs) { // here tagIDs gotta be sorted by asc, row.TagIDs already sorted
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

	return nil, nil
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
	setQuery = setQueryForUpdateModel(updateModel)

	// execute query only if updating something
	if setQuery != "" {
		_, err = tx.QueryxContext(ctx, fmt.Sprintf(`UPDATE content SET %v WHERE content_id = %v`, setQuery, contentIdStruct.ContentID))
		if err != nil {
			return err
		}
	}

	/* update tag ids in banner_tag table:
	to do this we need firstly delete all rows from banner_tag where banner_id = id,
	then add new rows in this table of form (banner_id = id, tag_id = updateModel.tagIds[i])
	*/
	if len(updateModel.TagIDs) != 0 {
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
