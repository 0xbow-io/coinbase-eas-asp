package reportdb

import "sync"

type ReportDB struct {
	data map[string]Report
	mut  sync.RWMutex

	// notifications on new reports for public IDs
	notifChans []chan string
}

func NewReportDB() *ReportDB {
	return &ReportDB{
		data: make(map[string]Report),
	}
}

func (r *ReportDB) SubscribeToNotif(notifChan chan string) {
	r.notifChans = append(r.notifChans, notifChan)
}

func (r *ReportDB) Get(publicID string) Report {
	r.mut.RLock()
	defer r.mut.RUnlock()
	if report, ok := r.data[publicID]; ok {
		return report
	}
	return Report{}
}

func (r *ReportDB) SendNotification(publicID string) {
	for _, notifChan := range r.notifChans {
		go func(c chan string) {
			c <- publicID
		}(notifChan)
	}
}

func (r *ReportDB) Set(publicID string, ts int64, report Report) {
	r.mut.Lock()
	defer r.mut.Unlock()

	// Get the latest report for the public ID
	// And check for timestamp
	latestReport, ok := r.data[publicID]
	if ok {
		// if the new report is older than the latest report, ignore it
		// if the new report is newer than the latest report, update it
		lastTs := latestReport.GetTimeStamp()
		if ts < lastTs {
			return
		}
		r.data[publicID] = report
	} else {
		r.data[publicID] = report
	}

	r.SendNotification(publicID)
}
