package service

const (
	// HttpServiceTypeStandard is the http
	// service type for standard http services.
	HttpServiceTypeStandard = "standard"

	// HttpServiceTypeConnectorFileServer is the http service
	// type for the connector's built-in file webserver.
	HttpServiceTypeConnectorFileServer = "connector_file_server"
)

// HttpServiceConfiguration represents service
// configuration for http services (fka sockets).
type HttpServiceConfiguration struct {
	HttpServiceType string `json:"http_service_type"`

	// mutually exclusive fields below
	StandardHttpServiceConfiguration   *StandardHttpServiceConfiguration   `json:"standard_http_service_configuration,omitempty"`
	FileServerHttpServiceConfiguration *FileServerHttpServiceConfiguration `json:"fileserver_http_service_configuration,omitempty"`
}

// StandardHttpServiceConfiguration represents service
// configuration for standard http services (fka sockets).
type StandardHttpServiceConfiguration struct {
	HostnameAndPort // inherited
	HostSniHeader   string
}

// FileServerHttpServiceConfiguration represents service
// configuration for the connector built-in file webserver.
type FileServerHttpServiceConfiguration struct {
	TopLevelDirectory string `json:"top_level_directory,omitempty"`
}

// Validate validates the HttpServiceConfiguration.
func (c *HttpServiceConfiguration) Validate() error {
	// TODO
	return nil
}
