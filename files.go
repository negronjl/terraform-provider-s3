package main

import (
	"github.com/hashicorp/terraform/helper/schema"
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
			"file_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type: schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceS3FileCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceS3FileRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceS3FileUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceS3FileDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}