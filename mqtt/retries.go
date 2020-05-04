package mqtt

type Retries interface {
	addRetry(retry Retry) error
	remRetry(clientId string, packetIdentifier int) error
	findRetriesByClientId(clientId string, packetIdentifier int) []Retry
}
