package main

import (
	"github.com/bhoriuchi/terraform-provider-windns/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.New})
}
