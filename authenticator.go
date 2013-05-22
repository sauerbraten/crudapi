package crudapi

import (
	"net/url"
)

const (
	ActionCreate    = "create"
	ActionGet       = "get"
	ActionGetAll    = "getAll"
	ActionUpdate    = "update"
	ActionDelete    = "delete"
	ActionDeleteAll = "deleteAll"
)

// An Response is returned by Storage's methods. It sets the HTTP status code of the response and describes what kind of error occured, if any.
type AuthenticatorResponse struct {
	Allowed      bool   // true if the client is allowed to perform the action
	ErrorMessage string // the error, if any
}

// Tries to authenticate a client's request using the action it wants to perform, which is one of the Action* constants, the kind of resource it wants to perform the action on, and the url parameters (e.g. using API keys or signed requests).
type AuthenticateFunction func(action string, kind string, params url.Values) AuthenticatorResponse

// default
func allowAll(action string, kind string, params url.Values) AuthenticatorResponse {
	return AuthenticatorResponse{true, ""}
}
