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
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
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

	ws.Route(ws.GET("/watch/{cluster}/{resources}").
		To(r.SSEHandler).
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

func (r *resources) eventsHandler(req *restful.Request, resp *restful.Response) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	resp.AddHeader("Access-Control-Allow-Origin", "*")
	resp.AddHeader("Access-Control-Expose-Headers", "Content-Type")

	// Set headers for Server-Sent Events
	resp.AddHeader("Content-Type", "text/event-stream")
	resp.AddHeader("Cache-Control", "no-cache")
	resp.AddHeader("Connection", "keep-alive")

	// Simulate sending events (you can replace this with real data)
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		fmt.Fprintf(resp, "data: %s\n\n", fmt.Sprintf("Event %d", i))
		time.Sleep(2 * time.Second)
		resp.Flush()
	}

	// Wait for the client to close the connection
	closeNotify := resp.ResponseWriter.(http.CloseNotifier).CloseNotify()
	<-closeNotify
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

// SSEWriter 是一个封装了 http.ResponseWriter 的结构体
type SSEWriter struct {
	writer http.ResponseWriter
}

// NewSSEWriter 创建一个新的 SSEWriter
func NewSSEWriter(w http.ResponseWriter) *SSEWriter {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	return &SSEWriter{writer: w}
}

// SendMessage 发送 SSE 消息
func (sse *SSEWriter) SendMessage(data watch.EventType) error {
	message := fmt.Sprintf("data: %s\n\n", data)
	_, err := sse.writer.Write([]byte(message))
	if err != nil {
		return err
	}
	sse.writer.(http.Flusher).Flush()
	return nil
}

// SSEHandler 是处理 SSE 消息的主函数
func (r *resources) SSEHandler(req *restful.Request, resp *restful.Response) {
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
	for event := range watcher.ResultChan() {
		// 通过 SendMessage 方法发送消息
		if err := sse.SendMessage(event.Type); err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}
