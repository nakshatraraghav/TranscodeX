package types

type RequestBodyContextKey string

type AuthedUserContextKey string

var ContextKey RequestBodyContextKey = "body"

var AuthContextKey AuthedUserContextKey = "authuser"
