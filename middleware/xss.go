package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/net/html"
)

func XssHandler(whitelistURLs []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, u := range whitelistURLs {
			if strings.HasPrefix(ctx.Request.URL.String(), u) {
				ctx.Next()
				return
			}
		}
		sanitizedQuery, err := XSSFilterQuery(ctx.Request.URL.RawQuery)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Request.URL.RawQuery = sanitizedQuery

		var sanitizedBody string
		bding := binding.Default(ctx.Request.Method, ctx.ContentType())
		body, err := ctx.GetRawData()
		if err != nil {
			ctx.Error(err)
			return
		}

		// XSSFilterJSON() will return error when body is empty.
		if len(body) == 0 {
			ctx.Next()
			return
		}

		switch bding {
		case binding.JSON:
			sanitizedBody, err = XSSFilterJSON(string(body))
			if err != nil {
				ctx.Error(err)
				return
			}
		case binding.FormMultipart:
			sanitizedBody = XSSFilterPlain(string(body))
		case binding.Form:
			sanitizedBody, err = XSSFilterQuery(string(body))
			if err != nil {
				ctx.Error(err)
				return
			}
		}
		if err != nil {
			ctx.Error(err)

			return
		}
		ctx.Request.Body = ioutil.NopCloser(strings.NewReader(sanitizedBody))

		ctx.Next()
	}
}

func XSSFilterQuery(s string) (string, error) {
	values, err := url.ParseQuery(s)
	if err != nil {
		return "", err
	}

	for k, v := range values {
		values.Del(k)
		for _, vv := range v {
			values.Add(k, XSSFilterPlain(vv))
		}
	}

	return values.Encode(), nil
}

func XSSFilterPlain(s string) string {
	return html.EscapeString(s)
}

func XSSFilterJSON(s string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	err = e.Encode(xssFilterJSONData(data))
	if err != nil {
		return "", err
	}
	// use `TrimSpace` to trim newline char add by `Encode`.
	return strings.TrimSpace(b.String()), nil
}

func xssFilterJSONData(data interface{}) interface{} {
	if s, ok := data.([]interface{}); ok {
		for i, v := range s {
			s[i] = xssFilterJSONData(v)
		}
		return s
	} else if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			m[k] = xssFilterJSONData(v)
		}
		return m
	} else if str, ok := data.(string); ok {
		return XSSFilterPlain(str)
	}
	return data
}
