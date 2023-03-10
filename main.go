package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

func main() {
	conf, err := ReadConfig("tmp/conf.toml")
	PrintAndExit(err)
	db, err := conf.Database.Open()
	PrintAndExit(err)
	defer db.Close()
	err = db.Ping()
	PrintAndExit(err)
	stmt, err := conf.Database.Prepare(db)
	PrintAndExit(err)

	logger := NewZapLogger()
	c := cron.New(cron.WithParser(cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month,
	)))
	// 每小时清除一次过期cookie
	_, err = c.AddFunc("0 * * *", func() {
		num, err := stmt.DeleteSessionExpired()
		if err == nil {
			logger.Info("cron", zap.Int("deleted_sessions", int(num)))
		} else {
			logger.Error("cron", zap.String("failed_delete_session", err.Error()))
		}
	})
	PrintAndExit(err)

	const CookieName = "savage_garden"
	const CookieExpireDur = 7 * 24 * time.Hour
	uuid.EnableRandPool()
	NewCookie := func() *http.Cookie {
		return &http.Cookie{
			Name:     CookieName,
			Value:    uuid.NewString(),
			Path:     "/",
			Domain:   conf.Server.Domain,
			Expires:  time.Now().Add(CookieExpireDur),
			Secure:   false,
			HttpOnly: true,
		}
	}

	e := echo.New()

	// cookie setter
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 如果没有cookie或cookie过期，则生成cookie并覆盖
			var refresh bool
			cookie, err := c.Cookie(CookieName)
			if refresh = err != nil; !refresh { // ↓找到了cookie↓
				notExpired, err := stmt.VerifySessionNotExpired(cookie.Value)
				if refresh = err != nil; !refresh { // ↓数据库中找到session↓
					if refresh = !notExpired; refresh { // ↓session已过期↓
						_ = stmt.DeleteSessionById(cookie.Value)
					}
				}
			}
			if refresh {
				cookie = NewCookie()
				c.SetCookie(cookie)
				_ = stmt.CreateSessionByCookie(cookie)
			}
			return next(c)
		}
	})

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

	secure := e.Group("/secure")

	const AuthRealm = "Savage Garden"
	const IdKey = "student_id"
	// cookie checker
	secure.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 只有cookie是已登录状态才可以访问后面内容
			cookie, err := c.Cookie(CookieName)
			if err == nil {
				sid, err := stmt.FindSessionSidById(cookie.Value)
				if err == nil {
					c.Set(IdKey, sid)
					return next(c)
				}
			}
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
	})

	e.GET("/login", func(c echo.Context) error {
		// 如果通过了Basic Auth，说明验证通过，新建一个已登录状态的cookie
		if cookie, err := c.Cookie(CookieName); err == nil {
			_ = stmt.DeleteSessionById(cookie.Value)
		}
		cookie := NewCookie()
		c.SetCookie(cookie)
		_ = stmt.CreateSessionByCookieWithLoggedIn(cookie, c.Get(IdKey).(string))
		return c.NoContent(http.StatusOK)
	}, middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(username, password string, c echo.Context) (bool, error) {
			hash, err := stmt.FindStudentPasswordById(username)
			if err != nil {
				return false, nil
			}
			rs := CheckPasswordHash(password, hash)
			if rs {
				c.Set(IdKey, username)
			}
			return rs, nil
		},
		Realm: AuthRealm,
	}))

	e.GET("/logout", func(c echo.Context) error {
		if cookie, err := c.Cookie(CookieName); err == nil {
			_ = stmt.DeleteSessionById(cookie.Value)
		}
		cookie := NewCookie()
		c.SetCookie(cookie)
		_ = stmt.CreateSessionByCookie(cookie)
		return c.NoContent(http.StatusOK)
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Public Path Test")
	})

	secure.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Secure Path Test")
	})

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
	conf.OutputPaths = []string{"stdout"}
	conf.DisableStacktrace = true
	conf.DisableCaller = true
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logger, _ := conf.Build()
	return logger
}
