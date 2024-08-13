package types

type RequestBodyContextKey string

type AuthedUserContextKey string

type ApiKeyContextKey string

var ContextKey RequestBodyContextKey = "body"

var AuthContextKey AuthedUserContextKey = "authuser"

var ApiContextKey ApiKeyContextKey = "apikey"
