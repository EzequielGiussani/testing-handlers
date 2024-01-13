package handlers_test

import (
	"app/internal/auth"
	"app/internal/product/handlers"
	"app/internal/product/repository"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetProducts(t *testing.T) {

	t.Run("should return a list of products", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{
			1: {Name: "Product 1", Quantity: 10, CodeValue: "code1", IsPublished: true, Expiration: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC), Price: 100.0},
		}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Get()
		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusOK
		expectedHeader := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		expectedBody := `{"data":[{"id":1,"name":"Product 1","quantity":10,"code_value":"code1","is_published":true,"expiration":"2021-12-31","price":100}],"message":"products"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

}

func TestGetProductById(t *testing.T) {

	t.Run("should return the product by the specified", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{
			1: {Name: "Product 1", Quantity: 10, CodeValue: "code1", IsPublished: true, Expiration: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC), Price: 100.0},
		}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.GetByID()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req := httptest.NewRequest("GET", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusOK
		expectedHeader := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		expectedBody := `{"data": {"id":1,"name":"Product 1","quantity":10,"code_value":"code1","is_published":true,"expiration":"2021-12-31","price":100},"message":"product"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("should return bad request due to invalid id", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.GetByID()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "test1234")
		req := httptest.NewRequest("GET", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusBadRequest
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Bad Request","message":"Invalid id"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("should return product not found", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.GetByID()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req := httptest.NewRequest("GET", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusNotFound
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Not Found","message":"Product not found"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

}

func TestCreateProduct(t *testing.T) {

	t.Run("should create the product", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Create()
		body := bytes.NewBufferString(`{"name": "Product 1", "quantity": 10, "code_value": "code1", "is_published": true, "expiration": "2020-02-02", "price": 100.0}`)
		req := httptest.NewRequest("POST", "/products", body)
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)
		// Assert

		expectedCode := http.StatusCreated
		expectedHeader := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		expectedBody := `{"data":{"id":1,"name":"Product 1","quantity":10,"code_value":"code1","is_published":true,"expiration":"2020-02-02","price":100},"message":"product created"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("should fail due to invalid token", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("test_token")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Create()
		req := httptest.NewRequest("POST", "/products", nil)
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)
		// Assert

		expectedCode := http.StatusUnauthorized
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Unauthorized","message":"Unauthorized"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

}

func TestDeleteProductById(t *testing.T) {

	t.Run("should return the product by the specified", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{
			1: {Name: "Product 1", Quantity: 10, CodeValue: "code1", IsPublished: true, Expiration: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC), Price: 100.0},
		}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Delete()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusNoContent
		expectedHeader := http.Header{}

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		require.Equal(t, expectedCode, res.Code)
		require.Empty(t, body)
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("should return bad request due to invalid id", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Delete()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "test1234")
		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusBadRequest
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Bad Request","message":"Invalid id"}`

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header())
		require.JSONEq(t, expectedBody, res.Body.String())

	})

	t.Run("should return product not found", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Delete()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusNotFound
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Not Found","message":"Product not found"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("should fail due to invalid token", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("test_token")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Delete()
		req := httptest.NewRequest("DELETE", "/products", nil)
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)
		// Assert

		expectedCode := http.StatusUnauthorized
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Unauthorized","message":"Unauthorized"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

}

func TestPutProduct(t *testing.T) {

	t.Run("should return bad request due to invalid id", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.UpdateOrCreate()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "test1234")
		req := httptest.NewRequest("PUT", "/products/1", nil)
		res := httptest.NewRecorder()

		// Act
		hdFunc(res, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Bad Request","message":"Invalid id"}`

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header())
		require.JSONEq(t, expectedBody, res.Body.String())

	})

	t.Run("should fail due to invalid token", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("test_token")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.UpdateOrCreate()
		req := httptest.NewRequest("PUT", "/products/1", nil)
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)
		// Assert

		expectedCode := http.StatusUnauthorized
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Unauthorized","message":"Unauthorized"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

}

func TestPatchProduct(t *testing.T) {

	t.Run("should return bad request due to invalid id", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Update()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "test1234")
		req := httptest.NewRequest("PATCH", "/products/1", nil)
		res := httptest.NewRecorder()

		// Act
		hdFunc(res, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Bad Request","message":"Invalid id"}`

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header())
		require.JSONEq(t, expectedBody, res.Body.String())

	})

	t.Run("should return product not found", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Update()
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		body := bytes.NewBufferString(`{"name": "Product 1"}`)
		req := httptest.NewRequest("PATCH", "/products/1", body)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)

		// Assert

		expectedCode := http.StatusNotFound
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Not Found","message":"Product not found"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("should fail due to invalid token", func(t *testing.T) {
		// Arrange
		db := map[int]repository.ProductAttributesMap{}
		st := repository.NewRepositoryProductMap(db, 0, "")
		au := auth.NewAuthTokenBasic("test_token")
		hd := handlers.NewHandlerProducts(st, au)
		hdFunc := hd.Update()
		req := httptest.NewRequest("PATCH", "/products/1", nil)
		res := httptest.NewRecorder()
		// Act

		hdFunc(res, req)
		// Assert

		expectedCode := http.StatusUnauthorized
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		expectedBody := `{"status":"Unauthorized","message":"Unauthorized"}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

}
