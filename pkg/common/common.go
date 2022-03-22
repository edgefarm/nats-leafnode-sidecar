package common

const (
	// RegisterSubject is the subject used to register a new user
	RegisterSubject = "register"
	// UnregisterSubject is the subject used to unregister a user
	UnregisterSubject = "unregister"
	// OkResponse is the response sent when a request is successful
	OkResponse = "ok"
	// ErrorResponse is the response sent when a request is unsuccessful
	ErrorResponse = "error"
	// CredentialsFile is the name of the file containing the credentials
	CredentialsFile = "/creds/edgefarm-sys.creds"

	// FixedNatsUser is the user used to connect to the NATS server
	FixedNatsUser = "nats-sidecar"
	// FixedNatsPassword is the user used to connect to the NATS server
	FixedNatsPassword = "nats-sidecar"
)
