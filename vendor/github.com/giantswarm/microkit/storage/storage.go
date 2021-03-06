// Package storage provides interface and error specifications. The storage sub
// packages provide specific storage implementations.
package storage

import (
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	microerror "github.com/giantswarm/microkit/error"
	"github.com/giantswarm/microkit/storage/etcd"
	"github.com/giantswarm/microkit/storage/etcdv2"
	"github.com/giantswarm/microkit/storage/memory"
)

const (
	// KindMemory is the kind to be used to create a memory storage service.
	KindMemory = "memory"
	// KindEtcd is the kind to be used to create an etcd storage service.
	KindEtcd = "etcd"
	// KindEtcdV2 is the kind to be used to create an etcd v2 storage service.
	KindEtcdV2 = "etcdv2"
)

// Config represents the configuration used to create a storage service.
type Config struct {
	// Settings.
	EtcdAddress string
	EtcdPrefix  string
	Kind        string
}

// DefaultConfig provides a default configuration to create a new storage
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		EtcdAddress: "",
		EtcdPrefix:  "",
		Kind:        KindMemory,
	}
}

// New creates a new configured storage service.
func New(config Config) (Service, error) {
	// Settings.
	if config.Kind == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "kind must not be empty")
	}
	if config.Kind != KindMemory && config.Kind != KindEtcd && config.Kind != KindEtcdV2 {
		return nil, microerror.MaskAnyf(invalidConfigError, "kind must be one of: %s, %s, %s", KindMemory, KindEtcd, KindEtcdV2)
	}

	var err error

	var storageService Service
	{
		switch config.Kind {
		case KindMemory:
			storageConfig := memory.DefaultConfig()
			storageService, err = memory.New(storageConfig)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		case KindEtcd:
			if config.EtcdAddress == "" {
				return nil, microerror.MaskAnyf(invalidConfigError, "etcd address must not be empty")
			}

			etcdConfig := clientv3.Config{
				Endpoints:   []string{config.EtcdAddress},
				DialTimeout: 5 * time.Second,
			}
			etcdClient, err := clientv3.New(etcdConfig)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			storageConfig := etcd.DefaultConfig()
			storageConfig.EtcdClient = etcdClient
			storageConfig.Prefix = config.EtcdPrefix
			storageService, err = etcd.New(storageConfig)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		case KindEtcdV2:
			if config.EtcdAddress == "" {
				return nil, microerror.MaskAnyf(invalidConfigError, "etcd address must not be empty")
			}

			etcdConfig := client.Config{
				Endpoints: []string{config.EtcdAddress},
				Transport: client.DefaultTransport,
			}
			etcdClient, err := client.New(etcdConfig)
			if err != nil {
				panic(err)
			}

			storageConfig := etcdv2.DefaultConfig()
			storageConfig.EtcdClient = etcdClient
			storageConfig.Prefix = config.EtcdPrefix
			storageService, err = etcdv2.New(storageConfig)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}
	}

	return storageService, nil
}
