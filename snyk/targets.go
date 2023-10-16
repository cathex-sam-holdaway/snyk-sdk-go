package snyk

import (
	"context"
	"fmt"
	"net/http"
)

// Use the beta version of the REST API as targets cannot be retrieved from the stable version.
// Limit the targets returned to those imported by the GitHub Enterprise integration so that
// we only see GitHub repos (need to refactor this or create separate functions to return different target types)
// I'm also overriding the default limit of 10 targets per page on the API to 100.
// TODO: Find a way to use a paginated query
const targetBasePath = "orgs/%v/targets?version=2023-09-29~beta&limit=100&origin=github-enterprise"

// TargetsService handles communication with the target related methods of the Snyk REST API.
type TargetsService service

// Target represents a Snyk target.
type Target struct {
	Type       string           `json:"type,omitempty"`
	ID         string           `json:"id,omitempty"`
	Attributes TargetAttributes `json:"attributes,omitempty"`
}

// TargetAttributes represents the attributes of a Snyk target, including the repo details.
type TargetAttributes struct {
	IsPrivate   bool   `json:"isPrivate,omitempty"`
	Origin      string `json:"origin,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	RemoteUrl   string `json:"remoteUrl,omitempty"`
}

// targetsRoot represents the root of the JSON response from the Snyk API.
type targetsRoot struct {
	Targets []Target `json:"data,omitempty"`
}

// List provides a list of all targets (for now those imported using the GitHub Enterprise integration)
// for the given organization.
func (s *TargetsService) List(ctx context.Context, organizationID string) ([]Target, *Response, error) {
	if organizationID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf(targetBasePath, organizationID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(targetsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Targets, resp, nil
}
