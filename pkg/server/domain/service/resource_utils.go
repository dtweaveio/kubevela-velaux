package service

import (
	"encoding/json"
	"fmt"
	"github.com/kubevela/velaux/pkg/server/interfaces/types/resp"
	"github.com/kubevela/velaux/pkg/server/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/klog/v2"
	"sort"
	"strings"
)

// CompareFunc return true is left greater than right
type CompareFunc func(runtime.Object, runtime.Object, resp.Field) bool

type FilterFunc func(runtime.Object, resp.Filter) bool

type TransformFunc func(runtime.Object) runtime.Object

func DefaultList(objects []runtime.Object, q *resp.Query, compareFunc CompareFunc, filterFunc FilterFunc, transformFuncs ...TransformFunc) ([]runtime.Object, int) {
	// selected matched ones
	var filtered []runtime.Object
	if len(q.Filters) != 0 {
		for _, object := range objects {
			selected := true
			for field, value := range q.Filters {
				if !filterFunc(object, resp.Filter{Field: field, Value: value}) {
					selected = false
					break
				}
			}

			if selected {
				for _, transform := range transformFuncs {
					object = transform(object)
				}
				filtered = append(filtered, object)
			}
		}
	} else {
		filtered = objects
	}

	// sort by sortBy field
	sort.Slice(filtered, func(i, j int) bool {
		if !q.Ascending {
			return compareFunc(filtered[i], filtered[j], q.SortBy)
		}
		return !compareFunc(filtered[i], filtered[j], q.SortBy)
	})

	total := len(filtered)

	if q.Pagination == nil {
		q.Pagination = resp.NoPagination
	}

	start, end := q.Pagination.GetValidPagination(total)
	//remainingItemCount := int64(total - end)

	return filtered[start:end], total
}

// DefaultObjectMetaCompare return true is left greater than right
func DefaultObjectMetaCompare(left, right metav1.Object, sortBy resp.Field) bool {
	switch sortBy {
	// ?sortBy=name
	case resp.FieldName:
		return strings.Compare(left.GetName(), right.GetName()) > 0
	//	?sortBy=creationTimestamp
	default:
		fallthrough
	case resp.FieldCreateTime:
		fallthrough
	case resp.FieldCreationTimeStamp:
		// compare by name if creation timestamp is equal
		ltime := left.GetCreationTimestamp()
		rtime := right.GetCreationTimestamp()
		if ltime.Equal(&rtime) {
			return strings.Compare(left.GetName(), right.GetName()) > 0
		}
		return left.GetCreationTimestamp().After(right.GetCreationTimestamp().Time)
	}
}

// Default metadata filter
func DefaultObjectMetaFilter(item metav1.Object, filter resp.Filter) bool {
	switch filter.Field {
	case resp.FieldNames:
		for _, name := range strings.Split(string(filter.Value), ",") {
			if item.GetName() == name {
				return true
			}
		}
		return false
	// /namespaces?page=1&limit=10&name=default
	case resp.FieldName:
		return strings.Contains(item.GetName(), string(filter.Value))
		// /namespaces?page=1&limit=10&uid=a8a8d6cf-f6a5-4fea-9c1b-e57610115706
	case resp.FieldUID:
		return strings.Compare(string(item.GetUID()), string(filter.Value)) == 0
		// /deployments?page=1&limit=10&namespace=kubesphere-system
	case resp.FieldNamespace:
		return strings.Compare(item.GetNamespace(), string(filter.Value)) == 0
		// /namespaces?page=1&limit=10&ownerReference=a8a8d6cf-f6a5-4fea-9c1b-e57610115706
	case resp.FieldOwnerReference:
		for _, ownerReference := range item.GetOwnerReferences() {
			if strings.Compare(string(ownerReference.UID), string(filter.Value)) == 0 {
				return true
			}
		}
		return false
		// /namespaces?page=1&limit=10&ownerKind=Workspace
	case resp.FieldOwnerKind:
		for _, ownerReference := range item.GetOwnerReferences() {
			if strings.Compare(ownerReference.Kind, string(filter.Value)) == 0 {
				return true
			}
		}
		return false
		// /namespaces?page=1&limit=10&annotation=openpitrix_runtime
	case resp.FieldAnnotation:
		return labelMatch(item.GetAnnotations(), string(filter.Value))
		// /namespaces?page=1&limit=10&label=kubesphere.io/workspace:system-workspace
	case resp.FieldLabel:
		return labelMatch(item.GetLabels(), string(filter.Value))
	case resp.ParameterFieldSelector:
		return contains(item.(runtime.Object), filter.Value)

	//case query.FieldNodeName:
	//	return strings.Contains(item..GetName(), string(filter.Value))
	//return pod.Spec.NodeName == string(filter.Value)
	default:
		return false
	}
}

func labelMatch(labels map[string]string, filter string) bool {
	fields := strings.SplitN(filter, "=", 2)
	var key, value string
	var opposite bool
	if len(fields) == 2 {
		key = fields[0]
		if strings.HasSuffix(key, "!") {
			key = strings.TrimSuffix(key, "!")
			opposite = true
		}
		value = fields[1]
	} else {
		key = fields[0]
		value = "*"
	}
	for k, v := range labels {
		if opposite {
			if (k == key) && v != value {
				return true
			}
		} else {
			if (k == key) && (value == "*" || v == value) {
				return true
			}
		}
	}
	return false
}

// implement a generic query filter to support multiple field selectors with "jsonpath.JsonPathLookup"
// https://github.com/oliveagle/jsonpath/blob/master/readme.md
func contains(object runtime.Object, queryValue resp.Value) bool {
	// call the ParseSelector function of "k8s.io/apimachinery/pkg/fields/selector.go" to validate and parse the selector
	fieldSelector, err := fields.ParseSelector(string(queryValue))
	if err != nil {
		klog.V(4).Infof("failed parse selector error: %s", err)
		return false
	}
	for _, requirement := range fieldSelector.Requirements() {
		var negative bool
		// supports '=', '==' and '!='.(e.g. ?fieldSelector=key1=value1,key2=value2)
		// fields.ParseSelector(FieldSelector) has handled the case where the operator is '==' and converted it to '=',
		// so case selection.DoubleEquals can be ignored here.
		switch requirement.Operator {
		case selection.NotEquals:
			negative = true
		case selection.Equals:
			negative = false
		}
		key := requirement.Field
		value := requirement.Value

		var input map[string]interface{}
		data, err := json.Marshal(object)
		if err != nil {
			klog.V(4).Infof("failed marshal to JSON string: %s", err)
			return false
		}
		if err = json.Unmarshal(data, &input); err != nil {
			klog.V(4).Infof("failed unmarshal to map object: %s", err)
			return false
		}
		rawValue, err := utils.JsonPathLookup(input, "$."+key)
		if err != nil {
			klog.V(4).Infof("failed to lookup jsonpath: %s", err)
			return false
		}
		if (negative && fmt.Sprintf("%v", rawValue) != value) || (!negative && fmt.Sprintf("%v", rawValue) == value) {
			continue
		} else {
			return false
		}
	}
	return true
}
