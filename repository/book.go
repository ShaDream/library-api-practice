package repository

import (
	"database/sql"
	"fmt"
	"github.com/ShaDream/library-api-practice/model"
	error2 "github.com/ShaDream/library-api-practice/model/error"
	"gorm.io/gorm"
)

type BookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepositoryImpl(db *gorm.DB) *BookRepositoryImpl {
	return &BookRepositoryImpl{db: db}
}

const getBooksRequest = `SELECT b.id, b.title, b.description, b.release_year, b.edition, a.name, p.name
FROM book AS b
         INNER JOIN author a on a.id = b.author_id
         INNER JOIN publisher p on p.id = b.publisher_id
ORDER BY b.id ASC
OFFSET @offset LIMIT @size`

const getBookRequest = `SELECT b.id, b.title, b.description, b.release_year, b.edition, a.name, p.name
FROM book AS b
         INNER JOIN author a on a.id = b.author_id
         INNER JOIN publisher p on p.id = b.publisher_id
WHERE b.id = ?`

const getPhysicalBooksRequest = `SELECT p.id,
       s.section,
       s.shelf,
       i.issue_date,
       i.complete_date,
       i.reader_id,
       r.name,
       i.issue_date IS NULL AS in_library
FROM physical_book AS p
         INNER JOIN section s on s.id = p.section_id
         LEFT JOIN issuance i on p.id = i.physical_book_id AND i.actual_complete_date IS NULL
         LEFT JOIN reader r on i.reader_id = r.id
WHERE book_id = ?;`

const IsExistsBookRequest = `SELECT EXISTS(SELECT * FROM book WHERE id = ?)`

func (b BookRepositoryImpl) GetBooksFromDB(page int, size int) ([]model.BookInfo, error) {
	rows, err := b.db.Raw(getBooksRequest, sql.Named("offset", (page-1)*size), sql.Named("size", size)).Rows()
	if err != nil {
		return nil, error2.NewInternalServerError(err)
	}
	books := make([]model.BookInfo, 0)
	for rows.Next() {
		var book model.BookInfo
		err = rows.Scan(
			&book.Id,
			&book.Title,
			&book.Description,
			&book.ReleaseYear,
			&book.Edition,
			&book.AuthorName,
			&book.PublisherName,
		)
		if err != nil {
			return nil, error2.NewInternalServerError(err)
		}
		books = append(books, book)
	}
	return books, nil
}

func (b BookRepositoryImpl) GetBookFromDB(id int) (*model.BookInfo, error) {
	exist, err := b.IsBookExistFromDB(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, error2.NewNotFoundError(fmt.Errorf("book with id %d not found", id))
	}

	row := b.db.Raw(getBookRequest, id).Row()
	if row.Err() != nil {
		return nil, error2.NewInternalServerError(row.Err())
	}
	var book model.BookInfo
	err = row.Scan(
		&book.Id,
		&book.Title,
		&book.Description,
		&book.ReleaseYear,
		&book.Edition,
		&book.AuthorName,
		&book.PublisherName,
	)
	if err != nil {
		return nil, error2.NewInternalServerError(err)
	}
	return &book, nil
}

func (b BookRepositoryImpl) GetPhysicalBooksFromDB(id int) ([]model.PhysicalBookInfo, error) {
	exist, err := b.IsBookExistFromDB(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, error2.NewNotFoundError(fmt.Errorf("book with id %d not found", id))
	}

	rows, err := b.db.Raw(getPhysicalBooksRequest, id).Rows()
	if err != nil {
		return nil, error2.NewInternalServerError(err)
	}

	physicalBooks := make([]model.PhysicalBookInfo, 0)
	for rows.Next() {
		var book model.PhysicalBookInfo
		err = rows.Scan(
			&book.Id,
			&book.Section.Section,
			&book.Section.Shelf,
			&book.IssueDate,
			&book.CompleteDate,
			&book.ReaderId,
			&book.ReaderName,
			&book.InLibrary,
		)
		if err != nil {
			return nil, error2.NewInternalServerError(err)
		}

		physicalBooks = append(physicalBooks, book)
	}
	return physicalBooks, nil
}

func (b BookRepositoryImpl) IsBookExistFromDB(id int) (bool, error) {
	row := b.db.Raw(IsExistsBookRequest, id).Row()
	if row.Err() != nil {
		return false, error2.NewInternalServerError(row.Err())
	}
	var exist bool
	err := row.Scan(&exist)
	if err != nil {
		return false, error2.NewInternalServerError(row.Err())
	}
	return exist, nil
}
