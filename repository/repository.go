package repository

import (
	"github.com/ShaDream/library-api-practice/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	GetBooksFromDB(page int, size int) ([]model.BookInfo, error)
	GetBookFromDB(id int) (*model.BookInfo, error)
	GetPhysicalBooksFromDB(id int) ([]model.PhysicalBookInfo, error)
	IsBookExistFromDB(id int) (bool, error)
}

type Repository struct {
	db *gorm.DB
	BookRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:             db,
		BookRepository: NewBookRepositoryImpl(db),
	}
}
