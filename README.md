## 任务要求

- [ ] 完成一个基于 net/http 库实现的 Web 框架
  
  - [ ] 路由支持 GET、POST、DELETE、PUT 功能 
  
  - [ ] 实现 Context 功能 
  
  - [ ] 嵌入log、cors、recovery 等 middleware

- [ ] 使用该库实现简单的 HTTP 的请求与响应

## 功能实现

- [x] 动态路由

- [x] 上下文包装

- [x] 前缀树

- [x] 分组控制

- [x] 中间件

主要是对于 gin 进行模仿

## 具体功能

### engine

使用 engine 作为处理请求的核心，将路由存储到一张句柄函数 map 上

```go
type Engine struct {
    *RouterGroup
    router *router
    groups []*RouterGroup // store all groups
}
```

### Context

添加了 Context，对运行参数进行了封装，

- 提供了访问Query和PostForm参数的方法。
- 提供了快速构造String/Data/JSON/HTML响应的方法。

```go
type Context struct {
    /* origin objects */
    Writer http.ResponseWriter
    Req    *http.Request

    /* request info */
    Path   string
    Method string
    Params map[string]string // store the router parameter

    /* response info */
    StatusCode int

    /* middleware */
    handlers []HandlerFunc
    index    int
}
```

Context 的实质就是对于请求的封装，让其携带上需要的信息。

### Router

将和路由相关的方法和结构提取出来，对 router 的功能进行增强，例如增加动态路由的支持

```go
type router struct {
    // Create one root for each method,
    // e.g. roots['GET'], roots['POST']
    roots map[string]*node
    // handlers key matched one method and one path
    // e.g. handlers['GET-/user/info']
    handlers map[string]HandlerFunc
}
```

路由这个结构体中储存了两个字段，一个是各方法路由与前缀树上节点的映射表，即 `roots`，一个是各路由与其句柄函数的 `map`，

#### Trie tree

前缀树实现动态路由，通过提供对树上路由节点的操作方法实现功能。

```go
type node struct {
    pattern  string  // 待匹配路由，例如 /p/:lang
    part     string  // 路由中的一部分，例如 :lang
    children []*node // 子节点，例如 [doc, tutorial, intro]
    isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}
```

#### RouterGroup

分组路由

```go
type RouterGroup struct {
        prefix      string
        middlewares []HandlerFunc // support middleware
        parent      *RouterGroup  // support nesting
        engine      *Engine       // all groups share a Engine instance
}
```

### Middleware

#### Cors

```go
func Cors() HandlerFunc {
    return func(c *Context) {
        origin := c.Req.Header["Origin"] //请求头部
        if len(origin) == 0 {
            // 当Access-Control-Allow-Credentials为true时，将*替换为指定的域名
            c.Header("Access-Control-Allow-Origin", "http://example.com")
            c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
            c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization")
            c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
            c.Header("Access-Control-Allow-Credentials", "true")
            c.Header("Access-Control-Max-Age", "86400") // 可选
        }
        c.Next() // 继续执行
    }
}
```

为 Head 添加跨域所需的信息即可。

#### Recovery

```go
func Recovery() HandlerFunc {
    return func(c *Context) {
        defer func() {
            if err := recover(); err != nil {
                message := fmt.Sprintf("%s", err)
                log.Printf("%s\n\n", trace(message))
                c.Fail(http.StatusInternalServerError, "Internal Server Error")
            }
        }()

        c.Next()
    }
}

// print stack trace for debug
func trace(message string) string {
    var pcs [32]uintptr
    n := runtime.Callers(3, pcs[:]) // skip first 3 caller

    var str strings.Builder
    str.WriteString(message + "\nTraceback:")
    for _, pc := range pcs[:n] {
        fn := runtime.FuncForPC(pc)
        file, line := fn.FileLine(pc)
        str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
    }
    return str.String()
}
```

#### logger

```go
func Logger() HandlerFunc {
    return func(c *Context) {
        // Start timer
        t := time.Now()
        // Process request
        c.Next()
        // Calculate resolution time
        log.Printf("[gtp] [%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
    }
}
```

Logger 实现简单控制台日志功能

## 测试

静态路由，静态页面返回

```shell
$ curl http://localhost:8080/        


StatusCode        : 200
StatusDescription : OK
Content           : <h1>Hello Gee</h1>
RawContent        : HTTP/1.1 200 OK
                    Content-Length: 18
                    Content-Type: text/html
                    Date: Sat, 16 Jul 2022 01:54:16 GMT

                    <h1>Hello Gee</h1>
Forms             : {}
Headers           : {[Content-Length, 18], [Content-Type, text/html], [Date, Sat, 16 Jul 2022 01:54:16 GMT]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 18
```

动态路由正常返回，不使用 Cors 中间件

```shell
$ curl http://127.0.0.1:8080/v2/Anxiu


StatusCode        : 200
StatusDescription : OK
Content           : Hello Anxiu

RawContent        : HTTP/1.1 200 OK
                    Content-Length: 12
                    Content-Type: text/plain
                    Date: Sat, 16 Jul 2022 01:30:01 GMT

                    Hello Anxiu

Forms             : {}
Headers           : {[Content-Length, 12], [Content-Type, text/plain], [Date,
                     Sat, 16 Jul 2022 01:30:01 GMT]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 12
```

动态路由使用 Cors 中间件，

```shell
$ curl http://127.0.0.1:8080/v1/Anxiu


StatusCode        : 200
Content           : Hello Anxiu

RawContent        : HTTP/1.1 200 OK                                                               Access-Control-Allow-Credentials: true                                        Access-Control-Allow-Headers: Origin, X-Requested-With, X                     -Extra-Header, Content-Type, Accept, Authorization        
Forms             : {}
Headers           : {[Access-Control-Allow-Credentials, true], [Access-Contro
                    l-Allow-Headers, Origin, X-Requested-With, X-Extra-Header
                    , Content-Type, Accept, Authorization], [Access-Control-A
                    llow-Methods, POST, GET, OPTIONS, PUT, DELETE, UPDATE], [
                    Access-Control-Allow-Origin, *]...}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 12
```

POST 方法测试

```go
$ curl -X POST -F 'username=Anxiu' -F 'password=123456' http://localhost:8080/v2/Anxiu
{"password":"123456","username":"Anxiu"}
```

```go
# logger
2022/07/16 11:01:43 [gtp] [200] /v2/Anxiu in 59.744µs
```

Panic

```shell
$ curl http://localhost:8080/panic
curl : {"message":"Internal Server Error"}
```

```shell
# logging
2022/07/16 10:01:26 runtime error: index out of range [100] with length 1
Traceback:
        D:/Go/src/runtime/panic.go:838
        D:/Go/src/runtime/panic.go:89
        D:/Desktop/7/gtp/cmd/main.go:39
        D:/Desktop/7/gtp/gtp/context.go:60
        D:/Desktop/7/gtp/gtp/recovery.go:22
        D:/Desktop/7/gtp/gtp/context.go:60
        D:/Desktop/7/gtp/gtp/logger.go:15
        D:/Desktop/7/gtp/gtp/context.go:60
        D:/Desktop/7/gtp/gtp/router.go:93
        D:/Desktop/7/gtp/gtp/gtp.go:49
        D:/Go/src/net/http/server.go:2917
        D:/Go/src/net/http/server.go:1967
        D:/Go/src/runtime/asm_amd64.s:1572

2022/07/16 10:01:26 [gtp] [500] /panic in 2.7332ms
```

na
