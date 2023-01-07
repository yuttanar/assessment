//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	var ep Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&ep)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, ep.ID)
	assert.Equal(t, "strawberry smoothie", ep.Title)
	assert.Equal(t, 79.00, ep.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", ep.Note)
	assert.ElementsMatch(t, []string{"food", "beverage"}, ep.Tags)
}

func TestGetExpenseByID(t *testing.T) {
	ep := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(ep.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, ep.ID, latest.ID)
	assert.Equal(t, ep.Title, latest.Title)
	assert.Equal(t, ep.Amount, latest.Amount)
	assert.Equal(t, ep.Note, latest.Note)
	assert.ElementsMatch(t, ep.Tags, latest.Tags)
}

func seedExpense(t *testing.T) Expense {
	var ep Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	err := request(http.MethodPost, uri("expenses"), body).Decode(&ep)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return ep
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
