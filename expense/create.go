package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateExpenseHandler(c echo.Context) error {
	ep := Expense{}
	err := c.Bind(&ep)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount , note , tags) values ($1, $2 , $3 , $4)  RETURNING id", ep.Title, ep.Amount, ep.Note, ep.Tags)
	err = row.Scan(&ep.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ep)
}
