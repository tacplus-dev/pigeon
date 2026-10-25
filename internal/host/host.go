// Copyright (c) 2026 Jody Scott
// SPDX-License-Identifier: Apache-2.0

// Package host is the generic pigeon: it lives on one thread of the coop's
// topic, says Ready, does the tasks its dog puts there with the kinds it was
// built with, and answers with results. Nothing calls it; it is spoken to only
// on the thread (A30).
package host

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	eb_pb "github.com/tacplus-dev/pigeon/tacplus/types/generated/eventbus_support"
	coop_pb "github.com/tacplus-dev/pigeon/tacplus/types/services/coop"
	pigeon_pb "github.com/tacplus-dev/pigeon/tacplus/types/shared/pigeon"
)

// Kind does one task and returns its result. The host's built-in kinds are
// the platform's; a product's kind is that product's host.
type Kind func(ctx context.Context, task *pigeon_pb.PigeonTask) (*pigeon_pb.PigeonResult, error)

// Echo is the kind the fleet proves the round trip with: the result is the
// request, byte for byte.
func Echo(_ context.Context, task *pigeon_pb.PigeonTask) (*pigeon_pb.PigeonResult, error) {
	return &pigeon_pb.PigeonResult{Kind: task.GetKind(), Status: pigeon_pb.ResultStatus_RESULT_STATUS_OK, Result: task.GetRequest()}, nil
}

// BuiltinKinds is what this host can do.
var BuiltinKinds = map[string]Kind{"echo.ping": Echo}

// Run is the pigeon's run, from the Job's environment the coop set.
type Run struct {
	RunID, ThreadID, Tenant, Kind, Caller string
}

// Host is one pigeon.
type Host struct {
	Run Run
	// Self is this pigeon's SPIFFE id, from its own token: the sender of its
	// messages.
	Self    string
	Publish func(context.Context, *coop_pb.PigeonMessage) error
	Kinds   map[string]Kind
	Now     func() time.Time

	mu    sync.Mutex
	seq   int64
	ready bool
	done  chan struct{}
	once  sync.Once
}

// Done is closed when the coop says the run has ended: the process should
// exit, its Job is being deleted.
func (h *Host) Done() <-chan struct{} { return h.doneCh() }

func (h *Host) doneCh() chan struct{} {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.done == nil {
		h.done = make(chan struct{})
	}
	return h.done
}

func (h *Host) now() time.Time {
	if h.Now == nil {
		return time.Now()
	}
	return h.Now()
}

// Ready says Ready on the thread, once: the coop marks the run RUNNING on it.
func (h *Host) Ready(ctx context.Context) error {
	h.mu.Lock()
	if h.ready {
		h.mu.Unlock()
		return nil
	}
	h.ready = true
	h.mu.Unlock()
	return h.say(ctx, &pigeon_pb.FromPigeon{Body: &pigeon_pb.FromPigeon_Ready{Ready: &pigeon_pb.Ready{}}})
}

// OnMessage is the subscription to the coop's topic. Only this run's thread
// matters, and on it only the dog's TO_PIGEON messages and the coop's word
// that the run ended.
func (h *Host) OnMessage(ctx context.Context, m *coop_pb.PigeonMessage, _ *eb_pb.RecordMeta) error {
	if m.GetRunId() != h.Run.RunID {
		return nil
	}
	switch m.GetDirection() {
	case coop_pb.Direction_DIRECTION_FROM_COOP:
		if e := m.GetError(); e != nil && (e.GetCode() == coop_pb.MessageErrorCode_MESSAGE_ERROR_CODE_PIGEON_GONE || strings.HasPrefix(e.GetDetail(), "ended:")) {
			zap.L().Info("run ended by the coop", zap.String("detail", e.GetDetail()))
			h.once.Do(func() { close(h.doneCh()) })
		}
		return nil
	case coop_pb.Direction_DIRECTION_TO_PIGEON:
	default:
		return nil
	}
	if m.GetSender() != h.Run.Caller {
		zap.L().Warn("dropping a task from a sender that is not this run's dog", zap.String("sender", m.GetSender()))
		return nil
	}
	var to pigeon_pb.ToPigeon
	if err := proto.Unmarshal(m.GetBody(), &to); err != nil {
		return h.fail(ctx, pigeon_pb.FailureCode_FAILURE_CODE_INVALID_TASK, "message body is not a ToPigeon: "+err.Error())
	}
	task := to.GetTask()
	if task == nil {
		return nil // an Answer or a Cancel: this host asks no questions and a cancel is the coop's to end
	}
	if err := h.Ready(ctx); err != nil {
		return err
	}
	kind, ok := h.Kinds[task.GetKind()]
	if !ok || task.GetKind() != h.Run.Kind {
		return h.fail(ctx, pigeon_pb.FailureCode_FAILURE_CODE_INVALID_TASK,
			fmt.Sprintf("task names kind %q; this pigeon is a %s", task.GetKind(), h.Run.Kind))
	}
	result, err := kind(ctx, task)
	if err != nil {
		return h.fail(ctx, pigeon_pb.FailureCode_FAILURE_CODE_SCRIPT_FAILED, err.Error())
	}
	if result.GetModeUsed() == pigeon_pb.RunMode_RUN_MODE_UNSPECIFIED {
		result.ModeUsed = pigeon_pb.RunMode_RUN_MODE_SCRIPT
	}
	return h.say(ctx, &pigeon_pb.FromPigeon{Body: &pigeon_pb.FromPigeon_Result{Result: result}})
}

func (h *Host) fail(ctx context.Context, code pigeon_pb.FailureCode, detail string) error {
	return h.say(ctx, &pigeon_pb.FromPigeon{Body: &pigeon_pb.FromPigeon_Result{Result: &pigeon_pb.PigeonResult{
		Kind: h.Run.Kind, ModeUsed: pigeon_pb.RunMode_RUN_MODE_SCRIPT, Status: pigeon_pb.ResultStatus_RESULT_STATUS_FAILED,
		Failure: &pigeon_pb.Failure{Code: code, Detail: detail},
	}}})
}

func (h *Host) say(ctx context.Context, body *pigeon_pb.FromPigeon) error {
	raw, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	h.mu.Lock()
	h.seq++
	seq := h.seq
	h.mu.Unlock()
	return h.Publish(ctx, &coop_pb.PigeonMessage{
		RunId: h.Run.RunID, ThreadId: h.Run.ThreadID, Direction: coop_pb.Direction_DIRECTION_FROM_PIGEON, Kind: h.Run.Kind,
		Seq: seq, Body: raw, TenantNamespace: h.Run.Tenant, SentMs: h.now().UnixMilli(), Sender: h.Self,
	})
}
