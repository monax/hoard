package stores

import (
	"encoding/base64"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v4/logging"
)

type loggingStore struct {
	store  NamedStore
	logger log.Logger
}

// Decorates a Store with some simple logging of method/address pairs
func NewLoggingStore(store NamedStore, logger log.Logger) *loggingStore {
	inv := &loggingStore{
		store:  store,
		logger: logging.TraceLogger(log.With(logger, "module", "storage")),
	}
	inv.logger = log.With(inv.logger, "store", inv.Name())
	return inv
}

var _ NamedStore = (*loggingStore)(nil)

func (inv *loggingStore) Put(address []byte, data []byte) ([]byte, error) {
	address, err := inv.store.Put(address, data)
	return address, logErrorOrSuccess(log.With(inv.logger, "method", "Put", "address",
		formatAddress(address)), err)
}

func (inv *loggingStore) Get(address []byte) ([]byte, error) {
	data, err := inv.store.Get(address)
	return data, logErrorOrSuccess(log.With(inv.logger, "method", "Get", "address",
		formatAddress(address)), err)
}

func (inv *loggingStore) Stat(address []byte) (*StatInfo, error) {
	statInfo, err := inv.store.Stat(address)
	return statInfo, logErrorOrSuccess(log.With(inv.logger, "method", "Stat", "address",
		formatAddress(address)), err)
}

func (inv *loggingStore) Location(address []byte) string {
	inv.logger.Log("method", "Location", "address", formatAddress(address))
	return inv.store.Location(address)
}

func (inv *loggingStore) Name() string {
	return fmt.Sprintf("loggingStore(%s)", inv.store.Name())
}

func logErrorOrSuccess(logger log.Logger, err error) error {
	err = logging.Err(logger, err)
	if err == nil {
		logging.Msg(logger, "Success")
	}
	return err
}

func formatAddress(address []byte) string {
	return base64.StdEncoding.EncodeToString(address)
}
