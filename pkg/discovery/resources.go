package discovery

import (
	"k8s.io/client-go/kubernetes"
	//cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

var ErrResourceNotFound = fmt.Errorf("resouce requested not found in server")

func GetResource(clientset *kubernetes.Clientset, resourceName string) (v1.APIResource, error) {
	resources, e := clientset.Discovery().ServerPreferredResources()
	if e != nil{
		return v1.APIResource{}, e
	}

	for _, list := range resources{
		for _, r := range list.APIResources{
			if r.Name == resourceName || contains(r.ShortNames, resourceName) || r.SingularName == resourceName{
				r.Version = list.GroupVersion
				return r, nil
			}
		}
	}
	return v1.APIResource{}, ErrResourceNotFound

}

func GetResourceTypes(clientset *kubernetes.Clientset, groupVersion string) (resources []v1.APIResource, err error){

	//this group should be different for different cluster version:
	//1.8  apps/v1beta2
	//1.9+ app/v1

	// We also want the resource to be in the category "all"
	var resourceList *v1.APIResourceList
	resourceList, err = clientset.Discovery().ServerResourcesForGroupVersion(groupVersion)

	if err != nil{
		return
	}

	for _, r := range resourceList.APIResources{

		if contains(r.Categories, "all"){
			r.Version = resourceList.GroupVersion
			resources = append(resources, r)
		} else{
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