package consts

// TraceIDKey traceID 的 key; 為何用Var? 因為這樣可以開放給使用者去變動
var TraceIDKey = "traceID"

// fx tags

const (
	RouterTags     = `group:"routers"`
	MiddlewareTags = `group:"middlewares"`
)
