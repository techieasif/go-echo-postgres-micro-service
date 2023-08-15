package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techieasif/wisdom/internal/dberrors"
	"github.com/techieasif/wisdom/internal/models"
)

func (s *EchoServer) GetAllProducts(ctx echo.Context) error {
	vendorID := ctx.QueryParam("vendorId")

	products, err := s.DB.GetAllProducts(ctx.Request().Context(), vendorID)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)

	}

	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) AddProduct(ctx echo.Context) error {
	product := new(models.Product)

	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	addResult, err := s.DB.AddProduct(ctx.Request().Context(), product)

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

func (s *EchoServer) GetProduct(ctx echo.Context) error {
	productId := ctx.Param("id")

	product, err := s.DB.GetProduct(ctx.Request().Context(), productId)

	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, product)
}


func (s *EchoServer) UpdateProduct(ctx echo.Context) error {
	productId := ctx.Param("id")
	product := new(models.Product)

	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if productId != product.ProductID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match with the id on body")
	}

	updatedProduct, err := s.DB.UpdateProduct(ctx.Request().Context(), product)

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

	return ctx.JSON(http.StatusOK, updatedProduct)
}


func (s *EchoServer) DeleteProduct(ctx echo.Context) error {
	ID := ctx.Param("id")
	err := s.DB.DeleteProduct(ctx.Request().Context(), ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusResetContent)
}