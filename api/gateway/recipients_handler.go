package gateway

import (
	"net/http"
	"sumup-notifications/pkg/model"
	"sumup-notifications/pkg/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type RecipientsHandler struct {
	db *pgxpool.Pool
}

func (handler RecipientsHandler) Mount(e *echo.Echo) {
	e.GET("/recipients/:recipient", echo.HandlerFunc(handler.Get))
	e.PUT("/recipients/:recipient", echo.HandlerFunc(handler.Put))
	e.DELETE("/recipients/:recipient", echo.HandlerFunc(handler.Delete))
}

func (handler RecipientsHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	recipient := c.Param("recipient")

	recipients := storage.Recipients{DB: handler.db}
	contact, err := recipients.Load(ctx, recipient)
	if err != nil {
		if err == storage.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Recipient not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, contact)
}

func (handler RecipientsHandler) Put(c echo.Context) error {
	ctx := c.Request().Context()
	recipient := c.Param("recipient")

	var contact model.RecipientContact
	err := c.Bind(&contact)
	if err != nil {
		return err
	}

	recipients := storage.Recipients{DB: handler.db}
	err = recipients.Store(ctx, recipient, contact)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (handler RecipientsHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	recipient := c.Param("recipient")

	recipients := storage.Recipients{DB: handler.db}
	err := recipients.Delete(ctx, recipient)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
