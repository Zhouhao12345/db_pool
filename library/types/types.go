package types

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc
type Context = gin.Context
type Connect interface {
	New() (err error)
	Close() (err error)
	GetUsed() (used bool)
	SetUsed(used bool) (err error)
	HandlerRequest(f HandlerFunc) HandlerFunc
}
type ConfigMeta map[string]interface{}
