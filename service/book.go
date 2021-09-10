package service

import (
	"github.com/ShaDream/library-api-practice/model"
	"github.com/ShaDream/library-api-practice/repository"
)

type BookServiceImpl struct {
	rep repository.BookRepository
}

func NewBookServiceImpl(service repository.BookRepository) *BookServiceImpl {
	return &BookServiceImpl{rep: service}
}

func (b BookServiceImpl) GetBooks(page int, size int) ([]model.BookInfo, error) {
	return b.rep.GetBooksFromDB(page, size)
}

func (b BookServiceImpl) GetBook(id int) (*model.BookInfo, error) {
	return b.rep.GetBookFromDB(id)
}

func (b BookServiceImpl) GetPhysicalBooks(id int) ([]model.PhysicalBookInfo, error) {
	return b.rep.GetPhysicalBooksFromDB(id)
}
