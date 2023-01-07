package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *Api) UpdateExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	ep := Expense{}
	err := c.Bind(&ep)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := a.Db.QueryRow("UPDATE expenses SET title = $1 , amount = $2 , note = $3 , tags = $4 WHERE id = $5 RETURNING id , title , amount , note , tags", ep.Title, ep.Amount, ep.Note, ep.Tags, id)
	err = row.Scan(&ep.ID, &ep.Title, &ep.Amount, &ep.Note, &ep.Tags)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, ep)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}
