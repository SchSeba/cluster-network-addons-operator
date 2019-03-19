package network

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/kubevirt/cluster-network-addons-operator/pkg/render"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	opv1alpha1 "github.com/kubevirt/cluster-network-addons-operator/pkg/apis/networkaddonsoperator/v1alpha1"
)

func changeSafeKubeMacPool(prev, next *opv1alpha1.NetworkAddonsConfigSpec) []error {
	if prev.KubeMacPool != nil && !reflect.DeepEqual(prev.KubeMacPool, next.KubeMacPool) {
		return []error{errors.Errorf("cannot modify KubeMacPool configuration once it is deployed")}
	}
	return nil
}

// renderLinuxBridge generates the manifests of Linux Bridge
func renderKubeMacPool(conf *opv1alpha1.NetworkAddonsConfigSpec, manifestDir string) ([]*unstructured.Unstructured, error) {
	if conf.KubeMacPool == nil {
		return nil, nil
	}

	// render the manifests on disk
	data := render.MakeRenderData()
	data.Data["KubeMacPoolImage"] = os.Getenv("KUBEMACPOOL_IMAGE")
	data.Data["ImagePullPolicy"] = os.Getenv("IMAGE_PULL_POLICY")

	if conf.KubeMacPool.StartPoolRange == "" || conf.KubeMacPool.EndPoolRange == "" {
		data.Data["DataExist"] = false
	} else {
		data.Data["EndPoolRange"] = conf.KubeMacPool.EndPoolRange
		data.Data["StartPoolRange"] = conf.KubeMacPool.StartPoolRange
		data.Data["DataExist"] = true

	}

	objs, err := render.RenderDir(filepath.Join(manifestDir, "kubemacpool"), &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render kubemacpool manifests")
	}

	return objs, nil
}
