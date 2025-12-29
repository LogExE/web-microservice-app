package boxes

import "fmt"

type Outbox struct {
	ID        int64
	Processed bool
	Type      string
	Payload   []byte
}

// Query all unprocessed outbox entries
func (b *BoxRepo) OutboxList() ([]Outbox, error) {
	fail := func(err error) ([]Outbox, error) {
		return nil, fmt.Errorf("OutboxList: %v", err)
	}

	var outboxItems []Outbox
	rows, err := b.DB.Query("SELECT * FROM outbox WHERE processed = FALSE")
	if err != nil {
		return fail(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ob Outbox

		if err := rows.Scan(
			&ob.ID, &ob.Processed, &ob.Type, &ob.Payload,
		); err != nil {
			return fail(err)
		}
		outboxItems = append(outboxItems, ob)
	}
	return outboxItems, nil
}

// Update outbox flag
func (b *BoxRepo) OutboxClear(id int64) error {
	fail := func(err error) error {
		return fmt.Errorf("OutboxClear: %v", err)
	}

	_, err := b.DB.Exec("UPDATE outbox SET processed = true WHERE id = ?", id)
	if err != nil {
		return fail(err)
	}
	return nil
}
