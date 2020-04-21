package main

import (
    "fmt"
    "github.com/BabyEngine/Backend/networking"
    "log"
    "net/http"
    "time"
)

type handler struct {

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}

func main()  {
    var h = &handler{}
    s, err := networking.HTTPListenAndServeWithClose("", h)
    if err != nil {
        log.Println(err)
        return
    }
    exit := make(chan int)
    time.AfterFunc(time.Second*10, func() {
        s.Close()
        exit<-1
    })
    <-exit
    fmt.Println("closed")
    select {

    }
}
