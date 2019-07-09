package pkg

import (
	"context"
	"fmt"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	DESTINATION_RULE_SUFFIX    = "-destination-rules"
	VIRTUALSERVICE_RULE_SUFFIX = "-virtualservice"
	INTRASERVICE_GATEWAY       = "mesh"
	INTERNAL_GATEWAY           = "istio-gateway-internal"
	EXTERNAL_GATEWAY           = "istio-gateway"
)

// CreateRouteResource deploys an virtualservice with a destinationrule for a given api struct.
// this function will create a virtualservice and  a destinationrule
// if they exist already, this function wil override with the custom definitions
func CreateRouteResource(api utils.ApiValues, cid string, parentCtx context.Context) error {
	ctx := context.WithValue(parentCtx, "cid", cid)

	//Handling destinationRule
	dr, created, err := retrieveDestinationRule(api, cid, ctx)
	if err != nil {
		return err
	}

	if created {
		utils.Info(fmt.Sprintf("Creating new destinationrule: %s", dr.Name), cid)
		_, err = istioClient.NetworkingV1alpha3().DestinationRules(api.Namespace).Create(dr)
	} else {
		utils.Info(fmt.Sprintf("Updating destinationrule: %s", dr.Name), cid)
		_, err = istioClient.NetworkingV1alpha3().DestinationRules(api.Namespace).Update(dr)
	}

	if err != nil {
		return err
	}

	//Handling virtualservice
	vs, created, err := retrieveVirtualService(api, cid, ctx)
	if err != nil {
		return err
	}
	if created {
		utils.Info(fmt.Sprintf("Creating new virtualservice: %s", vs.Name), cid)
		_, err = istioClient.NetworkingV1alpha3().VirtualServices(api.Namespace).Create(vs)
	} else {
		utils.Info(fmt.Sprintf("Updating virtualservice: %s", vs.Name), cid)
		_, err = istioClient.NetworkingV1alpha3().VirtualServices(api.Namespace).Update(vs)
	}

	if err != nil {
		return err
	}

	return nil
}

// retrieveVirtualService this function retrieves an virtualservice that already exists.
// if it doesn't exists it will return a new one in memory with the given apistruct parameters as the initialization values.
func retrieveVirtualService(api utils.ApiValues, cid string, parentCtx context.Context) (virtualService *v1alpha32.VirtualService, new bool, error error) {
	name := api.Name + VIRTUALSERVICE_RULE_SUFFIX
	utils.Info(fmt.Sprintf("Retrieving virtualservice: %s", name), cid)

	vs, err := GetVirtualService(cid, name, api.Namespace)
	if err != nil {
		customErr, ok := err.(*errors.StatusError)
		if !ok {
			return nil, false, err
		}

		if customErr.Status().Code != 404 {
			return nil, false, err
		}

		utils.Info(fmt.Sprintf("VirtualService: %s not found, will create one now", name), cid)
		vs := &v1alpha32.VirtualService{}
		vs.Name = api.Name + VIRTUALSERVICE_RULE_SUFFIX
		vs.Namespace = api.Namespace

		vs.Spec.Hosts = []string{api.Name, api.ApiHostName}
		vs.Spec.Gateways = []string{INTRASERVICE_GATEWAY, INTERNAL_GATEWAY}

		if api.GrpcPort > 0 {
			defaultMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Prefix{Prefix: "/"}}}
			defaultDestination := &v1alpha3.HTTPRouteDestination{Destination: &v1alpha3.Destination{Host: api.Name, Subset: "", Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: api.GrpcPort}}}}

			defaultRoute := &v1alpha3.HTTPRoute{}
			defaultRoute.Match = append(defaultRoute.Match, defaultMatch)
			defaultRoute.Route = append(defaultRoute.Route, defaultDestination)

			vs.Spec.Http = append(vs.Spec.Http, defaultRoute)
			return vs, true, nil
		}

		defaultMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}}
		defaultDestination := &v1alpha3.HTTPRouteDestination{Destination: &v1alpha3.Destination{Host: api.Name, Subset: "", Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: api.HttpPort}}}}

		defaultRoute := &v1alpha3.HTTPRoute{}
		defaultRoute.Match = append(defaultRoute.Match, defaultMatch)
		defaultRoute.Route = append(defaultRoute.Route, defaultDestination)

		vs.Spec.Http = append(vs.Spec.Http, defaultRoute)
		return vs, true, nil
	}

	return vs, false, nil
}

// retrieveDestinationRule this function retrieves an destinationrule that already exists.
// if it doesn't exists it will return a new one in memory with the given apistruct parameters as the initialization values.
func retrieveDestinationRule(api utils.ApiValues, cid string, parentCtx context.Context) (destinationRule *v1alpha32.DestinationRule, new bool, error error) {
	drName := api.Name + DESTINATION_RULE_SUFFIX
	utils.Info(fmt.Sprintf("Retrieving destinationrule: %s", drName), cid)

	dr, err := GetDestinationRule(cid, drName, api.Namespace)
	if err != nil {
		customErr, ok := err.(*errors.StatusError)
		if !ok {
			return nil, false, err
		}

		if customErr.Status().Code != 404 {
			return nil, false, err
		}

		utils.Info(fmt.Sprintf("DestinationRule: %s not found, will create one now", drName), cid)
		dr := &v1alpha32.DestinationRule{}
		dr.Name = api.Name + DESTINATION_RULE_SUFFIX
		dr.Namespace = api.Namespace
		dr.Spec.Host = api.Name
		dr.Spec.TrafficPolicy = &v1alpha3.TrafficPolicy{Tls: &v1alpha3.TLSSettings{}}
		dr.Spec.TrafficPolicy.Tls.Mode = v1alpha3.TLSSettings_ISTIO_MUTUAL
		dr.Spec.Subsets = make([]*v1alpha3.Subset, 0)
		return dr, true, nil
	}

	return dr, false, nil
}
