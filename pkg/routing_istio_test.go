package pkg

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedclientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"

	"log"
	"os"

	"github.com/stretchr/testify/assert"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

//func init() {
//	istioClient = versionedclientFake.NewSimpleClientset()
//}

func TestMain(m *testing.M) {
	TearUp()
	result := m.Run()
	os.Exit(result)

}

func TearUp(){
	istioClient = versionedclientFake.NewSimpleClientset()

	var err error

	err = CreateMockedDestinationRule()
	if err != nil {
		log.Fatal("Could not create mocked DestinationRule")
	}

	err = CreateMockedVirtualService()
	if err != nil {
		log.Fatal("Could not create mocked VirtualService")
	}
}


func CreateMockedDestinationRule() (error error) {
	mockedDr := &v1alpha32.DestinationRule{}
	mockedDr.Name = "api-unit-test-destinationrule"
	mockedDr.Namespace = "default"

	mockedDr.Spec.Subsets = append(mockedDr.Spec.Subsets, &v1alpha3.Subset{
		Name: "subset-test",
		Labels: map[string]string{
			"app": "api-xpto",
		},
	})

	_, err := istioClient.NetworkingV1alpha3().DestinationRules("default").Create(mockedDr)

	return err
}

func CreateMockedVirtualService() (error error) {
	mockedVs := &v1alpha32.VirtualService{}
	mockedVs.Name = "api-unit-test-virtualservice"
	mockedVs.Namespace = "default"

	mockedVs.Spec.Hosts = []string{"api-unit-test.domain.io"}
	mockedVs.Spec.Gateways = []string{"unit-test-gateway"}

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, &v1alpha3.HTTPRoute{})

	defaultMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}}
	defaultDestination := &v1alpha3.HTTPRouteDestination{Destination: &v1alpha3.Destination{Host: "api-xpto", Subset: "subset-name", Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: 8080}}}}
	defaultRoute := &v1alpha3.HTTPRoute{}
	defaultRoute.Match = append(defaultRoute.Match, defaultMatch)
	defaultRoute.Route = append(defaultRoute.Route, defaultDestination)

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, defaultRoute)

	_, err := istioClient.NetworkingV1alpha3().VirtualServices("default").Create(mockedVs)

	return err
}

func TestGetAllDestinationRules(t *testing.T) {
	listOptions := metav1.ListOptions{}
	mockedDrs, err := GetAllDestinationRules("random-cid", "default", listOptions)

	assert.NoError(t, err)
	assert.IsType(t, v1alpha32.DestinationRuleList{}, *mockedDrs)
	assert.EqualValues(t, "api-unit-test-destinationrule", mockedDrs.Items[0].Name)
}

func TestGetAllVirtualServices(t *testing.T) {
	listOptions := metav1.ListOptions{}
	mockedVss, err := GetAllVirtualServices("random-cid", "default", listOptions)

	assert.NoError(t, err)
	assert.IsType(t, v1alpha32.VirtualServiceList{}, *mockedVss)
	assert.EqualValues(t, "api-unit-test-virtualservice", mockedVss.Items[0].Name)
}


//
//func TestGetAllDestinationRules(t *testing.T) {
//	_, err := GetAllDestinationRules("random-cid", "default")
//	assert.NoError(t, err)
//}
//
//func TestGetResourcesToUpdate(t *testing.T) {
//	v := IstioValues{
//		"api-xpto",
//		"2.0.0",
//		123,
//		"default",
//	}
//	userLabelSelector := map[string]string{
//		"app": "api-xpto",
//	}
//	_, err := GetResourcesToUpdate("random-cid", v, userLabelSelector)
//	assert.NoError(t, err)
//
//}
