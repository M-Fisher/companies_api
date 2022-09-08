package companies

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/api/base"
	"github.com/M-Fisher/companies_api/app/config"
	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/server"
	"github.com/M-Fisher/companies_api/app/internal/services/auth"
	"github.com/M-Fisher/companies_api/app/internal/services/companies"
)

type CompaniesTestsSuite struct {
	suite.Suite
}

func TestCompaniesAPI(t *testing.T) {
	suite.Run(t, new(CompaniesTestsSuite))
}

func (suite *CompaniesTestsSuite) TestGetCompaniesEmptyResponse() {
	expResp := base.Response{
		"companies": []*models.Company{},
	}
	req, err := http.NewRequest("GET", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	compmocks := new(companies.MockCompaniesService)
	compmocks.On("GetCompanies", mock.Anything, models.Company{}).Return([]*models.Company{}, nil)
	srv := server.Server{
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.GetCompanies(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.NoError(gotErr)
}

func (suite *CompaniesTestsSuite) TestGetCompaniesInternalError() {
	expResp := base.Response{
		"err": `some error`,
	}
	req, err := http.NewRequest("GET", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	compmocks := new(companies.MockCompaniesService)
	compmocks.On("GetCompanies", mock.Anything, models.Company{}).Return([]*models.Company{}, errors.New("some error"))
	srv := server.Server{
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.GetCompanies(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.Error(gotErr)
}

func (suite *CompaniesTestsSuite) TestGetCompaniesOk() {
	expResp := base.Response{
		"companies": []*models.Company{
			{
				Name:  "test",
				Code:  "TST",
				Phone: "1234",
			},
		},
	}
	req, err := http.NewRequest("GET", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	compmocks := new(companies.MockCompaniesService)
	compmocks.On("GetCompanies", mock.Anything, models.Company{}).Return([]*models.Company{
		{
			Name:  "test",
			Code:  "TST",
			Phone: "1234",
		},
	}, nil)
	srv := server.Server{
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.GetCompanies(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.NoError(gotErr)
}

func (suite *CompaniesTestsSuite) TestCreateCompanyUnauthorized() {
	req, err := http.NewRequest("POST", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}

	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		Log: zap.NewExample(),
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.CreateCompany(context.Background(), res, req)

	suite.Nil(resp)
	suite.Equal(errors.New(`not authorized`), gotErr)
}

func (suite *CompaniesTestsSuite) TestCreateCompanyInvalidIP() {
	req, err := http.NewRequest("POST", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("CreateCompany", mock.Anything, models.Company{}).Return([]*models.Company{}, nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`US`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.CreateCompany(context.Background(), res, req)

	suite.Nil(resp)
	suite.Equal(errors.New(`not authorized`), gotErr)
}

func (suite *CompaniesTestsSuite) TestCreateCompanyRegionCheckingFailed() {
	req, err := http.NewRequest("POST", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("CreateCompany", mock.Anything, models.Company{}).Return([]*models.Company{}, nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(``, errors.New(`service unavailable`))
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.CreateCompany(context.Background(), res, req)

	suite.Nil(resp)
	suite.Equal(errors.New(`not authorized`), gotErr)
}

func (suite *CompaniesTestsSuite) TestCreateCompanyInternalError() {
	expResp := base.Response{
		"err": "some error",
	}
	req, err := http.NewRequest("POST", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("CreateCompany", mock.Anything, models.Company{}).Return(uint64(0), errors.New("some error"))

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.CreateCompany(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.Equal(errors.New("server error"), gotErr)
}

func (suite *CompaniesTestsSuite) TestCreateCompanyParamsError() {
	reqBody := `{"name":"Test1","code":"TST1",}`
	req, err := http.NewRequest("POST", "api/companies", strings.NewReader(reqBody))
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService: auth.NewService(clmock),
		Log:         zap.NewExample(),
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.CreateCompany(context.Background(), res, req)

	suite.Equal(base.Response{}, resp)
	suite.Equal(errors.New("incorrect params"), gotErr)
}

func (suite *CompaniesTestsSuite) TestCreateCompanyOk() {
	expResp := base.Response{
		"company_id": uint64(1),
	}
	req, err := http.NewRequest("GET", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("CreateCompany", mock.Anything, models.Company{}).Return(uint64(1), nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.CreateCompany(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.NoError(gotErr)
}

func (suite *CompaniesTestsSuite) TestDeleteCompanyUnauthorized() {
	req, err := http.NewRequest("DELETE", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}

	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		Log: zap.NewExample(),
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.DeleteCompany(context.Background(), res, req)
	suite.Nil(resp)
	suite.Equal(errors.New(`not authorized`), gotErr)
}

func (suite *CompaniesTestsSuite) TestDeleteCompanyInvalidIP() {
	req, err := http.NewRequest("DELETE", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("CreateCompany", mock.Anything, models.Company{}).Return([]*models.Company{}, nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`US`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.DeleteCompany(context.Background(), res, req)
	suite.Nil(resp)
	suite.Equal(errors.New(`not authorized`), gotErr)
}

func (suite *CompaniesTestsSuite) TestDeleteCompanyRegionCheckingFailed() {
	req, err := http.NewRequest("DELETE", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("DeleteCompany", mock.Anything, models.Company{}).Return([]*models.Company{}, nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(``, errors.New(`service unavailable`))
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.DeleteCompany(context.Background(), res, req)
	suite.Nil(resp)
	suite.Equal(errors.New(`not authorized`), gotErr)
}

func (suite *CompaniesTestsSuite) TestDeleteCompanyEmptyID() {
	req, err := http.NewRequest("DELETE", "api/companies/", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("DeleteCompany", mock.Anything, models.Company{}).Return([]*models.Company{}, nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.DeleteCompany(context.Background(), res, req)
	suite.Equal(base.Response{}, resp)
	suite.Equal(errors.New(`company id required`), gotErr)
}

func (suite *CompaniesTestsSuite) TestDeleteCompanyInternalError() {
	expResp := base.Response{
		"err": `some error`,
	}
	req, err := http.NewRequest("DELETE", "api/companies/12", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`
	req = mux.SetURLVars(req, map[string]string{
		"id": "12",
	})
	compmocks := new(companies.MockCompaniesService)
	compmocks.On("DeleteCompany", mock.Anything, uint64(12)).Return(errors.New(`some error`))

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Router: srv.Router,
			Srv:    &srv,
		},
	}

	res := httptest.NewRecorder()

	resp, gotErr := a.DeleteCompany(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.Equal(errors.New(`some error`), gotErr)

}

func (suite *CompaniesTestsSuite) TestDeleteCompanyOk() {
	req, err := http.NewRequest("DELETE", "api/companies/12", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`
	req = mux.SetURLVars(req, map[string]string{
		"id": "12",
	})
	compmocks := new(companies.MockCompaniesService)
	compmocks.On("DeleteCompany", mock.Anything, uint64(12)).Return(nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Router: srv.Router,
			Srv:    &srv,
		},
	}

	res := httptest.NewRecorder()

	resp, gotErr := a.DeleteCompany(context.Background(), res, req)

	suite.Equal(base.Response{}, resp)
	suite.NoError(gotErr)

}

func (suite *CompaniesTestsSuite) TestUpdateCompanyIncorrectParams() {
	req, err := http.NewRequest("GET", "api/companies", nil)
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("UpdateCompany", mock.Anything, models.Company{}).Return(nil, errors.New("some error"))

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.UpdateCompany(context.Background(), res, req)

	suite.Equal(base.Response{}, resp)
	suite.Equal(errors.New("company id required"), gotErr)
}

func (suite *CompaniesTestsSuite) TestUpdateCompanyInternalError() {
	expResp := base.Response{
		"err": "some error",
	}
	reqBody := `{"name":"Test1","code":"TST1"}`
	req, err := http.NewRequest("GET", "api/companies", strings.NewReader(reqBody))
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`
	req = mux.SetURLVars(req, map[string]string{
		"id": "12",
	})

	compmocks := new(companies.MockCompaniesService)
	compmocks.On(
		"UpdateCompany",
		mock.Anything,
		uint64(12),
		models.Company{Name: "Test1", Code: "TST1"},
	).Return(nil, errors.New("some error"))

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.UpdateCompany(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.Equal(errors.New("some error"), gotErr)
}

func (suite *CompaniesTestsSuite) TestUpdateCompanyParamsError() {
	reqBody := `{"name":"Test1","code":"TST1",}`
	req, err := http.NewRequest("PUT", "api/companies/1", strings.NewReader(reqBody))
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService: auth.NewService(clmock),
		Log:         zap.NewExample(),
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.UpdateCompany(context.Background(), res, req)

	suite.Equal(base.Response{}, resp)
	suite.Equal(errors.New("incorrect params"), gotErr)
}

func (suite *CompaniesTestsSuite) TestUpdateCompanyOk() {
	expResp := base.Response{
		"company": &models.Company{
			ID:   1,
			Name: "Updated",
		},
	}
	reqBody := `{"name":"Updated"}`
	req, err := http.NewRequest("PUT", "api/companies/1", strings.NewReader(reqBody))
	if err != nil {
		suite.FailNow(err.Error())
	}
	req.Header.Add(`Authorization`, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`)
	req.RemoteAddr = `192.168.1.1`
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	compmocks := new(companies.MockCompaniesService)
	compmocks.On("UpdateCompany", mock.Anything, uint64(1), models.Company{Name: "Updated"}).Return(&models.Company{
		ID:   uint64(1),
		Name: "Updated",
	}, nil)

	clmock := new(auth.MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.168.1.1`).Return(`CY`, nil)
	srv := server.Server{
		Config: &config.Config{
			JWTSecret: `test`,
		},
		AuthService:      auth.NewService(clmock),
		Log:              zap.NewExample(),
		CompaniesService: compmocks,
	}

	a := CompaniesAPI{
		API: base.API{
			Srv: &srv,
		},
	}
	res := httptest.NewRecorder()
	resp, gotErr := a.UpdateCompany(context.Background(), res, req)

	suite.Equal(expResp, resp)
	suite.NoError(gotErr)
}
