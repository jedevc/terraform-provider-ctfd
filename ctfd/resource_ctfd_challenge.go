package ctfd

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCTFdChallenge() *schema.Resource {
	return &schema.Resource{
		Create:   resourceCTFdChallengeCreate,
		Read:     resourceCTFdChallengeRead,
		Update:   nil,
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
				ForceNew: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"points": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"hidden": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},
			"max_attempts": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  0,
				Optional: true,
				ForceNew: true,
			},
			// "file": &schema.Schema{
			// 	Type:     schema.TypeSet,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"content": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"filename": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// },
			// "flag": &schema.Schema{
			// 	Type:     schema.TypeSet,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"pattern": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"regex": {
			// 				Type:     schema.TypeBool,
			// 				Default:  false,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },
			// "hint": &schema.Schema{
			// 	Type:     schema.TypeSet,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"text": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"cost": {
			// 				Type:     schema.TypeInt,
			// 				Default:  0,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },
			// "requirements": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeInt,
			// 	},
			// },
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

	chal := Challenge{
		Type:        "standard",
		Name:        d.Get("name").(string),
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		Value:       d.Get("points").(int),
		State:       state,
		MaxAttempts: uint(d.Get("max_attempts").(int)),
	}
	challengeID, err := client.CreateChallenge(chal)
	if err != nil {
		return err
	}

	d.Set("challenge_id", challengeID)
	d.SetId(strconv.FormatUint(uint64(challengeID), 10))
	return nil
}

func resourceCTFdChallengeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
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
