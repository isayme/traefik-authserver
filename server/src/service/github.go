package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/traefik-authserver/server/src/conf"
	"github.com/isayme/traefik-authserver/server/src/util"
)

const GithubOAuthAuthorizeUrl = "https://github.com/login/oauth/authorize"
const GithubOAuthAccessTokenUrl = "https://github.com/login/oauth/access_token"
const GithubApiUserUrl = "https://api.github.com/user"
const ScopeUserEmail = "user:email"

type Github struct {
	clientId     string
	clientSecret string
	redirectUrl  string
}

func NewGithub(config *conf.Github) *Github {
	return &Github{
		clientId:     config.ClientId,
		clientSecret: config.ClientSecret,
		redirectUrl:  config.RedirectUrl,
	}
}

func (g *Github) getRedirectUrl(nextUrl string) string {
	redirectUrl := g.redirectUrl
	if util.IsNotBlank(redirectUrl) && util.IsNotBlank(nextUrl) {
		u, err := url.Parse(redirectUrl)
		if err != nil {
			logger.Warnw("url.Parse fail", "redirectUrl", redirectUrl)
		} else {
			query := u.Query()
			query.Add("next_url", nextUrl)
			u.RawQuery = query.Encode()

			redirectUrl = u.String()
		}
	}

	return redirectUrl
}

func (g *Github) GenAuthorizeUrl(nextUrl string) string {
	redirectUrl := g.getRedirectUrl(nextUrl)

	state := util.UUID()
	if util.IsBlank(redirectUrl) {
		return fmt.Sprintf("%s?client_id=%s&scope=%s&state=%s", GithubOAuthAuthorizeUrl, g.clientId, ScopeUserEmail, state)
	} else {
		return fmt.Sprintf("%s?client_id=%s&scope=%s&state=%s&redirect_uri=%s", GithubOAuthAuthorizeUrl, g.clientId, ScopeUserEmail, state, url.QueryEscape(redirectUrl))
	}
}

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (g *Github) ExchangeAccessToken(ctx context.Context, code string) (*GithubAccessToken, error) {
	reqBody := map[string]interface{}{
		"client_id":     g.clientId,
		"client_secret": g.clientSecret,
		"code":          code,
	}

	var response GithubAccessToken

	headers := http.Header{}
	err := util.Request(ctx, http.MethodPost, GithubOAuthAccessTokenUrl, headers, reqBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type GithubUser struct {
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
	Name      string `json:"name"`
}

func (g *Github) GetUser(ctx context.Context, accessToken string) (*GithubUser, error) {
	var response GithubUser

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("token %s", accessToken))
	err := util.Request(ctx, http.MethodGet, GithubApiUserUrl, headers, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
