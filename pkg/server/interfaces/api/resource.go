package api

import (
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/kubevela/pkg/multicluster"
	"github.com/kubevela/velaux/pkg/server/domain/service"
	"github.com/kubevela/velaux/pkg/server/interfaces/types"
	"github.com/kubevela/velaux/pkg/server/interfaces/types/resp"
	"github.com/kubevela/velaux/pkg/server/utils/bcode"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ApiRootPath = "/kapis/"
	GroupName   = "resources"
)

type resources struct {
	RbacService      service.RBACService     `inject:""`
	ResourcesService service.ResourceService `inject:""`
}

func NewResources() Interface {
	return &resources{}
}

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

func (r *resources) GetWebServiceRoute() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(ApiRootPath + GroupVersion.String()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("api for user manage")

	tags := []string{"ClusteredResource"}

	ws.Route(ws.GET("/{cluster}/{resources}").
		To(r.listResources).
		Doc("Cluster level resources").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// TODO
		//Filter(c.RbacService.CheckPerm("user", "list")).
		Param(ws.PathParameter("cluster", "cluster name, e.g. local.").DefaultValue("local")).
		Param(ws.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(ws.QueryParameter(resp.ParameterName, "name used to do filtering").Required(false)).
		Param(ws.QueryParameter(resp.ParameterPage, "page").Required(false).DataFormat("page=%d").DefaultValue("page=1")).
		Param(ws.QueryParameter(resp.ParameterLimit, "limit").Required(false)).
		Param(ws.QueryParameter(resp.ParameterAscending, "sort parameters, e.g. reverse=true").Required(false).DefaultValue("ascending=false")).
		Param(ws.QueryParameter(resp.ParameterOrderBy, "sort parameters, e.g. orderBy=createTime")).
		//Filter(r.RbacService.CheckPerm("role", "list")).
		Returns(http.StatusOK, "OK", resp.ListResult{}).
		Returns(http.StatusBadRequest, "Bad Request", bcode.Bcode{}).
		Writes(resp.ListResult{}))

	ws.Route(ws.GET("/{cluster}/namespaces/{namespace}/{resources}").
		To(r.listResources).
		Metadata(restfulspec.KeyOpenAPITags, tags).Doc("Namespace level resources").
		Param(ws.PathParameter("cluster", "cluster name, e.g. local.").DefaultValue("local")).
		Param(ws.PathParameter("namespace", "the name of the namespace")).
		Param(ws.PathParameter("resources", "namespace level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(ws.QueryParameter(resp.ParameterName, "name used to do filtering").Required(false)).
		Param(ws.QueryParameter(resp.ParameterPage, "page").Required(false).DataFormat("page=%d").DefaultValue("page=1")).
		Param(ws.QueryParameter(resp.ParameterLimit, "limit").Required(false)).
		Param(ws.QueryParameter(resp.ParameterAscending, "sort parameters, e.g. reverse=true").Required(false).DefaultValue("ascending=false")).
		Param(ws.QueryParameter(resp.ParameterOrderBy, "sort parameters, e.g. orderBy=createTime")).
		//Filter(r.RbacService.CheckPerm("role", "list")).
		Returns(http.StatusOK, "OK", resp.ListResult{}).
		Returns(http.StatusBadRequest, "Bad Request", bcode.Bcode{}).
		Writes(resp.ListResult{}))

	ws.Route(ws.GET("/{cluster}/namespaces/{namespace}/{resources}/{name}").
		To(r.detailResource).
		Doc("Namespace level resource").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("cluster", "cluster name, e.g. local.").DefaultValue("local")).
		Param(ws.PathParameter("namespace", "the name of the namespace")).
		Param(ws.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(ws.PathParameter("name", "name of resource instance")))

	ws.Route(ws.POST("/{cluster}/{resources}").
		To(r.createResource).
		Doc("Resource create").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("cluster", "cluster name, e.g. local.").DefaultValue("local")).
		Param(ws.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Returns(http.StatusOK, "OK", nil).
		Returns(http.StatusBadRequest, "Bad Request", bcode.Bcode{}))

	ws.Route(ws.GET("/watch/{cluster}/{resources}").
		To(r.watchResource).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Cluster level resources watch").
		Param(ws.PathParameter("cluster", "cluster name, e.g. local.").DefaultValue("local")).
		Param(ws.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Returns(http.StatusOK, "OK", nil).
		Returns(http.StatusBadRequest, "Bad Request", bcode.Bcode{}))

	// TODO
	//ws.Filter(authCheckFilter)
	return ws
}

func (r *resources) listResources(req *restful.Request, res *restful.Response) {
	query := resp.ParseQueryParameter(req)
	resources := req.PathParameter("resources")
	namespace := req.PathParameter("namespace")
	cluster := req.PathParameter("cluster")

	resourceType := types.SupportedResources[resources]
	if resourceType.Empty() {
		klog.Errorf("%s, resource type: %s", bcode.ErrResourceNotSupported.Message, resourceType)
		bcode.ReturnError(req, res, bcode.ErrResourceNotSupported)
		return
	}

	ctx := req.Request.Context()
	if "" != cluster {
		ctx = multicluster.WithCluster(ctx, cluster)
	}

	result, err := r.ResourcesService.ListResources(ctx, resourceType, namespace, query)
	if err != nil {
		bcode.ReturnError(req, res, err)
		return
	}

	if err := res.WriteEntity(result); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
}

func (r *resources) detailResource(req *restful.Request, res *restful.Response) {
	ctx := req.Request.Context()
	cluster := req.PathParameter("cluster")
	resources := req.PathParameter("resources")
	namespace := req.PathParameter("namespace")
	name := req.PathParameter("name")

	resourceType := types.SupportedResources[resources]
	if resourceType.Empty() {
		klog.Errorf("%s, resource type: %s", bcode.ErrResourceNotSupported.Message, resourceType)
		bcode.ReturnError(req, res, bcode.ErrResourceNotSupported)
		return
	}

	if "" != cluster {
		ctx = multicluster.WithCluster(ctx, cluster)
	}

	result, err := r.ResourcesService.GetResource(ctx, resourceType, namespace, name)
	if err != nil {
		bcode.ReturnError(req, res, err)
		return
	}

	if err := res.WriteEntity(result); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
}

func (r *resources) createResource(req *restful.Request, res *restful.Response) {
	ctx := req.Request.Context()
	cluster := req.PathParameter("cluster")
	resources := req.PathParameter("resources")

	// 动态创建资源对象
	object, err := types.CreateResourceObject(resources)
	if err != nil {
		klog.Errorf("Failed to create object for resource type: %s, error: %v", resources, err)
		bcode.ReturnError(req, res, bcode.ErrResourceNotSupported)
		return
	}

	err = req.ReadEntity(object)
	if err != nil {
		klog.Errorf("%s, resource type: %s", bcode.ErrResourceReadEntity.Message, resources)
		bcode.ReturnError(req, res, bcode.ErrResourceReadEntity)
		return
	}

	if cluster != "" {
		ctx = multicluster.WithCluster(ctx, cluster)
	}

	err = r.ResourcesService.CreateResource(ctx, object)
	if err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
}

// SSEHandler 是处理 SSE 消息的主函数
func (r *resources) watchResource(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()
	cluster := req.PathParameter("cluster")
	resources := req.PathParameter("resources")

	resourceType := types.SupportedResources[resources]
	if resourceType.Empty() {
		klog.Errorf("%s, resource type: %s", bcode.ErrResourceNotSupported.Message, resourceType)
		bcode.ReturnError(req, resp, bcode.ErrResourceNotSupported)
		return
	}

	if "" != cluster {
		ctx = multicluster.WithCluster(ctx, cluster)
	}

	var opts []client.ListOption
	watcher, err := r.ResourcesService.WatchResource(ctx, resourceType, opts...)
	if err != nil {
		bcode.ReturnError(req, resp, err)
	}
	defer watcher.Stop()

	sse := NewSSEWriter(resp.ResponseWriter)
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			// 通过 SendMessage 方法发送消息
			if err := sse.SendMessage(event.Type); err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}
