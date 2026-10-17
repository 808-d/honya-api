package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/response"
)

const accessTokenCookie = "access_token"

type Handler struct {
	svc          *Service
	secureCookie bool
}

func NewHandler(svc *Service, secureCookie bool) *Handler {
	return &Handler{svc: svc, secureCookie: secureCookie}
}

// Login handles POST /auth/login. The service panics on failure and the global
// exception handler renders it, so there is no per-call error check here.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decode(r, &req); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	if req.Identifier == "" || req.Password == "" {
		panic(apperror.BadRequest("identifier and password are required"))
	}

	result := h.svc.Login(r.Context(), req)
	h.setAuthCookie(w, result)
	response.OK(w, "login successful", LoginResponse{
		ExpiresAt: result.ExpiresAt,
		User:      result.User,
	})
}

// Register handles POST /auth/register. It creates an unverified account and
// sends a verification link by email.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decode(r, &req); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Register(r.Context(), req)
	response.Created(w, "registration successful, check your email to verify", result)
}

// VerifyRegister handles GET /auth/register/verify?id=&token=. It consumes the
// token from Redis, activates the account and logs the user in.
func (h *Handler) VerifyRegister(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")

	result := h.svc.VerifyRegister(r.Context(), id, token)
	h.setAuthCookie(w, result)
	response.OK(w, "email verified", LoginResponse{
		ExpiresAt: result.ExpiresAt,
		User:      result.User,
	})
}

// ForgotPassword handles POST /auth/password/forgot.
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := decode(r, &req); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	h.svc.ForgotPassword(r.Context(), req)
	response.OK(w, "if the email is registered, a reset link has been sent", nil)
}

// ResetPassword handles POST /auth/password/reset.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := decode(r, &req); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	h.svc.ResetPassword(r.Context(), req)
	response.OK(w, "password has been reset", nil)
}

func (h *Handler) setAuthCookie(w http.ResponseWriter, result *LoginResult) {
	expiresAt := time.Unix(result.ExpiresAt, 0)
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    result.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
}

func decode(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
