package audit

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"
)

type Event struct {
	UserId    int64
	Action    string
	Metadata  string
	CreatedAt time.Time
}

type Service struct {
	db    *sql.DB
	queue chan Event
	wg    sync.WaitGroup
}

func NewService(db *sql.DB, bufferSize int) *Service {
	s := &Service{
		db:    db,
		queue: make(chan Event, bufferSize),
	}

	s.wg.Add(1)
	go s.worker()

	return s
}

func (s *Service) Dispatch(event Event) {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	select {
	case s.queue <- event:
	default:
		slog.Warn("audit log is full, dropping event", "action", event.Action)
	}
}

func (s *Service) worker() {
	defer s.wg.Done()
	query := `INSERT INTO audit_logs (user_id, action, metadata, created_at) VALUES (?, ?, ?, ?)`

	for event := range s.queue {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := s.db.ExecContext(ctx, query, event.UserId, event.Action, event.Metadata, event.CreatedAt)
		cancel()
		if err != nil {
			slog.Error("failed to write audit log in background", "error", err, "action", event.Action)
		} else {
			slog.Info("audit log recorded in background", "action", event.Action, "user_id", event.UserId)
		}
	}
}

func (s *Service) Close() {
	close(s.queue)
	s.wg.Wait()
}
