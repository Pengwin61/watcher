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

func GetPersonalView(user string) []ViewSession {
	tempusers := make([]ViewSession, 0)
	for _, v := range viewUsers {
		if v.Username == user {
			temp := &ViewSession{
				Username:     v.Username,
				Status:       v.Status,
				Hostname:     v.Hostname,
				StartSession: v.StartSession,
				StopSession:  v.StopSession,
				SessionID:    v.SessionID}

			tempusers = append(tempusers, *temp)
		}
	}
	return tempusers
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
