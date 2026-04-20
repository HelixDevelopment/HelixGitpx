package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/auth/v1"
	"github.com/helixgitpx/helixgitpx/services/auth/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/auth/internal/repo"
	hauth "github.com/helixgitpx/platform/auth"
)

// Server implements AuthService for the non-OIDC RPCs.
// OIDC code-exchange + refresh live in handler/http/router.go because they
// set cookies on the response.
type Server struct {
	pb.UnimplementedAuthServiceServer

	Users    *repo.UsersPG
	Sessions *repo.SessionsPG
	PATs     *repo.PATsPG
	MFA      *repo.MFAPG
	Issuer   *domain.TokensIssuer
}

func (s *Server) IssuePAT(ctx context.Context, req *pb.IssuePATRequest) (*pb.PAT, error) {
	uid, ok := hauth.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user id in context")
	}
	token, hashed, err := domain.IssuePAT()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue: %v", err)
	}
	var expiresAt *time.Time
	if req.ExpiresInSeconds > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresInSeconds) * time.Second)
		expiresAt = &t
	}
	row, err := s.PATs.Insert(ctx, uid, req.Name, hashed, req.Scopes, expiresAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "persist: %v", err)
	}
	return patProto(row, token), nil
}

func (s *Server) RevokePAT(ctx context.Context, req *pb.RevokePATRequest) (*emptypb.Empty, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	if err := s.PATs.Revoke(ctx, req.Id, uid); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListPATs(ctx context.Context, _ *emptypb.Empty) (*pb.ListPATsResponse, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	pats, err := s.PATs.List(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list: %v", err)
	}
	resp := &pb.ListPATsResponse{}
	for _, p := range pats {
		pp := p
		resp.Pats = append(resp.Pats, patProto(&pp, ""))
	}
	return resp, nil
}

func (s *Server) WhoAmI(ctx context.Context, _ *emptypb.Empty) (*pb.User, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	u, err := s.Users.GetBySubject(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user: %v", err)
	}
	return &pb.User{Id: u.ID.String(), Subject: u.Subject, Email: u.Email, DisplayName: u.DisplayName}, nil
}

func (s *Server) EnrollTOTP(ctx context.Context, _ *emptypb.Empty) (*pb.EnrollTOTPResponse, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	u, err := s.Users.GetBySubject(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user: %v", err)
	}
	otpURL, secret, err := domain.EnrollTOTP(u.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "enroll: %v", err)
	}
	if _, err := s.MFA.InsertTOTP(ctx, uid, secret); err != nil {
		return nil, status.Errorf(codes.Internal, "persist: %v", err)
	}
	return &pb.EnrollTOTPResponse{OtpauthUrl: otpURL, Secret: secret}, nil
}

func (s *Server) VerifyMFA(ctx context.Context, req *pb.VerifyMFARequest) (*pb.MFAVerification, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	if req.TotpCode != "" {
		f, err := s.MFA.GetTOTP(ctx, uid)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "no totp enrolled: %v", err)
		}
		if domain.VerifyTOTP(string(f.Secret), req.TotpCode) {
			return &pb.MFAVerification{Verified: true}, nil
		}
	}
	return &pb.MFAVerification{Verified: false}, nil
}

func (s *Server) ListSessions(ctx context.Context, _ *emptypb.Empty) (*pb.ListSessionsResponse, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	rows, err := s.Sessions.List(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list sessions: %v", err)
	}
	resp := &pb.ListSessionsResponse{}
	for _, r := range rows {
		resp.Sessions = append(resp.Sessions, &pb.Session{
			Id:        r.ID.String(),
			CreatedAt: timestamppb.New(r.CreatedAt),
			ExpiresAt: timestamppb.New(r.ExpiresAt),
			UserAgent: r.UserAgent,
			Ip:        r.IP,
		})
	}
	return resp, nil
}

func (s *Server) RevokeSession(ctx context.Context, req *pb.RevokeSessionRequest) (*emptypb.Empty, error) {
	uid, _ := hauth.UserIDFromContext(ctx)
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad session id: %v", err)
	}
	if err := s.Sessions.Revoke(ctx, id, uid); err != nil {
		return nil, status.Errorf(codes.NotFound, "revoke: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ExchangeOIDC(_ context.Context, _ *pb.ExchangeOIDCRequest) (*pb.Tokens, error) {
	return nil, status.Error(codes.Unimplemented, "use GET /v1/auth/callback (REST)")
}
func (s *Server) RefreshToken(_ context.Context, _ *pb.RefreshTokenRequest) (*pb.Tokens, error) {
	return nil, status.Error(codes.Unimplemented, "use POST /v1/auth/refresh (REST)")
}
func (s *Server) EnrollFIDO2(_ context.Context, _ *pb.EnrollFIDO2Request) (*pb.EnrollFIDO2Response, error) {
	return nil, status.Error(codes.Unimplemented, "FIDO2 handled via REST (webauthn needs request context)")
}

func patProto(p *repo.PAT, token string) *pb.PAT {
	out := &pb.PAT{
		Id:        p.ID.String(),
		Name:      p.Name,
		Token:     token,
		Scopes:    p.Scopes,
		CreatedAt: timestamppb.New(p.CreatedAt),
	}
	if p.ExpiresAt != nil {
		out.ExpiresAt = timestamppb.New(*p.ExpiresAt)
	}
	return out
}
