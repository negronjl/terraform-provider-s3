package main

import (
	"log"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/minio/minio-go/v6"
	"strings"
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
			"filepath": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contentType": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "application/octet-stream",
			},
			"debug": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceS3FileCreate(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	filepath := d.Get("filepath").(string)
	content := d.Get("content").(string)
	contentType := d.Get("contentType").(string)
	client := meta.(*s3Client).s3Client

	var err error
	if filepath != "" {
		if debug {
			log.Printf("[DEBUG] Creating object [%s] from file [%s] in bucket [%s]",
				name, filepath, bucket)
		}

		_, err = client.FPutObject(bucket, name, filepath,
			minio.PutObjectOptions{ContentType: contentType})
	} else {
		if debug {
			log.Printf("[DEBUG] Creating object [%s] from [%d] bytes of content in bucket [%s]",
				name, len(content), bucket)
		}

		reader := strings.NewReader(content)
		_, err = client.PutObject(bucket, name, reader, reader.Size(),
			minio.PutObjectOptions{ContentType: contentType})
	}

	if err != nil {
		log.Printf("[FATAL] Unable to create object [%s]. Error: %v", name, err)
		return fmt.Errorf("Unable to create object [%s].  Error: %v", name, err)
	}

	if debug {
		log.Printf("[DEBUG] Created object [%s] from file [%s] in bucket [%s]",
			name, filepath, bucket)
	}

	return nil
}

func resourceS3FileRead(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	filepath := d.Get("filepath").(string)
	client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Reading file [%s] from bucket [%s] into file [%s]", name, bucket, filepath)
	}

	err := client.FGetObject(bucket, name, filepath, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("[FATAL]  Unable to read file [%s] from bucket [%s] into file [%s].  Error: %v", name, bucket, filepath, err)
		return fmt.Errorf("Unable to read file [%s].  Error: %v", name, err)
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
	client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Deleting file [%s] from bucket [%s]", name, bucket)
	}

	err := client.RemoveObject(bucket, name)
	if err != nil {
		log.Printf("[FATAL] Unable to delete file [%s] from bucket [%s].  Error: %v", name, bucket, err)
		return fmt.Errorf("Unable to delete file [%s] from bucket [%s].  Error: %v", name, bucket, err)
	}

	if debug {
		log.Printf("[DEBUG] Deleted file [%s] from bucket [%s]", name, bucket)
	}
	return nil
}
