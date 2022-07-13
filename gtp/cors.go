package gtp

func Cors() HandlerFunc {
	return func(c *Context) {
		origin := c.Req.Header["Origin"] //请求头部
		if len(origin) == 0 {
			// 当Access-Control-Allow-Credentials为true时，将*替换为指定的域名
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400") // 可选
		}
		c.Next()
	}
}
