package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/dberrors"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/models"
)

func (s *EchoServer) GetAllCustomers(ctx echo.Context) error {
	emailAddress := ctx.QueryParam("emailAddress")

	customers, err := s.DB.GetAllCustomers(ctx.Request().Context(), emailAddress)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, customers)

}

func (s *EchoServer) AddCustomer(ctx echo.Context) error {
	customer := new(models.Customer)

	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	customer, err := s.DB.AddCustomer(ctx.Request().Context(), customer)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, customer)

}

func (s *EchoServer) GetCustomer(ctx echo.Context) error {
	customerId := ctx.Param("id")

	customer, err := s.DB.GetCustomer(ctx.Request().Context(), customerId)

	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, customer)
}

func (s *EchoServer) UpdateCustomer(ctx echo.Context) error {
	customerId := ctx.Param("id")
	customer := new(models.Customer)

	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if customerId != customer.CustomerID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match with the id on body")
	}

	updatedCustomer, err := s.DB.UpdateCustomer(ctx.Request().Context(), customer)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, updatedCustomer)
}

func (s *EchoServer) DeleteCustomer(ctx echo.Context) error {
	ID := ctx.Param("id")
	err := s.DB.DeleteCustomer(ctx.Request().Context(), ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusResetContent)
}
