/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package helper

import (
	"fmt"
	"time"

	"github.com/openblockchain/obc-peer/openchain/consensus"
	"github.com/openblockchain/obc-peer/openchain/chaincode"
	pb "github.com/openblockchain/obc-peer/protos"

	gp "google/protobuf"

	"github.com/op/go-logging"
	"golang.org/x/net/context"
)

// =============================================================================
// Init.
// =============================================================================

// Package-level logger.
var logger *logging.Logger

func init() {

	logger = logging.MustGetLogger("helper")
}

// =============================================================================
// Structure definitions go here.
// =============================================================================

// Helper data structure.
type Helper struct {
	consenter consensus.Consenter
}

// =============================================================================
// Constructs go here.
// =============================================================================

// New constructs the consensus helper object.
func New() consensus.CPI {

	if logger.IsEnabledFor(logging.DEBUG) {
		logger.Debug("Creating a new helper.")
	}

	return &Helper{}
}

// =============================================================================
// Interface implementations go here.
// =============================================================================

// SetConsenter is called from the system when a validating peer is booted. It
// creates the link between the consensus helper (in `helper`) and the specific
// algorithm implementation of the plugin package.
func (h *Helper) SetConsenter(c consensus.Consenter) {

	if logger.IsEnabledFor(logging.DEBUG) {
		logger.Debug("Setting the helper's consenter.")
	}

	h.consenter = c
}

// HandleMsg is called by the validating peer when a REQUEST or CONSENSUS
// message is received.
func (h *Helper) HandleMsg(msg *pb.OpenchainMessage) error {

	if logger.IsEnabledFor(logging.DEBUG) {
		logger.Debug("Handling an incoming message.")
	}

	switch msg.Type {
	case pb.OpenchainMessage_REQUEST:
		return h.consenter.RecvMsg(msg)
	case pb.OpenchainMessage_CONSENSUS:
		return h.consenter.RecvMsg(msg)
	default:
		return fmt.Errorf("Cannot process message type: %s", msg.Type)
	}
}

// Broadcast sends a message to all validating peers. during a consensus round.
// The argument is the serialized message that is specific to the consensus
// algorithm implementation. It is wrapped into a CONSENSUS message.
func (h *Helper) Broadcast(msgPayload []byte) error {

	if logger.IsEnabledFor(logging.DEBUG) {
		logger.Debug("Broadcasting a message.")
	}

	// Wrap as message of type OpenchainMessage_CONSENSUS.
	msgTime := &gp.Timestamp{Seconds: time.Now().Unix(), Nanos: 0}
	msg := &pb.OpenchainMessage{
		Type:      pb.OpenchainMessage_CONSENSUS,
		Timestamp: msgTime,
		Payload:   msgPayload,
	}

	// TODO: Call a function in the comms layer.
	// Waiting for Jeff's implementation.
	var _ = msg // Just to silence the compiler error.

	return nil
}

// ExecTXs will execute all the transactions listed in the `txs` array
// one-by-one. If all the executions are successful, it returns the candidate
// global state hash, and nil error array.
func (h *Helper) ExecTXs(ctxt context.Context, txs []*pb.Transaction) ([]byte, []error) {

	if logger.IsEnabledFor(logging.DEBUG) {
		logger.Debug("Executing the transactions.")
	}
	return chaincode.ExecuteTransactions(ctxt, chaincode.DefaultChain, txs)	

}

// Unicast is called by the validating peer to send a CONSENSUS message to a
// specified receiver. As is the case with `Broadcast()`, the argument is the
// serialized payload of an implementation-specific message.
func (h *Helper) Unicast(msgPayload []byte, receiver string) error {

	// Wrap as message of type OpenchainMessage_CONSENSUS.
	msgTime := &gp.Timestamp{Seconds: time.Now().Unix(), Nanos: 0}
	msg := &pb.OpenchainMessage{
		Type:      pb.OpenchainMessage_CONSENSUS,
		Timestamp: msgTime,
		Payload:   msgPayload,
	}

	// TODO: Call a function in the comms layer.
	// Waiting for Jeff's implementation.
	var _ = msg // Just to silence the compiler error.

	return nil
}
