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

var viewUsers = make([]ViewSession, 0)

func GetUsersView() []ViewSession {
	return viewUsers
}
func SetUserView(users []ViewSession) {
	viewUsers = users
}

func ShowSession(personsSession *[]PersonSession) {
	viewUsers = nil

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

			viewUsers = append(viewUsers, vTmp)
		}
	} else {
		viewUsers = nil
	}

}
