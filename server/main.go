package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/isayme/go-config"
	"github.com/isayme/traefik-authserver/server/src/conf"
	"github.com/isayme/traefik-authserver/server/src/service"
	"github.com/isayme/traefik-authserver/server/src/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

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

type sessionEntry struct {
	Config conf.SessionConfig
	Store  sessions.Store
}

var sessionEntries []sessionEntry

func main() {
	globalConfig := conf.Config{
		Session: conf.SessionConfig{
			Name:     "sid",
			MaxAge:   86400 * 7,
			HttpOnly: false,
			Secure:   false,
		},
	}
	config.Parse(&globalConfig)

	if len(globalConfig.Sessions) > 0 {
		for i, s := range globalConfig.Sessions {
			if util.IsBlank(s.Secret) {
				log.Errorf("config sessions[%d].secret is required", i)
				return
			}
			if util.IsBlank(s.LoginUrl) {
				log.Errorf("config sessions[%d].loginUrl is required", i)
				return
			}
			sessionEntries = append(sessionEntries, sessionEntry{
				Config: s,
				Store:  sessions.NewCookieStore([]byte(s.Secret)),
			})
		}
	} else {
		if util.IsBlank(globalConfig.Session.Secret) {
			log.Error("config session.secret is required")
			return
		}
		if util.IsBlank(globalConfig.Session.LoginUrl) {
			log.Error("config session.loginUrl is required")
			return
		}
		sessionEntries = append(sessionEntries, sessionEntry{
			Config: globalConfig.Session,
			Store:  sessions.NewCookieStore([]byte(globalConfig.Session.Secret)),
		})
	}

	e := echo.New()

	e.Use(middleware.Recover())

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))

	e.POST("/api/login", func(c echo.Context) error {
		reqBody := LoginReq{}
		if err := c.Bind(&reqBody); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		username := reqBody.Username
		if util.IsBlank(username) {
			return responseError(c, http.StatusBadRequest, "UsernameRequired", "username required")
		}

		for _, user := range globalConfig.Users {
			if user.Username == username {
				if verifyPassword(user.Password, reqBody.Password) {
					setSession(c, username)
					return c.JSON(http.StatusOK, LoginResp{Username: username})
				}
			}
		}
		return responseError(c, http.StatusForbidden, "Forbidden", "invalid username/password")
	})

	e.POST("/api/logout", func(c echo.Context) error {
		setSession(c, "")
		return c.JSON(http.StatusOK, LogoutResp{})
	})

	e.GET("/api/me", func(c echo.Context) error {
		username := getSession(c)
		if util.IsNotBlank(username) {
			return c.JSON(http.StatusOK, GetMeResp{Username: username})
		}
		return responseError(c, http.StatusUnauthorized, "Unauthorized", "unauthorized")
	})

	e.GET("/api/check-login", func(c echo.Context) error {
		username := getSession(c)
		if util.IsNotBlank(username) {
			return c.JSON(http.StatusOK, GetMeResp{Username: username})
		}

		entry := getEntry(c)
		if entry == nil {
			return responseError(c, http.StatusNotFound, "NoSessionConfig", "no matching session config")
		}

		location := entry.Config.LoginUrl
		uri, err := url.Parse(location)
		if err != nil {
			return c.Redirect(http.StatusFound, location)
		}

		forwardedProto := c.Request().Header.Get("X-Forwarded-Proto")
		forwardedHost := c.Request().Header.Get("X-Forwarded-Host")
		forwardedUri := c.Request().Header.Get("X-Forwarded-Uri")
		if util.IsNotBlank(forwardedProto) && util.IsNotBlank(forwardedHost) && util.IsNotBlank(forwardedUri) {
			nextUrl := fmt.Sprintf("%s://%s%s", forwardedProto, forwardedHost, forwardedUri)
			query := uri.Query()
			query.Add("next_url", nextUrl)
			uri.RawQuery = query.Encode()
			location = uri.String()
		}

		return c.Redirect(http.StatusFound, location)
	})

	githubConfig := globalConfig.Github
	if util.IsNotBlank(githubConfig.ClientId) && util.IsNotBlank(githubConfig.ClientSecret) {
		githubService := service.NewGithub(&githubConfig)

		e.GET("/oauth/github/login", func(c echo.Context) error {
			nextUrl := c.QueryParam("next_url")
			url := githubService.GenAuthorizeUrl(nextUrl)

			return c.Redirect(http.StatusFound, url)
		})

		e.GET("/oauth/github/redirect", func(c echo.Context) error {
			code := c.QueryParam("code")
			nextUrl := c.QueryParam("next_url")
			// state := c.QueryParam("state")

			ctx := c.Request().Context()

			accessTokenInfo, err := githubService.ExchangeAccessToken(ctx, code)
			if err != nil {
				return responseError(c, http.StatusBadRequest, "", err.Error())
			}

			githubUser, err := githubService.GetUser(ctx, accessTokenInfo.AccessToken)
			if err != nil {
				return responseError(c, http.StatusBadRequest, "", err.Error())
			}

			for _, user := range globalConfig.Users {
				if user.Github == githubUser.Login {
					setSession(c, user.Username)

					if util.IsNotBlank(nextUrl) {
						return c.Redirect(http.StatusFound, nextUrl)
					}
					return c.JSON(http.StatusOK, LoginResp{Username: user.Username})
				}
			}

			return responseError(c, http.StatusBadRequest, "NotFound", "no user bound with current github user")
		})
	}

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "public",
		Index: "index.html",
		HTML5: true,
	}))

	e.Logger.Fatal(e.Start(":1323"))
}

func getHost(c echo.Context) string {
	host := c.Request().Header.Get("X-Forwarded-Host")
	if util.IsNotBlank(host) {
		return host
	}
	return c.Request().Host
}

func matchSession(host string) *sessionEntry {
	for _, entry := range sessionEntries {
		domain := entry.Config.Domain
		if strings.HasSuffix(host, domain) {
			return &entry
		}
	}
	return nil
}

func getEntry(c echo.Context) *sessionEntry {
	host := getHost(c)
	return matchSession(host)
}

func setSession(c echo.Context, username string) {
	entry := getEntry(c)
	if entry == nil {
		return
	}

	sess, _ := entry.Store.Get(c.Request(), entry.Config.Name)
	sess.Options = &sessions.Options{
		Domain:   entry.Config.Domain,
		Path:     "/",
		MaxAge:   entry.Config.MaxAge,
		HttpOnly: entry.Config.HttpOnly,
		Secure:   entry.Config.Secure,
	}

	sess.Values["username"] = username
	sess.Save(c.Request(), c.Response())
}

func getSession(c echo.Context) string {
	entry := getEntry(c)
	if entry == nil {
		return ""
	}

	sess, _ := entry.Store.Get(c.Request(), entry.Config.Name)
	if v := sess.Values["username"]; v != nil {
		if username, ok := v.(string); ok {
			if util.IsNotBlank(username) {
				return username
			}
		}
	}

	return ""
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
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
