package client

import (
	"context"
	"fmt"
	"net/http"
)

// UserService is an interface for API client methods that interact with Border0 API to manage users.
type UserService interface {
	User(ctx context.Context, id string) (out *User, err error)
	Users(ctx context.Context) (out *Users, err error)
	CreateUser(ctx context.Context, in *User, opts ...UserOption) (out *User, err error)
	UpdateUser(ctx context.Context, in *User) (out *User, err error)
	DeleteUser(ctx context.Context, id string) (err error)
}

type userConfig struct {
	SkipNotification bool
}

// UserOption is a user creation option.
type UserOption func(*userConfig)

// WithSkipNotification is the UserOption to skip sending emails to notify added users of their addition.
func WithSkipNotification(skip bool) UserOption {
	return func(uc *userConfig) { uc.SkipNotification = skip }
}

// User fetches a user from your Border0 organization by UUID. User UUID is globally unique and immutable.
func (api *APIClient) User(ctx context.Context, id string) (out *User, err error) {
	out = new(User)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/organizations/iam/users/%s", id), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("user with ID [%s] not found: %w", id, err)
		}
		return nil, err
	}
	return out, nil
}

// Users fetches all users from your Border0 organization.
func (api *APIClient) Users(ctx context.Context) (out *Users, err error) {
	out = new(Users)
	_, err = api.request(ctx, http.MethodGet, "/organizations/iam/users", nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateUser creates a new user in your Border0 organization. User email must
// be unique within your organization, otherwise API will return an error.
func (api *APIClient) CreateUser(ctx context.Context, in *User, opts ...UserOption) (out *User, err error) {
	userConfig := &userConfig{}
	for _, opt := range opts {
		opt(userConfig)
	}
	url := "/organizations/iam/users"
	if userConfig.SkipNotification {
		url = fmt.Sprintf("%s?skip_notification=true", url)
	}
	out = new(User)
	_, err = api.request(ctx, http.MethodPost, url, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateUser updates an existing user in your Border0 organization.
func (api *APIClient) UpdateUser(ctx context.Context, in *User) (out *User, err error) {
	out = new(User)
	_, err = api.request(ctx, http.MethodPut, "/organizations/iam/users", in, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("user with ID [%s] not found: %w", in.ID, err)
		}
		return nil, err
	}
	return out, nil
}

// DeleteUser deletes an existing user from your Border0 organization.
func (api *APIClient) DeleteUser(ctx context.Context, id string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/organizations/iam/users/%s", id), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// User represents a user in your Border0 organization.
type User struct {
	// input and output fields
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`

	// output field
	ID               string            `json:"id"`
	UserType         string            `json:"user_type"`
	DirectoryService *DirectoryService `json:"directory_service,omitempty"`
}

// Users represents a list of users in your Border0 organization.
type Users struct {
	List []User `json:"list"`
}

// DirectoryService represents a directory service in your Border0 organization.
type DirectoryService struct {
	// input fields
	DisplayName string `json:"display_name"`
	ServiceType string `json:"service_type"`

	// output fields
	ID string `json:"id"`
}
