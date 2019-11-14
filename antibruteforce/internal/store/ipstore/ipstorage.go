package ipstore

import (
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"context"
	"github.com/jackc/pgx/pgtype"
	"github.com/jmoiron/sqlx"
	"net"
	"time"
)

type IpTable struct {
	ID          int64
	Kind        entities.KindIp
	IP          pgtype.Inet `db:"ip"`
	DateCreated time.Time   `db:"date_created"`
}

// DbRepo - ip repository
type DbRepo struct {
	db *sqlx.DB
}

// NewDbRepo - create event repository
func NewDbRepo(db *sqlx.DB) *DbRepo {
	return &DbRepo{db: db}
}

// Add adding new ip may by black or white
func (d *DbRepo) Add(ctx context.Context, ip *entities.IPList) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO ip_list (ip, kind, date_created) VALUES ($1,$2,$3)", ip.IP.String(), ip.Kind, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// GetByIP get by ip
func (d *DbRepo) GetByIP(ctx context.Context, ip *net.IPNet) (*entities.IPList, error) {
	dest := &IpTable{}
	err := d.db.GetContext(ctx, dest, "SELECT * FROM ip_list WHERE ip=$1", ip.String())
	if err != nil {
		return nil, err
	}
	return &entities.IPList{
		ID:          dest.ID,
		Kind:        dest.Kind,
		IP:          dest.IP.IPNet,
		DateCreated: dest.DateCreated,
	}, nil
}

// Delete delete by ip
func (d *DbRepo) DeleteByIp(ctx context.Context, ip *net.IPNet) error {
	result, err := d.db.ExecContext(ctx, "DELETE FROM ip_list WHERE ip =$1", ip.IP.String())
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return exceptions.ObjectNoteFound
	}
	return nil
}
