// Copyright (c) 2026 Jody Scott
// SPDX-License-Identifier: Apache-2.0

package host

import (
	"context"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	coop_pb "github.com/tacplus-dev/pigeon/tacplus/types/services/coop"
	pigeon_pb "github.com/tacplus-dev/pigeon/tacplus/types/shared/pigeon"
)

const (
	dog  = "spiffe://local.dev/tenants/acme/services/tenantecho"
	self = "spiffe://local.dev/system/services/coop/workloads/pigeon-01run"
)

func newHost() (*Host, *[]*coop_pb.PigeonMessage) {
	var sent []*coop_pb.PigeonMessage
	h := &Host{
		Run:     Run{RunID: "01run", ThreadID: "01run", Tenant: "/tenants/acme", Kind: "echo.ping", Caller: dog},
		Self:    self,
		Publish: func(_ context.Context, m *coop_pb.PigeonMessage) error { sent = append(sent, m); return nil },
		Kinds:   BuiltinKinds,
		Now:     func() time.Time { return time.Unix(1_700_000_000, 0) },
	}
	return h, &sent
}

func task(t *testing.T, sender, kind string, req []byte) *coop_pb.PigeonMessage {
	t.Helper()
	body, err := proto.Marshal(&pigeon_pb.ToPigeon{Body: &pigeon_pb.ToPigeon_Task{Task: &pigeon_pb.PigeonTask{Kind: kind, JobId: "j", Request: req}}})
	if err != nil {
		t.Fatal(err)
	}
	return &coop_pb.PigeonMessage{RunId: "01run", ThreadId: "01run", Direction: coop_pb.Direction_DIRECTION_TO_PIGEON, Kind: "echo.ping", Seq: 1, Body: body, TenantNamespace: "/tenants/acme", SpiffeId: sender}
}

func from(t *testing.T, m *coop_pb.PigeonMessage) *pigeon_pb.FromPigeon {
	t.Helper()
	var f pigeon_pb.FromPigeon
	if err := proto.Unmarshal(m.GetBody(), &f); err != nil {
		t.Fatal(err)
	}
	return &f
}

func TestEchoRoundTrip(t *testing.T) {
	ctx := context.Background()
	h, sent := newHost()
	if err := h.Ready(ctx); err != nil {
		t.Fatal(err)
	}
	if err := h.Ready(ctx); err != nil {
		t.Fatal(err)
	}
	if len(*sent) != 1 || from(t, (*sent)[0]).GetReady() == nil {
		t.Fatalf("Ready is said once: %v", *sent)
	}
	if err := h.OnMessage(ctx, task(t, dog, "echo.ping", []byte("hello")), nil); err != nil {
		t.Fatal(err)
	}
	if len(*sent) != 2 {
		t.Fatalf("want one result, got %d messages", len(*sent))
	}
	m := (*sent)[1]
	r := from(t, m).GetResult()
	if string(r.GetResult()) != "hello" || r.GetStatus() != pigeon_pb.ResultStatus_RESULT_STATUS_OK || r.GetModeUsed() != pigeon_pb.RunMode_RUN_MODE_SCRIPT {
		t.Fatalf("result %v", r)
	}
	if m.GetDirection() != coop_pb.Direction_DIRECTION_FROM_PIGEON || m.GetSpiffeId() != "" || m.GetSeq() != 2 || m.GetTenantNamespace() != "/tenants/acme" || m.GetRunId() != "01run" {
		t.Fatalf("envelope %v", m)
	}
}

func TestReadyPrecedesTheFirstResult(t *testing.T) {
	ctx := context.Background()
	h, sent := newHost()
	if err := h.OnMessage(ctx, task(t, dog, "echo.ping", []byte("x")), nil); err != nil {
		t.Fatal(err)
	}
	if len(*sent) != 2 || from(t, (*sent)[0]).GetReady() == nil || from(t, (*sent)[1]).GetResult() == nil {
		t.Fatalf("want Ready then Result: %v", *sent)
	}
}

func TestIgnoresWhatIsNotItsDogsTask(t *testing.T) {
	ctx := context.Background()
	h, sent := newHost()
	other := task(t, dog, "echo.ping", []byte("x"))
	other.RunId = "99other"
	_ = h.OnMessage(ctx, other, nil)
	_ = h.OnMessage(ctx, task(t, "spiffe://local.dev/tenants/acme/services/other", "echo.ping", []byte("x")), nil)
	mine := task(t, dog, "echo.ping", []byte("x"))
	mine.Direction = coop_pb.Direction_DIRECTION_FROM_PIGEON
	_ = h.OnMessage(ctx, mine, nil)
	if len(*sent) != 0 {
		t.Fatalf("answered something not its own: %v", *sent)
	}
	// A task for another kind fails INVALID_TASK (after Ready).
	if err := h.OnMessage(ctx, task(t, dog, "texas811.file", nil), nil); err != nil {
		t.Fatal(err)
	}
	r := from(t, (*sent)[len(*sent)-1]).GetResult()
	if r.GetStatus() != pigeon_pb.ResultStatus_RESULT_STATUS_FAILED || r.GetFailure().GetCode() != pigeon_pb.FailureCode_FAILURE_CODE_INVALID_TASK {
		t.Fatalf("wrong kind: %v", r)
	}
}

func TestEndedByTheCoopClosesDone(t *testing.T) {
	ctx := context.Background()
	h, _ := newHost()
	done := h.Done()
	select {
	case <-done:
		t.Fatal("done before the coop said anything")
	default:
	}
	_ = h.OnMessage(ctx, &coop_pb.PigeonMessage{RunId: "01run", Direction: coop_pb.Direction_DIRECTION_FROM_COOP,
		Error: &coop_pb.MessageError{Detail: "ended: done"}}, nil)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("the coop's end message did not close Done")
	}
}
