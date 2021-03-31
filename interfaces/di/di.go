package di

import (
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/domain/repository"
	"github.com/champon1020/argus/infra/persistence"
	"github.com/champon1020/argus/interfaces/handler"
	"github.com/champon1020/argus/usecase"
)

// DI is struct for dependency injection.
type DI interface {
	NewArticleRepository() repository.ArticleRepository
	NewArticleUseCase() usecase.ArticleUseCase
	NewAppHandler() handler.AppHandler
}

type di struct {
	config *config.Config
	ap     repository.ArticleRepository
	au     usecase.ArticleUseCase
	h      handler.AppHandler
}

// NewDI creates DI instance.
func NewDI(config *config.Config) DI {
	d := &di{config: config}
	d.ap = persistence.NewArticlePersistence()
	d.au = usecase.NewArticleUseCase(d.ap)
	d.h = handler.NewAppHandler(d.au, d.config)
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
