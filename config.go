package main

import (
	"errors"
	"log"

	"github.com/minio/minio-go/v6"
)

type Config struct {
	s3_server     string
	s3_region     string
	s3_access_key string
	s3_secret_key string
	api_signature string
	ssl           bool
	debug         bool
}

type s3Client struct {
	region   string
	s3Client *minio.Client
}

func (c *Config) NewClient() (interface{}, error) {

	// Debug
	if c.debug {
		log.Printf("[DEBUG] Debug enabled.")
	}

	// S3 Server
	if len(c.s3_server) < 1 {
		log.Println("[FATAL] S3 Server undefined!")
		return nil, errors.New("S3 Server not defined!")
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Server: [%s]", c.s3_server)
	}

	// S3 Region
	if len(c.s3_region) < 1 {
		if c.debug {
			log.Println("S3 Region not defined.  Using default value of [us-east-1]")
		}
		c.s3_region = "us-east-1"
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Region: [%s]", c.s3_region)
	}

	// S3 Access Key
	if len(c.s3_access_key) < 1 {
		log.Println("[FATAL] S3 Access Key not defined!")
		return nil, errors.New("S3 Access Key not defined!")
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Access Key: [%s]", c.s3_access_key)
	}

	// S3 Secret Key
	if len(c.s3_secret_key) < 1 {
		log.Println("[FATAL] S3 Secret Key not defined!")
		return nil, errors.New("S3 Secret Key not defined!")
	}
	if c.debug {
		log.Printf("[DEBUG] S3 Secret Key: [%s]", c.s3_secret_key)
	}

	// API Signature
	if len(c.api_signature) < 1 {
		if c.debug {
			log.Println("[DEBUG] API Signature not defined.  Using default value of [v4]")
		}
		c.api_signature = "v4"
	}
	if c.debug {
		log.Printf("[DEBUG] API Signature: [%s]", c.api_signature)
	}

	// SSL
	if c.debug {
		log.Printf("[DEBUG] SSL: %v", c.ssl)
	}

	// Initialize minio client object.
	minioClient := new(minio.Client)
	var err error
	if c.api_signature == "v2" {
		minioClient, err = minio.NewV2(c.s3_server, c.s3_access_key, c.s3_secret_key, c.ssl)
	} else if c.api_signature == "v4" {
		minioClient, err = minio.NewV4(c.s3_server, c.s3_access_key, c.s3_secret_key, c.ssl)
	} else {
		minioClient, err = minio.New(c.s3_server, c.s3_access_key, c.s3_secret_key, c.ssl)
	}
	if err != nil {
		log.Println("[FATAL] Error connecting to S3 server.")
		return nil, err
	} else {
		if c.debug {
			log.Printf("[DEBUG] S3 client initialized")
		}
	}

	return &s3Client{
		region:   c.s3_region,
		s3Client: minioClient,
	}, nil
}
