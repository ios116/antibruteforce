package ipstore

import (
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/domain/exceptions"
	"context"
	"net"
	"time"

	"github.com/jackc/pgx/pgtype"
	"github.com/jmoiron/sqlx"
)

type IPTable struct {
	ID          int64
	Kind        entities.IPKind
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

// CreateBucket adding new ip may by black or white
func (d *DbRepo) Add(ctx context.Context, ip *entities.IPItem) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO ip_list (ip, kind, date_created) VALUES ($1,$2,$3)", ip.IP.String(), ip.Kind, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// GetByIPWithMask get by ip
func (d *DbRepo) GetByIPWithMask(ctx context.Context, ip *net.IPNet) (*entities.IPItem, error) {
	dest := &IPTable{}
	err := d.db.GetContext(ctx, dest, "SELECT * FROM ip_list WHERE ip=$1", ip.String())
	if err != nil {
		return nil, err
	}
	return &entities.IPItem{
		ID:          dest.ID,
		Kind:        dest.Kind,
		IP:          dest.IP.IPNet,
		DateCreated: dest.DateCreated,
	}, nil
}

// GetByIPWithMask get by ip
func (d *DbRepo) GetByIP(ctx context.Context, ip net.IP) (*entities.IPItem, error) {
	dest := &IPTable{}
	err := d.db.GetContext(ctx, dest, "SELECT * FROM ip_list WHERE ip >> $1", ip.String())
	if err != nil {
		return nil, err
	}
	return &entities.IPItem{
		ID:          dest.ID,
		Kind:        dest.Kind,
		IP:          dest.IP.IPNet,
		DateCreated: dest.DateCreated,
	}, nil
}



// Delete delete by ip
func (d *DbRepo) DeleteByIP(ctx context.Context, ip *net.IPNet) error {
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
