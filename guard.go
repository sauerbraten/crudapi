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

// An GuardResponse is returned by the AuthenticateAndAuthorize method. It describes wether the client could be authenticated, the request could be authorized, and what kind of error occured, if any.
type GuardResponse struct {
	Authenticated bool   // true if the client could be authenticated (e.g. the API key is valid or the signed request checked out)
	Allowed       bool   // true if the client is allowed to perform the action
	ErrorMessage  string // the error, if any
}

// A guard authenticates users and authorizes their requests.
type Guard interface {
	// Tries to authenticate a client and authorize their request using the action (one of the Action* constants) they want to perform, the kind of resource it wants to perform the action on, and the url parameters (e.g. using API keys or signed requests).
	AuthenticateAndAuthorize(action string, kind string, params url.Values) GuardResponse
}

// default guard; allows everyone to do everything
type defaultGuard struct{}

func (d defaultGuard) AuthenticateAndAuthorize(action string, kind string, params url.Values) GuardResponse {
	return GuardResponse{true, true, ""}
}
