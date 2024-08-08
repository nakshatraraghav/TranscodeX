package util

import "net/http"

func ClearCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,                 // The name of the cookie to clear
		Value:    "",                   // Set an empty value
		Path:     "/",                  // The path of the cookie
		Domain:   "",                   // Domain if needed
		MaxAge:   -1,                   // Set max age to -1 to expire the cookie
		Secure:   true,                 // Set true if using HTTPS
		HttpOnly: true,                 // Set true if the cookie should be accessible only via HTTP(S)
		SameSite: http.SameSiteLaxMode, // Adjust based on your needs
	}
}
