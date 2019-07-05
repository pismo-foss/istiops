package pipeline

import (
	"context"

	"github.com/pismo/istiops/pkg"
	"github.com/pismo/istiops/utils"
)

func IstioRouting(api utils.ApiValues, cid string, parentCtx context.Context) error {
	var istioResult error
	headers := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	labelSelector := map[string]string{
		"fullname": "sec-bankaccounts-ext-pr-97-7",
		// "app": api.ApiFullname,
	}

	var istiops pkg.IstioOperationsInterface = pkg.IstioValues{"api-gateway", "2.0.0", 323, "default"}

	istioResult = istiops.Headers(cid, labelSelector, headers)
	if istioResult != nil {
		return istioResult
	}

	istioResult = istiops.Percentage(cid, labelSelector, 20)
	if istioResult != nil {
		return istioResult
	}

	return nil
}
