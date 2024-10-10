package request

import (
	"github.com/Rafael24595/go-api-core/src/domain"
)

type IRepositoryCommand interface {
	Insert(request domain.Request) *domain.Request
	Delete(request domain.Request) *domain.Request
}