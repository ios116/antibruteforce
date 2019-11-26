package main

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/godog"
	"google.golang.org/grpc"
	"integration_test/config"
	"integration_test/grpcserver"
)

type rpc struct {
	login    string
	pass     string
	ip       string
	response string
	config   *config.GrpcConf
}

func (r *rpc) grpcConnection() (*grpc.ClientConn, error) {
	conn, err := grpcserver.NewGrpcConnection(r.config)
	return conn, err
}

func (r *rpc) requestSet(arg1, arg2, arg3 string) error {
	r.login = arg1
	r.pass = arg2
	r.ip = arg3
	return nil
}

func (r *rpc) checkRequestTimes(arg1 int) error {
	conn, err := r.grpcConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	req := &grpcserver.CheckRequest{
		Login:    r.login,
		Password: r.pass,
		Ip:       r.ip,
	}
	server := grpcserver.NewAntiBruteForceClient(conn)
	for i := 0; i < arg1; i++ {
		status, err := server.Check(ctx, req)
		if err != nil && status == nil {
			r.response = "false"
			return nil
		}
	}
	r.response = "true"
	return nil
}

func (r *rpc) responseShouldBeMatch(arg1 string) error {
	if r.response == arg1 {
		return nil
	}
	return fmt.Errorf("error shoud be %s but we have %s", arg1, r.response)
}

func (r *rpc) resetBucket(arg1, arg2 string) error {
	//return godog.ErrPending
	conf := config.NewGrpcConf()
	conn, err := grpcserver.NewGrpcConnection(conf)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	server := grpcserver.NewAntiBruteForceClient(conn)
	req := new(grpcserver.ResetBucketRequest)
	switch arg1 {
	case "login":
		req.Kind = grpcserver.BucketKind_LOGIN
		req.Key = arg2
	default:
		return fmt.Errorf("bucket  %s not define", arg1)
	}
	status, err := server.ResetBucket(ctx, req)
	if err != nil {
		r.response = "false"
		return err
	}
	if status.Ok == true {
		r.response = "true"
		return nil
	}
	r.response = "false"
	return fmt.Errorf("unexpected error for reset bucket with parasms %s %s ", arg1, arg2)
}

func (r *rpc) addIpToList(arg1, arg2 string) error {
	conf := config.NewGrpcConf()
	conn, err := grpcserver.NewGrpcConnection(conf)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	server := grpcserver.NewAntiBruteForceClient(conn)
	var kind grpcserver.List
	switch arg2 {

	case "whitelist":
		kind = grpcserver.List_WHITE
	case "blacklist":
		kind = grpcserver.List_BLACK
	default:
		return fmt.Errorf("unexpected list parasms %s ", arg2)
	}
	req := &grpcserver.AddIpRequest{
		Net:  arg1,
		List: kind,
	}
	if _, err := server.AddIP(ctx, req); err != nil {
		r.response = "false"
		return err
	}
	r.response = "true"
	return nil
}

func (r *rpc) removeIpFromList(arg1 string) error {
	conf := config.NewGrpcConf()
	conn, err := grpcserver.NewGrpcConnection(conf)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	server := grpcserver.NewAntiBruteForceClient(conn)
	req := &grpcserver.DeleteIpRequest{Net: arg1}
	_, err = server.DeleteIP(ctx, req)
	if err != nil {
		r.response = "false"
		return err
	}
	r.response = "true"
	return nil
}

func FeatureContext(s *godog.Suite) {
	test := new(rpc)
	test.config = config.NewGrpcConf()
	s.Step(`^request set "([^"]*)", "([^"]*)", "([^"]*)"$`, test.requestSet)
	s.Step(`^check request (\d+) times$`, test.checkRequestTimes)
	s.Step(`^response should be match "([^"]*)"$`, test.responseShouldBeMatch)
	s.Step(`^reset Bucket "([^"]*)", "([^"]*)"$`, test.resetBucket)
	s.Step(`^add ip "([^"]*)" to list "([^"]*)"$`, test.addIpToList)
	s.Step(`^remove ip "([^"]*)" from list$`, test.removeIpFromList)
}
