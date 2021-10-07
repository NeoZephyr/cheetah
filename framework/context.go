package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	r *http.Request
	w http.ResponseWriter
	ctx context.Context

	timeout bool
	writeMux *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		r: r,
		w: w,
		ctx: r.Context(),
		writeMux: &sync.Mutex{},
	}
}

func (ctx *Context) WriteMux() *sync.Mutex {
	return ctx.writeMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.r
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.w
}

func (ctx *Context) GetTimeout() bool {
	return ctx.timeout
}

func (ctx *Context) SetTimeout(timeout bool) {
	ctx.timeout = timeout
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.r.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.r != nil {
		return ctx.r.URL.Query()
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, defaultValue int) int {
	params := ctx.QueryAll()

	if vs, ok := params[key]; ok {
		l := len(vs)

		if l > 0 {
			iv, err := strconv.Atoi(vs[l - 1])

			if err != nil {
				return defaultValue
			}

			return iv
		}
	}

	return defaultValue
}

func (ctx *Context) QueryString(key string, defaultValue string) string {
	params := ctx.QueryAll()

	if vs, ok := params[key]; ok {
		l := len(vs)

		if l > 0 {
			return vs[l - 1]
		}
	}

	return defaultValue
}

func (ctx *Context) QueryArray(key string, defaultValue []string) []string {
	params := ctx.QueryAll()

	if vs, ok := params[key]; ok {
		return vs
	}

	return defaultValue
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.r != nil {
		return ctx.r.PostForm
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, defaultValue int) int {
	params := ctx.FormAll()

	if vs, ok := params[key]; ok {
		l := len(vs)

		if l > 0 {
			iv, err := strconv.Atoi(vs[l - 1])

			if err != nil {
				return defaultValue
			}

			return iv
		}
	}

	return defaultValue
}

func (ctx *Context) FormString(key string, defaultValue string) string {
	params := ctx.FormAll()

	if vs, ok := params[key]; ok {
		l := len(vs)

		if l > 0 {
			return vs[l - 1]
		}
	}

	return defaultValue
}

func (ctx *Context) FormArray(key string, defaultValue []string) []string {
	params := ctx.FormAll()

	if vs, ok := params[key]; ok {
		return vs
	}

	return defaultValue
}

func (ctx *Context) BindJson(content interface{}) error {
	if ctx.r != nil {
		body, err := ioutil.ReadAll(ctx.r.Body)

		if err != nil {
			return err
		}

		ctx.r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, content)

		if err != nil {
			return err
		}
	}

	return errors.New("request empty")
}

func (ctx *Context) Json(status int, content interface{}) error {
	if ctx.GetTimeout() {
		return nil
	}

	ctx.w.Header().Set("Content-Type", "application/json")
	ctx.w.WriteHeader(status)

	byteArr, err := json.Marshal(content)

	if err != nil {
		ctx.w.WriteHeader(500)
		return err
	}

	ctx.w.Write(byteArr)
	return nil
}

func (ctx *Context) Html(status int, content interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, content string) error {
	return nil
}

