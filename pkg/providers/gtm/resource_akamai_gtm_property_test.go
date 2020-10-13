package gtm

import (
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"net/http"
	"regexp"
	"testing"
)

var prop = gtm.Property{
	BackupCName:            "",
	BackupIp:               "",
	BalanceByDownloadScore: false,
	CName:                  "www.boo.wow",
	Comments:               "",
	DynamicTTL:             300,
	FailbackDelay:          0,
	FailoverDelay:          0,
	HandoutMode:            "normal",
	HealthMax:              0,
	HealthMultiplier:       0,
	HealthThreshold:        0,
	Ipv6:                   false,
	LastModified:           "2019-04-25T14:53:12.000+00:00",
	Links: []*gtm.Link{
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property",
			Rel:  "self",
		},
	},
	LivenessTests: []*gtm.LivenessTest{
		{
			DisableNonstandardPortWarning: false,
			HttpError3xx:                  true,
			HttpError4xx:                  true,
			HttpError5xx:                  true,
			Name:                          "health check",
			RequestString:                 "",
			ResponseString:                "",
			SslClientCertificate:          "",
			SslClientPrivateKey:           "",
			TestInterval:                  60,
			TestObject:                    "/status",
			TestObjectPassword:            "",
			TestObjectPort:                80,
			TestObjectProtocol:            "HTTP",
			TestObjectUsername:            "",
			TestTimeout:                   25.0,
		},
	},
	LoadImbalancePercentage:   10.0,
	MapName:                   "",
	MaxUnreachablePenalty:     0,
	Name:                      "tfexample_prop_1",
	ScoreAggregationType:      "mean",
	StaticTTL:                 600,
	StickinessBonusConstant:   0,
	StickinessBonusPercentage: 50,
	TrafficTargets: []*gtm.TrafficTarget{
		{
			DatacenterId: 3131,
			Enabled:      true,
			HandoutCName: "",
			Name:         "",
			Servers: []string{
				"1.2.3.4",
				"1.2.3.5",
			},
			Weight: 50.0,
		},
	},
	Type:                 "weighted-round-robin",
	UnreachableThreshold: 0,
	UseComputedTargets:   false,
}

func TestResGtmProperty(t *testing.T) {

	t.Run("create property", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetProperty",
			mock.Anything, // ctx is irrelevant for this test
			prop.Name,
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.PropertyResponse{}
		resp.Resource = &prop
		resp.Status = &pendingResponseStatus
		client.On("CreateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("UpdateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("DeleteProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_property.tfexample_prop_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmProperty/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted-round-robin"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmProperty/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted-round-robin"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property failed", func(t *testing.T) {
		client := &mockgtm{}

		client.On("CreateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmProperty/create_basic.tf"),
						ExpectError: regexp.MustCompile("Property Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property denied", func(t *testing.T) {
		client := &mockgtm{}

		dr := gtm.PropertyResponse{}
		dr.Resource = &prop
		dr.Status = &deniedResponseStatus
		client.On("CreateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(&dr, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmProperty/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
