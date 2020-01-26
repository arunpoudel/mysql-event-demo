package event

import (
	"context"
	"github.com/siddontang/go-mysql/replication"
)

type Unmarshler interface {
	Unmarshal([]interface{}) error
}

type Listener interface {
	OnAdd(interface{})
	OnEdit(interface{}, interface{})
	OnDelete(interface{})
	OnError(error)
}

type ListenerConfig struct {
	Flavor      string
	Host        string
	Port        uint16
	Username    string
	Password    string
	Listeners   map[string]Listener
	Unmarshlers map[string]Unmarshler
}

func NewListener(config ListenerConfig) {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 1,
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "",
	}

	syncer := replication.NewBinlogSyncer(cfg)
	streamer, _ := syncer.StartSync(syncer.GetNextPosition())
	for {
		ev, _ := streamer.GetEvent(context.Background())
		switch ev.Event.(type) {
		case *replication.RowsEvent:
			event := ev.Event.(*replication.RowsEvent)
			if unmarshler, ok := config.Unmarshlers[string(event.Table.Table)]; ok {
				if listener, ok := config.Listeners[string(event.Table.Table)]; ok {
					if ev.Header.EventType == replication.WRITE_ROWS_EVENTv2 {
						for _, row := range event.Rows {
							added := unmarshler
							added.Unmarshal(row)
							listener.OnAdd(added)
						}
					} else if ev.Header.EventType == replication.UPDATE_ROWS_EVENTv2 {
						for i := 0; i < len(event.Rows); i += 2 {
							previous := unmarshler
							previous.Unmarshal(event.Rows[i])
							current := unmarshler
							current.Unmarshal(event.Rows[i+1])
							listener.OnEdit(previous, current)
						}
					} else if ev.Header.EventType == replication.DELETE_ROWS_EVENTv2 {
						for _, row := range event.Rows {
							deleted := unmarshler
							deleted.Unmarshal(row)
							listener.OnDelete(deleted)
						}
					}
				}
			}
		}
	}
}
