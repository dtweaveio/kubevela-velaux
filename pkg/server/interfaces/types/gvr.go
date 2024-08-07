package types

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
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
