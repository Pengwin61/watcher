package core

type ViewSession struct {
	Username     string
	Status       string
	Hostname     string
	StartSession string
	StopSession  string
}

var Tmp = make([]ViewSession, 0)

func ShowSession(personsSession *[]PersonSession) {
	Tmp = nil

	if personsSession != nil {

		for _, v := range *personsSession {

			switch v.SessionState {
			case "S":
				v.SessionState += "toped"

			case "R":
				v.SessionState += "unning"
			}

			vTmp := ViewSession{
				Username:     v.UserSession,
				Status:       v.SessionState,
				Hostname:     v.Hostname,
				StartSession: viewTimeFormat(v.StartDateSession),
				StopSession:  viewTimeFormat(v.StopDateSession)}
			Tmp = append(Tmp, vTmp)
		}
	} else {
		Tmp = nil
	}

}
