package main

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
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
				Default:  false,
			},
		},
	}
}

func resourceS3BucketCreate(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	region := meta.(*s3Client).region
	s3_client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Creating bucket: [%s] in region: [%s]", bucket, region)
	}

	if s3_client.MakeBucket(bucket, region) != nil {
		log.Printf("[FATAL] Unable to create bucket [%s] in region [%s]", bucket, region)
		return errors.New(fmt.Sprintf("Unable to create bucket [%s] in region [%s]", bucket, region))
	}
	if debug {
		log.Printf("[DEBUG] Created bucket: [%s] in region: [%s]", bucket, region)
	}
	return nil
}

func resourceS3BucketRead(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	region := meta.(*s3Client).region
	s3_client := meta.(*s3Client).s3Client
	if debug {
		log.Printf("[DEBUG] Reading bucket [%s] in region [%s]", bucket, region)
	}
	found, err := s3_client.BucketExists(bucket)
	if !found {
		return errors.New(fmt.Sprintf("[FATAL] Unable to find bucket [%s] in region [%s].  Error: %v",
			bucket, region, err))
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
	s3_client := meta.(*s3Client).s3Client
	if debug {
		log.Printf("[DEBUG] Deleting bucket [%s] from region [%s]", bucket, region)
	}
	if s3_client.RemoveBucket(bucket) != nil {
		log.Printf("[FATAL]  Unable to remove bucket [%s].")
		return errors.New(fmt.Sprintf("[FATAL] Unable to remove bucket [%s]", bucket))
	}
	return nil
}
