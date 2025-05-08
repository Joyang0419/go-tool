package ginx_middleware

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	/*
		AllowAllOrigins:
		這設置允許來自任何來源(域名)的請求訪問您的API。
		當設為 true 時，伺服器回應會包含 Access-Control-Allow-Origin: * 頭部。
		這意味著任何網站都可以調用您的 API，沒有域名限制。
	*/
	AllowAllOrigins bool `mapstructure:"allowAllOrigins" yaml:"allowAllOrigins"`
	/*
		AllowHeaders:
		這設置允許請求中包含任何HTTP請求頭部。
		瀏覽器在進行跨域請求時，可以帶上任何自定義頭部，如 Authorization、Content-Type 等。
		通常，CORS 只允許一部分標準頭部，但設置為 * 後允許所有頭部。
	*/
	AllowHeaders []string `mapstructure:"allowHeaders" yaml:"allowHeaders"`
	/*
		AllowMethods:
		這設置允許使用任何HTTP方法來訪問您的API。
		包括 GET、POST、PUT、DELETE、OPTIONS 等所有HTTP方法。
		瀏覽器可以使用任何這些方法向您的API發送跨域請求。
	*/
	AllowMethods []string `mapstructure:"allowMethods" yaml:"allowMethods"`
}

type CORSMiddleware struct {
	config CORSConfig
}

func NewCORSMiddleware(config CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: config,
	}
}

func (receiver *CORSMiddleware) HandlerFunc() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	slog.Info("CORSMiddleware.HandlerFunc", slog.Any("config", receiver.config))
	corsConfig.AllowAllOrigins = receiver.config.AllowAllOrigins
	corsConfig.AllowHeaders = receiver.config.AllowHeaders
	corsConfig.AllowMethods = receiver.config.AllowMethods
	return cors.New(corsConfig)
}
