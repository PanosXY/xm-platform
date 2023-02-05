package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/PanosXY/xm-platform/api/v1/company"
	"github.com/PanosXY/xm-platform/api/v1/misc"
	"github.com/PanosXY/xm-platform/config"
	companyRepoPG "github.com/PanosXY/xm-platform/store/postgres/company"
	miscRepoPG "github.com/PanosXY/xm-platform/store/postgres/misc"
	"github.com/PanosXY/xm-platform/utils/logger"
	"github.com/PanosXY/xm-platform/utils/postgres"
)

func NewRouter(configuration *config.Configuration, log *logger.Logger, dbClient *postgres.PostgresClient) *chi.Mux {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	miscHandler := misc.NewMiscHandler(
		configuration,
		log,
		misc.NewMiscService(miscRepoPG.NewPostgresMiscRepository(dbClient)),
	)

	companyHandler := company.NewCompanyHandler(
		log,
		company.NewCompanyuService(companyRepoPG.NewPostgresCompanyRepository(dbClient)),
	)

	router.Get("/v1/health", miscHandler.HealthCheck)
	router.Post("/v1/login", miscHandler.Login)

	router.Get("/v1/company/{id}", companyHandler.GetCompanyByID)

	router.Route("/v1", func(r chi.Router) {
		r.Use(miscHandler.JWTAuthToken)
		r.Post("/company", companyHandler.CreateCompany)
		r.Delete("/company/{id}", companyHandler.DeleteCompanyByID)
		r.Patch("/company/{id}", companyHandler.PatchCompanyByID)
	})

	return router
}
