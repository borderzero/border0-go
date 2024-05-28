package client

import (
	"context"
	"fmt"
	"net/http"
)

// GroupService is an interface for API client methods that interact with Border0 API to manage groups.
type GroupService interface {
	Group(ctx context.Context, id string) (out *Group, err error)
	CreateGroup(ctx context.Context, in *Group) (out *Group, err error)
	UpdateGroup(ctx context.Context, in *Group) (out *Group, err error)
	DeleteGroup(ctx context.Context, id string) (err error)
}

// Group fetches a group from your Border0 organization by UUID. Group UUID is globally unique and immutable.
func (api *APIClient) Group(ctx context.Context, id string) (out *Group, err error) {
	out = new(Group)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/organizations/iam/groups/%s", id), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("group with ID [%s] not found: %w", id, err)
		}
		return nil, err
	}
	return out, nil
}

// CreateGroup creates a new group in your Border0 organization.
func (api *APIClient) CreateGroup(ctx context.Context, in *Group) (out *Group, err error) {
	out = new(Group)
	_, err = api.request(ctx, http.MethodPost, "/organizations/iam/groups", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateGroup updates an existing group in your Border0 organization.
func (api *APIClient) UpdateGroup(ctx context.Context, in *Group) (out *Group, err error) {
	out = new(Group)
	_, err = api.request(ctx, http.MethodPut, "/organizations/iam/groups", in, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("group with ID [%s] not found: %w", in.ID, err)
		}
		return nil, err
	}
	return out, nil
}

// DeleteGroup deletes an existing group from your Border0 organization.
func (api *APIClient) DeleteGroup(ctx context.Context, id string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/organizations/iam/groups/%s", id), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// Group represents a group in your Border0 organization.
type Group struct {
	// input fields
	DisplayName string `json:"display_name"`

	// output fields
	ID               string            `json:"id"`
	GroupType        string            `json:"group_type"`
	DirectoryService *DirectoryService `json:"directory_service,omitempty"`
	Members          []User            `json:"members,omitempty"`
}
