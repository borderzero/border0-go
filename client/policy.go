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
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	OrgID       string    `json:"org_id"`
	OrgWide     bool      `json:"org_wide"`
	PolicyData  any       `json:"policy_data"`
	CreatedAt   time.Time `json:"created_at"`
	SocketIDs   []string  `json:"socket_ids"`
	Deleted     bool      `json:"deleted"`
}

// PolicyData represents the policy data schema for v1 policies. A policy can have multiple actions, and its condition determines when the
// actions are applied. See [PolicyCondition] for more details about the policy condition schema.
type PolicyData struct {
	Version   string          `json:"version,omitempty"`
	Action    []string        `json:"action"`
	Condition PolicyCondition `json:"condition"`
}

// PolicyDataV2 represents the policy data schema for v2 policies. A policy can have multiple actions, and its condition determines when the
// actions are applied. See [PolicyCondition] for more details about the policy condition schema.
type PolicyDataV2 struct {
	Permissions PolicyPermissions `json:"permissions"`
	Condition   PolicyConditionV2 `json:"condition"`
}

// PolicyCondition represents the policy condition schema for v1 policies. A policy condition can define "who", "where" and "when" conditions.
// See [PolicyWho], [PolicyWhere] and [PolicyWhen] for more details about the policy condition schema.
type PolicyCondition struct {
	Who   PolicyWho   `json:"who,omitempty"`
	Where PolicyWhere `json:"where,omitempty"`
	When  PolicyWhen  `json:"when,omitempty"`
}

// PolicyConditionV2 represents the policy condition schema for v2 policies. A policy condition can define "who", "where" and "when" conditions.
// See [PolicyWho], [PolicyWhere] and [PolicyWhen] for more details about the policy condition schema.
type PolicyConditionV2 struct {
	Who   PolicyWhoV2 `json:"who,omitempty"`
	Where PolicyWhere `json:"where,omitempty"`
	When  PolicyWhen  `json:"when,omitempty"`
}

// PolicyWho represents the policy condition "who" schema. It specifies who the policy applies to, based on allowed email
// addresses, domains, groups and service accounts.
type PolicyWho struct {
	Email          []string `json:"email,omitempty"`
	Domain         []string `json:"domain,omitempty"`
	Group          []string `json:"group,omitempty"`
	ServiceAccount []string `json:"service_account,omitempty"`
}

// PolicyWhoV2 represents the policy condition "who" schema, for v2 policies. It specifies who the policy applies to, based on allowed email
// addresses, groups and service accounts.
type PolicyWhoV2 struct {
	Email          []string `json:"email,omitempty"`
	Group          []string `json:"group,omitempty"`
	ServiceAccount []string `json:"service_account,omitempty"`
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

// PolicyPermissions represents permissions for policy (v2).
type PolicyPermissions struct {
	Database   *DatabasePermissions   `json:"database,omitempty"`
	SSH        *SSHPermissions        `json:"ssh,omitempty"`
	HTTP       *HTTPPermissions       `json:"http,omitempty"`
	TLS        *TLSPermissions        `json:"tls,omitempty"`
	VNC        *VNCPermissions        `json:"vnc,omitempty"`
	RDP        *RDPPermissions        `json:"rdp,omitempty"`
	VPN        *VPNPermissions        `json:"vpn,omitempty"`
	Kubernetes *KubernetesPermissions `json:"kubernetes,omitempty"`
}

// DatabasePermissions represents database permissions for policy (v2).
type DatabasePermissions struct {
	AllowedDatabases          *[]DatabasePermission `json:"allowed_databases,omitempty"`
	MaxSessionDurationSeconds *int                  `json:"max_session_duration_seconds,omitempty"`
}

// DatabasePermission represents a single database permission for policy (v2).
type DatabasePermission struct {
	Database          string    `json:"database"`
	AllowedQueryTypes *[]string `json:"allowed_query_types,omitempty"`
}

// SSHPermissions represents ssh service permissions for policy (v2).
type SSHPermissions struct {
	Shell                     *SSHShellPermission         `json:"shell,omitempty"`
	Exec                      *SSHExecPermission          `json:"exec,omitempty"`
	SFTP                      *SSHSFTPPermission          `json:"sftp,omitempty"`
	TCPForwarding             *SSHTCPForwardingPermission `json:"tcp_forwarding,omitempty"`
	KubectlExec               *SSHKubectlExecPermission   `json:"kubectl_exec,omitempty"`
	DockerExec                *SSHDockerExecPermission    `json:"docker_exec,omitempty"`
	MaxSessionDurationSeconds *int                        `json:"max_session_duration_seconds,omitempty"`
	AllowedUsernames          *[]string                   `json:"allowed_usernames,omitempty"`
}

// SSHShellPermission represents the shell ssh permission for policy (v2).
type SSHShellPermission struct{}

// SSHExecPermission represents the exec ssh permission for policy (v2).
type SSHExecPermission struct {
	Commands *[]string `json:"commands,omitempty"`
}

// SSHSFTPPermission represents the sftp ssh permission for policy (v2).
type SSHSFTPPermission struct{}

// SSHTCPForwardingPermission represents the tcp forwarding ssh permission for policy (v2).
type SSHTCPForwardingPermission struct {
	AllowedConnections *[]SSHTcpForwardingConnection `json:"allowed_connections,omitempty"`
}

// SSHTcpForwardingConnection represents data regarding a tcp forwarding ssh permission for policy (v2).
type SSHTcpForwardingConnection struct {
	DestinationAddress *string `json:"destination_address,omitempty"`
	DestinationPort    *string `json:"destination_port,omitempty"`
}

// SSHKubectlExecPermission represents the kubectl exec ssh permission for policy (v2).
type SSHKubectlExecPermission struct {
	AllowedNamespaces *[]KubectlExecNamespace `json:"allowed_namespaces,omitempty"`
}

// KubectlExecNamespace represents a single namespace and pod selector for a kubectl exec ssh permission for policy (v2).
type KubectlExecNamespace struct {
	Namespace   string             `json:"namespace"`
	PodSelector *map[string]string `json:"pod_selector,omitempty"`
}

// SSHDockerExecPermission represents the docker exec ssh permission for policy (v2).
type SSHDockerExecPermission struct {
	AllowedContainers *[]string `json:"allowed_containers,omitempty"`
}

// HTTPPermissions represents http service permissions for policy (v2).
type HTTPPermissions struct{}

// TLSPermissions represents tls service permissions for policy (v2).
type TLSPermissions struct{}

// VNCPermissions represents vnc service permissions for policy (v2).
type VNCPermissions struct{}

// RDPPermissions represents rdp service permissions for policy (v2).
type RDPPermissions struct{}

// VPNPermissions represents vpn service permissions for policy (v2).
type VPNPermissions struct{}

// KubernetesPermissions represents kubernetes service permissions for policy (v2).
type KubernetesPermissions struct {
	Rules *[]KubernetesRule `json:"rules,omitempty"`
}

// KubernetesRule represents a single kubernetes rule for kubernetes service permissions for policy (v2).
type KubernetesRule struct {
	APIGroups     []string `json:"api_groups,omitempty"`
	Namespaces    []string `json:"namespaces,omitempty"`
	Verbs         []string `json:"verbs,omitempty"`
	Resources     []string `json:"resources,omitempty"`
	ResourceNames []string `json:"resource_names,omitempty"`
	// Note: support for selectors will come later...
}
