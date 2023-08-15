package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/dberrors"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/models"
)

func (s *EchoServer) GetAllVendors(ctx echo.Context) error {
	vendors, err := s.DB.GetAllVendors(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, vendors)

}

func (s *EchoServer) AddVendor(ctx echo.Context) error {
	vendor := new(models.Vendor)

	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	addResult, err := s.DB.AddVendor(ctx.Request().Context(), vendor)

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

func (s *EchoServer) GetVendor(ctx echo.Context) error {
	vendorId := ctx.Param("id")

	vendor, err := s.DB.GetVendor(ctx.Request().Context(), vendorId)

	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, vendor)
}



func (s *EchoServer) UpdateVendor(ctx echo.Context) error {
	vendorId := ctx.Param("id")
	vendor := new(models.Vendor)

	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if vendorId != vendor.VendorID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match with the id on body")
	}

	updatedVendor, err := s.DB.UpdateVendor(ctx.Request().Context(), vendor)

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

	return ctx.JSON(http.StatusOK, updatedVendor)
}

func (s *EchoServer) DeleteVendor(ctx echo.Context) error {
	ID := ctx.Param("id")
	err := s.DB.DeleteVendor(ctx.Request().Context(), ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusResetContent)
}