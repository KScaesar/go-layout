package utility

import (
	"time"
)

// ReadProxy 用於管理多個資料節點之間的讀取和寫入操作,
// 允許從 Replica 中快速獲取資料, 如果無法取得，則嘗試從 Primary 中讀取.
//
// 避免同一個 key 的多個併發請求到達 Primary or Replica, 解決 Hotspot Invalid, Cache Avalanche 問題
//
// ( Primary, Replica ) 可以分別代表不同的資料存取方式,
// 例如 ( Database, Cache ) 或 ( RemoteCache, LocalCache )
type ReadProxy[ViewModel any, Read func(key string) (ViewModel, error), Write func(string, *ViewModel) error] struct {
	ReadReplica  Read
	ReadPrimary  Read
	WriteReplica Write
	SingleFlight *Singleflight
}

// SafeReadPrimaryNode 併發時, 只會保護 Primary Node, 不會保護 Replica Node
func (proxy ReadProxy[ViewModel, Read, Write]) SafeReadPrimaryNode(key string) (val ViewModel, err error) {
	val, err = proxy.ReadReplica(key)
	if err == nil {
		return val, nil
	}
	return proxy.SafeReadPrimaryAndReplicaNode(key)
}

// SafeReadPrimaryAndReplicaNode 併發時, 會保護 Primary Node and Replica Node
func (proxy ReadProxy[ViewModel, Read, Write]) SafeReadPrimaryAndReplicaNode(key string) (val ViewModel, err error) {
	value, err, _ := proxy.SingleFlight.Do(key, func() (any, error) {
		return proxy.Read(key)
	})
	proxy.SingleFlight.Expire(key, time.Second)
	// proxy.SingleFlight.Forget(key)
	return value.(ViewModel), err
}

func (proxy ReadProxy[ViewModel, Read, Write]) Read(key string) (val ViewModel, err error) {
	val, err = proxy.ReadReplica(key)
	if err == nil {
		return val, nil
	}

	val, err = proxy.ReadPrimary(key)
	if err != nil {
		return val, err
	}

	err = proxy.WriteReplica(key, &val)
	if err != nil {
		return val, err
	}

	return val, nil
}
