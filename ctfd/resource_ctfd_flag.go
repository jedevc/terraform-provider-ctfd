package ctfd

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jedevc/terraform-provider-ctfd/api"
)

func resourceCTFdFlag() *schema.Resource {
	return &schema.Resource{
		Create:   resourceCTFdFlagCreate,
		Read:     resourceCTFdFlagRead,
		Update:   resourceCTFdFlagUpdate,
		Delete:   resourceCTFdFlagDelete,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"flag_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"challenge": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"regex": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceCTFdFlagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	var flagType string
	if d.Get("regex").(bool) {
		flagType = "regex"
	} else {
		flagType = "static"
	}

	flag := api.Flag{
		Challenge: uint(d.Get("challenge").(int)),
		Type:      flagType,
		Content:   d.Get("pattern").(string),
	}
	flag, err := client.CreateFlag(flag)
	if err != nil {
		return err
	}
	d.Set("flag_id", flag.ID)

	d.SetId(strconv.Itoa(d.Get("flag_id").(int)))
	return nil
}

func resourceCTFdFlagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	flagID := uint(d.Get("flag_id").(int))
	flag, err := client.GetFlag(flagID)
	if err != nil {
		return err
	}

	d.Set("regex", flag.Type == "regex")
	d.Set("challenge", flag.Challenge)
	d.Set("pattern", flag.Content)

	return nil
}

func resourceCTFdFlagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	var flagType string
	if d.Get("regex").(bool) {
		flagType = "regex"
	} else {
		flagType = "static"
	}

	flag := api.Flag{
		ID:        uint(d.Get("flag_id").(int)),
		Challenge: uint(d.Get("challenge").(int)),
		Type:      flagType,
		Content:   d.Get("pattern").(string),
	}
	_, err := client.ModifyFlag(flag)
	return err
}

func resourceCTFdFlagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	flagID := uint(d.Get("flag_id").(int))

	err := client.DeleteFlag(flagID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
