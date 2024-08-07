package api

import (
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
	ws.Path(ApiRootPath+GroupVersion.String()).
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("api for user manage")

	tags := []string{"ClusteredResource"}

	ws.Route(ws.GET("/{cluster}/{resources}").
		To(r.listResources).
		Doc("Cluster level resources").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// TODO
		//Filter(c.RbacService.CheckPerm("user", "list")).
		Param(ws.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(ws.PathParameter("cluster", "cluster name, e.g. local.").DefaultValue("local")).
		Param(ws.QueryParameter(resp.ParameterName, "name used to do filtering").Required(false)).
		Param(ws.QueryParameter(resp.ParameterPage, "page").Required(false).DataFormat("page=%d").DefaultValue("page=1")).
		Param(ws.QueryParameter(resp.ParameterLimit, "limit").Required(false)).
		Param(ws.QueryParameter(resp.ParameterAscending, "sort parameters, e.g. reverse=true").Required(false).DefaultValue("ascending=false")).
		Param(ws.QueryParameter(resp.ParameterOrderBy, "sort parameters, e.g. orderBy=createTime")).
		//Filter(r.RbacService.CheckPerm("role", "list")).
		Returns(http.StatusOK, "OK", resp.ListResult{}).
		Returns(http.StatusBadRequest, "Bad Request", bcode.Bcode{}).
		Writes(resp.ListResult{}))

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
