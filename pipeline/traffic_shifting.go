package pipeline

import (
	"context"
	"github.com/pismo/istiops/pkg"
	"github.com/pismo/istiops/utils"
)

func IstioRouting(api utils.ApiValues, cid string, parentCtx context.Context) error {
	var istioResult error
	//headers := map[string]string{
	//	"app":   "api-xpto",
	//	"build": "123",
	//}

	labelSelector := map[string]string{
		"environment": "pipeline-go",
	}

	var istiops pkg.IstioOperationsInterface = pkg.IstioValues{"sec-bankaccounts", "2.0.0", 323, "default"}

	//istioResult = istiops.SetLabelsDestinationRule(cid, "sec-bankaccounts-destination-rules", labelSelector)
	//if istioResult != nil {
	//	return istioResult
	//}
	//
	//virtualServices := []string{"sec-bankaccounts-virtualservice", "sec-bankaccounts-internal-virtualservice"}
	//for _, virtualService := range virtualServices {
	//	istioResult = istiops.SetLabelsVirtualService(cid, virtualService, labelSelector)
	//	if istioResult != nil {
	//		return istioResult
	//	}
	//}
	//
	//subsetName, istioResult := istiops.SetHeaders(cid, labelSelector, "api-xpto", headers, 8080)
	//if (istioResult != nil) || (subsetName == "") {
	//	return istioResult
	//}
	//
	//istioResult = istiops.SetPercentage(cid, "sec-bankaccounts-virtualservice", subsetName, 85)
	//if istioResult != nil {
	//	return istioResult
	//}

	istioResult = istiops.ClearRules(cid, "uri", labelSelector)
	if istioResult != nil {
		return istioResult
	}

	return nil
}
