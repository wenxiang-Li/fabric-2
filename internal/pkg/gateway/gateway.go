/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
)

var logger = flogging.MustGetLogger("gateway")

// Server represents the GRPC server for the Gateway.
type Server struct {
	registry *registry
	options  Options
}

type EndorserServerAdapter struct {
	Server peer.EndorserServer
}

func (e *EndorserServerAdapter) ProcessProposal(ctx context.Context, req *peer.SignedProposal, _ ...grpc.CallOption) (*peer.ProposalResponse, error) {
	return e.Server.ProcessProposal(ctx, req)
}

// CreateServer creates an embedded instance of the Gateway.
func CreateServer(localEndorser peer.EndorserClient, discovery Discovery, selfEndpoint string, options Options) *Server {
	gwServer := &Server{
		registry: &registry{
			localEndorser:       &endorser{client: localEndorser, endpointConfig: &endpointConfig{address: selfEndpoint}},
			discovery:           discovery,
			logger:              logger,
			endpointFactory:     &endpointFactory{timeout: options.EndorsementTimeout},
			remoteEndorsers:     map[string]*endorser{},
			broadcastClients:    map[string]*orderer{},
			tlsRootCerts:        map[string][][]byte{},
			channelsInitialized: map[string]bool{},
		},
		options: options,
	}

	return gwServer
}
