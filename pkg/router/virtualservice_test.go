package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/stretchr/testify/assert"
	"istio.io/api/networking/v1alpha3"
	"testing"
)

func TestValidateVirtualServiceList_Unit(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: []v1alpha32.VirtualService{
				{},
			},
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: []v1alpha32.DestinationRule{
				{},
			},
		},
	}

	err := ValidateVirtualServiceList(&irl)
	assert.NoError(t, err)
}

func TestValidateVirtualServiceList_Unit_EmptyItems(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: nil,
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: nil,
		},
	}

	err := ValidateVirtualServiceList(&irl)
	assert.EqualError(t, err, "empty virtualServices")
}

func TestRemove_Unit(t *testing.T) {
	var routes []*v1alpha3.HTTPRoute
	var destinations []*v1alpha3.HTTPRouteDestination

	destinations = append(destinations, &v1alpha3.HTTPRouteDestination{
		Weight: 10,
	})

	newMatch := &v1alpha3.HTTPMatchRequest{
		Headers: map[string]*v1alpha3.StringMatch{},
	}

	newMatch.Headers["x-email"] = &v1alpha3.StringMatch{
		MatchType: &v1alpha3.StringMatch_Exact{
			Exact: "somebody@domain.io",
		},
	}

	routes = append(routes, &v1alpha3.HTTPRoute{
		Match:   nil,
		Route:   nil,
		Headers: nil,
	})

	routes = append(routes, &v1alpha3.HTTPRoute{
		Match:   nil,
		Route:   destinations,
		Headers: nil,
	})

	newRoute := &v1alpha3.HTTPRoute{}
	newRoute.Match = append(newRoute.Match, newMatch)
	routes = append(routes, newRoute)

	updatedRoutes := Remove(routes, 1)
	assert.Equal(t, 2, len(updatedRoutes))
	assert.Equal(t, "somebody@domain.io", updatedRoutes[1].Match[0].Headers["x-email"].GetExact())
}

func TestVirtualService_Validate_Unit_ErrorCases(t *testing.T) {
	failureCases := []struct{
		vs VirtualService
		shift Shift
		want string
	}{
		{
			VirtualService{},
			Shift{
				Port:     0,
				Hostname: "",
				Selector: nil,
				Traffic:  Traffic{},
			},
			"could not update route without 'weight' or 'headers'",
		},
		{
			VirtualService{},
			Shift{
				Port:     0,
				Hostname: "",
				Selector: nil,
				Traffic:  Traffic{
					RequestHeaders: map[string]string {
						"header-key": "header-value",
					},
					Weight: 10,
				},
			},
			"a route needs to be served with a 'weight' or 'request headers', not both",
		},
	}

	for _, tt := range failureCases {
		err := tt.vs.Validate(tt.shift)
		assert.EqualError(t, err, tt.want)
	}
}

func TestVirtualService_Validate_Unit_Success(t *testing.T) {

	sucCases := []struct {
		vs VirtualService
		shift Shift
	}{
		{
			VirtualService{},
			Shift{
				Traffic: Traffic{
					Weight: 10,
				},
			},
		},
		{
			VirtualService{},
			Shift{
				Traffic: Traffic{
					RequestHeaders: map[string]string{
						"x-email": "somebody@domain.io",
					},
				},
			},
		},
	}

	for _, tt := range sucCases {
		err := tt.vs.Validate(tt.shift)
		assert.NoError(t, err)
	}
}