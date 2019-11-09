package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/jedevc/terraform-provider-ctfd/ctfd"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ctfd.Provider,
	})
}
