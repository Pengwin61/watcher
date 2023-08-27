package core

type ViewSession struct {
	Username     string
	Status       string
	Hostname     string
	StartSession string
	StopSession  string
}

var ViewData = make([]ViewSession, 0)

func ShowSession(personsSession *[]PersonSession) {
	ViewData = nil

	if personsSession != nil {

		for _, v := range *personsSession {

			switch v.State {
			case "S":
				v.State += "toped"

			case "R":
				v.State += "unning"
			}

			vTmp := ViewSession{
				Username:     v.Username,
				Status:       v.State,
				Hostname:     v.Hostname,
				StartSession: viewTimeFormat(v.StartDateSession),
				StopSession:  viewTimeFormat(v.StopDateSession)}
			ViewData = append(ViewData, vTmp)
		}
	} else {
		ViewData = nil
	}

}
