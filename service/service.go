package service

import (
	"github.com/ShaDream/library-api-practice/model"
	"github.com/ShaDream/library-api-practice/repository"
)

type BookService interface {
	GetBooks(page int, size int) ([]model.BookInfo, error)
	GetBook(id int) (*model.BookInfo, error)
	GetPhysicalBooks(id int) ([]model.PhysicalBookInfo, error)
}

type Service struct {
	BookService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{BookService: NewBookServiceImpl(repository)}
}
