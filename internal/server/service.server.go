package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/dberrors"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/models"
)

func (s *EchoServer) GetAllServices(ctx echo.Context) error {
	services, err := s.DB.GetAllServices(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)

	}

	return ctx.JSON(http.StatusOK, services)
}

func (s *EchoServer) AddService(ctx echo.Context) error {
	service := new(models.Service)

	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	addResult, err := s.DB.AddService(ctx.Request().Context(), service)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, addResult)

}

func (s *EchoServer) GetService(ctx echo.Context) error {
	serviceId := ctx.Param("id")

	service, err := s.DB.GetService(ctx.Request().Context(), serviceId)

	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, service)
}

func (s *EchoServer) UpdateService(ctx echo.Context) error {
	serviceId := ctx.Param("id")
	service := new(models.Service)

	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if serviceId != service.ServiceID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match with the id on body")
	}

	updatedService, err := s.DB.UpdateService(ctx.Request().Context(), service)

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

	return ctx.JSON(http.StatusOK, updatedService)
}

func (s *EchoServer) DeleteService(ctx echo.Context) error {
	ID := ctx.Param("id")
	err := s.DB.DeleteService(ctx.Request().Context(), ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusResetContent)
}