package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	result := m.Run()
	os.Exit(result)
}

func TestValidateDestinationRuleList_Unit(t *testing.T) {
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

	err := ValidateDestinationRuleList(&irl)
	assert.NoError(t, err)
}

func TestValidateDestinationRuleList_Unit_EmptyItems(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: nil,
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: nil,
		},
	}

	err := ValidateDestinationRuleList(&irl)
	assert.EqualError(t, err, "empty destinationRules")
}

func TestDestinationRule_List_Integrated_Empty(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Name:       "api-testing",
		Namespace:  "default",
		Build:      10000,
		Istio:      fakeIstioClient,
	}

	irl, err := dr.List(map[string]string{"environment": "integration-tests"})
	assert.EqualError(t, err, "could not find any destinationRules which matched label-selector 'environment=integration-tests'")
	assert.Nil(t, irl)
}

func TestDestinationRule_List_Integrated(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Name:       "api-testing",
		Namespace:  "default",
		Build:      10000,
		Istio:      fakeIstioClient,
	}

	d := v1alpha32.DestinationRule{
		Spec: v1alpha32.DestinationRuleSpec{},
	}

	labelSelector := map[string]string{
		"environment": "integration-tests",
	}

	d.Name = "custom-dr"
	d.Namespace = dr.Namespace
	d.Labels = labelSelector

	_, _ = dr.Istio.NetworkingV1alpha3().DestinationRules(dr.Namespace).Create(&d)
	irl, err := dr.List(labelSelector)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(irl.DList.Items))
	assert.Equal(t, "custom-dr", irl.DList.Items[0].Name)
}

func TestDestinationRule_Validate_Unit(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	cases := []struct {
		dr    DestinationRule
		shift Shift
		want  string
	}{
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: nil,
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty label-selector",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     0,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty port",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     1000,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"port not in range 1024 - 65535",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     66000,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"port not in range 1024 - 65535",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic:  Traffic{},
			},
			"empty pod selector",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'name' attribute",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-test",
			Namespace:  "",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'namespace' attribute",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-test",
			Namespace:  "arrow",
			Build:      0,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'build' attribute",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-test",
			Namespace:  "arrow",
			Build:      1,
			Istio:      nil,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"nil istioClient object",
		},
		{DestinationRule{
			TrackingId: "",
			Name:       "api-test",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'trackingId' attribute",
		},
	}

	for _, tt := range cases {
		err := tt.dr.Validate(tt.shift)
		assert.EqualError(t, err, tt.want)
	}
}

func TestDestinationRule_Create_Integrated(t *testing.T) {
	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Name:       "api-testing",
		Namespace:  "arrow",
		Build:      10000,
		Istio:      fakeIstioClient,
	}

	shift := Shift{
		Traffic: Traffic{
			PodSelector: map[string]string{
				"environment": "test",
				"app":         "api-testing",
			},
		},
	}

	ir, err := dr.Create(shift)
	assert.NotNil(t, ir)
	assert.NoError(t, err)
	assert.Equal(t, "api-testing-10000-arrow", ir.Subset.Name)
}

func TestDestinationRule_Clear_Integrated_EmptyVirtualServiceRoutes(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	labelSelector := map[string]string{
		"app":         "api-test",
		"environment": "integration-tests",
	}

	// create a destinationRule object in memory
	tdr := v1alpha32.DestinationRule{
		Spec: v1alpha32.DestinationRuleSpec{},
	}

	tdr.Name = "integration-testing-dr"
	tdr.Namespace = dr.Namespace
	tdr.Labels = labelSelector

	// create a virtualService object in memory
	tvs := v1alpha32.VirtualService{
		Spec: v1alpha32.VirtualServiceSpec{},
	}
	tvs.Labels = labelSelector

	_, err := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Create(&tdr)
	_, err = fakeIstioClient.NetworkingV1alpha3().VirtualServices(dr.Namespace).Create(&tvs)

	shift := Shift{
		Selector: map[string]string{
			"app":         "api-test",
			"environment": "integration-tests",
		},
	}

	err = dr.Clear(shift)
	re, _ := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Get(dr.Name, metav1.GetOptions{})

	assert.NoError(t, err)
	assert.Equal(t, 0, len(re.Spec.Subsets))
}

func TestDestinationRule_Clear_Integrated_ExistentVirtualServiceRoutes(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()

	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	labelSelector := map[string]string{
		"app":         "api-test",
		"environment": "integration-tests",
	}

	// create a destinationRule object in memory
	tdr := v1alpha32.DestinationRule{
		Spec: v1alpha32.DestinationRuleSpec{},
	}

	tdr.Name = "integration-testing-dr"
	tdr.Namespace = dr.Namespace
	tdr.Labels = labelSelector
	tdr.Spec.Subsets = append(tdr.Spec.Subsets, &v1alpha3.Subset{
		Name: "existent-subset",
		Labels: map[string]string{
			"label":   "value",
			"version": "PR-integrated",
		},
	})

	tdr.Spec.Subsets = append(tdr.Spec.Subsets, &v1alpha3.Subset{
		Name: "subset-to-be-removed",
		Labels: map[string]string{
			"label":   "value2",
			"version": "1.3.2",
		},
	})

	// create a virtualService object in memory
	tvs := v1alpha32.VirtualService{
		Spec: v1alpha32.VirtualServiceSpec{},
	}
	tvs.Labels = labelSelector
	tvs.Spec.Http = append(tvs.Spec.Http, &v1alpha3.HTTPRoute{})
	tvs.Spec.Http[0].Match = append(tvs.Spec.Http[0].Match, &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}})
	tvs.Spec.Http[0].Route = append(tvs.Spec.Http[0].Route, &v1alpha3.HTTPRouteDestination{
		Destination: &v1alpha3.Destination{
			Host:   "api-integrated-test",
			Subset: "existent-subset",
		},
	})

	_, err := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Create(&tdr)
	_, err = fakeIstioClient.NetworkingV1alpha3().VirtualServices(dr.Namespace).Create(&tvs)

	shift := Shift{
		Selector: map[string]string{
			"app":         "api-test",
			"environment": "integration-tests",
		},
	}

	err = dr.Clear(shift)
	re, _ := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Get(dr.Name, metav1.GetOptions{})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(re.Spec.Subsets))
	assert.Equal(t, "existent-subset", re.Spec.Subsets[0].Name)
	assert.Equal(t, "PR-integrated", re.Spec.Subsets[0].Labels["version"])
}

func TestDestinationRule_Update_Integrated(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()
	dr := DestinationRule{
		Namespace:  "integration",
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	// create a destinationRule object in memory
	tdr := v1alpha32.DestinationRule{
		Spec: v1alpha32.DestinationRuleSpec{},
	}

	tdr.Name = "integration-testing-dr"
	tdr.Namespace = dr.Namespace
	labelSelector := map[string]string{
		"app":         "api-test",
		"environment": "integration-tests",
	}
	tdr.Labels = labelSelector

	_, err := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Create(&tdr)

	shift := Shift{
		Port:     8080,
		Hostname: "api-domain",
		Selector: labelSelector,
		Traffic: Traffic{
			PodSelector: map[string]string{"version": "1.2.3"},
		},
	}

	err = dr.Update(shift)

	mockedDr, _ := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Get(tdr.Name, metav1.GetOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "integration-testing-dr", mockedDr.Name)
	assert.Equal(t, "integration", mockedDr.Namespace)
	assert.Equal(t, "integration-tests", mockedDr.Labels["environment"])
}
