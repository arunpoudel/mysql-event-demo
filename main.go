package main

import (
	"github.com/arunpoudel/mysql-event-demo/event"
	"github.com/siddontang/go-log/log"
)

type HolaListener struct{}

func (*HolaListener) OnAdd(h interface{}) {
	log.Infof("New entry added: %+v", h)
}

func (*HolaListener) OnEdit(old interface{}, new interface{}) {
	log.Info("Entry Edited.")
	log.Infof("Old Value: %+v", old)
	log.Infof("New Value: %+v", new)
}

func (*HolaListener) OnDelete(h interface{}) {
	log.Infof("Entry Deleted: %+v", h)
}

func (*HolaListener) OnError(err error) {
	// Should not hit here, but the idea is, if there is some error, propagate it?
	log.Infof("Error: %+v", err)
}

type Hola struct {
	Column1 uint32
	Column2 string
	Column3 string
	Column4 string
	Column5 string
}

func (h *Hola) Unmarshal(b []interface{}) error {
	h.Column1 = uint32(b[0].(int32))
	if b[1] != nil {
		h.Column2 = b[1].(string)
	}
	if b[2] != nil {
		h.Column3 = b[2].(string)
	}
	if b[3] != nil {
		h.Column4 = b[3].(string)
	}
	if b[4] != nil {
		h.Column5 = b[4].(string)
	}

	return nil
}

func main() {
	// Listen for hola table changes
	// hola table structure as shown
	//CREATE TABLE `hola` (
	//	`Column1` int NOT NULL AUTO_INCREMENT,
	//	`Column2` varchar(100) DEFAULT NULL,
	//	`Column3` varchar(100) DEFAULT NULL,
	//	`Column4` varchar(100) DEFAULT NULL,
	//	`Column5` varchar(100) DEFAULT NULL,
	//	PRIMARY KEY (`Column1`)
	//) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	event.NewListener(event.ListenerConfig{
		Host:     "127.0.01",
		Port:     3306,
		Username: "root",
		Password: "",
		Listeners: map[string]event.Listener{
			"hola": &HolaListener{},
		},
		Unmarshlers: map[string]event.Unmarshler{
			"hola": &Hola{},
		},
	})
}
