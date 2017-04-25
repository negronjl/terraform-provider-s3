package main

import (
	"log"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
)

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
				Required:    true,
				Description: "S3 Server Access Key",
			},
			"s3_secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "S3 Server Secret Key",
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
		log.Printf("[DEBUG] Initializing the S3 Provider")
	}
	config := Config{
		s3_server:     d.Get("s3_server").(string),
		s3_region:     d.Get("s3_region").(string),
		s3_access_key: d.Get("s3_access_key").(string),
		s3_secret_key: d.Get("s3_secret_key").(string),
		api_signature: d.Get("s3_api_signature").(string),
		ssl:           d.Get("s3_ssl").(bool),
		debug:         d.Get("s3_debug").(bool),
	}
	return config.NewClient()
}
