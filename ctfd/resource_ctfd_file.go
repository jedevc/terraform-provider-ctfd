package ctfd

import (
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jedevc/terraform-provider-ctfd/api"
)

func resourceCTFdFile() *schema.Resource {
	return &schema.Resource{
		Create:   resourceCTFdFileCreate,
		Read:     resourceCTFdFileRead,
		Update:   nil,
		Delete:   resourceCTFdFileDelete,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"file_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"challenge": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hash": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCTFdFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	config := meta.(*TerraformCTFdContext).config

	fd, err := os.Open(d.Get("filename").(string))
	if err != nil {
		log.Print("could not open file for reading")
		return err
	}
	defer fd.Close()

	spec := api.FileSpec{
		Challenge: uint(d.Get("challenge").(int)),
		Type:      "challenge",
		File:      fd,
	}
	files, err := client.CreateFile(spec)
	if err != nil {
		return err
	}
	file := files[0]

	d.Set("location", config.URL+"/files/"+file.Location)
	d.Set("file_id", file.ID)

	d.SetId(strconv.Itoa(d.Get("file_id").(int)))
	return nil
}

func resourceCTFdFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	config := meta.(*TerraformCTFdContext).config

	fileID := uint(d.Get("file_id").(int))
	file, err := client.GetFile(fileID)
	if err != nil {
		return err
	}

	d.Set("location", config.URL+"/files/"+file.Location)

	return nil
}

func resourceCTFdFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	fileID := uint(d.Get("file_id").(int))
	err := client.DeleteFile(fileID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
