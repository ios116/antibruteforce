package ipstore

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"context"
	"log"
	"net"
	"testing"
	"time"
)

func TestDbRepo_Add(t *testing.T) {
	ipv4Addr, ipv4Net, err := net.ParseCIDR("192.4.4.88/32")
	if err != nil {
		log.Fatal(err)
	}
	t.Log(ipv4Addr, ipv4Net)

	dbConf := config.NewDateBaseConf()
	db, err := config.DBConnection(dbConf)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewDbRepo(db)
	ctx := context.Background()

	ip := &entities.IPItem{
		Kind:        entities.Black,
		IP:          ipv4Net,
		DateCreated: time.Time{},
	}

	t.Run("add IP", func(t *testing.T) {
		err = repo.Add(ctx, ip)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetBucketByHash by ip", func(t *testing.T) {
		result, err := repo.GetByIP(ctx, ipv4Net)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
		if result.IP.String() != ipv4Net.String() {
			t.Fatal("ip is not equal", result.IP, ipv4Net)
		}
	})

	t.Run("Delete by ip", func(t *testing.T) {
		err = repo.DeleteByIP(ctx, ipv4Net)
		if err != nil {
			t.Fatal(err)
		}
	})

}
