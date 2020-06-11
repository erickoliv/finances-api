// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import categories "github.com/erickoliv/finances-api/categories"
import context "context"
import mock "github.com/stretchr/testify/mock"
import rest "github.com/erickoliv/finances-api/pkg/http/rest"
import uuid "github.com/google/uuid"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, pk, owner
func (_m *Repository) Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error {
	ret := _m.Called(ctx, pk, owner)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, pk, owner)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, pk, owner
func (_m *Repository) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*categories.Category, error) {
	ret := _m.Called(ctx, pk, owner)

	var r0 *categories.Category
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *categories.Category); ok {
		r0 = rf(ctx, pk, owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*categories.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, pk, owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx, filters
func (_m *Repository) Query(ctx context.Context, filters *rest.Query) ([]categories.Category, error) {
	ret := _m.Called(ctx, filters)

	var r0 []categories.Category
	if rf, ok := ret.Get(0).(func(context.Context, *rest.Query) []categories.Category); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]categories.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rest.Query) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, row
func (_m *Repository) Save(ctx context.Context, row *categories.Category) error {
	ret := _m.Called(ctx, row)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *categories.Category) error); ok {
		r0 = rf(ctx, row)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
