package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/helixgitpx/helixgitpx/services/auth/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/auth/internal/repo"
)

type Router struct {
	OIDC       *oidc.Provider
	OAuth      *oauth2.Config
	Users      *repo.UsersPG
	Sessions   *repo.SessionsPG
	Issuer     *domain.TokensIssuer
	RefreshTTL time.Duration
}

func (r *Router) Register(g *gin.Engine) {
	g.GET("/v1/auth/callback", r.callback)
	g.POST("/v1/auth/refresh", r.refresh)
	g.GET("/.well-known/jwks.json", r.jwks)
}

func (r *Router) callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	tok, err := r.OAuth.Exchange(ctx, code)
	if err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "oauth exchange: " + err.Error()})
		return
	}
	rawIDToken, ok := tok.Extra("id_token").(string)
	if !ok {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "no id_token"})
		return
	}
	verifier := r.OIDC.Verifier(&oidc.Config{ClientID: r.OAuth.ClientID})
	idTok, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"error": "id_token verify: " + err.Error()})
		return
	}
	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	_ = idTok.Claims(&claims)

	u, err := r.Users.UpsertBySubject(ctx, claims.Sub, claims.Email, claims.Name)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": "upsert user: " + err.Error()})
		return
	}

	persist := func(sid uuid.UUID, exp time.Time) error {
		return r.Sessions.Create(ctx, sid, u.ID.String(), exp, c.Request.UserAgent(), c.ClientIP())
	}
	tokens, err := r.Issuer.Issue(ctx, u.ID.String(), persist)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": "issue: " + err.Error()})
		return
	}

	c.SetCookie("access_token", tokens.AccessToken, int(tokens.ExpiresIn.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(r.RefreshTTL.Seconds()), "/", "", true, true)
	c.JSON(nethttp.StatusOK, gin.H{"user": u.Email})
}

func (r *Router) refresh(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"error": "no refresh cookie"})
		return
	}
	sid, err := uuid.Parse(rt)
	if err != nil {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"error": "bad refresh"})
		return
	}
	sess, err := r.Sessions.Active(c.Request.Context(), sid)
	if err != nil || sess == nil {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"error": "expired or revoked"})
		return
	}
	_ = r.Sessions.Revoke(c.Request.Context(), sid, sess.UserID)

	persist := func(newID uuid.UUID, exp time.Time) error {
		return r.Sessions.Create(c.Request.Context(), newID, sess.UserID, exp, c.Request.UserAgent(), c.ClientIP())
	}
	tokens, err := r.Issuer.Issue(c.Request.Context(), sess.UserID, persist)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": "issue: " + err.Error()})
		return
	}
	c.SetCookie("access_token", tokens.AccessToken, int(tokens.ExpiresIn.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(r.RefreshTTL.Seconds()), "/", "", true, true)
	c.JSON(nethttp.StatusOK, gin.H{"ok": true})
}

func (r *Router) jwks(c *gin.Context) {
	// Minimal JWKS — production key rotation lives in the signer; M3 exposes only the
	// current public key. Populated from app.Run via a setter, for simplicity keep empty.
	c.Header("Content-Type", "application/json")
	_ = json.NewEncoder(c.Writer).Encode(map[string]any{"keys": []any{}})
}
