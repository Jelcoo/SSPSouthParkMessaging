package scheduler

import (
	"log"
	"math/rand/v2"
	"sp-messaging/api/internal/core/domain"
	"sp-messaging/api/internal/core/ports"

	"github.com/robfig/cron/v3"
)

var predefinedMessages = []domain.Message{
	{Author: "Cartman", Body: "Respect my authoritah!"},
	{Author: "Cartman", Body: "Screw you guys, I'm going home."},
	{Author: "Kyle", Body: "You know, I learned something today..."},
	{Author: "Stan", Body: "Oh my God, they killed Kenny!"},
	{Author: "Kyle", Body: "You bastards!"},
	{Author: "Kenny", Body: "Mmph mmmph mmph!"},
	{Author: "Butters", Body: "Oh, hamburgers."},
	{Author: "Randy", Body: "I thought this was America!"},
	{Author: "Mr. Mackey", Body: "Drugs are bad, m'kay?"},
	{Author: "Chef", Body: "Hello there, children!"},
	{Author: "Cartman", Body: "Whatever, whatever, I do what I want!"},
	{Author: "Cartman", Body: "Seriously, guys!"},
	{Author: "Cartman", Body: "But Mooom!"},
	{Author: "Cartman", Body: "Sweet! Kewl!"},
	{Author: "Kyle", Body: "Dude, this is pretty f'd up right here."},
	{Author: "Kyle", Body: "Goddammit, Cartman!"},
	{Author: "Stan", Body: "Dude, that's messed up."},
	{Author: "Stan", Body: "I'm not your buddy, guy!"},
	{Author: "Kenny", Body: "Mmmph mmph mmmmph mph!"},
	{Author: "Butters", Body: "Oh, geez."},
	{Author: "Butters", Body: "Professor Chaos will have his revenge!"},
	{Author: "Randy", Body: "I'm not gonna lie to you, Stan."},
	{Author: "Randy", Body: "Lorde ya ya ya ya!"},
	{Author: "Mr. Garrison", Body: "There are no stupid questions, just stupid people."},
	{Author: "Mr. Mackey", Body: "Okay, m'kay?"},
	{Author: "Towelie", Body: "Don't forget to bring a towel!"},
	{Author: "Towelie", Body: "Wanna get high?"},
	{Author: "Jimmy", Body: "W-w-what a terrific audience!"},
	{Author: "Timmy", Body: "Timmy!"},
	{Author: "Officer Barbrady", Body: "Okay, people, move along. Nothing to see here."},
	{Author: "Mr. Hankey", Body: "Howdy ho!"},
	{Author: "Terrance", Body: "Hey Phillip, pull my finger!"},
	{Author: "Scott Tenorman", Body: "Hey, give me back my sixteen dollars!"},
}

type CronScheduler struct {
	cron *cron.Cron
	svc  ports.MessengerService
}

func NewCronScheduler(svc ports.MessengerService) *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
		svc:  svc,
	}
}

func (s *CronScheduler) ScheduleRandomMessages(spec string) error {
	_, err := s.cron.AddFunc(spec, s.sendRandomMessage)
	return err
}

func (s *CronScheduler) sendRandomMessage() {
	message := predefinedMessages[rand.IntN(len(predefinedMessages))]
	posted, err := s.svc.PostMessage(message)
	if err != nil {
		log.Printf("cron: failed to send message: %v", err)
		return
	}
	log.Printf("cron: sent message %s from %s", posted.ID, posted.Author)
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}

func (s *CronScheduler) Stop() {
	<-s.cron.Stop().Done()
}
