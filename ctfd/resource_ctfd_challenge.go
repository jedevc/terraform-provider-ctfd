package ctfd

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jedevc/terraform-provider-ctfd/api"
)

func resourceCTFdChallenge() *schema.Resource {
	return &schema.Resource{
		Create:   resourceCTFdChallengeCreate,
		Read:     resourceCTFdChallengeRead,
		Update:   resourceCTFdChallengeUpdate,
		Delete:   resourceCTFdChallengeDelete,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"challenge_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"points": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"hidden": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"max_attempts": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  0,
				Optional: true,
			},
		},
	}
}

func resourceCTFdChallengeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	var state string
	if d.Get("hidden").(bool) {
		state = "hidden"
	} else {
		state = "visible"
	}

	chal := api.Challenge{
		Type:        "standard",
		Name:        d.Get("name").(string),
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		Value:       d.Get("points").(int),
		State:       state,
		MaxAttempts: uint(d.Get("max_attempts").(int)),
	}
	chal, err := client.CreateChallenge(chal)
	if err != nil {
		return err
	}

	d.Set("challenge_id", chal.ID)
	d.SetId(strconv.Itoa(d.Get("challenge_id").(int)))
	return nil
}

func resourceCTFdChallengeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	challengeID := uint(d.Get("challenge_id").(int))
	chal, err := client.GetChallenge(challengeID)
	if err != nil {
		return err
	}

	d.Set("name", chal.Name)
	d.Set("category", chal.Category)
	d.Set("description", chal.Description)
	d.Set("points", chal.Value)
	d.Set("hidden", chal.State == "hidden")
	d.Set("max_attempts", chal.MaxAttempts)

	return nil
}

func resourceCTFdChallengeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	var state string
	if d.Get("hidden").(bool) {
		state = "hidden"
	} else {
		state = "visible"
	}

	chal := api.Challenge{
		ID:          uint(d.Get("challenge_id").(int)),
		Type:        "standard",
		Name:        d.Get("name").(string),
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		Value:       d.Get("points").(int),
		State:       state,
		MaxAttempts: uint(d.Get("max_attempts").(int)),
	}
	_, err := client.ModifyChallenge(chal)
	return err
}

func resourceCTFdChallengeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformCTFdContext).client
	// config := meta.(*TerraformCTFdContext).config

	challengeID := uint(d.Get("challenge_id").(int))

	err := client.DeleteChallenge(challengeID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
