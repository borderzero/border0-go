package connector

// Metadata represents informational data about a connector.
type Metadata struct {
	AwsEc2IdentityMetadata    *AwsEc2IdentityMetadata    `json:"aws_ec2_identity_metadata,omitempty"`
	ConnectorInternalMetadata *ConnectorInternalMetadata `json:"connector_internal_metadata,omitempty"`
}

// AwsEc2IdentityMetadata represents metadata for connectors running on AWS EC2 instances.
type AwsEc2IdentityMetadata struct {
	AwsAccountId        string `json:"aws_account_id,omitempty"`
	AwsRegion           string `json:"aws_region,omitempty"`
	AwsAvailabilityZone string `json:"aws_availability_zone,omitempty"`
	Ec2InstanceId       string `json:"ec2_instance_id,omitempty"`
	Ec2InstanceType     string `json:"ec2_instance_type,omitempty"`
	Ec2ImageId          string `json:"ec2_image_id,omitempty"`
	KernelId            string `json:"kernel_id,omitempty"`
	RamdiskId           string `json:"ramdisk_id,omitempty"`
	Architecture        string `json:"architecture,omitempty"`
	PrivateIpAddress    string `json:"private_ip_address,omitempty"`
}

// ConnectorInternalMetadata represents metadata for connector internal data. This includes the version
// of the connector, the date the connector was built, the IP address of the connector, and the metadata
// for the IP address.
type ConnectorInternalMetadata struct {
	Version    string      `json:"version,omitempty"`
	BuiltDate  string      `json:"built_date,omitempty"`
	IPAddress  string      `json:"ip_address,omitempty"`
	IPMetadata *IPMetadata `json:"ip_metadata,omitempty"`
}

// IPMetadata represents metadata for an IP address. This includes the country name, country code, region
// name, region code, city name, latitude, longitude, and ISP. This data is retrieved from IP Geolocation
// database with the IP address of the connector.
type IPMetadata struct {
	CountryName string  `json:"country_name,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
	RegionName  string  `json:"region_name,omitempty"`
	RegionCode  string  `json:"region_code,omitempty"`
	CityName    string  `json:"city_name,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	ISP         string  `json:"isp,omitempty"`
}
