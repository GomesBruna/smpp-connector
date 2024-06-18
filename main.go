package main

import (
	//"context"
	"encoding/json"
    "log"
    "net/http"
    "io/ioutil"
    "github.com/GomesBruna/smpp-connector/connector"
	"github.com/julienschmidt/httprouter"
	"github.com/linxGnu/gosmpp"
)


type Request struct{
    SenderName     string `json:"senderName"`
    SenderAddress  string `json:"senderAddress"`
    ReceiveAddress string `json:"receiveAddress"`
    Message        string `json:"message"`
}

  
func SendSMS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    
    var req Request
    body, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(body, &req)
    connector.SendingAndReceiveSMS(transmitter, req.SenderAddress, req.ReceiveAddress, req.Message)
	retorno := connector.GetRetorno()
	if retorno == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sucesso\n"))
	}
	if retorno == 1 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Erro sem retentativa\n"))
	}
	if retorno == 2 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Erro com retentativa\n"))
	}
}

var transmitter *gosmpp.Session
func main() {
	transmitter = connector.ConnectSmpp()
	defer connector.Close(transmitter)
    router := httprouter.New()
    router.POST("/", SendSMS)
    log.Fatal(http.ListenAndServe(":8080", router))
}