package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/minio/minio-go/v6"
)

// Config represents the s3 client connection string information.
type Config struct {
	server       string
	region       string
	accessKey    string
	secretKey    string
	apiSignature string
	ssl          bool
	insecure     bool
	issuerPEM    string
	debug        bool
}

type s3Client struct {
	region   string
	debug    bool
	s3Client *minio.Client
}

func (c *Config) getRoundTripper() http.RoundTripper {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.insecure,
	}

	if c.issuerPEM != "" {

		// Get the SystemCertPool, continue with an empty pool on error
		log.Printf("[DEBUG] loading system cert pool.")
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			log.Printf("[DEBUG] system cert pool not found.  Creating empty cert pool.")
			rootCAs = x509.NewCertPool()
		}

		// Read in the cert file

		certs, err := ioutil.ReadFile(c.issuerPEM)
		if err != nil {
			log.Fatalf("Failed to append %q to RootCAs: %v", c.issuerPEM, err)
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Println("No certs appended, using system certs only")
		}

		tlsConfig.RootCAs = rootCAs
	}

	var h http.RoundTripper = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return h
}

// NewClient generates a new s3 client from a Config struct.
func (c *Config) NewClient() (interface{}, error) {

	// Debug
	if c.debug {
		log.Printf("[DEBUG] Debug enabled.")
	}

	// S3 Server
	if len(c.server) < 1 {
		log.Println("[FATAL] S3 Server undefined!")
		return nil, errors.New("no S3 Server not defined")
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Server: [%s]", c.server)
	}

	// S3 Region
	if len(c.region) < 1 {
		if c.debug {
			log.Println("S3 Region not defined.  Using default value of [us-east-1]")
		}
		c.region = "us-east-1"
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Region: [%s]", c.region)
	}

	// S3 Access Key
	if len(c.accessKey) < 1 {
		log.Println("[FATAL] S3 Access Key not defined!")
		return nil, errors.New("no S3 Access Key defined")
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Access Key: [%s]", c.accessKey)
	}

	// S3 Secret Key
	if len(c.secretKey) < 1 {
		log.Println("[FATAL] S3 Secret Key not defined!")
		return nil, errors.New("no S3 Secret Key not defined")
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Secret Key: [%s]", c.secretKey)
	}

	// API Signature
	if len(c.apiSignature) < 1 {
		if c.debug {
			log.Println("[DEBUG] API Signature not defined.  Using default value of [v4]")
		}
		c.apiSignature = "v4"
	}
	if c.debug {
		log.Printf("[DEBUG] API Signature: [%s]", c.apiSignature)
	}

	// SSL
	if c.debug {
		log.Printf("[DEBUG] SSL: %v", c.ssl)
	}

	// Initialize minio client object.
	minioClient := new(minio.Client)
	var err error
	if c.apiSignature == "v2" {
		minioClient, err = minio.NewV2(c.server, c.accessKey, c.secretKey, c.ssl)
	} else if c.apiSignature == "v4" {
		minioClient, err = minio.NewV4(c.server, c.accessKey, c.secretKey, c.ssl)
	} else {
		minioClient, err = minio.New(c.server, c.accessKey, c.secretKey, c.ssl)
	}
	if err != nil {
		log.Println("[FATAL] Error connecting to S3 server.")
		return nil, err
	}

	minioClient.SetCustomTransport(c.getRoundTripper())

	if c.debug {
		log.Printf("[DEBUG] S3 client initialized")
	}

	return &s3Client{
		region:   c.region,
		debug:    c.debug,
		s3Client: minioClient,
	}, nil
}
