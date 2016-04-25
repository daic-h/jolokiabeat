package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/daichirata/jolokiabeat/config"
)

type Jolokiabeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration

	client *Jolokiabeat
}

// Creates beater
func New() *Jolokiabeat {
	return &Jolokiabeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Jolokiabeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	// Setting default period if not set
	if bt.beatConfig.Jolokiabeat.Period == "" {
		bt.beatConfig.Jolokiabeat.Period = "1s"
	}

	bt.period, err = time.ParseDuration(bt.beatConfig.Jolokiabeat.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Jolokiabeat) Setup(b *beat.Beat) error {
	conf := bt.beatConfig.Jolokiabeat

	var err error
	bt.client, err = NewJolokiaClient(conf.Hosts, &conf.Jolokia.Proxy)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Jolokiabeat) Run(b *beat.Beat) error {
	logp.Info("jolokiabeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		timerStart := time.Now()

		events := bt.client.GetJMXEvents()
		for _, event := range events {
			b.Events.PublishEvent(event)
		}

		timerEnd := time.Now()
		duration := timerEnd.Sub(timerStart)
		if duration.Nanoseconds() > bt.period.Nanoseconds() {
			logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
		}
	}
}

func (bt *Jolokiabeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Jolokiabeat) Stop() {
	close(bt.done)
}
