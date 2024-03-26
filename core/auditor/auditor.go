package auditor

import (
	cr "github.com/0xBow-io/base-eas-asp/pkg/change_request"

	eas "github.com/0xBow-io/base-eas-asp/pkg/base_eas"
	proofOfAudit "github.com/0xBow-io/base-eas-asp/pkg/proof-of-audit"
	reportDB "github.com/0xBow-io/base-eas-asp/pkg/reportDB"
	sDB "github.com/0xBow-io/base-eas-asp/pkg/statedb"
)

type Verifier interface {
	SubmitChangeRequest(cr cr.ChangeRequest) error
}

type ReportDB interface {
	Get(publicID string) reportDB.Report
	Set(publicID string, ts int64, report reportDB.Report)
	SubscribeToNotif(notifChan chan string)
}

type StateDB interface {
	NsExists(ns string) bool
	GetMembership(ns string) sDB.MEMBERSHIP_TYPE
}

type Auditor struct {
	v        Verifier
	rDb      ReportDB
	sDB      sDB.StateDB
	rDbNotif chan string
}

func NewAuditor(v Verifier, rDb ReportDB) *Auditor {
	rDbNotif := make(chan string)
	rDb.SubscribeToNotif(rDbNotif)

	return &Auditor{v: v, rDb: rDb, rDbNotif: rDbNotif}
}

// Async handle new reports
func (a *Auditor) HandleIncomingReports() {
	for i := 0; i < 10; i++ {
		publicID := <-a.rDbNotif
		r := a.rDb.Get(publicID)
		easRes, _, err := r.Parse()
		if err != nil {
			panic(err)
		}

		for _, e := range easRes {
			// check that we have any associated events for those EAS
			// by checking the existence of a namespace (EAS account)
			if a.sDB.NsExists(e.Account.Hex()) {
				// check if their membership in stateDB is valid
				expectedMembership := eas.EasTypeToMembership(e.Type)
				if expectedMembership != a.sDB.GetMembership(e.Account.Hex()) {
					// create a change request (cr)
					// cr will provide proof based on the report that a change in membership for the namesapce is needed
					proof, err := proofOfAudit.VerifyCommitment()
					if err != nil {
						panic(err)
					}
					a.v.SubmitChangeRequest(cr.ChangeRequest{
						Ns:         e.Account.Bytes(),
						Membership: expectedMembership,
						Proof:      *proof,
					})
				}
			}

			// generate change request

		}
	}
}
