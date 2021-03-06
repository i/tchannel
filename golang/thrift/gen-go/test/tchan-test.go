// Autogenerated by thrift-gen. Do not modify.
package test

import (
	"fmt"

	athrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/uber/tchannel/golang/thrift"
)

// Interfaces for the service and client for the services defined in the IDL.

type TChanSecondService interface {
	Echo(ctx thrift.Context, arg string) (string, error)
}

type TChanSimpleService interface {
	Call(ctx thrift.Context, arg *Data) (*Data, error)
	Simple(ctx thrift.Context) error
}

// Implementation of a client and service handler.

type tchanSecondServiceClient struct {
	client thrift.TChanClient
}

func NewTChanSecondServiceClient(client thrift.TChanClient) TChanSecondService {
	return &tchanSecondServiceClient{client: client}
}

func (c *tchanSecondServiceClient) Echo(ctx thrift.Context, arg string) (string, error) {
	var resp EchoResult
	args := EchoArgs{
		Arg: arg,
	}
	success, err := c.client.Call(ctx, "SecondService", "Echo", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

type tchanSecondServiceServer struct {
	handler TChanSecondService
}

func NewTChanSecondServiceServer(handler TChanSecondService) thrift.TChanServer {
	return &tchanSecondServiceServer{handler}
}

func (s *tchanSecondServiceServer) Service() string {
	return "SecondService"
}

func (s *tchanSecondServiceServer) Methods() []string {
	return []string{
		"Echo",
	}
}

func (s *tchanSecondServiceServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "Echo":
		return s.handleEcho(ctx, protocol)
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanSecondServiceServer) handleEcho(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req EchoArgs
	var res EchoResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.Echo(ctx, req.Arg)

	if err != nil {
		return false, nil, err
	}
	res.Success = &r

	return err == nil, &res, nil
}

type tchanSimpleServiceClient struct {
	client thrift.TChanClient
}

func NewTChanSimpleServiceClient(client thrift.TChanClient) TChanSimpleService {
	return &tchanSimpleServiceClient{client: client}
}

func (c *tchanSimpleServiceClient) Call(ctx thrift.Context, arg *Data) (*Data, error) {
	var resp CallResult
	args := CallArgs{
		Arg: arg,
	}
	success, err := c.client.Call(ctx, "SimpleService", "Call", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

func (c *tchanSimpleServiceClient) Simple(ctx thrift.Context) error {
	var resp SimpleResult
	args := SimpleArgs{}
	success, err := c.client.Call(ctx, "SimpleService", "Simple", &args, &resp)
	if err == nil && !success {
		if e := resp.SimpleErr; e != nil {
			err = e
		}
	}

	return err
}

type tchanSimpleServiceServer struct {
	handler TChanSimpleService
}

func NewTChanSimpleServiceServer(handler TChanSimpleService) thrift.TChanServer {
	return &tchanSimpleServiceServer{handler}
}

func (s *tchanSimpleServiceServer) Service() string {
	return "SimpleService"
}

func (s *tchanSimpleServiceServer) Methods() []string {
	return []string{
		"Call",
		"Simple",
	}
}

func (s *tchanSimpleServiceServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "Call":
		return s.handleCall(ctx, protocol)
	case "Simple":
		return s.handleSimple(ctx, protocol)
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanSimpleServiceServer) handleCall(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req CallArgs
	var res CallResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.Call(ctx, req.Arg)

	if err != nil {
		return false, nil, err
	}
	res.Success = r

	return err == nil, &res, nil
}

func (s *tchanSimpleServiceServer) handleSimple(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req SimpleArgs
	var res SimpleResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	err :=
		s.handler.Simple(ctx)

	if err != nil {
		switch v := err.(type) {
		case *SimpleErr:
			res.SimpleErr = v
		default:
			return false, nil, err
		}
	}

	return err == nil, &res, nil
}
