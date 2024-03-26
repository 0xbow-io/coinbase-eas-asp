package aggregator

import (
	"fmt"

	reportDB "github.com/0xBow-io/base-eas-asp/pkg/reportDB"
	"github.com/pkg/errors"
)

var (
	quiknodeFeedIDs = map[string]string{
		"base-eas-attest": "604ab59e-f362-413e-a04e-64723176595a", // Coinbase EAS Feed
	}
)

type ReportFeed interface {
	SubscribeTo(NotificationID string, feedChan chan<- reportDB.Report) error
}

type ReportDB interface {
	Get(publicID string) reportDB.Report
	Set(publicID string, ts int64, report reportDB.Report)
}

type ReportAggregator struct {
	feed ReportFeed

	// map of latest report for a public ID
	rDb ReportDB
}

func NewReportAggregator(feed ReportFeed, rDb ReportDB) *ReportAggregator {
	a := &ReportAggregator{feed: feed, rDb: rDb}
	go a.aggregate()
	return a
}

/*
Collect all the reports from the feed and store them in the reportDB
*/
func (a *ReportAggregator) collect(notificationID string, errChan chan error) {
	feed := make(chan reportDB.Report)
	err := a.feed.SubscribeTo(notificationID, feed)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to subscribe to feed for id: "+notificationID)
		return
	}
	fmt.Println("subscribed to feed for id: ", notificationID)

	for report := range feed {
		publicIDs, ts, err := report.Parse()
		if err != nil {
			errChan <- errors.Wrap(err, "failed to parse report for id: "+notificationID)
		}
		for _, id := range publicIDs {
			a.rDb.Set(id.Account.String(), ts, report)
		}
	}

}

func (a *ReportAggregator) aggregate() {
	errChan := make(chan error)
	for _, id := range quiknodeFeedIDs {
		go a.collect(id, errChan)
	}

	// simply panic if eror occurs
	for err := range errChan {
		panic(err)
	}
}

func (a *ReportAggregator) GetLatestReportForPublicID(publicID string) reportDB.Report {
	return a.rDb.Get(publicID)
}
