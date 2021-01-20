package consul

import (
	"bytes"
	"encoding/json"

	"github.com/hashicorp/consul/api"
)

// NewKV creates a KV.
func NewKV() (*KV, error) {
	cfg := api.DefaultConfig()
	cl, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &KV{inner: cl.KV()}, nil
}

// KV represents a simplified key-value store on top of Hashicorp Consul.
type KV struct {
	inner *api.KV
}

// Get returns a list of strings found at key.
func (k *KV) Get(key string) ([]string, error) {
	pair, _, err := k.inner.Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	var items []string
	if err = json.NewDecoder(bytes.NewReader(pair.Value)).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// Put inserts a list of strings at key.
func (k *KV) Put(key string, values []string) error {
	var data bytes.Buffer
	if err := json.NewEncoder(&data).Encode(values); err != nil {
		return err
	}

	pair := &api.KVPair{Key: key, Value: data.Bytes()}
	_, err := k.inner.Put(pair, nil)
	return err
}

// Delete removes the record at key.
func (k *KV) Delete(key string) error {
	_, err := k.inner.Delete(key, nil)
	return err
}
