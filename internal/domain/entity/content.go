package entity

type Content struct {
	ID    int    `db:"content_id"`
	Title string `db:"title"`
	Text  string `db:"text"`
	Url   string `db:"url"`
}
