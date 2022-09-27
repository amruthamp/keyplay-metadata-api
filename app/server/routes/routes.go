package routes

const (
	Healthcheck           = "/keyplay-metadata-api/healthcheck"
	SampleRoute           = "/keyplay-metadata-api/sample"
	Swagger               = "/keyplay-metadata-api/swaggerui"
	GetKeyplayMetadata    = "/keyplay/metadata/:programId/:channelId/:id"
	CreateKeyplayMetadata = "/keyplay/metadata/:programId/:channelId/:id"
	UpdateKeyplayMetadata = "/keyplay/metadata/:programId/:channelId/:id"
	DeleteKeyplayMetadata = "/keyplay/metadata/:programId/:channelId/:id"
)
