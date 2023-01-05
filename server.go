package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yuttanar/assessment/expense"
)

func main() {
	var PORT string = os.Getenv("PORT")
	var DATABASE_URL string = os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		fmt.Println("Please set ENV 'DATABASE_URL' before start the server.")
		return
	}
	if PORT == "" {
		PORT = ":2565"
	}

	expense.InitDB()

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", PORT)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start go server
	go func() {
		if err := e.Start(PORT); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down the server, bye...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
