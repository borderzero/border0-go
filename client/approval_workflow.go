package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ApprovalWorkflowService is an interface for API client methods that interact with
// Border0 API to manage approval workflows. Approval workflows are surfaced to users
// as "approval flows".
type ApprovalWorkflowService interface {
	ApprovalWorkflow(ctx context.Context, id string) (out *ApprovalWorkflow, err error)
	ApprovalWorkflows(ctx context.Context) (out *ApprovalWorkflows, err error)
	CreateApprovalWorkflow(ctx context.Context, in *ApprovalWorkflow) (out *ApprovalWorkflow, err error)
	UpdateApprovalWorkflow(ctx context.Context, in *ApprovalWorkflow) (out *ApprovalWorkflow, err error)
	DeleteApprovalWorkflow(ctx context.Context, id string) (err error)
}

// ApprovalWorkflow fetches an approval workflow from your Border0 organization by UUID.
// Approval workflow UUID is globally unique and immutable.
func (api *APIClient) ApprovalWorkflow(ctx context.Context, id string) (out *ApprovalWorkflow, err error) {
	out = new(ApprovalWorkflow)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/approval_workflows/%s", id), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("approval workflow with ID [%s] not found: %w", id, err)
		}
		return nil, err
	}
	return out, nil
}

// ApprovalWorkflows fetches all approval workflows from your Border0 organization.
func (api *APIClient) ApprovalWorkflows(ctx context.Context) (out *ApprovalWorkflows, err error) {
	out = new(ApprovalWorkflows)
	_, err = api.request(ctx, http.MethodGet, "/approval_workflows", nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateApprovalWorkflow creates a new approval workflow in your Border0 organization.
func (api *APIClient) CreateApprovalWorkflow(ctx context.Context, in *ApprovalWorkflow) (out *ApprovalWorkflow, err error) {
	out = new(ApprovalWorkflow)
	_, err = api.request(ctx, http.MethodPost, "/approval_workflows", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateApprovalWorkflow updates an existing approval workflow in your Border0 organization.
func (api *APIClient) UpdateApprovalWorkflow(ctx context.Context, in *ApprovalWorkflow) (out *ApprovalWorkflow, err error) {
	out = new(ApprovalWorkflow)
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/approval_workflows/%s", in.ID), in, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("approval workflow with ID [%s] not found: %w", in.ID, err)
		}
		return nil, err
	}
	return out, nil
}

// DeleteApprovalWorkflow deletes an existing approval workflow from your Border0 organization.
func (api *APIClient) DeleteApprovalWorkflow(ctx context.Context, id string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/approval_workflows/%s", id), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// ApprovalWorkflow represents an approval workflow ("approval flow") in your Border0
// organization. It defines who may request just-in-time access to a set of sockets,
// and who may approve those requests.
type ApprovalWorkflow struct {
	// input fields
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	SocketIDs         []string          `json:"socket_ids"`
	SocketTags        map[string]string `json:"socket_tags"`
	RequesterUserIDs  []string          `json:"requester_user_ids"`
	RequesterGroupIDs []string          `json:"requester_group_ids"`
	ApproverUserIDs   []string          `json:"approver_user_ids"`
	ApproverGroupIDs  []string          `json:"approver_group_ids"`
	AllowSelfApproval bool              `json:"allow_self_approval"`

	// output fields
	ID        string    `json:"approval_workflow_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ApprovalWorkflows represents a list of approval workflows in your Border0 organization.
type ApprovalWorkflows struct {
	List []ApprovalWorkflow `json:"list"`
}
