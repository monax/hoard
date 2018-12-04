package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type ipfsStore struct {
	Protocol        string
	Address         string
	Port            string
	addressEncoding AddressEncoding
}

func NewIPFSStore(proto, address, port string, addressEncoding AddressEncoding) (*ipfsStore, error) {
	if proto == "http://" {
		fmt.Println("Warning: IPFS connection not secure.")
	}
	_, err := http.Get(proto + address + ":" + port + "/api/v0/")
	if err != nil {
		return nil, err
	}
	return &ipfsStore{
		Protocol:        proto,
		Address:         address,
		Port:            port,
		addressEncoding: addressEncoding,
	}, nil
}

func (ipfss *ipfsStore) Put(address []byte, data []byte) ([]byte, error) {
	uri := ipfss.Protocol + ipfss.Address + ":" + ipfss.Port + "/api/v0/add"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField("arg")
	if _, err = fw.Write((data)[:]); err != nil {
		return address, nil
	}
	w.Close()

	req, err := http.NewRequest("POST", uri, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return address, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return address, err
	}
	var m map[string]interface{}
	json.Unmarshal(body.Bytes(), &m)

	// TODO deterministically generate IPFS addresses natively
	// currently, we `add` the blob and read the return address
	return []byte(m["Name"].(string)), nil
}

func (ipfss *ipfsStore) Get(address []byte) ([]byte, error) {
	resp, err := http.Get(ipfss.Protocol + ipfss.Address + ":" + ipfss.Port + "/api/v0/cat?arg=" + string(address))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (ipfss *ipfsStore) Stat(address []byte) (*StatInfo, error) {
	resp, err := http.Get(ipfss.Protocol + ipfss.Address + ":" + ipfss.Port + "/api/v0/cat?arg=" + string(address[:]))
	if err != nil || resp.StatusCode != 200 {
		return &StatInfo{
			Exists: false,
		}, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &StatInfo{
		Exists: true,
		Size:   uint64(len(body)),
	}, nil
}

func (ipfss *ipfsStore) Location(address []byte) string {
	return string(address)
}

func (ipfss *ipfsStore) Name() string {
	return fmt.Sprintf("ipfsStore[api=%s:%s]", ipfss.Address, ipfss.Port)
}
