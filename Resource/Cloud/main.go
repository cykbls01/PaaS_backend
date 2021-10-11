//go:generate statik -src=./assets -dest=.
package main

import (
	"Cloud/pkg/setting"
	"Cloud/routers"
	"Cloud/routers/api"
	"fmt"
	"net/http"
	"time"
)

func main(){

	router := routers.InitROuter()

	s := &http.Server{
		Addr:			fmt.Sprintf(":%d",setting.HTTPPort),
		Handler:		router,
		ReadTimeout:	setting.ReadTimeout,
		WriteTimeout:	setting.WriteTimeout,
		MaxHeaderBytes:	1 << 20,
	}

	d := time.Duration(time.Minute * 1)

	t := time.NewTicker(d)
	defer t.Stop()
	go func(){

		for {
			<- t.C

			api.GetTotalResource()
		}
	}()

	d1 := time.Duration(time.Hour * 24)

	t2 := time.NewTicker(d1)
	defer t2.Stop()

	go func(){
		for {
			<- t2.C

			api.DealWithEndedClass()
		}
	}()

	s.ListenAndServe()
}