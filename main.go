package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

func main() {
	fmt.Println("你好，世界！")
	conf, err := ReadConfig("tmp/conf.toml")
	PrintAndExit(err)
	db, err := conf.Database.Open()
	PrintAndExit(err)
	defer db.Close()
	err = db.Ping()
	PrintAndExit(err)
	stmt, err := conf.Database.Prepare(db)
	PrintAndExit(err)

	e := echo.New()

	e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper: func(c echo.Context) bool {

			return false
		},
		Validator: func(username, password string, c echo.Context) (bool, error) {
			hash, err := stmt.FindStudentPasswordById(username)
			if err != nil {
				return false, nil
			}
			return CheckPasswordHash(password, hash), nil
		},
		Realm: "Savage Garden",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	logger := NewZapLogger()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP: true, LogMethod: true, LogURI: true, LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("remote_ip", v.RemoteIP),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	e.Logger.Fatal(e.Start(":1323"))
}

func PrintAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func NewZapLogger() *zap.Logger {
	conf := zap.NewDevelopmentConfig()
	conf.DisableStacktrace = true
	conf.DisableCaller = true
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	logger, _ := conf.Build()
	return logger
}
