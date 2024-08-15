package leader_elector

import (
	"context"
	"database/sql"
)

func Lock(ctx context.Context, db *sql.DB, lockName, id string) error {
	query := `
	update locks (id)
		set id = $1
		where name == lockName and timestamp < $1 returning id
	on conflict
		insert into locks(id, name) values($1, $2);`

	row := db.QueryRowContext(ctx, query, id)

	result := ""
	err := row.Scan(&result)
	if err != nil {
		return err
	}

	return nil
}
