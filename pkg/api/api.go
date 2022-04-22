package api

// Credentials is used to store credentials, including the network and component
type Credentials struct {
	// Nats Address
	NatsAddress string `json:"natsAddress"`
	// Component is the network participant
	Component string `json:"component"`
	// Network is the network
	Network string `json:"network"`
	// Creds is the credentials for the network
	Creds string `json:"creds"`
	// AccountPublicKey is the public key of the account the user belongs to
	AccountPublicKey string `json:"accountPublicKey"`
}
