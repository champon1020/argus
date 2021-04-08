package di

import (
	"github.com/champon1020/argus"
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/domain/gcp"
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
	logger *argus.Logger
	aP     repository.ArticleRepository
	tP     repository.TagRepository
	iP     gcp.CloudStorage
	aU     usecase.ArticleUseCase
	tU     usecase.TagUseCase
	iU     usecase.ImageUseCase
	h      handler.AppHandler
}

// NewDI creates DI instance.
func NewDI(config *config.Config, logger *argus.Logger) DI {
	d := &di{config: config, logger: logger}
	d.aP = persistence.NewArticlePersistence()
	d.tP = persistence.NewTagPersistence()
	d.iP = persistence.NewImagePersistence()

	d.aU = usecase.NewArticleUseCase(d.aP, d.tP)
	d.tU = usecase.NewTagUseCase(d.tP)
	d.iU = usecase.NewImageUseCase(d.iP)

	d.h = handler.NewAppHandler(d.aU, d.tU, d.iU, d.config, d.logger)
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
