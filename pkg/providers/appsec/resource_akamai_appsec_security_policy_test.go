package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSecurityPolicy_res_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getSecurityPolicyResponse := appsec.GetSecurityPolicyResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"), &getSecurityPolicyResponse)
		require.NoError(t, err)

		createSecurityPolicyResponse := appsec.CreateSecurityPolicyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicyCreate.json"), &createSecurityPolicyResponse)
		require.NoError(t, err)

		removeSecurityPolicyResponse := appsec.RemoveSecurityPolicyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"), &removeSecurityPolicyResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicy",
			mock.Anything,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&getSecurityPolicyResponse, nil)

		client.On("CreateSecurityPolicy",
			mock.Anything,
			appsec.CreateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyName: "PLE Cloned Test for Launchpad 15", PolicyPrefix: "PLE", DefaultSettings: true},
		).Return(&createSecurityPolicyResponse, nil)

		client.On("RemoveSecurityPolicy",
			mock.Anything,
			appsec.RemoveSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&removeSecurityPolicyResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy.test", "id", "43253:PLE_114049"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
