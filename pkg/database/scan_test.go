package database

import (
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreScan(t *testing.T) {
	db, closer := makeDb()
	defer closer()

	db.Set(&schema.SetRequest{KVs: []*schema.KeyValue{{Key: []byte(`aaa`), Value: []byte(`item1`)}}})
	db.Set(&schema.SetRequest{KVs: []*schema.KeyValue{{Key: []byte(`bbb`), Value: []byte(`item2`)}}})

	txID, err := db.Set(&schema.SetRequest{KVs: []*schema.KeyValue{{Key: []byte(`abc`), Value: []byte(`item3`)}}})
	require.NoError(t, err)

	item, err := db.Get(&schema.KeyRequest{Key: []byte(`abc`), FromTx: int64(txID.Id)})
	require.Equal(t, []byte(`abc`), item.Key)
	require.NoError(t, err)

	scanOptions := schema.ScanRequest{
		Prefix:  []byte(`a`),
		Offset:  nil,
		Limit:   0,
		Reverse: true,
		Deep:    false,
	}

	db.WaitForIndexingUpto(txID.Id)

	list, err := db.Scan(&scanOptions)

	assert.NoError(t, err)
	assert.Exactly(t, 2, len(list.Items))
	assert.Equal(t, list.Items[0].Key, []byte(`aaa`))
	assert.Equal(t, list.Items[0].Value, []byte(`item1`))
	assert.Equal(t, list.Items[1].Key, []byte(`abc`))
	assert.Equal(t, list.Items[1].Value, []byte(`item3`))

	scanOptions1 := schema.ScanRequest{
		Prefix:  []byte(`a`),
		Offset:  nil,
		Limit:   0,
		Reverse: false,
		Deep:    false,
	}

	list1, err1 := db.Scan(&scanOptions1)
	assert.NoError(t, err1)
	assert.Exactly(t, 2, len(list1.Items))
	assert.Equal(t, list1.Items[0].Key, []byte(`abc`))
	assert.Equal(t, list1.Items[0].Value, []byte(`item3`))
	assert.Equal(t, list1.Items[1].Key, []byte(`aaa`))
	assert.Equal(t, list1.Items[1].Value, []byte(`item1`))
}
