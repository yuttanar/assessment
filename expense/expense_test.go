package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var ep = &Expense{
	ID:     1,
	Title:  "strawberry smoothie",
	Amount: 79,
	Note:   "night market promotion discount 10 bath",
	Tags:   []string{"food", "beverage"},
}

var expenseJSON string = `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`

func TestShouldGetExpense(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(ep.ID, ep.Title, ep.Amount, ep.Note, ep.Tags)

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	mock.ExpectQuery("SELECT id, title, amount , note , tags FROM expenses WHERE id = (.+)").WithArgs("1").WillReturnRows(rows)

	app := &Api{Db}

	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, app.GetExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expenseJSON, strings.TrimSpace(rec.Body.String()))
	}
}

func TestShouldCreateExpense(t *testing.T) {
	var expenseReqJSON string = `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expenseReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	row := sqlmock.NewRows([]string{"id"}).AddRow(ep.ID)

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	mock.ExpectQuery("INSERT INTO expenses (.+) values (.+)  RETURNING id").WithArgs(ep.Title, ep.Amount, ep.Note, ep.Tags).WillReturnRows(row)

	app := &Api{Db}

	c := e.NewContext(req, rec)

	if assert.NoError(t, app.CreateExpenseHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expenseJSON, strings.TrimSpace(rec.Body.String()))
	}
}
