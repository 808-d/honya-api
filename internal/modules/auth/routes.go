package auth

import "github.com/go-chi/chi/v5"

// RegisterRoutes mounts the auth endpoints on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Get("/register/verify", h.VerifyRegister)
	r.Post("/password/forgot", h.ForgotPassword)
	r.Post("/password/reset", h.ResetPassword)
	r.Post("/logout", h.Logout)
	r.Get("/me", h.Me)
}
