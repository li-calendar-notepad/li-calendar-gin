package initialize

import "calendar-note-gin/ability"

func Start() {
	ability.EventReminder.Start()
}
