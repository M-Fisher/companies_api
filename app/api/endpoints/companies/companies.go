package companies

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/api/base"
	"github.com/M-Fisher/companies_api/app/internal/logger"
	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/services/auth"
)

// CompaniesAPI - common API struct for companies
type CompaniesAPI struct {
	base.API
}

// SetRoutes initial routing
func (a *CompaniesAPI) SetRoutes() {
	a.API.SetJSONHandler("", a.GetCompanies).Methods("GET")
	a.API.SetJSONHandler("/{id}", a.GetCompany).Methods("GET")
	a.API.SetJSONHandler("/{id}", a.UpdateCompany).Methods("PUT")
	a.API.SetJSONHandler("/{id}", a.DeleteCompany).Methods("DELETE")
	a.API.SetJSONHandler("", a.CreateCompany).Methods("POST")
}

func (a *CompaniesAPI) GetCompanies(rw http.ResponseWriter, r *http.Request) (any, error) {
	log := logger.FromContext(r.Context()).With(zap.String("method", "GetCompanies"))
	var data models.GetCompanyRequest
	err := base.DecodeQuery(&data, r)
	if err != nil {
		log.Error("Failed to decode query params", zap.Error(err))
		return base.Response{}, errors.New(`incorrect params`)
	}

	cmp, err := a.Srv.CompaniesService.GetCompanies(r.Context(), models.Company{
		Name:    data.Name,
		Code:    data.Code,
		Country: data.Country,
		Website: data.Website,
		Phone:   data.Phone,
	})
	if err != nil {
		log.Error("Failed to get companies", zap.Error(err))
		return base.Response{
			"err": err.Error(),
		}, errors.New(`global error`)
	}
	return base.Response{
		"companies": cmp,
	}, nil
}

func (a *CompaniesAPI) GetCompany(rw http.ResponseWriter, r *http.Request) (any, error) {
	log := logger.FromContext(r.Context()).With(zap.String("method", "GetCompanies"))
	compID, err := base.GetVarInt(r, `id`)
	if err != nil {
		log.Error("Failed to parse company id from request url", zap.Error(err))
		return base.Response{}, errors.New(`company id required`)
	}

	cmp, err := a.Srv.CompaniesService.GetCompany(r.Context(), uint64(compID))
	if err != nil {
		log.Error("Failed to get company", zap.Error(err))
		return base.Response{
			"err": err.Error(),
		}, errors.New(`global error`)
	}
	if cmp == nil {
		return nil, ErrCompanyNotFound
	}
	return base.Response{
		"company": cmp,
	}, nil
}

func (a *CompaniesAPI) UpdateCompany(rw http.ResponseWriter, r *http.Request) (any, error) {
	log := logger.FromContext(r.Context()).With(zap.String("method", "UpdateCompany"))
	compID, err := base.GetVarInt(r, `id`)
	if err != nil {
		log.Error("Failed to parse company id from request url", zap.Error(err))
		return base.Response{}, errors.New(`company id required`)
	}

	var data models.Company
	err = base.DecodeBody(&data, r)
	if err != nil {
		log.Error("Failed to parse request body", zap.Error(err))
		return base.Response{}, errors.New(`incorrect params`)
	}

	newComp, err := a.Srv.CompaniesService.UpdateCompany(r.Context(), uint64(compID), data)
	if err != nil {
		log.Error("Failed to update company", zap.Error(err))
		return base.Response{
			"err": err.Error(),
		}, err
	}
	return base.Response{
		"company": newComp,
	}, nil
}

func (a *CompaniesAPI) DeleteCompany(rw http.ResponseWriter, r *http.Request) (any, error) {
	log := logger.FromContext(r.Context()).With(zap.String("method", "DeleteCompany"))
	_, err := a.VerifyUser(r)
	if err != nil {
		log.Error("Failed to verify user")
		return nil, err
	}
	stat, err := a.Srv.AuthService.IsActionAllowed(auth.ActionCompanyDelete, r.RemoteAddr)
	if err != nil {
		log.Error("Failed to verify user location", zap.Error(err))
	}
	if !a.Srv.Config.DevMode && !stat {
		log.Info("Action is not allowed due to region restrictions", zap.String("remote_addr", r.RemoteAddr))
		return nil, base.ErrUnauthorized
	}

	compID, err := base.GetVarInt(r, `id`)
	if err != nil {
		log.Error("Failed to parse company id from request", zap.Error(err))
		return base.Response{}, errors.New(`company id required`)
	}

	err = a.Srv.CompaniesService.DeleteCompany(r.Context(), uint64(compID))
	if err != nil {
		log.Error("Failed to delete company", zap.Error(err))
		return base.Response{
			"err": err.Error(),
		}, err
	}
	return base.Response{}, nil
}

func (a *CompaniesAPI) CreateCompany(rw http.ResponseWriter, r *http.Request) (any, error) {
	log := logger.FromContext(r.Context()).With(zap.String("method", "CreateCompany"))
	_, err := a.VerifyUser(r)
	if err != nil {
		log.Error("Failed to verify user", zap.Error(err))
		return nil, err
	}
	stat, err := a.Srv.AuthService.IsActionAllowed(auth.ActionCompanyCreate, r.RemoteAddr)
	if err != nil {
		log.Error("Failed to verify user location", zap.String("remote_addr", r.RemoteAddr), zap.Error(err))
	}
	if !a.Srv.Config.DevMode && !stat {
		log.Info("Action is not allowed due to region restrictions", zap.String("remote_addr", r.RemoteAddr))
		return nil, base.ErrUnauthorized
	}
	var userData models.Company
	err = base.DecodeBody(&userData, r)
	if err != nil {
		log.Error("Failed to decode request body", zap.Error(err))
		return base.Response{}, errors.New(`incorrect params`)
	}

	compID, err := a.Srv.CompaniesService.CreateCompany(r.Context(), userData)
	if err != nil {
		log.Error("Failed to create company", zap.Error(err))
		return base.Response{
			"err": err.Error(),
		}, errors.New("server error")
	}

	return base.Response{
		"company_id": compID,
	}, nil
}
