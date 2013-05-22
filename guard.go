package crudapi

import (
	"net/url"
)

// A CRUD action. See constants for predefined actions.
type Action string

const (
	ActionCreate    Action = "create"
	ActionGet       Action = "get"
	ActionGetAll    Action = "get all"
	ActionUpdate    Action = "update"
	ActionDelete    Action = "delete"
	ActionDeleteAll Action = "delete all"
)

// A guard authenticates users and authorizes their requests.
type Guard interface {
	// Tries to authenticate a client (e.g. using API keys or signed requests). Returns wether the client could be authenticated, a string which will be passed to Guard.Authorize() and should be used to set levels of permissions, or per-user-permissions, and an error message in case of an error or in case the client could not be authenticated.
	Authenticate(params url.Values) (ok bool, client string, errorMessage string)

	// Tries to authorize the action (one of the Action* constants) to be performed on the kind of resource by the client. Returns wether the action is authorized, and if not, an error message to be sent back to the client.
	Authorize(client string, action Action, kind string) (ok bool, errorMessage string)
}

// default guard; allows everyone to do everything
type defaultGuard struct{}

func (d defaultGuard) Authenticate(params url.Values) (ok bool, client string, errorMessage string) {
	ok = true
	return
}

func (d defaultGuard) Authorize(client string, action Action, kind string) (ok bool, errorMessage string) {
	ok = true
	return
}
