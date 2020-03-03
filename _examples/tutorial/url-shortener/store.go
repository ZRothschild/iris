package main

import (
	"bytes"

	bolt "github.com/etcd-io/bbolt"
)

//Panic，如果您不想因严重的INITIALIZE-ONLY-ERRORS 异常而将其更改

// Panic panics, change it if you don't want to panic on critical INITIALIZE-ONLY-ERRORS
var Panic = func(v interface{}) {
	panic(v)
}

// Store是网址的存储接口
//注意：没有Del功能

// Store is the store interface for urls.
// Note: no Del functionality.
type Store interface {
	Set(key string, value string) error // 如果出了问题返回错误 | error if something went wrong
	Get(key string) string              // 如果找不到，则为空值 | empty value if not found
	Len() int                           // 应该返回所有记录/表/桶的数量 | should return the number of all the records/tables/buckets
	Close()                             // 释放存储或忽略 | release the store or ignore
}

var tableURLs = []byte("urls")

//Store的数据库表示形式。
//只有一个表/存储桶包含网址，因此它不是完整的数据库，
//它仅适用于单个存储桶，因为我们需要这些。

// DB representation of a Store.
// Only one table/bucket which contains the urls, so it's not a fully Database,
// it works only with single bucket because that all we need.
type DB struct {
	db *bolt.DB
}

var _ Store = &DB{}

// openDatabase打开一个新的数据库连接并返回其实例

// openDatabase open a new database connection
// and returns its instance.
func openDatabase(stumb string) *bolt.DB {
	//打开当前工作目录下的data（base）文件，
	// 如果不存在该文件将被创建

	// Open the data(base) file in the current working directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(stumb, 0600, nil)
	if err != nil {
		Panic(err)
	}
	//在此处创建存储桶

	// create the buckets here
	tables := [...][]byte{
		tableURLs,
	}

	db.Update(func(tx *bolt.Tx) (err error) {
		for _, table := range tables {
			_, err = tx.CreateBucketIfNotExists(table)
			if err != nil {
				Panic(err)
			}
		}

		return
	})

	return db
}

// NewDB返回一个新的数据库实例，其连接已打开
// DB实现Store

// NewDB returns a new DB instance, its connection is opened.
// DB implements the Store.
func NewDB(stumb string) *DB {
	return &DB{
		db: openDatabase(stumb),
	}
}

// Set设置一个缩短的网址及其键
//注意：调用方负责生成密钥

// Set sets a shorten url and its key
// Note: Caller is responsible to generate a key.
func (d *DB) Set(key string, value string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(tableURLs)
		//生成url的ID注意：我们可以使用它代替随机的字符串键
		// 但是我们想模拟一个实际的url缩短器，因此我们跳过它
		// id, _ := b.NextSequence()

		// Generate ID for the url
		// Note: we could use that instead of a random string key
		// but we want to simulate a real-world url shortener
		// so we skip that.
		// id, _ := b.NextSequence()
		if err != nil {
			return err
		}

		k := []byte(key)
		valueB := []byte(value)
		c := b.Cursor()

		found := false
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(valueB, v) {
				found = true
				break
			}
		}
		//如果值已经存在，请不要重新输入

		// if value already exists don't re-put it.
		if found {
			return nil
		}

		return b.Put(k, []byte(value))
	})
}

//Clear将清除表URL的所有数据库条目

// Clear clears all the database entries for the table urls.
func (d *DB) Clear() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(tableURLs)
	})
}

// Get通过其键返回一个URL
//
//如果找不到，则返回一个空字符串

// Get returns a url by its key.
//
// Returns an empty string if not found.
func (d *DB) Get(key string) (value string) {
	keyB := []byte(key)
	d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableURLs)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(keyB, k) {
				value = string(v)
				break
			}
		}

		return nil
	})

	return
}

// GetByValue返回特定（原始）URL值的所有键

// GetByValue returns all keys for a specific (original) url value.
func (d *DB) GetByValue(value string) (keys []string) {
	valueB := []byte(value)
	d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableURLs)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		//首先为存储区的表"urls"
		// first for the bucket's table "urls"
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(valueB, v) {
				keys = append(keys, string(k))
			}
		}

		return nil
	})

	return
}

// Len返回所有“短”网址的长度
// Len returns all the "shorted" urls length
func (d *DB) Len() (num int) {
	d.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(tableURLs)
		if b == nil {
			return nil
		}

		b.ForEach(func([]byte, []byte) error {
			num++
			return nil
		})
		return nil
	})
	return
}

//关闭将关闭数据库）连接

// Close shutdowns the data(base) connection.
func (d *DB) Close() {
	if err := d.db.Close(); err != nil {
		Panic(err)
	}
}
