package main

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/awoodbeck/gnp/ch12-2026/housework/v1"
)

type Rosie struct {
	housework.UnimplementedRobotMaidServer
	mu     sync.Mutex
	chores []*housework.Chore
}

func (r *Rosie) Add(_ context.Context, chores *housework.Chores) (
	*housework.Response, error) {
	r.mu.Lock()
	r.chores = append(r.chores, chores.Chores...)
	r.mu.Unlock()

	return &housework.Response{Message: "ok"}, nil
}

func (r *Rosie) Complete(_ context.Context,
	req *housework.CompleteRequest) (*housework.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil || req.ChoreNumber < 1 ||
		int(req.ChoreNumber) > len(r.chores) {
		return nil, status.Errorf(codes.OutOfRange, "chore %d not found",
			req.ChoreNumber)
	}

	r.chores[req.ChoreNumber-1].Complete = true

	return &housework.Response{Message: "ok"}, nil
}

func (r *Rosie) List(_ context.Context, _ *housework.Empty) (
	*housework.Chores, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil {
		r.chores = make([]*housework.Chore, 0)
	}

	return &housework.Chores{Chores: r.chores}, nil
}
