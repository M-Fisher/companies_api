package api

import (
	"github.com/gorilla/mux"

	"github.com/M-Fisher/companies_api/app/api/base"
	ca "github.com/M-Fisher/companies_api/app/api/endpoints/companies"
	"github.com/M-Fisher/companies_api/app/internal/server"
)

func SetAPIRouter(s *server.Server, p string) {
	r := s.CreateDefaultRouter()
	apiRouter := r.PathPrefix(p).Subrouter()
	setSubRouterByStruct(apiRouter, "/companies", &ca.CompaniesAPI{
		API: base.API{
			Srv: s,
		},
	})
}

func setSubRouterByStruct(r *mux.Router, p string, a base.APIInterface) {
	subRouter := r.PathPrefix(p).Subrouter()
	a.SetRouter(subRouter)
	a.SetRoutes()
}
