// Package app is the composition root for auth.
package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	nethttp "net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	pb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/auth/v1"
	"github.com/helixgitpx/helixgitpx/services/auth/internal/domain"
	grpchandler "github.com/helixgitpx/helixgitpx/services/auth/internal/handler/grpc"
	httphandler "github.com/helixgitpx/helixgitpx/services/auth/internal/handler/http"
	"github.com/helixgitpx/helixgitpx/services/auth/internal/repo"
	hauth "github.com/helixgitpx/platform/auth"
	"github.com/helixgitpx/platform/config"
	hgin "github.com/helixgitpx/platform/gin"
	hgrpc "github.com/helixgitpx/platform/grpc"
	"github.com/helixgitpx/platform/health"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
	"github.com/helixgitpx/platform/telemetry"
)

type cfg struct {
	HTTPAddr     string `env:"HTTP_ADDR" default:":8002"`
	GRPCAddr     string `env:"GRPC_ADDR" default:":9002"`
	HealthAddr   string `env:"HEALTH_ADDR" default:":8082"`
	PostgresDSN  string `env:"POSTGRES_DSN" vault:"kv/auth#pg_dsn" required:"true"`
	JWTPrivPEM   string `env:"JWT_PRIVATE_PEM" vault:"kv/auth/jwt#private_pem" required:"true"`
	OIDCIssuer   string `env:"OIDC_ISSUER" default:"https://keycloak.helix.local/realms/helixgitpx"`
	OIDCClient   string `env:"OIDC_CLIENT_ID" default:"auth-service"`
	OIDCSecret   string `env:"OIDC_CLIENT_SECRET" vault:"kv/auth#oidc_client_secret"`
	OIDCRedirect string `env:"OIDC_REDIRECT" default:"https://auth.helix.local/v1/auth/callback"`
	OTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Version      string `env:"VERSION" default:"m3-dev"`
}

// Run wires dependencies and serves until ctx is done.
func Run(ctx context.Context, lg *log.Logger) error {
	var c cfg
	if err := config.Load(&c, config.Options{Prefix: "AUTH"}); err != nil {
		return err
	}

	shutdownTel, _ := telemetry.Start(ctx, telemetry.Options{Service: "auth", Version: c.Version, OTLPEndpoint: c.OTLPEndpoint})
	defer func() {
		sh, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTel(sh)
	}()

	pool, err := pg.Open(ctx, pg.Options{DSN: c.PostgresDSN})
	if err != nil {
		return err
	}
	defer pool.Close()

	block, _ := pem.Decode([]byte(c.JWTPrivPEM))
	if block == nil {
		return errors.New("auth: bad private key PEM")
	}
	var priv *rsa.PrivateKey
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		priv = k
	} else {
		k2, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return fmt.Errorf("auth: parse private key: %w / %v", err, err2)
		}
		var ok bool
		priv, ok = k2.(*rsa.PrivateKey)
		if !ok {
			return errors.New("auth: private key is not RSA")
		}
	}
	signer := hauth.NewSigner(priv, "kid-1", c.OIDCIssuer)

	provider, err := oidc.NewProvider(ctx, c.OIDCIssuer)
	if err != nil {
		return fmt.Errorf("auth: OIDC discovery: %w", err)
	}
	oauthCfg := &oauth2.Config{
		ClientID:     c.OIDCClient,
		ClientSecret: c.OIDCSecret,
		RedirectURL:  c.OIDCRedirect,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	users := &repo.UsersPG{Pool: pool}
	sessions := &repo.SessionsPG{Pool: pool}
	pats := &repo.PATsPG{Pool: pool}
	mfa := &repo.MFAPG{Pool: pool}
	issuer := &domain.TokensIssuer{
		Signer:     signer,
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 14 * 24 * time.Hour,
	}

	validator := hauth.NewValidatorFromKey(&priv.PublicKey, "kid-1", c.OIDCIssuer)
	_ = validator // will be wired when interceptor surface settles

	grpcSrv, err := hgrpc.NewServer(hgrpc.Options{})
	if err != nil {
		return err
	}
	pb.RegisterAuthServiceServer(grpcSrv, &grpchandler.Server{
		Users: users, Sessions: sessions, PATs: pats, MFA: mfa, Issuer: issuer,
	})

	router := hgin.NewRouter(hgin.Options{Service: "auth", Version: c.Version})
	(&httphandler.Router{
		OIDC: provider, OAuth: oauthCfg, Users: users, Sessions: sessions, Issuer: issuer,
		RefreshTTL: 14 * 24 * time.Hour,
	}).Register(router)

	hh := health.New()
	hh.Register("pg", pg.Probe(pool))
	hmux := nethttp.NewServeMux()
	hh.Routes(hmux)
	telemetry.RegisterPprof(hmux)

	grpcL, err := net.Listen("tcp", c.GRPCAddr)
	if err != nil {
		return err
	}
	httpSrv := &nethttp.Server{Addr: c.HTTPAddr, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	healthSrv := &nethttp.Server{Addr: c.HealthAddr, Handler: hmux, ReadHeaderTimeout: 5 * time.Second}

	errCh := make(chan error, 3)
	go func() { errCh <- grpcSrv.Serve(grpcL) }()
	go func() { errCh <- httpSrv.ListenAndServe() }()
	go func() { errCh <- healthSrv.ListenAndServe() }()

	lg.Info("auth serving", "grpc", c.GRPCAddr, "http", c.HTTPAddr, "health", c.HealthAddr)

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			lg.Error("server exited", "err", err.Error())
		}
	}

	sh, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcSrv.GracefulStop()
	_ = httpSrv.Shutdown(sh)
	_ = healthSrv.Shutdown(sh)
	return nil
}
