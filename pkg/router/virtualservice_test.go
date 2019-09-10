package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	failureCases := []struct {
		vs    VirtualService
		shift Shift
		want  string
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
				Traffic: Traffic{
					RequestHeaders: map[string]string{
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
		vs    VirtualService
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

func TestUpdateVirtualService_Integrated(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()
	vs := VirtualService{
		TrackingId: "unit-testing-uuid",
		Namespace:  "integration",
		Istio:      fakeIstioClient,
	}

	v := v1alpha32.VirtualService{
		Spec: v1alpha32.VirtualServiceSpec{},
	}

	v.Name = "updated-virtualservice"
	v.Namespace = vs.Namespace

	_, err := fakeIstioClient.NetworkingV1alpha3().VirtualServices(vs.Namespace).Create(&v)

	// updating resource
	v.Labels = map[string]string{"label-key": "label-value"}

	err = UpdateVirtualService(&vs, &v)
	assert.NoError(t, err)

	mockedVs, _ := fakeIstioClient.NetworkingV1alpha3().VirtualServices(vs.Namespace).Get(v.Name, metav1.GetOptions{})

	assert.Equal(t, "label-value", mockedVs.Labels["label-key"])
}

func TestVirtualService_Clear_Integrated_EmptyRoutes(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	vs := VirtualService{
		TrackingId: "unit-testing-uuid",
		Name:       "api-testing",
		Namespace:  "integration",
		Build:      1,
		Istio:      fakeIstioClient,
	}

	shift := Shift{
		Port:     0,
		Hostname: "",
		Selector: map[string]string{
			"environment": "integration-tests",
		},
		Traffic: Traffic{},
	}

	// create a virtualService object in memory
	tvs := v1alpha32.VirtualService{
		Spec: v1alpha32.VirtualServiceSpec{},
	}

	tvs.Name = "integration-testing-dr"
	tvs.Namespace = vs.Namespace
	labelSelector := map[string]string{
		"app":         "api-test",
		"environment": "integration-tests",
	}
	tvs.Labels = labelSelector

	_, err := fakeIstioClient.NetworkingV1alpha3().VirtualServices(vs.Namespace).Create(&tvs)

	err = vs.Clear(shift)
	assert.EqualError(t, err, "empty routes when cleaning virtualService's rules")
}

func TestVirtualService_Clear_Integrated(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	vs := VirtualService{
		TrackingId: "unit-testing-uuid",
		Name:       "api-testing",
		Namespace:  "integration",
		Build:      1,
		Istio:      fakeIstioClient,
	}

	shift := Shift{
		Port:     0,
		Hostname: "",
		Selector: map[string]string{
			"environment": "integration-tests",
		},
		Traffic: Traffic{},
	}

	// create a virtualService object in memory
	tvs := v1alpha32.VirtualService{
		Spec: v1alpha32.VirtualServiceSpec{},
	}

	tvs.Name = "integration-testing-vs"
	tvs.Namespace = vs.Namespace
	labelSelector := map[string]string{
		"app":         "api-test",
		"environment": "integration-tests",
	}
	tvs.Labels = labelSelector
	tvs.Spec.Http = append(tvs.Spec.Http, &v1alpha3.HTTPRoute{
		Match: nil,
		Route: nil,
	})

	tvs.Spec.Http[0].Match = append(tvs.Spec.Http[0].Match, &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}})
	tvs.Spec.Http = append(tvs.Spec.Http, &v1alpha3.HTTPRoute{})

	_, err := fakeIstioClient.NetworkingV1alpha3().VirtualServices(vs.Namespace).Create(&tvs)

	err = vs.Clear(shift)
	assert.NoError(t, err)

	mockedVs, _ := fakeIstioClient.NetworkingV1alpha3().VirtualServices(vs.Namespace).Get(tvs.Name, metav1.GetOptions{})

	assert.Equal(t, "integration-testing-vs", mockedVs.Name)
	assert.Equal(t, "integration", mockedVs.Namespace)
	assert.Equal(t, 1, len(mockedVs.Spec.Http))
	assert.Equal(t, ".+", mockedVs.Spec.Http[0].Match[0].GetUri().GetRegex())
}

func TestBalance_Unit_PartialPercent(t *testing.T) {
	shift := Shift{
		Port:     8080,
		Hostname: "host",
		Selector: nil,
		Traffic:  Traffic{
			Weight: 40,
		},
	}

	balancedRoutes, err := Balance("current-subset","new-subset", shift)
	assert.NoError(t, err)
	assert.Equal(t, int32(60), balancedRoutes[0].Weight)
	assert.Equal(t, int32(40), balancedRoutes[1].Weight)
	assert.Equal(t, "current-subset", balancedRoutes[0].Destination.GetSubset())
	assert.Equal(t, "new-subset", balancedRoutes[1].Destination.GetSubset())
	assert.Equal(t, "host", balancedRoutes[0].Destination.GetHost())
	assert.Equal(t, "host", balancedRoutes[1].Destination.GetHost())
	assert.Equal(t, uint32(8080), balancedRoutes[0].Destination.GetPort().GetNumber())
	assert.Equal(t, uint32(8080), balancedRoutes[1].Destination.GetPort().GetNumber())
	assert.Equal(t, 2, len(balancedRoutes))
}

func TestBalance_Unit_FullPercent(t *testing.T) {
	shift := Shift{
		Port:     9090,
		Hostname: "host",
		Selector: nil,
		Traffic:  Traffic{
			Weight: 100,
		},
	}

	balancedRoutes, err := Balance("current","new", shift)
	assert.NoError(t, err)
	assert.Equal(t, int32(100), balancedRoutes[0].Weight)
	assert.Equal(t, "new", balancedRoutes[0].Destination.GetSubset())
	assert.Equal(t, "host", balancedRoutes[0].Destination.GetHost())
	assert.Equal(t, uint32(9090), balancedRoutes[0].Destination.GetPort().GetNumber())
	assert.Equal(t, 1, len(balancedRoutes))

}