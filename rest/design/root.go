package design

import (
	"regexp"

	"goa.design/goa.v2/design"
)

var (
	// Root holds the root expression built on process initialization.
	Root = &RootExpr{}

	// WildcardRegex is the regular expression used to capture path parameters.
	WildcardRegex = regexp.MustCompile(`/(?::|\*)([a-zA-Z0-9_]+)`)

	// ErrorMedia is the built-in media type for error responses.
	ErrorMedia = design.ErrorMedia
)

const (
	// DefaultView is the name of the default view.
	DefaultView = "default"
)

const (
	// SchemeHTTP denotes an API that uses HTTP.
	SchemeHTTP APIScheme = "http"

	// SchemeHTTPS denotes an API that uses HTTPS.
	SchemeHTTPS = "https"

	// SchemeWS denotes an API that uses websocket.
	SchemeWS = "ws"

	// SchemeWSS denotes an API that uses secure websocket.
	SchemeWSS = "wss"
)

type (
	// RootExpr is the data structure built by the top level HTTP DSL.
	RootExpr struct {
		// Path is the common request path prefix to all the service
		// HTTP endpoints.
		Path string
		// Params defines common request parameters to all the service
		// HTTP endpoints.
		Params *AttributeMapExpr
		// Schemes is the supported API URL schemes
		Schemes []APIScheme
		// Consumes lists the mime types supported by the API controllers
		Consumes []*EncodingExpr
		// Produces lists the mime types generated by the API controllers
		Produces []*EncodingExpr
		// Resources contains the resources created by the DSL.
		Resources []*ResourceExpr
		// Responses available to all API actions.
		Responses []*HTTPResponseExpr
		// Built-in responses
		DefaultResponses []*HTTPResponseExpr
	}

	// APIScheme lists the possible values for Scheme
	APIScheme string
)

// EvalName is the expression name used by the evaluation engine to display
// error messages.
func (r *RootExpr) EvalName() string {
	return "API HTTP"
}

// Response returns the response with the given name if any.
func (r *RootExpr) Response(name string) *HTTPResponseExpr {
	for _, resp := range r.Responses {
		if resp.Name == name {
			return resp
		}
	}
	return nil
}

// DefaultResponse returns the default response with the given name if any.
func (r *RootExpr) DefaultResponse(name string) *HTTPResponseExpr {
	for _, resp := range r.DefaultResponses {
		if resp.Name == name {
			return resp
		}
	}
	return nil
}

// Resource returns the resource with the given name if any.
func (r *RootExpr) Resource(name string) *ResourceExpr {
	for _, res := range r.Resources {
		if res.Name == name {
			return res
		}
	}
	return nil
}

// ResourceFor creates a new or returns the existing resource definition for the
// given service.
func (r *RootExpr) ResourceFor(s *design.ServiceExpr) *ResourceExpr {
	if res := r.Resource(s.Name); res != nil {
		return res
	}
	res := &ResourceExpr{
		ServiceExpr: s,
		Actions:     make([]*ActionExpr, len(s.Endpoints)),
	}
	for i, e := range s.Endpoints {
		res.Actions[i] = &ActionExpr{
			EndpointExpr: e,
			Resource:     res,
		}
	}
	r.Resources = append(r.Resources, res)
	return res
}

// ExtractWildcards returns the names of the wildcards that appear in path.
func ExtractWildcards(path string) []string {
	matches := WildcardRegex.FindAllStringSubmatch(path, -1)
	wcs := make([]string, len(matches))
	for i, m := range matches {
		wcs[i] = m[1]
	}
	return wcs
}
