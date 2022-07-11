## 任务要求

- [ ] 完成一个基于 net/http 库实现的 Web 框架
  
  - [ ] 路由支持 GET、POST、DELETE、PUT 功能 
  
  - [ ] 实现 Context 功能 
  
  - [ ] 嵌入log、cors、recovery 等 middleware

- [ ] 使用该库实现简单的 HTTP 的请求与响应

- [ ] 使用 benchmark 进行测试，至少保证 10k 的并发读写量。



## 功能实现

- [ ] 动态路由

- [ ] 上下文包装

- [ ] 前缀树

- [ ] 分组控制

- [ ] 中间件



## gtp 设计思路

主要是对于 gin 进行模仿

```puml
package gtp{
    class Context {

    }

    class Engine {
        RouterGroup
        EngineConfig...: Options
    }

    interface IRoutes {
        Use(...HandlerFunc) IRoutes

        Handle(pattern, string, ...HandlerFunc) IRoutes
        GET(pattern, ...HandlerFunc) IRoutes
        POST(pattern, ...HandlerFunc) IRoutes
        DELETE(pattern, ...HandlerFunc) IRoutes
        PUT(pattern, ...HandlerFunc) IRoutes

        StaticFile(pattern, string) IRoutes
        Static(pattern, string) IRoutes
        StaticFS(pattern, http.FileSystem) IRoutes
    }

    class RouterGroup {
        Handlers []HandlerFunc
        basePath string
        engine   *Engine
        root     bool

        Use(...HandlerFunc) IRoutes

        Handle(pattern, string, ...HandlerFunc) IRoutes
        GET(pattern, ...HandlerFunc) IRoutes
        POST(pattern, ...HandlerFunc) IRoutes
        DELETE(pattern, ...HandlerFunc) IRoutes
        PUT(pattern, ...HandlerFunc) IRoutes

        StaticFile(pattern, string) IRoutes
        Static(pattern, string) IRoutes
        StaticFS(pattern, http.FileSystem) IRoutes
    }

    RouterGroup --o Engine
    RouterGroup *-- Engine
    RouterGroup --|> IRouters
}

package gtp_test{

}
```

## 具体功能

### engine

使用 engine 作为处理请求的核心，将路由存储到一张句柄函数 map 上

### Context

添加了 Context，对运行参数进行了封装，

- 提供了访问Query和PostForm参数的方法。
- 提供了快速构造String/Data/JSON/HTML响应的方法。

### Router

将和路由相关的方法和结构提取出来，对 router 的功能进行增强，例如增加动态路由的支持

### trie tree

前缀树实现动态路由
