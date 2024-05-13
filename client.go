package temporal_tests

import (
	"crypto/tls"
	"log"

	"go.temporal.io/sdk/client"
)

func NewCloudClient(namespace string) (client.Client, error) {
	if namespace == "" {
		namespace = "shahab-test2"
	}
	namespace += ".temporal-dev"

	//cloudHostPort := os.Getenv("TEMPORAL_CLOUD")
	cloudHostPort := "tmprl-test.cloud:7233"

	if cloudHostPort == "" {
		log.Fatalln("TEMPORAL_CLOUD env var is not set")
	}

	// Get the key and cert from your env or local machine
	clientKeyPath := "/Users/shahab/GolandProjects/certs/shahab-test.key"
	clientCertPath := "/Users/shahab/GolandProjects/certs/shahab-test.pem"
	// Specify the host and port of your Temporal Cloud Namespace
	// Host and port format: namespace.unique_id.tmprl.cloud:port
	hostPort := namespace + "." + cloudHostPort
	// Use the crypto/tls package to create a cert object
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalln("Unable to load cert and key pair.", err)
	}
	// Add the cert to the tls certificates in the ConnectionOptions of the Client
	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
		},
	})

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return c, err
}

func NewLocalClient(namespace string) (client.Client, error) {
	c, err := client.Dial(client.Options{
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return c, err
}
