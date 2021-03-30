package di

import (
	"github.com/champon1020/argus/domain/repository"
	"github.com/champon1020/argus/infra/persistence"
	"github.com/champon1020/argus/interfaces/handler"
	"github.com/champon1020/argus/usecase"
)

type DI interface {
	NewArticleRepository() repository.ArticleRepository
	NewArticleUseCase() usecase.ArticleUseCase
	NewAppHandler() handler.AppHandler
}

type di struct {
	ap repository.ArticleRepository
	au usecase.ArticleUseCase
	h  handler.AppHandler
}

func NewDI() DI {
	d := &di{}
	d.ap = persistence.NewArticlePersistence()
	d.au = usecase.NewArticleUseCase(d.ap)
	d.h = handler.NewAppHandler(d.au)
	return d
}

func (d di) NewArticleRepository() repository.ArticleRepository {
	return d.ap
}

func (d di) NewArticleUseCase() usecase.ArticleUseCase {
	return d.au
}

func (d di) NewAppHandler() handler.AppHandler {
	return d.h
}
