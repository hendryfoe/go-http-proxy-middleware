package middlewares

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpProxyConfig struct {
	Target         string
	ErrorHandler   func(*gin.Context, http.ResponseWriter, *http.Request, error)
	ModifyResponse func(*http.Response) error
	PathRewrite    map[string]string
}

func rewriteRulesRegex(rewrite map[string]string) map[*regexp.Regexp]string {
	rulesRegex := map[*regexp.Regexp]string{}
	for k, v := range rewrite {
		k = regexp.QuoteMeta(k)
		if strings.HasPrefix(k, `\^`) {
			k = strings.Replace(k, `\^`, "^", -1)
		}
		rulesRegex[regexp.MustCompile(k)] = v
	}
	return rulesRegex
}

func HttpProxy(config *HttpProxyConfig) gin.HandlerFunc {
	remote, err := url.Parse(config.Target)
	if err != nil {
		panic(err)
	}

	// Initialize Proxy
	proxy := httputil.NewSingleHostReverseProxy(remote)

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"message": "Bad Gateway"})
		}

		if config != nil {

			if config.PathRewrite != nil {
				for regexPath, destPath := range rewriteRulesRegex(config.PathRewrite) {
					if isMatchPath := regexPath.MatchString(path); isMatchPath {
						path = regexPath.ReplaceAllString(path, destPath)
					}
				}
			}

			if config.ModifyResponse != nil {
				proxy.ModifyResponse = config.ModifyResponse
			}

			if config.ErrorHandler != nil {
				proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
					config.ErrorHandler(c, w, r, e)
				}
			}

		}

		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = path
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
