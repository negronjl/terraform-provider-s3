package main

import (
	"log"

	"bytes"
	"crypto/sha512"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/minio/minio-go/v6"
	"io/ioutil"
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
			"file_path": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "application/octet-stream",
			},
			"debug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceS3FileCreate(d *schema.ResourceData, meta interface{}) error {
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	filepath := d.Get("file_path").(string)
	content := d.Get("content").(string)
	contentType := d.Get("content_type").(string)
	client := meta.(*s3Client).s3Client

	_, debugExists := d.GetOk("debug")
	debug := meta.(*s3Client).debug
	if debugExists {
		debug = d.Get("debug").(bool)
	}
	d.Set("debug", debug)

	var err error
	if filepath != "" {
		if debug {
			log.Printf("[DEBUG] Creating object [%s] from file [%s] in bucket [%s]",
				name, filepath, bucket)
		}

		buf, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Printf("[FATAL] Unable to read file [%s].  Error %v", filepath, err)
			return fmt.Errorf("[FATAL] Unable to read file [%s].  Error %v", filepath, err)
		}
		content := string(buf)
		d.Set("content", content)
	}

	reader := strings.NewReader(content)
	_, err = client.PutObject(bucket, name, reader, reader.Size(),
		minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		log.Printf("[FATAL] Unable to create object [%s]. Error: %v", name, err)
		return fmt.Errorf("Unable to create object [%s].  Error: %v", name, err)
	}

	if debug {
		log.Printf("[DEBUG] Created object [%s] from file [%s] in bucket [%s]",
			name, filepath, bucket)
	}
	region, err := client.GetBucketLocation(bucket)
	if err != nil {
		log.Printf("[DEBUG] Could not retrieve bucket location for bucket [%s] at host [%s]", bucket, client.EndpointURL())
		return fmt.Errorf("[DEBUG] Could not retrieve bucket location for bucket [%s] at host [%s]", bucket, client.EndpointURL())
	}

	idkeysource := fmt.Sprintf("ObjectKey [%s] Bucket: [%s] Region: [%s] Host: [%s]", name, bucket, region, client.EndpointURL())
	id := fmt.Sprintf("%x", sha512.Sum512([]byte(idkeysource)))

	d.SetId(id)
	d.Set("endpointURL", client.EndpointURL())
	d.Set("region", region)
	d.Set("debug", debug)

	return nil
}

func resourceS3FileRead(d *schema.ResourceData, meta interface{}) error {
	debug := d.Get("debug").(bool)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	client := meta.(*s3Client).s3Client

	if debug {
		log.Printf("[DEBUG] Reading file [%s] from bucket [%s] into memory", name, bucket)
	}

	reader, err := client.GetObject(bucket, name, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("[FATAL]  Unable to read content [%s] from bucket [%s] into memory.  Error: %v", name, bucket, err)
		return fmt.Errorf("Unable to read file [%s].  Error: %v", name, err)
	}

	defer reader.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)

	if err != nil {
		log.Printf("[FATAL] Unable to read content from reader for content [%s] from bucket [%s] into memory.  Error %v", name, bucket, err)
		return fmt.Errorf("[FATAL] Unable to read content from reader for content [%s] from bucket [%s] into memory.  Error %v", name, bucket, err)
	}

	region, err := client.GetBucketLocation(bucket)
	if err != nil {
		log.Printf("[FATAL] Unable to read region from bucket [%s].  Error %v", bucket, err)
		return fmt.Errorf("[FATAL] Unable to read region from bucket [%s].  Error %v", bucket, err)
	}

	idkeysource := fmt.Sprintf("ObjectKey [%s] Bucket: [%s] Region: [%s] Host: [%s]", name, bucket, region, client.EndpointURL())
	id := fmt.Sprintf("%x", sha512.Sum512([]byte(idkeysource)))

	d.SetId(id)
	d.Set("region", region)
	d.Set("endpointURL", client.EndpointURL())
	d.Set("content", buf.String())

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
