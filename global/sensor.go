package global

import (
	. "github.com/pspaces/gospace/space"
	"strconv"
)

var passTime = 0

func sensor(ts *Space, batiment int, porte int) bool {
	var idPorte int
	var idBatiment int
	socket, errLaser := ts.QueryP("Laser detected", &idBatiment, &idPorte)
	socket1, errdoorClose := ts.QueryP("DoorClose", &idBatiment, &idPorte)

	if errLaser == nil {
		idPorte = (socket.GetFieldAt(2)).(int)
		idBatiment = (socket.GetFieldAt(1)).(int)
		if batiment == idBatiment && porte == idPorte {
			passTime += 1
			if passTime > 1 {
				println("Déclenche Alarme")
				SendToConn("{\"message\": \"Alarme Door\",\"idPorte\": " + strconv.Itoa(porte) + ",\"idBatiment\": " + strconv.Itoa(batiment) + "}")
			}
			ts.Get("Laser detected", idBatiment, idPorte)
		}
	}
	if errdoorClose == nil {
		idPorte = (socket1.GetFieldAt(2)).(int)
		idBatiment = (socket1.GetFieldAt(1)).(int)
		if batiment == idBatiment && porte == idPorte {
			ts.Get("DoorClose", idBatiment, idPorte)
			var res bool = passTime >= 1
			passTime = 0
			ts.Put("Human Pass", idBatiment, idPorte, res)
			return res
		}
	}

	return sensor(ts, batiment, porte)
}
