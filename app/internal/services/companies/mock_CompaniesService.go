// Code generated by mockery v2.10.6. DO NOT EDIT.

package companies

import (
	context "context"

	models "github.com/M-Fisher/companies_api/app/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockCompaniesService is an autogenerated mock type for the CompaniesService type
type MockCompaniesService struct {
	mock.Mock
}

// CreateCompany provides a mock function with given fields: ctx, company
func (_m *MockCompaniesService) CreateCompany(ctx context.Context, company models.Company) (uint64, error) {
	ret := _m.Called(ctx, company)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, models.Company) uint64); ok {
		r0 = rf(ctx, company)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Company) error); ok {
		r1 = rf(ctx, company)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCompany provides a mock function with given fields: ctx, companyID
func (_m *MockCompaniesService) DeleteCompany(ctx context.Context, companyID uint64) error {
	ret := _m.Called(ctx, companyID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, companyID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCompanies provides a mock function with given fields: ctx, params
func (_m *MockCompaniesService) GetCompanies(ctx context.Context, params models.Company) ([]*models.Company, error) {
	ret := _m.Called(ctx, params)

	var r0 []*models.Company
	if rf, ok := ret.Get(0).(func(context.Context, models.Company) []*models.Company); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Company)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Company) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCompany provides a mock function with given fields: ctx, compID, company
func (_m *MockCompaniesService) UpdateCompany(ctx context.Context, compID uint64, company models.Company) (*models.Company, error) {
	ret := _m.Called(ctx, compID, company)

	var r0 *models.Company
	if rf, ok := ret.Get(0).(func(context.Context, uint64, models.Company) *models.Company); ok {
		r0 = rf(ctx, compID, company)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Company)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, models.Company) error); ok {
		r1 = rf(ctx, compID, company)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}