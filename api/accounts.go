package api

import (
	"net/http"
	"strings"
	"time"
)

const (
	stateCookie   = "bebop_oauth_state"
	resultCookie  = "bebop_oauth_result"
	clientTimeout = 10 * time.Second
)

func (h *Handler) handleAccountLogin(w http.ResponseWriter, r *http.Request) {
	currentUser := h.currentUser(r)

	if currentUser != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Already logged in.")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || strings.ToUpper(authHeader[:6]) != "BEARER" {
		return
	}
	parts := strings.Split(authHeader[7:], ":")
	name := parts[0]
	pass := parts[1]

	user, err := h.Store.Users().GetByNamePass(name, pass)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "That user does not exist.")
		return
	}

	authToken, err := h.JWTService.Create(user.ID)
	if err != nil {
		h.logError("Failed to create auth token: $s", err)
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Failed to register.")
		return
	}
	
	http.SetCookie(w, &http.Cookie{ 
		Name: resultCookie,
		Value: "success:"+authToken,
		Path: h.Config.CookiePath,
		Secure: strings.HasPrefix(h.Config.MountURL, "https"),
		MaxAge: 10 * 60,
	})
}

func (h *Handler) handleAccountRegister(w http.ResponseWriter, r *http.Request) {
	currentUser := h.currentUser(r)

	if currentUser != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Already have an account.")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || strings.ToUpper(authHeader[:6]) != "BEARER" {
		return
	}
	parts := strings.Split(authHeader[7:], ":")
	name := parts[0]
	pass := parts[1]
	_, err := h.Store.Users().NewLocal(name, pass)
	if err != nil {
		h.logError("Failed to register new user: %s", err)
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Failed to register.")
	}
	
	response := struct {
		message string
	}{
		message: "Signup successful. You may now login.",
	}
	h.render(w, http.StatusOK, response)
}
