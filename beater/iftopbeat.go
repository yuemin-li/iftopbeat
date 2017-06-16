package beater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/yuemin-li/iftopbeat/config"
)

type Iftopbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

type IftopEvent struct {
	Interface string `json:"interface"`
	Interval  string `json:"interval"`
	Src       string `json:"source"`
	Dest      string `json:"destination"`
	Upload    string `json:"upload"`
	Download  string `json:"download"`
}

type Pair struct {
	Src  string `json:"source"`
	Dest string `json:"destination"`
}

// New beater creation function
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	localConfig := config.Config{}
	if err := cfg.Unpack(&localConfig); err != nil {
		fmt.Errorf("Error reading config file: %v, using DefaultConfig instead", err)
		localConfig = config.DefaultConfig
	}

	bt := &Iftopbeat{
		done:   make(chan struct{}),
		config: localConfig,
	}
	return bt, nil
}

func (bt *Iftopbeat) getEvents() ([]IftopEvent, error) {
	//TODO: LookPath happends at creating beater.
	iftop, err := exec.LookPath("iftop")
	if err != nil {
		log.Fatal(err)
	}
	logp.Info("iftop is available at %s\n", iftop)

	logp.Info(fmt.Sprintf(
		"iftopbeat configs: period: %v; interval: %d; listenOn: %s; numLines: %d",
		bt.config.Period,
		bt.config.Interval,
		bt.config.ListenOn,
		bt.config.NumLines,
	))
	interval := bt.config.Interval
	listenOn := bt.config.ListenOn
	numLines := bt.config.NumLines
	args := []string{"-t", "-s", strconv.Itoa(interval), "-L", strconv.Itoa(numLines), "-i", listenOn, "-n"}
	cmd := exec.Command(iftop, args...)
	log.Print(cmd.Path, cmd.Args)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}
	lines := strings.Split(out.String(), "\n")
	var ret []Pair
	for _, l := range lines {
		log.Print(l)
		if len(l) != 0 && strings.Contains(l, "=>") {
			ret = append(ret, Pair{})
			ret[len(ret)-1].Src = l
		} else if len(l) != 0 && strings.Contains(l, "<=") {
			ret[len(ret)-1].Dest = l
		}
	}

	events := []IftopEvent{}

	for _, value := range ret {
		event := IftopEvent{Interval: strconv.Itoa(interval), Interface: listenOn}
		// TODO(yuemin): use more verbose FieldsFunc
		log.Print(value)
		uploadRecord := strings.Fields(value.Src)
		event.Src = uploadRecord[1]
		event.Upload = uploadRecord[4]
		downRecord := strings.Fields(value.Dest)
		event.Dest = downRecord[0]
		event.Download = downRecord[3]

		output, _ := json.Marshal(event)
		//log.Print(string(output))
		logp.Info(string(output))
		events = append(events, event)
	}
	return events, nil
}

func (bt *Iftopbeat) Run(b *beat.Beat) error {
	logp.Info("iftopbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		events, err := bt.getEvents()
		if err != nil {
			log.Fatal(err)
		}
		for _, record := range events {
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"counter":    counter,
				"event": common.MapStr{
					"interface":   record.Interface,
					"interval":    record.Interval,
					"source":      record.Src,
					"destination": record.Dest,
					"upload":      record.Upload,
					"download":    record.Download,
				},
			}
			bt.client.PublishEvent(event)
			logp.Info("Event sent")
			counter++
		}

	}
}

func (bt *Iftopbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
