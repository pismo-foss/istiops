package router

import (
	"istio.io/api/networking/v1alpha3"
)

type Route struct {
	Port     uint32
	Hostname string
	Selector *Selector
	Traffic  *Traffic
}

type Traffic struct {
	RequestHeaders map[string]string
	Weight         int32
}

type Metadata struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
}

type Selector struct {
	ResourceSelector map[string]string
	PodSelector      map[string]string
}

type Router interface {
	Validate(route Route) error
	Update(route Route) error
	Delete(route Route) error
}

type Subset struct {
	Subset *v1alpha3.Subset
}

type TrafficShift struct {
	Headers map[string]string
	Percent int32
}

//GetResourcesToUpdate returns a slice of all DestinationRules and/or VirtualServices (based on given labelSelectors to a posterior update
//func GetResourcesToUpdate(labelSelector map[string]string) (*IstioRouteList, error) {
//	StringifyLabelSelector, _ := utils.StringifyLabelSelector(ips.TrackingId, labelSelector)
//
//	listOptions := metav1.ListOptions{
//		LabelSelector: StringifyLabelSelector,
//	}
//
//	matchedDrs, err := GetAllDestinationRules(ips, listOptions)
//	if err != nil {
//		utils.Fatal(fmt.Sprintf("%s", err), "")
//		return nil, err
//	}
//
//	matchedVss, err := GetAllVirtualServices(ips, listOptions)
//	if err != nil {
//		utils.Fatal(fmt.Sprintf("%s", err), "")
//		return nil, err
//	}
//
//	if len(matchedDrs.Items) == 0 || len(matchedVss.Items) == 0 {
//		utils.Fatal(fmt.Sprintf("Couldn't find any istio resources based on given labelSelector '%s' to update. ", StringifyLabelSelector), "")
//		return nil, err
//	}
//
//	matchedResourcesList := &IstioRouteList{
//		matchedVss,
//		matchedDrs,
//	}
//
//	return matchedResourcesList, nil
//}
