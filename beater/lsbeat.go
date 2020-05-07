package beater

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/liu-xiao-guo/lsbeat/config"
)

// lsbeat configuration.
type lsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	lastIndexTime time.Time
}

// New creates an instance of lsbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &lsbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *lsbeat) Run(b *beat.Beat) error {
    logp.Info("lsbeat is running! Hit CTRL-C to stop it.")

    var err error
    bt.client, err = b.Publisher.Connect()
    if err != nil {
	return err
    }

    // fmt.Printf("%+v\n", b)

    ticker := time.NewTicker(bt.config.Period)
    for {
        now := time.Now()
        bt.listDir(bt.config.Path, b.Info.Beat) // call listDir
        bt.lastIndexTime = now                  // mark Timestamp
        logp.Info("Event sent")
        select {
        case <-bt.done:
            return nil
        case <-ticker.C:
        }
    }
    return nil
}

// Stop stops lsbeat.
func (bt *lsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *lsbeat) listDir(dirFile string, beatname string) {
    files, _ := ioutil.ReadDir(dirFile)
    for _, f := range files {
        t := f.ModTime()
        path := filepath.Join(dirFile, f.Name())
        if t.After(bt.lastIndexTime) {

	    event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr {
			"type":       beatname,
			"modtime":    common.Time(t),
			"filename":   f.Name(),
			"path":       path,
			"directory":  f.IsDir(),
			"filesize":   f.Size(),
		},
	    }

	    bt.client.Publish(event)
        }
        if f.IsDir() {
            bt.listDir(path, beatname)
        }
    }
}
