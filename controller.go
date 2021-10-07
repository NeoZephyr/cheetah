package main

import (
	"cheetah/framework"
	"context"
	"log"
	"math/rand"
	"time"
)

func FooController(ctx *framework.Context) error {
	finishChan := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	timeoutCtx, cancelFunc := context.WithTimeout(ctx.BaseContext(), 1 * time.Second)
	defer cancelFunc()

	go func() {
		// 每个 Goroutine 创建的时候，使用 defer 和 recover 为当前 Goroutine 捕获 panic 异常，并进行处理
		// 否则，任意一处 panic 就会导致整个进程崩溃
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		secondCount := rand.Intn(5)
		log.Printf("sleep %d seconds\n", secondCount)
		time.Sleep(time.Duration(secondCount) * time.Second)
		ctx.Json(200, "ok")
		finishChan <- struct{}{}
	}()

	select {
	case p := <- panicChan:
		log.Println("panic", p)
		ctx.WriteMux().Lock()
		defer ctx.WriteMux().Unlock()
		log.Println(p)
		ctx.Json(500, "panic")
	case <- finishChan:
		log.Println("finish")
	case <- timeoutCtx.Done():
		log.Println("timeout")
		ctx.WriteMux().Lock()
		defer ctx.WriteMux().Unlock()
		ctx.Json(500, "timeout")
		ctx.SetTimeout(true)
	}
	return nil
}
