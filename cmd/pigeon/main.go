// Copyright (c) 2026 Jody Scott
// SPDX-License-Identifier: Apache-2.0

// pigeon is the generic pigeon host: the Job the coop creates for one run.
// It bootstraps from the identity the coop minted (KERNEL_BOOTSTRAP_TOKEN_FILE),
// learns its run from the environment the coop set, and lives on its thread.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/tacplus-dev/pigeon/internal/host"
	loggerchassis "github.com/tacplus-dev/pigeon/tacplus/chassis/logger"
	"github.com/tacplus-dev/pigeon/tacplus/chassis/servicemanager"
	"github.com/tacplus-dev/pigeon/tacplus/chassis/token"
	sm "github.com/tacplus-dev/pigeon/tacplus/s2sm"
	fe "github.com/tacplus-dev/pigeon/tacplus/s2sm/grpc/servers/services/pigeon"
	pb "github.com/tacplus-dev/pigeon/tacplus/types/services/pigeon"
	healthreport_pb "github.com/tacplus-dev/pigeon/tacplus/types/sys/healthreport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func env(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("%s is not set; a pigeon runs only as the Job the coop creates", name)
	}
	return v, nil
}

func run() error {
	logger, err := loggerchassis.InstallFromConfig(nil)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync() //nolint:errcheck

	var r host.Run
	for name, dst := range map[string]*string{
		"PIGEON_RUN_ID": &r.RunID, "PIGEON_THREAD_ID": &r.ThreadID, "PIGEON_TENANT": &r.Tenant, "PIGEON_KIND": &r.Kind, "PIGEON_CALLER": &r.Caller,
	} {
		if *dst, err = env(name); err != nil {
			return err
		}
	}

	s, err := sm.New()
	if err != nil {
		return fmt.Errorf("servicemanager: %w", err)
	}
	defer s.Close()

	raw, err := s.GetToken()
	if err != nil {
		return fmt.Errorf("own token: %w", err)
	}
	own, err := token.ParseToken(raw)
	if err != nil {
		return fmt.Errorf("own token: %w", err)
	}
	if own.GetSubject() == "" {
		return errors.New("own token has no subject")
	}
	h := &host.Host{Run: r, Self: own.GetSubject(), Publish: s.CreateCoopPigeonMessageBusStreamPublisher(), Kinds: host.BuiltinKinds}
	s.SetOnStreamCoopPigeonMessageBusStreamFunc(h.OnMessage)
	zap.L().Info("pigeon", zap.String("self", own.GetSubject()), zap.String("run", r.RunID), zap.String("kind", r.Kind), zap.String("dog", r.Caller))

	// Ready once the thread is being read, so the coop's RUNNING is not
	// said before a task could be heard; and leave when the coop ends the run.
	go func() {
		time.Sleep(2 * time.Second)
		if err := h.Ready(context.Background()); err != nil {
			zap.L().Error("saying Ready", zap.Error(err))
		}
		<-h.Done()
		zap.L().Info("run ended; leaving")
		s.Close()
		os.Exit(0)
	}()
	return s.Run(server{})
}

// server is the generated handler: a pigeon serves no RPC of its own, and its
// config state is the config store's (kernel's), not its.
type server struct{}

func (server) PigeonConfigStateConfig(*pb.Config) error { return nil }
func (server) PigeonConfigStateGrpcImpl() fe.PigeonConfigStateServiceHandler {
	return notHosted{}
}
func (server) Shutdown(context.Context) {}
func (server) CheckHealth() healthreport_pb.InstanceStatus {
	return healthreport_pb.InstanceStatus_INSTANCE_STATUS_HEALTHY
}

var errNotHosted = status.Error(codes.Unimplemented, "the pigeon's config state is served by the config store (kernel), not by the pigeon")

type notHosted struct{}

func (notHosted) GetConfig(context.Context, *servicemanager.Caller) (*pb.Config, error) {
	return nil, errNotHosted
}
func (notHosted) SetConfig(context.Context, *servicemanager.Caller, *pb.Config) (*pb.Config, error) {
	return nil, errNotHosted
}
func (notHosted) WatchConfig(*servicemanager.Caller, grpc.ServerStreamingServer[pb.Config]) error {
	return errNotHosted
}
