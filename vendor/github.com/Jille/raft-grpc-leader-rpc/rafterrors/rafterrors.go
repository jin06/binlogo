// Package rafterrors annotates Raft errors with gRPC status codes.
//
// Use MarkRetriable/MarkUnretriable to add a gRPC status code.
//
// Use MarkRetriable for atomic operations like Apply, ApplyLog, Barrier, changing configuration/voters.
//
// Use MarkUnretriable if your own application already made changes that it didn't roll back and for Restore.
// Restore does multiple operation, and errors could be from the first or second and it's unsafe to distinguish.
package rafterrors

import (
	"github.com/hashicorp/raft"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MarkRetriable annotates a Raft error with a gRPC status code, given that the entire operation is retriable.
func MarkRetriable(err error) error {
	return status.Error(RetriableCode(err), err.Error())
}

// MarkUnretriable annotates a Raft error with a gRPC status code, given that the entire operation is not retriable.
func MarkUnretriable(err error) error {
	return status.Error(UnretriableCode(err), err.Error())
}

// RetriableCode returns a gRPC status code for a given Raft error, given that the entire operation is retriable.
func RetriableCode(err error) codes.Code {
	return code(err, true)
}

// Code returns a gRPC status code for a given Raft error, given that the entire operation is not retriable.
func UnretriableCode(err error) codes.Code {
	return code(err, false)
}

func code(err error, retriable bool) codes.Code {
	switch err {
	case raft.ErrLeader, raft.ErrNotLeader, raft.ErrLeadershipLost, raft.ErrRaftShutdown, raft.ErrLeadershipTransferInProgress:
		if retriable {
			return codes.Unavailable
		}
		return codes.Unknown
	case raft.ErrAbortedByRestore:
		return codes.Aborted
	case raft.ErrEnqueueTimeout:
		if retriable {
			return codes.Unavailable
		}
		// DeadlineExceeded is generally considered not safe to be retried, because (part of) the mutation might have been applied already.
		// In hashicorp/raft, there is one place in which ErrEnqueueTimeout doesn't mean nothing happened: during a Restore does two actions (restore + noop) and i f the latter failed the restore might've gone through.
		// So sadly we can't return a more retriable error code here.
		return codes.DeadlineExceeded
	case raft.ErrNothingNewToSnapshot, raft.ErrCantBootstrap:
		return codes.FailedPrecondition
	case raft.ErrUnsupportedProtocol:
		return codes.Unimplemented
	default:
		return codes.Unknown
	}
}
