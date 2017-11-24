package dslog

import (
	"context"
	"strings"
	"sync"

	"go.mercari.io/datastore"
)

var _ datastore.CacheStrategy = &logger{}

func NewLogger(prefix string, logf func(ctx context.Context, format string, args ...interface{})) datastore.CacheStrategy {
	return &logger{Prefix: prefix, Logf: logf, counter: 1}
}

type logger struct {
	Prefix string
	Logf   func(ctx context.Context, format string, args ...interface{})

	m       sync.Mutex
	counter int
}

func (l *logger) KeysToString(keys []datastore.Key) string {
	keyStrings := make([]string, 0, len(keys))
	for _, key := range keys {
		keyStrings = append(keyStrings, key.String())
	}

	return strings.Join(keyStrings, ", ")
}

func (l *logger) PutMultiWithoutTx(info *datastore.CacheInfo, keys []datastore.Key, psList []datastore.PropertyList) ([]datastore.Key, error) {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"PutMultiWithoutTx #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))

	keys, err := info.Next.PutMultiWithoutTx(info, keys, psList)

	if err == nil {
		l.Logf(info.Context, l.Prefix+"PutMultiWithoutTx #%d, keys=[%s]", cnt, l.KeysToString(keys))
	} else {
		l.Logf(info.Context, l.Prefix+"PutMultiWithoutTx #%d, err=%s", cnt, err.Error())
	}

	return keys, err
}

func (l *logger) PutMultiWithTx(info *datastore.CacheInfo, keys []datastore.Key, psList []datastore.PropertyList) ([]datastore.PendingKey, error) {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"PutMultiWithTx #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))

	pKeys, err := info.Next.PutMultiWithTx(info, keys, psList)

	if err != nil {
		l.Logf(info.Context, l.Prefix+"PutMultiWithTx #%d, err=%s", cnt, err.Error())
	}

	return pKeys, err
}

func (l *logger) GetMultiWithoutTx(info *datastore.CacheInfo, keys []datastore.Key, psList []datastore.PropertyList) error {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"GetMultiWithoutTx #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))

	err := info.Next.GetMultiWithoutTx(info, keys, psList)

	if err != nil {
		l.Logf(info.Context, l.Prefix+"GetMultiWithoutTx #%d, err=%s", cnt, err.Error())
	}

	return err
}

func (l *logger) GetMultiWithTx(info *datastore.CacheInfo, keys []datastore.Key, psList []datastore.PropertyList) error {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"GetMultiWithTx #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))

	err := info.Next.GetMultiWithTx(info, keys, psList)

	if err != nil {
		l.Logf(info.Context, l.Prefix+"GetMultiWithTx #%d, err=%s", cnt, err.Error())
	}

	return err
}

func (l *logger) DeleteMultiWithoutTx(info *datastore.CacheInfo, keys []datastore.Key) error {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"DeleteMultiWithoutTx #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))

	err := info.Next.DeleteMultiWithoutTx(info, keys)

	if err != nil {
		l.Logf(info.Context, l.Prefix+"DeleteMultiWithoutTx #%d, err=%s", cnt, err.Error())
	}

	return err
}

func (l *logger) DeleteMultiWithTx(info *datastore.CacheInfo, keys []datastore.Key) error {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"DeleteMultiWithTx #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))

	err := info.Next.DeleteMultiWithTx(info, keys)

	if err != nil {
		l.Logf(info.Context, l.Prefix+"DeleteMultiWithTx #%d, err=%s", cnt, err.Error())
	}

	return err
}

func (l *logger) PostCommit(info *datastore.CacheInfo, commit datastore.Commit) error {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"PostCommit #%d", cnt)

	return nil
}

func (l *logger) PostRollback(info *datastore.CacheInfo) error {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"PostRollback #%d", cnt)

	return nil
}

func (l *logger) Run(info *datastore.CacheInfo, q datastore.Query, qDump *datastore.QueryDump) datastore.Iterator {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"Run #%d, q=%s", cnt, qDump.String())

	return info.Next.Run(info, q, qDump)
}

func (l *logger) GetAll(info *datastore.CacheInfo, q datastore.Query, qDump *datastore.QueryDump, psList *[]datastore.PropertyList) ([]datastore.Key, error) {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"GetAll #%d, q=%s", cnt, qDump.String())

	keys, err := info.Next.GetAll(info, q, qDump, psList)

	if err == nil {
		l.Logf(info.Context, l.Prefix+"GetAll #%d, len(keys)=%d, keys=[%s]", cnt, len(keys), l.KeysToString(keys))
	} else {
		l.Logf(info.Context, l.Prefix+"GetAll #%d, err=%s", cnt, err.Error())
	}

	return keys, err
}

func (l *logger) Next(info *datastore.CacheInfo, q datastore.Query, qDump *datastore.QueryDump, iter datastore.Iterator, ps *datastore.PropertyList) (datastore.Key, error) {
	l.m.Lock()
	cnt := l.counter
	l.counter += 1
	l.m.Unlock()

	l.Logf(info.Context, l.Prefix+"Next #%d, q=%s", cnt, qDump.String())

	key, err := info.Next.Next(info, q, qDump, iter, ps)

	if err == nil {
		l.Logf(info.Context, l.Prefix+"Next #%d, key=%s", cnt, key.String())
	} else {
		l.Logf(info.Context, l.Prefix+"Next #%d, err=%s", cnt, err.Error())
	}

	return key, err
}
