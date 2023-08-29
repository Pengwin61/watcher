package core

type ViewSession struct {
	Username     string
	Status       string
	Hostname     string
	StartSession string
	StopSession  string
	SessionID    string
	DbID         int
}

var ViewData = make([]ViewSession, 0)

func ShowSession(personsSession *[]PersonSession) {
	ViewData = nil

	if personsSession != nil {

		for _, v := range *personsSession {

			vTmp := ViewSession{
				Username:     v.Username,
				Status:       *viewStatusFormat(&v.State),
				Hostname:     viewHostname(v.Hostname),
				StartSession: viewTimeFormat(v.StartDateSession),
				StopSession:  viewTimeFormat(v.StopDateSession),
				SessionID:    v.SessionID,
				DbID:         v.DbID}

			ViewData = append(ViewData, vTmp)
		}
	} else {
		ViewData = nil
	}

}
