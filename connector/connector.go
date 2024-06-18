package connector

import (
    "fmt"
    "log"
    "time"
	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
)

var retorno int

func ConnectSmpp() *gosmpp.Session{
	auth := gosmpp.Auth{
		SMSC:       "127.0.0.1:2775",
		SystemID:   "smppclient1",
		Password:   "password",
		SystemType: "",
	}

	trans, err := gosmpp.NewSession(
		gosmpp.TRXConnector(gosmpp.NonTLSDialer, auth),
		//gosmpp.TRXConnector(TLSDialer, auth),
		gosmpp.Settings{
			EnquireLink: 5 * time.Second,

			ReadTimeout: 10 * time.Second,

			OnSubmitError: func(_ pdu.PDU, err error) {
				log.Fatal("SubmitPDU error:", err)
			},

			OnReceivingError: func(err error) {
				log.Fatal("Receiving PDU/Network error:", err)
			},

			OnRebindingError: func(err error) {
				log.Fatal("Rebinding but error:", err)
			},

			OnPDU: handlePDU(),

			OnClosed: func(state gosmpp.State) {
				fmt.Println(state)
			},
		}, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	return trans
}

func Close(trans *gosmpp.Session) {
	_ = trans.Close()
}

func SendingAndReceiveSMS(trans *gosmpp.Session, senderAddress, receiveAddress, message string) {
	// sending SMS(s)
	if err := trans.Transmitter().Submit(newSubmitSM(senderAddress, receiveAddress, message)); err != nil {
		fmt.Println(err)
	}
}

func handlePDU() func(pdu.PDU, bool)  {
	return func(p pdu.PDU, _ bool) {
		switch pd := p.(type) {
		case *pdu.SubmitSMResp:
			retorno = -1
			fmt.Printf("SubmitSMResp:%+v\n", pd)
			if pd.CommandLength == 16 {
				retorno = 2
			}
			if pd.IsOk(){
				retorno = 0
			}
		}
		if retorno == -1 {
			retorno = 1
		}
	}
}

func GetRetorno() int{
	return retorno
}

func newSubmitSM(senderAddress, receiveAddress, message string) *pdu.SubmitSM {
	// build up submitSM
	srcAddr := pdu.NewAddress()
	srcAddr.SetTon(5)
	srcAddr.SetNpi(0)
	_ = srcAddr.SetAddress(senderAddress)

	destAddr := pdu.NewAddress()
	destAddr.SetTon(1)
	destAddr.SetNpi(1)
	_ = destAddr.SetAddress(receiveAddress)

	submitSM := pdu.NewSubmitSM().(*pdu.SubmitSM)
	submitSM.SourceAddr = srcAddr
	submitSM.DestAddr = destAddr
	_ = submitSM.Message.SetMessageWithEncoding(message, data.UCS2)
	submitSM.ProtocolID = 0
	submitSM.RegisteredDelivery = 1
	submitSM.ReplaceIfPresentFlag = 0
	submitSM.EsmClass = 0

	return submitSM
}
