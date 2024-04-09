package global

import (
	. "github.com/pspaces/gospace/space"
	"strconv"
	"sync"
	"time"
)

func openDoor(idPorte int, idBatiment int, WgPorte *sync.WaitGroup, ts *Space) {
	SendToConn("{\"message\":\"Open door\",\"idBatiment\":" + strconv.Itoa(idBatiment) + ",\"idPorte\":" + strconv.Itoa(idPorte) + "}")
	time.Sleep(5 * time.Second)
	SendToConn("{\"message\":\"TurnOffLight\",\"idBatiment\":" + strconv.Itoa(idBatiment) + ",\"idPorte\":" + strconv.Itoa(idPorte) + "}")
	time.Sleep(25 * time.Second)
	ts.Put("DoorClose", idBatiment, idPorte)
	SendToConn("{\"message\":\"Close Door\",\"idBatiment\":" + strconv.Itoa(idBatiment) + ",\"idPorte\":" + strconv.Itoa(idPorte) + "}")
	WgPorte.Done()
}

func closeDoor(idPorte int, idBatiment int, WgPorte *sync.WaitGroup) {
	SendToConn("{\"message\":\"Block door\",\"idBatiment\":" + strconv.Itoa(idBatiment) + ",\"idPorte\":" + strconv.Itoa(idPorte) + "}")
	time.Sleep(5 * time.Second)
	SendToConn("{\"message\":\"TurnOffLight\",\"idBatiment\":" + strconv.Itoa(idBatiment) + ",\"idPorte\":" + strconv.Itoa(idPorte) + "}")
	SendToConn("{\"message\":\"Unlock door\",\"idBatiment\":" + strconv.Itoa(idBatiment) + ",\"idPorte\":" + strconv.Itoa(idPorte) + "}")
	WgPorte.Done()
}
