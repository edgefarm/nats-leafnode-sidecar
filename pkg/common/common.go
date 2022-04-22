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
	CredentialsFile = "/creds/nats-sidecar.creds"
)

var (
	// Remote is the address of the remote nats server
	Remote string
)
