// scheduler
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/morpheusxaut/lootsheeter/models"
)

var (
	scheduler *Scheduler
)

type Scheduler struct {
	memberImportStop   chan struct{}
	memberImportTicker *time.Ticker
}

func NewScheduler() *Scheduler {
	scheduler := &Scheduler{}

	return scheduler
}

func InitialiseScheduler() {
	scheduler = NewScheduler()

	scheduler.StartMemberImport(4 * time.Hour)
}

func (s *Scheduler) StartMemberImport(interval time.Duration) {
	logger.Debugf("Starting member import scheduling...")

	s.memberImportStop = make(chan struct{})
	s.memberImportTicker = time.NewTicker(interval)

	err := s.ImportMembers()
	if err != nil {
		logger.Errorf("Failed to import members: [%v]", err)
	}

	go func() {
		for {
			select {
			case <-s.memberImportTicker.C:
				err = s.ImportMembers()
				if err != nil {
					logger.Errorf("Failed to import members: [%v]", err)
				}
			case <-s.memberImportStop:
				s.memberImportTicker.Stop()
				return
			}
		}
	}()

	logger.Debugf("Finished member import scheduling...")
}

func (s *Scheduler) ImportMembers() error {
	corporations, err := database.LoadAllCorporations()
	if err != nil {
		return err
	}

	for _, corporation := range corporations {
		apiURL := fmt.Sprintf("https://api.eveonline.com/corp/MemberTracking.xml.aspx?KeyID=%d&vCode=%s", corporation.APIID, corporation.APICode)

		resp, err := http.Get(apiURL)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		xmlContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var memberTracking models.MemberTracking

		err = xml.Unmarshal(xmlContent, &memberTracking)
		if err != nil {
			return err
		}

		for _, row := range memberTracking.Rows {
			_, err := database.SavePlayer(models.NewPlayer(-1, row.CharacterID, row.Name, corporation, models.AccessMaskMember))
			if !strings.Contains(err.Error(), "Duplicate entry") {
				return err
			}
		}
	}

	return nil
}
