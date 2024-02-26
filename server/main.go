package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/isayme/go-config"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

type SessionConfig struct {
	Name     string `json:"name" yaml:"name"`
	Secret   string `json:"secret" yaml:"secret"`
	LoginUrl string `json:"loginUrl" yaml:"loginUrl"`
}

type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Config struct {
	Session SessionConfig `json:"session" yaml:"session"`
	Users   []User        `json:"users" yaml:"users"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Username string `json:"username"`
}

type LogoutReq struct {
}

type LogoutResp struct {
}

type GetMeResp struct {
	Username string `json:"username"`
}

func main() {
	var globalConfig Config
	config.Parse(&globalConfig)

	e := echo.New()

	e.Use(middleware.Recover())

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(globalConfig.Session.Secret))))

	e.POST("/login", func(c echo.Context) error {
		reqBody := LoginReq{}
		if err := c.Bind(&reqBody); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		username := reqBody.Username
		if username == "" {
			return responseError(c, http.StatusBadRequest, "UsernameRequired", "username required")
		}

		for _, user := range globalConfig.Users {
			if user.Username == username {
				if verifyPassword(user.Password, reqBody.Password) {
					sess, _ := session.Get(globalConfig.Session.Name, c)
					sess.Options = &sessions.Options{
						Domain:   "localtest.me",
						Path:     "/",
						MaxAge:   86400 * 7,
						HttpOnly: true,
					}

					sess.Values["username"] = username
					sess.Save(c.Request(), c.Response())

					return c.JSON(http.StatusOK, LoginResp{Username: username})
				}
			}
		}
		return responseError(c, http.StatusForbidden, "Forbidden", "invalid username/password")
	})

	e.POST("/logout", func(c echo.Context) error {
		sess, _ := session.Get(globalConfig.Session.Name, c)
		sess.Values["username"] = ""
		sess.Save(c.Request(), c.Response())
		return c.JSON(http.StatusOK, LogoutResp{})
	})

	e.GET("/me", func(c echo.Context) error {
		sess, _ := session.Get(globalConfig.Session.Name, c)
		if v := sess.Values["username"]; v != nil {
			if username, ok := v.(string); ok {
				if username != "" {
					return c.JSON(http.StatusOK, GetMeResp{Username: username})
				}
			}
		}

		return responseError(c, http.StatusUnauthorized, "Unauthorized", "unauthorized")
	})

	e.GET("/check-login", func(c echo.Context) error {
		sess, _ := session.Get(globalConfig.Session.Name, c)
		if v := sess.Values["username"]; v != nil {
			if username, ok := v.(string); ok {
				if username != "" {
					return c.JSON(http.StatusOK, GetMeResp{Username: username})
				}
			}
		}

		return c.Redirect(http.StatusFound, globalConfig.Session.LoginUrl)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func responseError(c echo.Context, statusCode int, errCode, errMessage string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"code":    errCode,
		"message": errMessage,
	})
}

// func hashPassword(password string) (string, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(hash), nil
// }

func verifyPassword(hash, password string) bool {
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
