package usecase

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/email"
)

type EventReminderScheduler struct {
	userRepo   repository.UserRepository
	ticketRepo repository.TicketRepository
	mailSender email.Sender
	location   *time.Location
	nowFunc    func() time.Time
}

func NewEventReminderScheduler(
	userRepo repository.UserRepository,
	ticketRepo repository.TicketRepository,
	mailSender email.Sender,
	timezone string,
) *EventReminderScheduler {
	location, err := time.LoadLocation(strings.TrimSpace(timezone))
	if err != nil {
		location = time.Local
	}

	return &EventReminderScheduler{
		userRepo:   userRepo,
		ticketRepo: ticketRepo,
		mailSender: mailSender,
		location:   location,
		nowFunc:    time.Now,
	}
}

func (s *EventReminderScheduler) Start(ctx context.Context) {
	if s.mailSender == nil {
		log.Println("event reminder scheduler skipped: mail sender is not configured")
		return
	}

	for {
		now := s.nowFunc().In(s.location)
		nextRun := nextDailyRunAtNine(now)
		waitDuration := time.Until(nextRun)
		if waitDuration < 0 {
			waitDuration = 0
		}

		log.Printf("event reminder scheduler next run at %s", nextRun.Format(time.RFC3339))

		timer := time.NewTimer(waitDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("event reminder scheduler stopped")
			return
		case <-timer.C:
		}

		s.runOnce(ctx)
	}
}

func nextDailyRunAtNine(now time.Time) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func (s *EventReminderScheduler) runOnce(ctx context.Context) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		log.Printf("event reminder scheduler failed to fetch users: %v", err)
		return
	}

	now := s.nowFunc().In(s.location)
	tomorrowStart := startOfDay(now.AddDate(0, 0, 1))
	tomorrowEnd := tomorrowStart.Add(24 * time.Hour)

	sentCount := 0
	for _, user := range users {
		if ctx.Err() != nil {
			return
		}

		events, err := s.ticketRepo.FindDistinctEventsByUserIDAndDateRange(user.ID.String(), tomorrowStart, tomorrowEnd)
		if err != nil {
			log.Printf("event reminder scheduler failed to fetch events for user %s: %v", user.ID.String(), err)
			continue
		}

		if len(events) == 0 {
			continue
		}

		subject, body := s.buildReminderEmail(user.Name, events)
		if err := s.mailSender.Send(user.Email, subject, body); err != nil {
			log.Printf("event reminder scheduler failed to send email to %s: %v", user.Email, err)
			continue
		}

		sentCount++
	}

	log.Printf("event reminder scheduler finished: %d reminder emails sent", sentCount)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (s *EventReminderScheduler) buildReminderEmail(userName string, events []model.Event) (string, string) {
	subject := "Reminder Event Besok"
	var body strings.Builder

	body.WriteString(fmt.Sprintf("<p>Halo %s,</p>", html.EscapeString(userName)))
	body.WriteString("<p>Ini pengingat bahwa kamu memiliki event yang akan berlangsung besok.</p>")
	body.WriteString("<table style=\"border-collapse:collapse; width:100%; max-width:680px;\">")
	body.WriteString("<tr><th style=\"text-align:left; padding:8px 0;\">Event</th><th style=\"text-align:left; padding:8px 0;\">Tanggal</th><th style=\"text-align:left; padding:8px 0;\">Waktu</th><th style=\"text-align:left; padding:8px 0;\">Lokasi</th></tr>")

	for _, event := range events {
		eventDate := "-"
		if event.StartDate != nil {
			eventDate = event.StartDate.In(s.location).Format("02 Jan 2006")
		}

		eventTime := "-"
		if event.Time != nil && strings.TrimSpace(*event.Time) != "" {
			eventTime = strings.TrimSpace(*event.Time)
		}

		venue := "-"
		if event.Venue != nil && strings.TrimSpace(*event.Venue) != "" {
			venue = strings.TrimSpace(*event.Venue)
		}

		body.WriteString("<tr>")
		body.WriteString(fmt.Sprintf("<td style=\"padding:6px 0;\">%s</td>", html.EscapeString(event.Title)))
		body.WriteString(fmt.Sprintf("<td style=\"padding:6px 0;\">%s</td>", html.EscapeString(eventDate)))
		body.WriteString(fmt.Sprintf("<td style=\"padding:6px 0;\">%s</td>", html.EscapeString(eventTime)))
		body.WriteString(fmt.Sprintf("<td style=\"padding:6px 0;\">%s</td>", html.EscapeString(venue)))
		body.WriteString("</tr>")
	}

	body.WriteString("</table>")
	body.WriteString("<p style=\"margin-top:16px;\">Pastikan datang tepat waktu ya. Sampai jumpa di acara!</p>")

	return subject, body.String()
}
