package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// PolicyService is an interface for API client methods that interact with Border0 API to manage policies and policy socket attachments.
type PolicyService interface {
	Policy(ctx context.Context, id string) (out *Policy, err error)
	Policies(ctx context.Context) (out []Policy, err error)
	PoliciesByNames(ctx context.Context, names ...string) (out []Policy, err error)
	CreatePolicy(ctx context.Context, in *Policy) (out *Policy, err error)
	UpdatePolicy(ctx context.Context, id string, in *Policy) (out *Policy, err error)
	DeletePolicy(ctx context.Context, id string) (err error)
	AttachPolicyToSocket(ctx context.Context, policyID string, socketID string) (err error)
	RemovePolicyFromSocket(ctx context.Context, policyID string, socketID string) (err error)
	AttachPoliciesToSocket(ctx context.Context, policyIDs []string, socketID string) (err error)
	RemovePoliciesFromSocket(ctx context.Context, policyIDs []string, socketID string) (err error)
}

// Policy fetches a policy from your Border0 organization by policy ID.
func (api *APIClient) Policy(ctx context.Context, id string) (out *Policy, err error) {
	out = new(Policy)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/policy/%s", id), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("policy [%s] not found: %w", id, err)
		}
		return nil, err
	}
	return out, nil
}

// Policies fetches all policies in your Border0 organization.
func (api *APIClient) Policies(ctx context.Context) (out []Policy, err error) {
	_, err = api.request(ctx, http.MethodGet, "/policies", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PoliciesByNames finds policies in your Border0 organization by policy names. If any of the policies does not exist,
// an error will be returned. When only one policy name is provided, this method will use the /policies/find endpoint,
// otherwise it will fetch all policies and filter them by name.
func (api *APIClient) PoliciesByNames(ctx context.Context, names ...string) (out []Policy, err error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no policy names provided")
	}

	if len(names) == 1 {
		var found Policy
		_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/policies/find?name=%s", names[0]), nil, &found)
		if err != nil {
			if NotFound(err) {
				return nil, fmt.Errorf("policy [%s] does not exist, please create the policy first", names[0])
			}
			return nil, err
		}
		return []Policy{found}, nil
	}

	policies, err := api.Policies(ctx)
	if err != nil {
		return nil, err
	}
	policiesMap := make(map[string]Policy)
	for _, policy := range policies {
		policiesMap[policy.Name] = policy
	}
	for _, name := range names {
		policy, ok := policiesMap[name]
		if !ok {
			return nil, fmt.Errorf("policy [%s] does not exist, please create the policy first", name)
		}
		out = append(out, policy)
	}
	return out, nil
}

// CreatePolicy creates a new policy in your Border0 organization. Policy name must be unique within your organization,
// otherwise API will return an error. Policy name must contain only lowercase letters, numbers and dashes.
func (api *APIClient) CreatePolicy(ctx context.Context, in *Policy) (out *Policy, err error) {
	out = new(Policy)
	_, err = api.request(ctx, http.MethodPost, "/policies", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdatePolicy updates an existing policy in your Border0 organization.
func (api *APIClient) UpdatePolicy(ctx context.Context, id string, in *Policy) (out *Policy, err error) {
	out = new(Policy)
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/policy/%s", id), in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeletePolicy deletes a policy from your Border0 organization by policy ID.
func (api *APIClient) DeletePolicy(ctx context.Context, id string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/policy/%s", id), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// AttachPolicyToSocket attaches a policy to a socket by policy ID and socket ID.
func (api *APIClient) AttachPolicyToSocket(ctx context.Context, policyID string, socketID string) (err error) {
	in := PolicySocketAttachments{
		Actions: []PolicySocketAttachment{
			{Action: "add", ID: socketID},
		},
	}
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/policy/%s/socket", policyID), &in, nil)
	return err
}

// RemovePolicyFromSocket detaches a policy from a socket with policy ID and socket ID.
func (api *APIClient) RemovePolicyFromSocket(ctx context.Context, policyID string, socketID string) (err error) {
	in := PolicySocketAttachments{
		Actions: []PolicySocketAttachment{
			{Action: "remove", ID: socketID},
		},
	}
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/policy/%s/socket", policyID), &in, nil)
	return err
}

// AttachPoliciesToSocket attaches multiple policies to a socket by policy IDs and socket ID.
func (api *APIClient) AttachPoliciesToSocket(ctx context.Context, policyIDs []string, socketID string) (err error) {
	in := PolicySocketAttachments{
		Actions: []PolicySocketAttachment{},
	}
	for _, policyID := range policyIDs {
		in.Actions = append(in.Actions, PolicySocketAttachment{Action: "add", ID: policyID})
	}
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/socket/%s/policy", socketID), &in, nil)
	return err
}

// RemovePoliciesFromSocket detaches multiple policies from a socket by policy IDs and socket ID.
func (api *APIClient) RemovePoliciesFromSocket(ctx context.Context, policyIDs []string, socketID string) (err error) {
	in := PolicySocketAttachments{
		Actions: []PolicySocketAttachment{},
	}
	for _, policyID := range policyIDs {
		in.Actions = append(in.Actions, PolicySocketAttachment{Action: "remove", ID: policyID})
	}
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/socket/%s/policy", socketID), &in, nil)
	return err
}

// Policy represents a Border0 policy in your organization. See [PolicyData] for more details about the policy data schema.
// A policy can be set to be organization-wide, in which case it will be applied to all sockets in your organization. If
// a policy is not organization-wide, it can be attached to individual sockets. See [AttachPolicyToSocket] and [RemovePolicyFromSocket]
// for more details.
type Policy struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrgID       string     `json:"org_id"`
	OrgWide     bool       `json:"org_wide"`
	PolicyData  PolicyData `json:"policy_data"`
	CreatedAt   time.Time  `json:"created_at"`
	SocketIDs   []string   `json:"socket_ids"`
	Deleted     bool       `json:"deleted"`
}

// PolicyData represents the policy data schema. A policy can have multiple actions, and its condition determines when the
// actions are applied. See [PolicyCondition] for more details about the policy condition schema.
type PolicyData struct {
	Version   string          `json:"version"`
	Action    []string        `json:"action"`
	Condition PolicyCondition `json:"condition"`
}

// PolicyCondition represents the policy condition schema. A policy condition can define "who", "where" and "when" conditions.
// See [PolicyWho], [PolicyWhere] and [PolicyWhen] for more details about the policy condition schema.
type PolicyCondition struct {
	Who   PolicyWho   `json:"who,omitempty"`
	Where PolicyWhere `json:"where,omitempty"`
	When  PolicyWhen  `json:"when,omitempty"`
}

// PolicyWho represents the policy condition "who" schema. It specifies who the policy applies to, based on allowed email
// addresses and allowed email domains.
type PolicyWho struct {
	Email  []string `json:"email,omitempty"`
	Domain []string `json:"domain,omitempty"`
	Group  []string `json:"group,omitempty"`
}

// PolicyWhere represents the policy condition "where" schema. It specifies where the policy applies to, based on allowed
// IP addresses, allowed countries and countries not allowed.
type PolicyWhere struct {
	AllowedIP  []string `json:"allowed_ip,omitempty"`
	Country    []string `json:"country,omitempty"`
	CountryNot []string `json:"country_not,omitempty"`
}

// PolicyWhen represents the policy condition "when" schema. It specifies when the policy applies to, based on allowed
// dates and allowed times of day.
type PolicyWhen struct {
	After           string `json:"after,omitempty"`
	Before          string `json:"before,omitempty"`
	TimeOfDayAfter  string `json:"time_of_day_after,omitempty"`
	TimeOfDayBefore string `json:"time_of_day_before,omitempty"`
}

// PolicySocketAttachments represents a list of policy socket attachments. Border0 API client uses this schema to attach
// or detach a policy to/from a socket.
type PolicySocketAttachments struct {
	Actions []PolicySocketAttachment `json:"actions"`
}

// PolicySocketAttachment represents a single policy socket attachment. The action can be "add" or "remove", and the ID
// is the socket ID.
type PolicySocketAttachment struct {
	Action string `json:"action" binding:"required"`
	ID     string `json:"id" binding:"required"`
}
