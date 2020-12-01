package hoard

import (
	"context"

	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
)

const DefaultChunkSize = 1 << 20 // 1MiB

// How many refs to place in a single PutSeal Grant before introducing an intermediate LINK ref
const MaxRefsBeforeLinking = 16

// Service implements the GRPC Hoard service. It should mostly be plumbing to
// a DeterministicEncryptedStore (for which hoard.hoard is the canonical example)
// and also to Grants.
type Service struct {
	streaming *StreamingService
}

func NewService(grantService GrantService, chunkSize int) *Service {
	if chunkSize == 0 {
		chunkSize = DefaultChunkSize
	}
	return &Service{
		streaming: NewStreamingService(grantService, chunkSize, func(refs []*reference.Ref) bool {
			// TODO: using LINK refs with a nonce may have a role to play in our deletion system - we could delete a
			// noncified LINK ref without deleting the underlying linked refs. That would mean we could safely delete
			// while at the same time getting the benefit of only storing the same data once (as the underlying refs)
			// in effect data would become inaccessible when the last reference was deleted. We may still want to introduce
			// some form of reference counting in a stateful 'link layer' if we want to be able to truly delete the data
			// when it falls out of reference.1
			return len(refs) > MaxRefsBeforeLinking
		}),
	}
}

// PutSeal encrypts and seals plaintext
func (service *Service) PutSeal(srv api.Grant_PutSealServer) error {
	return service.streaming.PutSeal(srv.SendAndClose, srv.Recv)
}

func (service *Service) UnsealGet(grt *grant.Grant, srv api.Grant_UnsealGetServer) error {
	return service.streaming.UnsealGet(grt, srv.Send)
}

func (service *Service) Seal(srv api.Grant_SealServer) error {
	return service.streaming.Seal(srv.SendAndClose, srv.Recv)
}

func (service *Service) Unseal(grt *grant.Grant, srv api.Grant_UnsealServer) error {
	return service.streaming.Unseal(grt, srv.Send)
}

func (service *Service) Reseal(ctx context.Context, grts *api.GrantAndGrantSpec) (*grant.Grant, error) {
	return service.streaming.Reseal(grts)
}

func (service *Service) UnsealDelete(grt *grant.Grant, srv api.Grant_UnsealDeleteServer) error {
	return service.streaming.UnsealDelete(grt, srv.Send)
}

func (service *Service) Put(srv api.Cleartext_PutServer) error {
	return service.streaming.Put(srv.Send, srv.Recv)
}

func (service *Service) Get(srv api.Cleartext_GetServer) error {
	return service.streaming.Get(srv.Send, srv.Recv)
}

func (service *Service) Encrypt(srv api.Encryption_EncryptServer) error {
	return service.streaming.Encrypt(srv.Send, srv.Recv)
}

func (service *Service) Decrypt(srv api.Encryption_DecryptServer) error {
	return service.streaming.Decrypt(srv.Send, srv.Recv)
}

func (service *Service) Push(srv api.Storage_PushServer) error {
	return service.streaming.Push(srv.Send, srv.Recv)
}

func (service *Service) Pull(srv api.Storage_PullServer) error {
	return service.streaming.Pull(srv.Send, srv.Recv)
}

// Delete removes the data located at the address
func (service *Service) Delete(ctx context.Context, address *api.Address) (*api.Address, error) {
	return address, service.streaming.Delete(address.Address)
}

// Stat checks the data stored at the given address
func (service *Service) Stat(ctx context.Context, address *api.Address) (*stores.StatInfo, error) {
	return service.streaming.Stat(address)
}
