func convertTime(t string) time.Time {

	layout := "2006-01-02T15:04:05"
	timeSession, err := time.Parse(layout, t)
	if err != nil {
		log.Println(err)
	}
	return timeSession
}

func checkExpirationSession(t, state string,
	durationSession time.Duration) (bool, time.Duration) {

	var msk, _ = time.ParseDuration("3h")

	stopTimeSession := convertTime(t)
	delta := time.Since(stopTimeSession)
	delta = delta.Truncate(time.Second)

	delta += msk

	if delta >= durationSession && state != "R" {

		return true, delta
	}

	return false, delta
}