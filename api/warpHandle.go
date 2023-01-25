package api

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/yddeng/utils/task"
)

var taskQueue *task.TaskPool = task.NewTaskPool(1, 1024)

// ctx.request.body格式解析至参数, 并调用对应函数句柄
func WarpHandle(fn interface{}) gin.HandlerFunc {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	switch typ.NumIn() {
	case 1: // func(done *WaitConn)
		return func(ctx *gin.Context) {
			transBegin(ctx, fn)
		}
	case 2: // func(done *WaitConn, req struct)
		return func(ctx *gin.Context) {
			inValue, err := getJsonBody(ctx, typ.In(1))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Json unmarshal failed!",
					"error":   err.Error(),
				})
				return
			}

			transBegin(ctx, fn, inValue)
		}
	default:
		panic("func symbol error")
	}
}

func transBegin(ctx *gin.Context, fn interface{}, args ...reflect.Value) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != len(args)+1 {
		panic("func argument error")
	}

	route := getCurrentRoute(ctx)
	wait := NewWaitConn(ctx, route)
	if err := taskQueue.SubmitTask(webTask(func() {
		val.Call(append([]reflect.Value{reflect.ValueOf(wait)}, args...))
	})); err != nil {
		wait.SetResult("访问人数过多", nil)
		wait.Done()
	}
	wait.Wait()

	ctx.JSON(wait.code, wait.result)
}

func getCurrentRoute(ctx *gin.Context) string {
	return ctx.FullPath()
}

func getJsonBody(ctx *gin.Context, inType reflect.Type) (inValue reflect.Value, err error) {
	if inType.Kind() == reflect.Ptr {
		inValue = reflect.New(inType.Elem())
	} else {
		inValue = reflect.New(inType)
	}
	// 常用解析JSON数据的方法
	// 解析ctx.request.body到inValue
	if err = ctx.ShouldBindJSON(inValue.Interface()); err != nil {
		return
	}
	if inType.Kind() != reflect.Ptr {
		inValue = inValue.Elem()
	}
	return
}
