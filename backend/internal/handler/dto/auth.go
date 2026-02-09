package dto

// LdapLoginRequest represents LDAP login request payload
type LdapLoginRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	TurnstileToken string `json:"turnstile_token"`
}
