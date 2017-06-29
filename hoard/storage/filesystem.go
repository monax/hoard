package storage

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

type fileSystemStore struct {
	rootDirectory    string
	addressSegmenter AddressSegmenter
}

// Address splitter is used to create a directory hierarchy by splitting
// addresses into segments, the mapping should be bijective and the segments will
// be joined to form the path to file system
type AddressSegmenter interface {
	Segment(address []byte) (segments []string)
	Combine(segments []string) (address []byte)
}

func NewFileSystemStore(rootDirectory string,
	addressSegmenter AddressSegmenter) Store {
	return &fileSystemStore{
		rootDirectory:    rootDirectory,
		addressSegmenter: addressSegmenter,
	}
}

func NewFlatAddressSegmenter() AddressSegmenter {
	return &flatAddressSegmenter{}
}

type flatAddressSegmenter struct{}

func (fas *flatAddressSegmenter) Segment(address []byte) []string {
	return []string{string(address)}
}

func (fas *flatAddressSegmenter) Combine(segments []string) []byte {
	return ([]byte)(segments[0])
}

func (fss *fileSystemStore) Put(address, data []byte) error {
	return ioutil.WriteFile(fss.Path(address), data, 0644)
}

func (fss *fileSystemStore) Get(address []byte) ([]byte, error) {
	return ioutil.ReadFile(fss.Path(address))
}

func (fss *fileSystemStore) Stat(address []byte) (*StatInfo, error) {
	fileInfo, err := os.Stat(fss.Path(address))
	statInfo := new(StatInfo)
	// Any kind of error means we should set exists false
	statInfo.Exists = err == nil
	if statInfo.Exists {
		statInfo.Size = uint64(fileInfo.Size())
	}
	// Don't treat not existing as an error
	if os.IsNotExist(err) {
		return statInfo, nil
	}
	return statInfo, err
}

func (fss *fileSystemStore) Location(address []byte) string {
	filePath := fss.Path(address)
	uri, err := url.Parse(filePath)
	if err != nil {
		return filePath
	}
	return uri.String()
}

func (fss *fileSystemStore) Path(address []byte) string {
	return path.Join(fss.rootDirectory,
		path.Join(fss.addressSegmenter.Segment(address)...))
}
