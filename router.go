package nullgo

import (
	"fmt"
	"github.com/unknwon/com"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const (
	GET = iota
	POST
	PUT
	DELETE
)

type Router struct {
	ctx *Context
	regex  map[string]*regexp.Regexp
	mu     sync.RWMutex
	router []handlerMap
	wsMap  map[string]WebSocketConfig
}

var key = ""
var k = 1


type handlerMap map[string]HandlerFunc

type HandlerFunc func(*Context)

var orderMap = make(map[string]int)

var methodMap = map[string]int{
	"GET":    GET,
	"POST":   POST,
	"PUT":    PUT,
	"DELETE": DELETE,
}

var methodStringMap = map[int]string{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	DELETE: "DELETE",
}

func (r *Router) Init() {
	r.router = []handlerMap{
		GET:    make(map[string]HandlerFunc),
		POST:   make(map[string]HandlerFunc),
		PUT:    make(map[string]HandlerFunc),
		DELETE: make(map[string]HandlerFunc),
	}
	//r.params = make(Params, 10)
	r.wsMap = make(map[string]WebSocketConfig, 1)
	r.regex = make(map[string]*regexp.Regexp, 10)
	r.ctx = &Context{Params:make(Params,20)}
}

func (r *Router) add(method int, uri string, handle HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if uri == "" {
		//TODO log
		panic("http: invalid path")
		return
	}

	if handle == nil {
		//TODO log
		panic("http: nil handler")
		return
	}

	//①检查uri中是否需要进行参数或者正则匹配
	//	如：	/user/:id([1-9]+)
	parts := strings.Split(uri, "/")
	var params Params
	var expr string

	params = make(Params, 10)

	for i, part := range parts {

		if strings.HasPrefix(part, ":") {
			key = com.ToStr(parts[1:len(parts)-1])
			expr = "([^/]+)"
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
				fmt.Println(part)
			}

			params[k].key = strings.TrimLeft(part, ":")
			parts[i] = ""

		} else {
			key = com.ToStr(parts[1:])
		}

	}

	//②对uri进行重组,并组装正则
	// 如:   /user/([1-9]+)                        /user/([^/]+)
	uri = strings.Join(parts, "/")
	pattern := uri + expr
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
		return
	}
	regexMap := r.regex
	regexMap[key] = regex
	//r.params = params
	ctx := r.ctx
	ctx.Params[k].key = params[k].key
	//③对uri再次重组
	//如： /user/:
	length := len(uri)
	if uri[length-1] == '/' {
		uri += ":"
	}
	order := orderMap
	order[key] = k
	k++

	//④检查uri是否已经被注册
	for _, handlerMaps := range r.router {
		if _, exist := handlerMaps[uri]; exist {
			panic("http: multiple registrations for " + uri)
			return
		}
	}

	//⑤把uri和路由绑定
	handlerMaps := r.router[method]
	handlerMaps[uri] = handle
	Trace("Registered route： %s      %s\n", methodStringMap[method], uri)
}

func (r *Router) forward(w http.ResponseWriter, req *http.Request) {
	ctx := r.ctx







	ctx.ResponseWriter = w
	ctx.Request = req
	rawURI := req.RequestURI
	method := strings.ToUpper(req.Method)
	upgrade := req.Header.Get("Upgrade")
	//Trace("Access route： %s     %s\n",method, rawURI)
	methodToInt := methodMap[method]
	arr := strings.Split(rawURI, "?")
	uri := arr[0]
	array := strings.Split(uri,"/")
	key := com.ToStr(array[1:len(array)-1])
	k := orderMap[key]
    if regex, ok := r.regex[key]; ok {
		if ok := regex.MatchString(uri); ok {
			matches := regex.FindStringSubmatch(uri)
			if len(matches[0]) == len(uri) {
				for _, match := range matches[1:] {

					ctx.Params[k].value = match
					uri = strings.Replace(uri, match, ":", -1)
				}
			}
		}
	}

	var find = false
	for i, handlerMaps := range r.router {
		if handle, ok := handlerMaps[uri]; ok {
			find = true
			if methodToInt != i {
				w.WriteHeader(405)
				w.Write(QuickStringToBytes("Method not allowed"))
				printError("| 405 |          | %s            %s \n", method, rawURI)
			}
			//ctx.Params = r.params

			//如果http协议需要升级为websocket
			if upgrade == "websocket" {
				c, ok := r.wsMap[uri]
				if !ok {
					w.Write(QuickStringToBytes("404 pages not found"))
					//TODO log
					return
				}
				ctx.config = c
				handle(ctx)
				return
			}
			printDebug("| 200 |          | %s            %s \n", method, rawURI)
			handle(ctx)
		}
	}

	if !find {
		w.WriteHeader(404)
		w.Write(QuickStringToBytes("404 not found"))
		printError("| 404 |          | %s            %s \n", method, rawURI)
	}

}
