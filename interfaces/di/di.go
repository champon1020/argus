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
	aP     repository.ArticleRepository
	aU     usecase.ArticleUseCase
	tP     repository.TagRepository
	tU     usecase.TagUseCase
	h      handler.AppHandler
}

// NewDI creates DI instance.
func NewDI(config *config.Config) DI {
	d := &di{config: config}
	d.aP = persistence.NewArticlePersistence()
	d.tP = persistence.NewTagPersistence()

	d.aU = usecase.NewArticleUseCase(d.aP, d.tP)
	d.tU = usecase.NewTagUseCase(d.tP)

	d.h = handler.NewAppHandler(d.aU, d.tU, d.config)
	return d
}

func (d di) NewArticleRepository() repository.ArticleRepository {
	return d.aP
}

func (d di) NewArticleUseCase() usecase.ArticleUseCase {
	return d.aU
}

func (d di) NewAppHandler() handler.AppHandler {
	return d.h
}
