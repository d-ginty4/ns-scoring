package local

import (
	"log/slog"

	e "apis/echo"
	"apis/logger"

	"github.com/joho/godotenv"
)

func Handler() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("No .env file found")
	}

	echoInstance := e.NewEchoInstance()

	slog.Info("starting local server", "addr", ":3000")
	if err := echoInstance.Start(":3000"); err != nil {
		slog.Error("local server failed", logger.Err(err))
		panic(err)
	}
}
