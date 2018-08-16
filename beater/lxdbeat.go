package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/neophenix/lxdbeat/config"

	"github.com/lxc/lxd/client"
)

// Lxdbeat configuration.
type Lxdbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of lxdbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Lxdbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts lxdbeat.
func (bt *Lxdbeat) Run(b *beat.Beat) error {
	logp.Info("lxdbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		for _, host := range bt.config.Hosts {
			conn, err := bt.connect(host)
			if err != nil {
				return err
			}

			containers, err := conn.GetContainers()
			if err != nil {
				return err
			}

			for _, container := range containers {
				state, _, err := conn.GetContainerState(container.Name)
				if err != nil {
					return err
				}

				network := make(map[string]map[string]int64)

				if eth, ok := state.Network["eth0"]; ok {
					network["eth0"] = make(map[string]int64)
					network["eth0"]["bytes_received"] = eth.Counters.BytesReceived
					network["eth0"]["bytes_sent"] = eth.Counters.BytesSent
					network["eth0"]["packets_received"] = eth.Counters.PacketsReceived
					network["eth0"]["packets_sent"] = eth.Counters.PacketsSent
				}

				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":      b.Info.Name,
						"lxdhost":   host,
						"container": container.Name,
						"cpu":       state.CPU.Usage,
						"memory":    state.Memory.Usage,
						"processes": state.Processes,
						"network":   network,
					},
				}

				bt.client.Publish(event)
			}
		}
	}

	return nil
}

// Stop stops lxdbeat.
func (bt *Lxdbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Lxdbeat) connect(host string) (lxd.ContainerServer, error) {
	args := &lxd.ConnectionArgs{
		TLSClientCert: bt.config.ClientCert,
		TLSClientKey:  bt.config.ClientKey,
		TLSServerCert: bt.config.ServerCert,
	}

	conn, err := lxd.ConnectLXD("https://"+host, args)
	if err != nil {
		return conn, err
	}

	return conn, nil
}
