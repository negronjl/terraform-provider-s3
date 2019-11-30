package main

import (
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	homedir "github.com/mitchellh/go-homedir"
	"log"
)

// Provider defines the contract for the provider definition.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"s3_server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "S3 Server",
			},
			"s3_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "us-east-1",
				Description: "S3 Server region (default: us-east-1)",
			},
			"s3_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "S3 Server Access Key",
			},
			"s3_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "S3 Server Secret Key",
			},
			"s3_shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "S3 Shared Credentials File",
			},
			"s3_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "S3 Shared Credentials Profile Name",
			},
			"s3_api_signature": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "v4",
				Description: "S3 Server API Signature (type: string, options: v2 or v4, default: v4)",
			},
			"s3_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use SSL to connect to the S3 Server? (default: false)",
			},
			"s3_ssl_insecure_ssl_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip SSL Host Verification.  (default: false)",
			},
			"s3_ssl_issuer_pem": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to PEM encoded Issuer CA Chain for the S3 Server SSL Certificate.",
			},
			"s3_debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Print debugging informatioin (default: false)",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"s3_bucket": resourceS3Bucket(),
			"s3_object": resourceS3Object(),
			"s3_file":   resourceS3File(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	debug := d.Get("s3_debug").(bool)
	if debug {
		log.Printf("[INFO] Initializing the S3 Provider")
	}

	accessKey := ""
	secretKey := ""
	path := ""
	profile := ""

	if d.Get("s3_acccess_key") != nil {
		accessKey = d.Get("s3_acccess_key").(string)
	}

	if d.Get("s3_secret_key") != nil {
		secretKey = d.Get("s3_secret_key").(string)
	}

	if d.Get("s3_shared_credentials_file") != nil {
		path = d.Get("s3_shared_credentials_file").(string)
	}

	if d.Get("s3_profile") != nil {
		profile = d.Get("s3_profile").(string)
	}

	if debug {
		log.Printf("[INFO] ACCESS_KEY: %s", accessKey)
		log.Printf("[INFO] SECRET_KEY: %s", secretKey)
	}

	if accessKey == "" || secretKey == "" {
		if debug {
			log.Printf("[INFO] s3_access_key or s3_secret_key is the empty string.  Looking for shared credentials file.")
		}

		credsPath, err := homedir.Expand(path)
		if err != nil {
			return nil, err
		}

		sharedCreds := awsCredentials.NewSharedCredentials(credsPath, profile)
		creds, err := sharedCreds.Get()
		if err != nil {
			log.Printf("[ERROR] Error encountered retrieving profile `%s` from `%s`\n%s", profile, credsPath, err)
		}
		accessKey = creds.AccessKeyID
		secretKey = creds.SecretAccessKey

		if debug {
			log.Printf("[INFO] ACCESS_KEY: %s", accessKey)
			log.Printf("[INFO] SECRET_KEY: %s", secretKey)
		}
	}

	issuerPEM, _ := homedir.Expand(d.Get("s3_ssl_issuer_pem").(string))
	config := Config{
		server:       d.Get("s3_server").(string),
		region:       d.Get("s3_region").(string),
		accessKey:    accessKey,
		secretKey:    secretKey,
		apiSignature: d.Get("s3_api_signature").(string),
		ssl:          d.Get("s3_ssl").(bool),
		insecure:     d.Get("s3_ssl_insecure_ssl_skip_verify").(bool),
		issuerPEM:    issuerPEM,
		debug:        d.Get("s3_debug").(bool),
	}
	return config.NewClient()
}
