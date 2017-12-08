package gelf

import (
	"encoding/json"
	"errors"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/gliderlabs/logspout/router"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	router.AdapterFactories.Register(NewGelfAdapter, "gelf")
}

// GelfAdapter is an adapter that streams UDP JSON to Graylog
type GelfAdapter struct {
	writer *gelf.Writer
	route  *router.Route
}

// NewGelfAdapter creates a GelfAdapter with UDP as the default transport.
func NewGelfAdapter(route *router.Route) (router.LogAdapter, error) {
	_, found := router.AdapterTransports.Lookup(route.AdapterTransport("udp"))
	if !found {
		return nil, errors.New("unable to find adapter: " + route.Adapter)
	}

	gelfWriter, err := gelf.NewWriter(route.Address)
	if err != nil {
		return nil, err
	}

	return &GelfAdapter{
		route:  route,
		writer: gelfWriter,
	}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *GelfAdapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		m := &GelfMessage{message}
		level := gelf.LOG_INFO
		if m.Source == "stderr" {
			level = gelf.LOG_ERR
		}
		extra, err := m.getExtraFields()
		if err != nil {
			log.Println("Graylog:", err)
			continue
		}

		msg := gelf.Message{
			Version:  "1.1",
			Host:     m.Container.Config.Hostname,
			Short:    m.Message.Data,
			TimeUnix: float64(m.Message.Time.UnixNano()/int64(time.Millisecond)) / 1000.0,
			Level:    level,
			RawExtra: extra,
		}

		// here be message write.
		if err := a.writer.WriteMessage(&msg); err != nil {
			log.Println("Graylog:", err)
			continue
		}
	}
}

type GelfMessage struct {
	*router.Message
}

func (m GelfMessage) getExtraFields() (json.RawMessage, error) {
	nameParts := strings.Split(m.Container.Name, "_")

	extra := map[string]interface{}{
		"_kube_namespace": nameParts[3],
		"_kube_container": nameParts[1],
	}

	for _, assignment := range os.Environ() {
		if strings.HasPrefix(assignment, "KUBE_") {
			pair := strings.SplitN(assignment, "=", 2)
			key, value := pair[0], pair[1]
			extra["_" + strings.ToLower(key)] = value
		}
	}

	return json.Marshal(extra)
}
