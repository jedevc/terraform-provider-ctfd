package ctfd

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jedevc/terraform-provider-ctfd/api"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: providerDataSourcesMap(),
		ResourcesMap:   providerResources(),
		ConfigureFunc:  providerConfigure,
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The admin username",
		},
		"password": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The admin password",
		},
		"url": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL of your CTFd server",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"ctfd_challenge": resourceCTFdChallenge(),
		"ctfd_flag":      resourceCTFdFlag(),
		"ctfd_file":      resourceCTFdFile(),
	}
}

func providerDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := api.Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
	}
	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return &client, nil
}
