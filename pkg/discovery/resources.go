package discovery

import (
	"k8s.io/client-go/kubernetes"
	//cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetResourceTypes(clientset *kubernetes.Clientset) (resources []v1.APIResource, err error){

	//this group should be different for different cluster version:
	//1.8  apps/v1beta2
	//1.9+ app/v1

	// We also want the resource to be in the category "all"
	var resourceList *v1.APIResourceList
	resourceList, err = clientset.Discovery().ServerResourcesForGroupVersion("apps/v1beta2")

	if err != nil{
		return
	}

	for _, r := range resourceList.APIResources{
		if contains(r.Categories, "all"){
			resources = append(resources, r)
		}
	}

	return
}

func contains(list []string, probe string) bool{
	for _, s := range list{
		if s == probe{
			return true
		}
	}
	return false
}