package model

type BookInfo struct {
	Id            int
	Title         string
	Description   *string
	AuthorName    string
	ReleaseYear   *int
	Edition       *string
	PublisherName string
}
