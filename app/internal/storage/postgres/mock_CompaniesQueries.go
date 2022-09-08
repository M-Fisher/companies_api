// Code generated by mockery v2.10.6. DO NOT EDIT.

package postgres

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCompaniesQueries is an autogenerated mock type for the CompaniesQueries type
type MockCompaniesQueries struct {
	mock.Mock
}

// CreateCompany provides a mock function with given fields: ctx, params
func (_m *MockCompaniesQueries) CreateCompany(ctx context.Context, params Company) (uint64, error) {
	ret := _m.Called(ctx, params)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, Company) uint64); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, Company) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCompany provides a mock function with given fields: ctx, compID
func (_m *MockCompaniesQueries) DeleteCompany(ctx context.Context, compID uint64) error {
	ret := _m.Called(ctx, compID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, compID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCompanies provides a mock function with given fields: ctx, params
func (_m *MockCompaniesQueries) GetCompanies(ctx context.Context, params Company) ([]*Company, error) {
	ret := _m.Called(ctx, params)

	var r0 []*Company
	if rf, ok := ret.Get(0).(func(context.Context, Company) []*Company); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Company)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, Company) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCompanyByID provides a mock function with given fields: ctx, compID
func (_m *MockCompaniesQueries) GetCompanyByID(ctx context.Context, compID uint64) (*Company, error) {
	ret := _m.Called(ctx, compID)

	var r0 *Company
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *Company); ok {
		r0 = rf(ctx, compID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Company)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, compID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCompany provides a mock function with given fields: ctx, compID, data
func (_m *MockCompaniesQueries) UpdateCompany(ctx context.Context, compID uint64, data Company) (uint64, error) {
	ret := _m.Called(ctx, compID, data)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, uint64, Company) uint64); ok {
		r0 = rf(ctx, compID, data)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, Company) error); ok {
		r1 = rf(ctx, compID, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
