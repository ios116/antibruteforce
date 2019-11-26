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

func TestDbRepo(t *testing.T) {
	ipv4Addr, ipv4Net, err := net.ParseCIDR("192.168.0.1/24")
	// 192.168.0.254
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

	ip := &entities.IPListRow{
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
	t.Run("Get by ip with mask 32", func(t *testing.T) {
		ipv4Net := &net.IPNet{
			IP:   ipv4Addr,
			Mask: net.CIDRMask(32, 32),
		}
		result, err := repo.GetSubnetBySubnet(ctx, ipv4Net)
		if err != nil {
			t.Fatal(err)
		}
		if len(result) == 0 {
			t.Fatal("ip address not found")
		}
	})

	t.Run("Delete by ip", func(t *testing.T) {
		t.Log("delete", ipv4Net.String())
		err = repo.DeleteByIP(ctx, ipv4Net)
		if err != nil {
			t.Fatal(err)
		}
	})
	// ip:=net.ParseIP("192.168.0.254")

}
