package nullgo

import (
	"net/http"
	"net/url"
)

type Context struct {
	Params         Params
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	//query查询的缓存
	queryCache     url.Values

	//post表单的缓存
	formCache      url.Values
	config 			WebSocketConfig
}

type Params []Param

type Param struct {
	key   string
	value string
}


func (c *Context) Param(name string) string {
	return c.Params.getValue(name)
}

func (p Params) getValue(name string) string {
	value, _ := p.getV(name)
	return value
}

func (p Params) getV(name string) (string, bool) {
	for _, param := range p {
		if param.key == name {
			value := param.value
			return value, true
		}
	}
	return "", false
}

func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

func (c *Context) GetQuery(key string) (string,bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.getQueryCache()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *Context) getQueryCache() {
	if c.queryCache == nil {
		c.queryCache = make(url.Values)
		c.queryCache, _ = url.ParseQuery(c.Request.URL.RawQuery)
	}
}

func (c *Context) PostV(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

func (c *Context) GetPostForm(key string) (string,bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.getFormCache()
	if values, ok := c.formCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *Context) getFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.Request
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			if err != http.ErrNotMultipart {
				//TODO log print
				print("error on parse multipart form array: %v", err)
				//debugPrint("error on parse multipart form array: %v", err)
			}
		}
		c.formCache = req.PostForm
	}
}

func (c *Context)String(format string)  {
	c.ResponseWriter.Write(QuickStringToBytes(format))
	c.ResponseWriter.Write(QuickStringToBytes("\n"))
}
