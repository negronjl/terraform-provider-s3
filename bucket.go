package main

import (
	"log"

	"crypto/sha512"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceS3Bucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3BucketCreate,
		Read:   resourceS3BucketRead,
		Update: resourceS3BucketUpdate,
		Delete: resourceS3BucketDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"debug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceS3BucketCreate(d *schema.ResourceData, meta interface{}) error {
	bucket := d.Get("bucket").(string)
	region := meta.(*s3Client).region
	client := meta.(*s3Client).s3Client
	_, debugExists := d.GetOk("debug")
	debug := meta.(*s3Client).debug
	if debugExists {
		debug = d.Get("debug").(bool)
	}
	d.Set("debug", debug)

	if debug {
		log.Printf("[DEBUG] Creating bucket: [%s] in region: [%s]", bucket, region)
	}

	err := client.MakeBucket(bucket, region)
	if err != nil {
		log.Printf("[FATAL] Unable to create bucket [%s] in region [%s].  Failed with error: %v", bucket, region, err)
		return fmt.Errorf("Unable to create bucket [%s] in region [%s].  Failed with error: %v", bucket, region, err)
	}

	idkeysource := fmt.Sprintf("Bucket: [%s] Region: [%s] Host: [%s]", bucket, region, client.EndpointURL())
	id := fmt.Sprintf("%x", sha512.Sum512([]byte(idkeysource)))

	d.SetId(id)
	d.Set("endpointURL", client.EndpointURL())
	d.Set("region", region)

	if debug {
		log.Printf("[DEBUG] Created bucket: [%s] in region: [%s]", bucket, region)
	}
	return nil
}

func resourceS3BucketRead(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	region := meta.(*s3Client).region
	client := meta.(*s3Client).s3Client
	if debug {
		log.Printf("[DEBUG] Reading bucket [%s] in region [%s]", bucket, region)
	}
	found, err := client.BucketExists(bucket)
	if !found {
		return fmt.Errorf("[FATAL] Unable to find bucket [%s] in region [%s].  Error: %v",
			bucket, region, err)
	}
	return nil
}

func resourceS3BucketUpdate(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	region := meta.(*s3Client).region
	if debug {
		log.Printf("[DEBUG] Bucket update operation not implemented. Bucket: [%s], Region: [%s]",
			bucket, region)
	}
	return nil
}

func resourceS3BucketDelete(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	region := meta.(*s3Client).region
	client := meta.(*s3Client).s3Client
	if debug {
		log.Printf("[DEBUG] Deleting bucket [%s] from region [%s]", bucket, region)
	}
	if client.RemoveBucket(bucket) != nil {
		log.Printf("[FATAL]  Unable to remove bucket [%s].", bucket)
		return fmt.Errorf("[FATAL] Unable to remove bucket [%s]", bucket)
	}
	return nil
}
