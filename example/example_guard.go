package crudapi

import (
	"net/url"
)

type MapGuard struct {
	ValidActions map[string][]string
}

// no authentication; only valid actions are authorized
func (a MapGuard) AuthenticateAndAuthorize(action string, kind string, params url.Values) GuardResponse {
	for _, validAction := range a.ValidActions[kind] {
		if validAction == action {
			return GuardResponse{true, true, ""}
		}
	}

	return GuardResponse{true, false, "action not allowed for this kind of resource"}
}
