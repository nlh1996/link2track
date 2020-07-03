package ws

import (
	"log"
)

// HandlerFunc 路由处理方法
type HandlerFunc func(*Request)

// Router 路由.
type Router struct {
	Apis map[string]HandlerFunc
}

// NewRouter .
func NewRouter() *Router {
	return &Router{
		Apis: make(map[string]HandlerFunc),
	}
}

// DoMsgHandle .
func (mh *Router) DoMsgHandle(req *Request) {
	handlerFunc, ok := mh.Apis[req.ID]
	if !ok {
		log.Println(req.ID, "api not found!")
		return
	}
	handlerFunc(req)
}

// AddRouter .
func (mh *Router) AddRouter(id string, handlerFunc HandlerFunc) {
	if _, ok := mh.Apis[id]; ok {
		log.Panicln("repeat api, msgID = ", id)
	}
	mh.Apis[id] = handlerFunc
}

// BeforeHandle .
func (mh *Router) BeforeHandle(r *Request) {
	// fmt.Println("BeforeHandle call...")
	// req.Send(gin.H{"msg": "BeforeHandle call..."})
}

// AfterHandle .
func (mh *Router) AfterHandle(r *Request) {
	// fmt.Println("AfterHandle call...")
	// req.Send(gin.H{"msg": "AfterHandle call..."})
}
