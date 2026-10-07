package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SSHPermissions_AgentForwardingJSON(t *testing.T) {
	out, err := json.Marshal(SSHPermissions{AgentForwarding: &SSHAgentForwardingPermission{}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"agent_forwarding":{}}`, string(out))

	var in SSHPermissions
	require.NoError(t, json.Unmarshal([]byte(`{"sftp":{},"agent_forwarding":{}}`), &in))
	assert.NotNil(t, in.SFTP)
	assert.NotNil(t, in.AgentForwarding)

	var sftpOnly SSHPermissions
	require.NoError(t, json.Unmarshal([]byte(`{"sftp":{}}`), &sftpOnly))
	assert.Nil(t, sftpOnly.AgentForwarding, "absent key must leave the permission unset")
}
