package main

import (
	"github.com/sauerbraten/crudapi"
	"net/url"
)

// no authentication; only valid actions are authorized
type MapGuard struct {
	ValidActions map[string][]crudapi.Action
}

func (mg MapGuard) Authenticate(params url.Values) crudapi.AuthenticationResponse {
	return crudapi.AuthenticationResponse{true, "", ""}
}

func (mg MapGuard) Authorize(client string, action crudapi.Action, kind string) crudapi.AuthorizationResponse {
	for _, validAction := range mg.ValidActions[kind] {
		if validAction == action {
			return crudapi.AuthorizationResponse{true, true, ""}
		}
	}

	return crudapi.GuardResponse{true, false, "action not allowed for this kind of resource"}
}
