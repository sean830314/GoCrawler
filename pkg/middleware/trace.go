package middleware

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

const JaegerOpen = 1

func OpenTracing() gin.HandlerFunc {
	jaegerOpen := viper.GetBool("jaeger.open")
	jaegerHostPort := fmt.Sprintf("%s:%d", viper.GetString("jaeger.host"), viper.GetInt("jaeger.port"))
	return func(c *gin.Context) {
		if jaegerOpen {

			var parentSpan opentracing.Span
			tracer, closer := NewJaegerTracer(jaegerHostPort)
			defer closer.Close()

			spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err != nil {
				logrus.Warning(err)
				logrus.Info("opentracing:start")
				parentSpan = tracer.StartSpan(c.Request.URL.Path)
				defer parentSpan.Finish()
			} else {
				logrus.Info("opentracing:extract")
				parentSpan = opentracing.StartSpan(
					c.Request.URL.Path,
					opentracing.ChildOf(spCtx),
					opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
					ext.SpanKindRPCServer,
				)
				defer parentSpan.Finish()
			}
			// subspan used
			c.Set("Tracer", tracer)
			c.Set("ParentSpanContext", parentSpan.Context())
			// 記錄請求 Url
			ext.HTTPUrl.Set(parentSpan, c.Request.URL.Path)
			// Http Method
			ext.HTTPMethod.Set(parentSpan, c.Request.Method)
			// 記錄元件名稱
			ext.Component.Set(parentSpan, "Gin-Http")
			parentSpan.LogFields(
				log.String("Path", c.Request.URL.Path),
				log.String("Method", c.Request.Method))
			// 在 header 中加上當前程序的上下文資訊
			c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), parentSpan))
			ext.HTTPStatusCode.Set(parentSpan, uint16(c.Writer.Status()))
			c.Next()
		} else {
			c.Next()
		}
	}
}

func NewJaegerTracer(jaegerHostPort string) (opentracing.Tracer, io.Closer) {

	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: "GoGrawler",
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
