package main

import (
	"github.com/sauerbraten/crudapi"
	"net/url"
)

// no authentication; only valid actions are authorized
type MapGuard struct {
	ValidActions map[string][]crudapi.Action
}

func (mg MapGuard) Authenticate(params url.Values) (ok bool, client string, errorMessage string) {
	ok = true
	return
}

func (mg MapGuard) Authorize(client string, action crudapi.Action, vars map[string]string) (ok bool, errorMessage string) {
	for _, validAction := range mg.ValidActions[vars["kind"]] {
		if validAction == action {
			ok = true
			return
		}
	}

	errorMessage = "action not allowed for this kind of resource"
	return
}
