package message

import (
	"context"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
)

type Message struct {
	WaID     string
	From     string
	To       string
	MediaUrl string
	Received time.Time
	Replied  time.Time
}

func store(ctx context.Context, event *TranscribeEvent) error {
	_, err := sqldb.Exec(ctx, `
	INSERT INTO transcription_requests (
		wa_id,
		from_number,
		to_number,
		media_url,
		received_time,
	)
	VALUES ($1, $2, $3, $4, $5)`,
		event.WaID,
		event.From,
		event.To,
		event.MediaUrl,
		time.Now().UTC())
	return err
}

func get(ctx context.Context, waID string) (*Message, error) {
	query := fmt.Sprintf(`
		SELECT
		wa_id,
		from_number,
		to_number,
		media_url,
		received_time,
		replied_time
		
		FROM transcription_requests
		WHERE wa_id = '%s'
		LIMIT 1`, waID)
	message := &Message{}
	err := sqldb.QueryRow(ctx, query).Scan(
		&message.WaID,
		&message.From,
		&message.To,
		&message.MediaUrl,
		&message.Received,
		&message.Replied,
	)
	if err != nil {
		return nil, err
	}
	return message, nil
}
