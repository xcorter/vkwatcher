package observable

import (
	"database/sql"
	"fmt"
)

type Provider struct {
	db *sql.DB
}

func (p *Provider) Save(observable Observable) {
	_, e := p.db.Exec(
		`INSERT INTO observable (owner, "value", type, last_scan, chat_id) VALUES (?, ?, ?, ?, ?)`,
		observable.Owner,
		observable.Value,
		1,
		0,
		observable.ChatId,
	)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func (p *Provider) GetCountByChatId(chatId int64) int {
	var count int
	row := p.db.QueryRow("SELECT count(*) FROM observable where chat_id = ?", chatId)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (p *Provider) UpdateLastScan(observable Observable) {
	_, e := p.db.Exec(
		`UPDATE observable SET last_scan = ? WHERE owner = ? AND "value" = ?`,
		observable.LastScan,
		observable.Owner,
		observable.Value,
	)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func (p *Provider) GetData() []*Observable {
	var result []*Observable
	rows, err := p.db.Query("SELECT * FROM observable WHERE chat_id IS NOT NULL AND is_active = 1")
	if err != nil {
		fmt.Println("Error")
		return result
	}

	for rows.Next() {
		ob := &Observable{}
		err := rows.Scan(&ob.Owner, &ob.Value, &ob.ObservableType.Value, &ob.LastScan, &ob.ChatId, &ob.IsActive)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, ob)
	}
	return result
}

func NewProvider(db *sql.DB) *Provider {
	return &Provider{
		db: db,
	}
}
