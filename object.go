package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)


func resourceS3Object() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3ObjectCreate,
		Read:   resourceS3ObjectRead,
		Update: resourceS3ObjectUpdate,
		Delete: resourceS3ObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
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

func resourceS3ObjectCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceS3ObjectRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceS3ObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceS3ObjectDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}