package main

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)


func resourceS3File() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3FileCreate,
		Read:   resourceS3FileRead,
		Update: resourceS3FileUpdate,
		Delete: resourceS3FileDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_path": {
				Type: schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type: schema.TypeString,
				Required: false,
				Default: "application/octet-stream",
			},
			"debug": {
				Type:     schema.TypeBool,
				Required: false,
			},
		},
	}
}

func resourceS3FileCreate(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	file_path := d.Get("file_path").(string)
	content_type := d.Get("content_type").(string)
	s3_client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Creating object [%s] from file [%s] in bucket [%s]",
			name, file_path, bucket)
	}

	_, err := s3_client.FPutObject(bucket, name, file_path,content_type)
	if err != nil {
		log.Printf("[FATAL] Unable to create object [%s]", name)
		return errors.New(fmt.Sprintf("Unable to create object [%s]", name))
	}

	if debug {
		log.Printf("[DEBUG] Created object [%s] from file [%s] in bucket [%s]",
		name, file_path, bucket)
	}

	return nil
}

func resourceS3FileRead(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	file_path := d.Get("file_path").(string)
	s3_client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Reading file [%s] from bucket [%s] into file [%s]", name, bucket, file_path)
	}

	err := s3_client.FGetObject(bucket, name, file_path)
	if err != nil {
		log.Printf("[FATAL]  Unable to read file [%s] from bucket [%s] into file [%s]", name, bucket, file_path)
		return errors.New(fmt.Sprintf("Unable to read file [%s]", name))
	}

	if debug {
		log.Printf("[DEBUG] Read file [%s] from bucket [%s]", name, bucket)
	}
	return nil
}

func resourceS3FileUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceS3BucketCreate(d, meta)
}

func resourceS3FileDelete(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	s3_client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Deleting file [%s] from bucket [%s]", name, bucket)
	}

	err := s3_client.RemoveObject(bucket, name)
	if err != nil {
		log.Printf("[FATAL] Unable to delete file [%s] from bucket [%s]", name, bucket)
		return errors.New(fmt.Sprintf("Unable to delete file [%s] from bucket [%s]", name, bucket))
	}

	if debug {
		log.Printf("[DEBUG] Deleted file [%s] fro bucket [%s]", name, bucket)
	}
	return nil
}