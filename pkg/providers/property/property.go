package property

import "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/registry"

// SubproviderName defines name of the property subprovider
const SubproviderName = "property"

func init() {
	registry.RegisterPluginSubprovider(NewPluginSubprovider())
	registry.RegisterFrameworkSubprovider(NewFrameworkSubprovider())
}
