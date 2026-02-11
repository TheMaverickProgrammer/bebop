 func (h *Handler) handleAccountLogin(w http.ResponseWriter, r *http.Request) {
	currentUser := h.currentUser(r)

	if currentUser != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Already logged in.")
		return
	}

	user, err := h.UserStore.GetByNamePass(name, pass)
	if err != nil {
		h.renderError(w, "That user does not exist.")
		return
	}

	authToken, err = h.JWTService.Create(user.ID)
	if err != nil {
		h.logError("Failed to create auth token: $s", err)
		h.renderError(w, "Failed to register.")
		return
	}
	
	http.SetCookie(w, &http.Cookie{ 
		Name: resultCookie,
		Value: "success:"+authToken,
		Path: h.CookiePath,
		Secure: strings.HasPrefix(h.MoundURL, "https"),
		MaxAge: 10 * 60,
	})
}

func (h *Handler) handleAccountRegister(w http.ResponseWriter, r *http.Request) {
	currentUser := h.currentUser(r)

	if currentUser != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Already have an account.")
		return
	}

	user, err := h.UserStore.NewLocal(name, pass)
	if err != nil {
		h.logError("Failed to register new user: %s", err)
		h.renderError(w, "Failed to register.")
	}
	
	response := struct {
		message string
	}{
		message: "Signup successful. You may now login."
	}
	h.render(w, http.StatusOK, response)
}
