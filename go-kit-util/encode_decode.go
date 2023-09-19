package go_kit_util

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	gohttp "github.com/likearthian/go-http"
	"github.com/likearthian/go-http/router"
)

type requestDecoderOption struct {
	acceptedFields  map[string]struct{}
	urlParamsGetter func(context.Context) map[string]string
}

type RequestDecoderOption func(d *requestDecoderOption)

func MakeCommonGetRequestDecoder(output reflect.Type, options ...RequestDecoderOption) httptransport.DecodeRequestFunc {
	modelTyp := output
	if modelTyp.Kind() == reflect.Ptr {
		modelTyp = modelTyp.Elem()
	}

	opts := requestDecoderOption{
		acceptedFields: make(map[string]struct{}),
		urlParamsGetter: func(ctx context.Context) map[string]string {
			return make(map[string]string)
		},
	}

	for _, op := range options {
		op(&opts)
	}

	for i := 0; i < modelTyp.NumField(); i++ {
		tag := modelTyp.Field(i).Tag.Get("query")
		if tag != "" {
			opts.acceptedFields[tag] = struct{}{}
		}
	}

	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		obj := reflect.Zero(output)
		pv := reflect.New(obj.Type())
		pv.Elem().Set(obj)

		params := opts.urlParamsGetter(ctx)
		query := r.URL.Query()

		//include params into query to be parsed
		for k, v := range params {
			query.Set(k, v)
		}

		for field := range query {
			if _, ok := opts.acceptedFields[field]; !ok {
				return nil, fmt.Errorf("%w: unknown field '%s'", fmt.Errorf("bad request"), field)
			}
		}

		if err := gohttp.BindURLQuery(pv.Interface(), query); err != nil {
			return nil, err
		}

		return pv.Elem().Interface(), nil
	}
}

func MakeCommonPostRequestDecoder(output reflect.Type, options ...RequestDecoderOption) httptransport.DecodeRequestFunc {
	opts := requestDecoderOption{
		acceptedFields: make(map[string]struct{}),
		urlParamsGetter: func(ctx context.Context) map[string]string {
			return make(map[string]string)
		},
	}

	for _, op := range options {
		op(&opts)
	}

	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		obj := reflect.Zero(output)
		pv := reflect.New(obj.Type())
		pv.Elem().Set(obj)

		params := opts.urlParamsGetter(ctx)
		query := r.URL.Query()
		//include params into query to be parsed
		for k, v := range params {
			query.Set(k, v)
		}
		err := json.NewDecoder(r.Body).Decode(pv.Interface())
		if err != nil {
			return nil, fmt.Errorf("%w: %s", fmt.Errorf("bad request"), err)
		}

		if err := gohttp.BindURLQuery(pv.Interface(), query); err != nil {
			return nil, err
		}

		return pv.Elem().Interface(), nil
	}
}

func CommonJSONResponseEncoder(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(gohttp.HeaderContentType, gohttp.HttpContentTypeJson)
	var gw io.Writer = w
	if needGzipped(ctx) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gw = gz
	}

	return json.NewEncoder(gw).Encode(response)
}

func CommonByteResponseEncoder(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	buf, ok := response.([]byte)
	if !ok {
		return fmt.Errorf("response format for commonByteResponseEncoder is not []byte")
	}

	var gw io.Writer = w

	if needGzipped(ctx) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gw = gz
	}

	_, err := gw.Write(buf)
	return err
}

func MakeCommonByteResponseEncoder(contentType string) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set(gohttp.HeaderContentType, contentType)
		buf, ok := response.([]byte)
		if !ok {
			return fmt.Errorf("response format for commonByteResponseEncoder is not []byte")
		}

		var gw io.Writer = w

		if needGzipped(ctx) {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gw = gz
		}

		_, err := gw.Write(buf)
		return err
	}
}

func needGzipped(ctx context.Context) bool {
	val := ctx.Value(router.ContextKeyRequestAcceptEncoding)
	enc, ok := val.(string)
	var gzipped = false
	if ok {
		encodings := strings.Split(strings.ToLower(enc), ",")
		for _, e := range encodings {
			if strings.TrimSpace(e) == "gzip" {
				gzipped = true
			}
		}
	}

	return gzipped
}

func WithAcceptedQueryFields(acceptedFields []string) RequestDecoderOption {
	return func(d *requestDecoderOption) {
		if len(acceptedFields) == 0 {
			return
		}

		if d.acceptedFields == nil {
			d.acceptedFields = make(map[string]struct{})
		}

		for _, f := range acceptedFields {
			d.acceptedFields[f] = struct{}{}
		}
	}
}

func WithURLParamsGetter(fn func(context.Context) map[string]string) RequestDecoderOption {
	return func(d *requestDecoderOption) {
		d.urlParamsGetter = fn
	}
}
