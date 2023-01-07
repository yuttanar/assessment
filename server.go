package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yuttanar/assessment/expense"
)

func main() {
	var PORT string = os.Getenv("PORT")
	var DATABASE_URL string = os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		fmt.Println("Please set ENV 'DATABASE_URL' before start the server.")
		return
	}
	if PORT == "" || !checkPortPatternMatch(PORT) {
		PORT = ":2565"
	}
	if !checkPortIsBeingUsed(PORT) {
		fmt.Printf("Port %s is being used , Please use another port.", PORT)
		return
	}

	expense.InitDB()

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", PORT)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app := &expense.Api{Db: expense.Db}
	e.POST("/expenses", app.CreateExpenseHandler)
	e.GET("/expenses/:id", app.GetExpenseHandler)
	e.PUT("/expenses/:id", app.UpdateExpenseHandler)

	// Start go server
	go func() {
		if err := e.Start(PORT); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down the go server, bye bye...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

func checkPortIsBeingUsed(port string) bool {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return false
	}
	defer ln.Close()
	return true
}

func checkPortPatternMatch(port string) bool {
	re := regexp.MustCompile(`^:((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$`)
	return re.MatchString(port)
}
