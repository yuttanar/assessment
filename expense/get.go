package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *Api) GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	row := a.Db.QueryRow("SELECT id, title, amount , note , tags FROM expenses WHERE id = $1", id)

	ep := Expense{}
	err := row.Scan(&ep.ID, &ep.Title, &ep.Amount, &ep.Note, &ep.Tags)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, ep)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}
