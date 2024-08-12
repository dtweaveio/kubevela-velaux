package types

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClientType string
type ResourceType string

const (
	ClientKubernetes ClientType = "Kubernetes"
)

var SupportedGroupVersionResources = map[ClientType][]schema.GroupVersionResource{
	ClientKubernetes: {
		{Group: "", Version: "v1", Resource: "namespaces"},
		{Group: "", Version: "v1", Resource: "nodes"},
		{Group: "", Version: "v1", Resource: "resourcequotas"},
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		{Group: "", Version: "v1", Resource: "secrets"},
		{Group: "", Version: "v1", Resource: "configmaps"},
		{Group: "", Version: "v1", Resource: "serviceaccounts"},

		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},

		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
		{Group: "apps", Version: "v1", Resource: "replicasets"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "controllerrevisions"},

		{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},

		{Group: "batch", Version: "v1", Resource: "jobs"},
		{Group: "batch", Version: "v1", Resource: "cronjobs"},

		{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},

		//{Group: "autoscaling", Version: "v2beta2", Resource: "horizontalpodautoscalers"},
	},
}

var SupportedResources = map[string]schema.GroupVersionResource{
	"namespaces":             {Group: "", Version: "v1", Resource: "namespaces"},
	"nodes":                  {Group: "", Version: "v1", Resource: "nodes"},
	"pods":                   {Group: "", Version: "v1", Resource: "pods"},
	"services":               {Group: "", Version: "v1", Resource: "services"},
	"secrets":                {Group: "", Version: "v1", Resource: "secrets"},
	"configmaps":             {Group: "", Version: "v1", Resource: "configmaps"},
	"serviceaccounts":        {Group: "", Version: "v1", Resource: "serviceaccounts"},
	"resourcequotas":         {Group: "", Version: "v1", Resource: "resourcequotas"},
	"persistentvolumeclaims": {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	//
	"roles":               {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
	"rolebindings":        {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
	"clusterroles":        {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
	"clusterrolebindings": {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
	//
	"deployments":         {Group: "apps", Version: "v1", Resource: "deployments"},
	"daemonsets":          {Group: "apps", Version: "v1", Resource: "daemonsets"},
	"replicasets":         {Group: "apps", Version: "v1", Resource: "replicasets"},
	"statefulsets":        {Group: "apps", Version: "v1", Resource: "statefulsets"},
	"controllerrevisions": {Group: "apps", Version: "v1", Resource: "controllerrevisions"},
	//
	"storageclasses":        {Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
	"persistentvolumes":     {Group: "", Version: "v1", Resource: "persistentvolumes"},
	"volumesnapshotclasses": {Group: "", Version: "v1", Resource: "volumesnapshotclasses"},
	//
	"customresourcedefinitions": {Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"},
	//
	"jobs":     {Group: "batch", Version: "v1", Resource: "jobs"},
	"cronjobs": {Group: "batch", Version: "v1", Resource: "cronjobs"},
	//
	"ingresses": {Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	//
	"horizontalpodautoscalers": {Group: "autoscaling", Version: "v2beta2", Resource: "horizontalpodautoscalers"},

	"guestbooks": {Group: "samples.dtweave.io", Version: "v1", Resource: "guestbooks"},
}

func CreateResourceObject(resource string) (client.Object, error) {
	switch resource {
	case "namespaces":
		return &corev1.Namespace{}, nil
	case "nodes":
		return &corev1.Node{}, nil
	case "pods":
		return &corev1.Pod{}, nil
	case "services":
		return &corev1.Service{}, nil
	case "secrets":
		return &corev1.Secret{}, nil
	case "configmaps":
		return &corev1.ConfigMap{}, nil
	case "serviceaccounts":
		return &corev1.ServiceAccount{}, nil
	case "resourcequotas":
		return &corev1.ResourceQuota{}, nil
	case "persistentvolumeclaims":
		return &corev1.PersistentVolumeClaim{}, nil
	//
	case "roles":
		return &rbacv1.Role{}, nil
	case "rolebindings":
		return &rbacv1.RoleBinding{}, nil
	case "clusterroles":
		return &rbacv1.ClusterRole{}, nil
	case "clusterrolebindings":
		return &rbacv1.ClusterRoleBinding{}, nil
	//
	case "deployments":
		return &appsv1.Deployment{}, nil
	case "daemonsets":
		return &appsv1.DaemonSet{}, nil
	case "replicasets":
		return &appsv1.ReplicaSet{}, nil
	case "statefulsets":
		return &appsv1.StatefulSet{}, nil
	case "controllerrevisions":
		return &appsv1.ControllerRevision{}, nil
	//
	case "storageclasses":
		return &storagev1.StorageClass{}, nil
	case "persistentvolumes":
		return &corev1.PersistentVolume{}, nil
	//
	case "customresourcedefinitions":
		return &apiextv1.CustomResourceDefinition{}, nil
	//
	case "jobs":
		return &batchv1.Job{}, nil
	case "cronjobs":
		return &batchv1.CronJob{}, nil
	//
	case "ingresses":
		return &networkingv1.Ingress{}, nil
	//
	case "horizontalpodautoscalers":
		return &autoscalingv2beta2.HorizontalPodAutoscaler{}, nil
	default:
		return nil, fmt.Errorf("unsupported resource type: %v", resource)
	}
}

const (
	ResourceKindConfigMap                = "configmaps"
	ResourceKindDaemonSet                = "daemonsets"
	ResourceKindDeployment               = "deployments"
	ResourceKindEvent                    = "events"
	ResourceKindHorizontalPodAutoscaler  = "horizontalpodautoscalers"
	ResourceKindIngress                  = "ingresses"
	ResourceKindJob                      = "jobs"
	ResourceKindCronJob                  = "cronjobs"
	ResourceKindLimitRange               = "limitranges"
	ResourceKindNamespace                = "namespaces"
	ResourceKindNode                     = "nodes"
	ResourceKindPersistentVolumeClaim    = "persistentvolumeclaims"
	ResourceKindPersistentVolume         = "persistentvolumes"
	ResourceKindCustomResourceDefinition = "customresourcedefinitions"
	ResourceKindPod                      = "pods"
	ResourceKindReplicaSet               = "replicasets"
	ResourceKindResourceQuota            = "resourcequota"
	ResourceKindSecret                   = "secrets"
	ResourceKindService                  = "services"
	ResourceKindStatefulSet              = "statefulsets"
	ResourceKindStorageClass             = "storageclasses"
	ResourceKindClusterRole              = "clusterroles"
	ResourceKindClusterRoleBinding       = "clusterrolebindings"
	ResourceKindRole                     = "roles"
	ResourceKindRoleBinding              = "rolebindings"
	ResourceKindWorkspace                = "workspaces"
	ResourceKindS2iBinary                = "s2ibinaries"
	ResourceKindStrategy                 = "strategy"
	ResourceKindServicePolicy            = "servicepolicies"
	ResourceKindS2iBuilderTemplate       = "s2ibuildertemplates"
	ResourceKindeS2iRun                  = "s2iruns"
	ResourceKindS2iBuilder               = "s2ibuilders"
	ResourceKindApplication              = "applications"

	WorkspaceNone = ""
	ClusterNone   = ""
)
