package service

import (
	"context"
	"github.com/kubevela/velaux/pkg/server/interfaces/types/resp"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type ResourceService interface {
	//CreateObjectFromRawData(gvr schema.GroupVersionResource, rawData []byte) (client.Object, error)
	CreateResource(ctx context.Context, object client.Object) error
	UpdateResource(ctx context.Context, object client.Object) error
	PatchResource(ctx context.Context, object client.Object) error
	DeleteResource(ctx context.Context, gvr schema.GroupVersionResource, namespace string, name string) error
	GetResource(ctx context.Context, gvr schema.GroupVersionResource, namespace string, name string) (client.Object, error)
	ListResources(ctx context.Context, gvr schema.GroupVersionResource, namespace string, query *resp.Query) (*resp.ListResult, error)
	WatchResource(ctx context.Context, gvr schema.GroupVersionResource, opts ...client.ListOption) (watch.Interface, error)

	Get(ctx context.Context, namespace, name string, object client.Object) error
	List(ctx context.Context, namespace string, query *resp.Query, object client.ObjectList) (*resp.ListResult, error)
	Create(ctx context.Context, object client.Object) error
	Delete(ctx context.Context, object client.Object) error
	Update(ctx context.Context, old, new client.Object) error
	Patch(ctx context.Context, old, new client.Object) error
	Watch(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) (watch.Interface, error)
}

type resourceService struct {
	KubeClient  client.Client    `inject:"kubeClient"`
	WatchClient client.WithWatch `inject:"watchClient"`
}

func NewResourceService() ResourceService {
	return &resourceService{}
}

func (h *resourceService) GetResource(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string) (client.Object, error) {
	var obj client.Object
	gvk, err := h.getGVK(gvr)
	if err != nil {
		return nil, err
	}

	if h.KubeClient.Scheme().Recognizes(gvk) {
		gvkObject, err := h.KubeClient.Scheme().New(gvk)
		if err != nil {
			return nil, err
		}
		obj = gvkObject.(client.Object)
	} else {
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(gvk)
		obj = u
	}

	if err := h.Get(ctx, namespace, name, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

//func (h *resourceRepository) CreateObjectFromRawData(gvr schema.GroupVersionResource, rawData []byte) (client.Object, error) {
//	var obj client.Object
//	gvk, err := h.getGVK(gvr)
//	if err != nil {
//		return nil, err
//	}
//
//	if h.client.Scheme().Recognizes(gvk) {
//		gvkObject, err := h.client.Scheme().New(gvk)
//		if err != nil {
//			return nil, err
//		}
//		obj = gvkObject.(client.Object)
//	} else {
//		u := &unstructured.Unstructured{}
//		u.SetGroupVersionKind(gvk)
//		obj = u
//	}
//
//	err = json.Unmarshal(rawData, obj)
//	if err != nil {
//		return nil, err
//	}
//
//	// The object`s GroupVersionKind could be overridden if apiVersion and kind of rawData are different
//	// with GroupVersionKind from url, so that we should check GroupVersionKind after Unmarshal rawDate.
//	if obj.GetObjectKind().GroupVersionKind().String() != gvk.String() {
//		return nil, errors.NewBadRequest("wrong resource GroupVersionKind")
//	}
//
//	return obj, nil
//}

func (h *resourceService) ListResources(ctx context.Context, gvr schema.GroupVersionResource, namespace string, query *resp.Query) (*resp.ListResult, error) {
	var obj client.ObjectList

	gvk, err := h.getGVK(gvr)
	if err != nil {
		return nil, err
	}

	gvk = convertGVKToList(gvk)

	if h.KubeClient.Scheme().Recognizes(gvk) {
		gvkObject, err := h.KubeClient.Scheme().New(gvk)
		if err != nil {
			return nil, err
		}
		obj = gvkObject.(client.ObjectList)
	} else {
		u := &unstructured.UnstructuredList{}
		u.SetGroupVersionKind(gvk)
		obj = u
	}

	if res, err := h.List(ctx, namespace, query, obj); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h *resourceService) DeleteResource(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string) error {
	resource, err := h.GetResource(ctx, gvr, namespace, name)
	if err != nil {
		return err
	}
	return h.Delete(ctx, resource)
}

func (h *resourceService) UpdateResource(ctx context.Context, object client.Object) error {
	old := object.DeepCopyObject().(client.Object)
	err := h.Get(ctx, object.GetNamespace(), object.GetName(), old)
	if err != nil {
		return err
	}

	return h.Update(ctx, old, object)
}

func (h *resourceService) PatchResource(ctx context.Context, object client.Object) error {
	old := object.DeepCopyObject().(client.Object)
	err := h.Get(ctx, object.GetNamespace(), object.GetName(), old)
	if err != nil {
		return err
	}

	return h.Patch(ctx, old, object)
}

func (h *resourceService) CreateResource(ctx context.Context, object client.Object) error {
	return h.Create(ctx, object)
}

func (h *resourceService) WatchResource(ctx context.Context, gvr schema.GroupVersionResource, opts ...client.ListOption) (watch.Interface, error) {
	var listObj client.ObjectList

	gvk, err := h.getGVK(gvr)
	if err != nil {
		return nil, err
	}

	gvk = convertGVKToList(gvk)
	if h.KubeClient.Scheme().Recognizes(gvk) {
		gvkObject, err := h.KubeClient.Scheme().New(gvk)
		if err != nil {
			return nil, err
		}
		listObj = gvkObject.(client.ObjectList)
	} else {
		u := &unstructured.UnstructuredList{}
		u.SetGroupVersionKind(gvk)
		listObj = u
	}

	return h.WatchClient.Watch(ctx, listObj, opts...)
}

func convertGVKToList(gvk schema.GroupVersionKind) schema.GroupVersionKind {
	if strings.HasSuffix(gvk.Kind, "List") {
		return gvk
	}
	gvk.Kind = gvk.Kind + "List"
	return gvk
}

func (h *resourceService) getGVK(gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	var (
		gvk schema.GroupVersionKind
		err error
	)
	gvk, err = h.KubeClient.RESTMapper().KindFor(gvr)
	if err != nil {
		return gvk, err
	}
	return gvk, nil
}

func (h *resourceService) Get(ctx context.Context, namespace, name string, object client.Object) error {
	return h.KubeClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, object)
}

func (h *resourceService) List(ctx context.Context, namespace string, query *resp.Query, list client.ObjectList) (*resp.ListResult, error) {
	listOpt := &client.ListOptions{
		LabelSelector: query.Selector(),
		Namespace:     namespace,
	}

	err := h.KubeClient.List(ctx, list, listOpt)
	if err != nil {
		return nil, err
	}

	extractList, err := meta.ExtractList(list)
	if err != nil {
		return nil, err
	}

	filtered, total := DefaultList(extractList, query, compareList, filter)
	return &resp.ListResult{
		TotalItems: total,
		Items:      ObjectsToInterfaces(filtered),
	}, nil

	//list.SetRemainingItemCount(remainingItemCount)
	//if err := meta.SetList(list, filtered); err != nil {
	//	return nil, err
	//}
	//return nil, nil
}

func (h *resourceService) Create(ctx context.Context, object client.Object) error {
	return h.KubeClient.Create(ctx, object)
}

func (h *resourceService) Delete(ctx context.Context, object client.Object) error {
	return h.KubeClient.Delete(ctx, object)
}

func (h *resourceService) Update(ctx context.Context, old, new client.Object) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return h.KubeClient.Update(ctx, new)
}

func (h *resourceService) Patch(ctx context.Context, old, new client.Object) error {
	return h.KubeClient.Patch(ctx, new, client.MergeFrom(old))
}

func (h *resourceService) Watch(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) (watch.Interface, error) {
	return h.WatchClient.Watch(ctx, obj, opts...)
}

func compareList(left, right runtime.Object, field resp.Field) bool {
	l, err := meta.Accessor(left)
	if err != nil {
		return false
	}
	r, err := meta.Accessor(right)
	if err != nil {
		return false
	}
	return DefaultObjectMetaCompare(l, r, field)
}

func filter(object runtime.Object, filter resp.Filter) bool {
	o, err := meta.Accessor(object)
	if err != nil {
		return false
	}
	return DefaultObjectMetaFilter(o, filter)
}

func ObjectsToInterfaces(objs []runtime.Object) []runtime.Object {
	res := make([]runtime.Object, 0)
	for _, obj := range objs {
		res = append(res, obj)
	}
	return res
}
