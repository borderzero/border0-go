package recordings

import "github.com/borderzero/border0-go/lib/types/set"

const (
	// RecordingTypeAsciinema is the recording type for recordings in https://asciinema.org/ format.
	RecordingTypeAsciinema = "asciinema"

	// RecordingTypeKubernetesAPIRequestLog is the recording type for kubernetes api events in table format.
	RecordingTypeKubernetesAPIRequestLog = "kubernetes_api_request_log"

	// RecordingTypeDatabaseQueryLog is the recording type for database queries in table format.
	RecordingTypeDatabaseQueryLog = "database_query_log"

	// RecordingTypeBrowserSessionStream is the recording type for http DOM snapshot recordings e.g. rrweb.
	RecordingTypeBrowserSessionStream = "browser_session_stream"
)

// ValidRecordingTypes represents allowed values for recording types.
var ValidRecordingTypes = set.New(
	RecordingTypeAsciinema,
	RecordingTypeKubernetesAPIRequestLog,
	RecordingTypeDatabaseQueryLog,
	RecordingTypeBrowserSessionStream,
)
