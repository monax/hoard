package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type ipfsStore struct {
	host     string
	encoding AddressEncoding
}

func NewIPFSStore(host string, encoding AddressEncoding) (*ipfsStore, error) {
	host = fmt.Sprintf("%s/api/v0", strings.TrimRight(host, "/"))
	_, err := http.Get(host)
	if err != nil {
		return nil, err
	}
	return &ipfsStore{
		host:     host,
		encoding: encoding,
	}, nil
}

func (inv *ipfsStore) Put(address []byte, data []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/add", inv.host)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// https://docs.ipfs.io/reference/api/http/
	fw, err := w.CreateFormField("arg")
	if _, err = fw.Write((data)[:]); err != nil {
		return address, nil
	}
	w.Close()

	// pinning is true by default
	req, err := http.NewRequest("POST", url, &b)
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

	// currently, we `add` the blob and read the return address
	return []byte(m["Name"].(string)), nil
}

func (inv *ipfsStore) Get(address []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/cat?arg=%s", inv.host, string(address))
	resp, err := http.Get(url)
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

func (inv *ipfsStore) Stat(address []byte) (*StatInfo, error) {
	url := fmt.Sprintf("%s/cat?arg=%s", inv.host, string(address))
	resp, err := http.Get(url)
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

func (inv *ipfsStore) Location(address []byte) string {
	return string(address)
}

func (inv *ipfsStore) Name() string {
	return fmt.Sprintf("ipfsStore[api=%s]", inv.host)
}
