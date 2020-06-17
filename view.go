package razlink

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// View is something that a PageHandler returns and is capable of rendering a page
type View struct {
	StatusCode int
	Error      error
	Data       interface{}
	renderer   func(w http.ResponseWriter)
}

// Render ...
func (view *View) Render(w http.ResponseWriter) {
	view.renderer(w)
}

// ViewOption ...
type ViewOption func(view *View)

// WithError ...
func WithError(err error, errcode int) ViewOption {
	return func(view *View) {
		view.Error = err
		view.StatusCode = errcode
	}
}

// WithErrorMessage ...
func WithErrorMessage(errmsg string, errcode int) ViewOption {
	return WithError(fmt.Errorf("%s", errmsg), errcode)
}

// WithData ...
func WithData(data interface{}) ViewOption {
	return func(view *View) {
		view.Data = data
	}
}

var errViewRenderer, _ = BindLayout("<strong>{{.}}</strong>", nil, nil, nil)

// ErrorView returns a View that represents an error
func ErrorView(r *http.Request, errmsg string, errcode int, opts ...ViewOption) *View {
	renderer := func(w http.ResponseWriter) {
		w.WriteHeader(errcode)
		errViewRenderer(w, r, errmsg, errmsg)
	}
	v := &View{
		StatusCode: errcode,
		Error:      fmt.Errorf("%s", errmsg),
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// ErrorView ...
func (r *PageRequest) ErrorView(errmsg string, errcode int, opts ...ViewOption) *View {
	return ErrorView(r.Request, errmsg, errcode, opts...)
}

// EmbedView returns a View that embeds the given website
func EmbedView(url string, opts ...ViewOption) *View {
	renderer := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<iframe src="%s" style="position:fixed; top:0; left:0; bottom:0; right:0; width:100%%; height:100%%; border:none; margin:0; padding:0; overflow:hidden; z-index:999999;"></iframe>`, url)
	}
	v := &View{
		StatusCode: http.StatusOK,
		Data:       url,
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// EmbedView ...
func (r *PageRequest) EmbedView(url string, opts ...ViewOption) *View {
	return EmbedView(url, opts...)
}

// RedirectView ...
func RedirectView(r *http.Request, url string, opts ...ViewOption) *View {
	renderer := func(w http.ResponseWriter) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	v := &View{
		StatusCode: http.StatusSeeOther,
		Data:       url,
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// RedirectView ...
func (r *PageRequest) RedirectView(url string, opts ...ViewOption) *View {
	return RedirectView(r.Request, url, opts...)
}

// CookieAndRedirectView ...
func CookieAndRedirectView(r *http.Request, cookie *http.Cookie, url string, opts ...ViewOption) *View {
	renderer := func(w http.ResponseWriter) {
		http.SetCookie(w, cookie)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	v := &View{
		StatusCode: http.StatusSeeOther,
		Data:       cookie,
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// CookieAndRedirectView ...
func (r *PageRequest) CookieAndRedirectView(cookie *http.Cookie, url string, opts ...ViewOption) *View {
	return CookieAndRedirectView(r.Request, cookie, url, opts...)
}

// CopyView ...
func CopyView(resp *http.Response, opts ...ViewOption) *View {
	bytes, _ := ioutil.ReadAll(resp.Body)
	renderer := func(w http.ResponseWriter) {
		for k, v := range resp.Header {
			w.Header().Set(k, v[0])
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(bytes)
	}
	v := &View{
		StatusCode: http.StatusOK,
		Data:       resp,
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// CopyView ...
func (r *PageRequest) CopyView(resp *http.Response, opts ...ViewOption) *View {
	return CopyView(resp, opts...)
}

// AsyncCopyView ...
func AsyncCopyView(resp *http.Response, opts ...ViewOption) *View {
	renderer := func(w http.ResponseWriter) {
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header().Set(k, v[0])
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
	v := &View{
		StatusCode: http.StatusOK,
		Data:       resp,
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// AsyncCopyView ...
func (r *PageRequest) AsyncCopyView(resp *http.Response, opts ...ViewOption) *View {
	return AsyncCopyView(resp, opts...)
}

// HandlerView ...
func HandlerView(r *http.Request, handler http.HandlerFunc, opts ...ViewOption) *View {
	renderer := func(w http.ResponseWriter) {
		handler(w, r)
	}
	v := &View{
		StatusCode: http.StatusOK,
		renderer:   renderer,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// HandlerView ...
func (r *PageRequest) HandlerView(handler http.HandlerFunc, opts ...ViewOption) *View {
	return HandlerView(r.Request, handler, opts...)
}
