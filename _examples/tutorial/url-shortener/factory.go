package main

import (
	"net/url"

	"github.com/iris-contrib/go.uuid"
)

//生成类型以生成密钥（短网址）

// Generator the type to generate keys(short urls)
type Generator func() string

// DefaultGenerator是默认的URL生成器

// DefaultGenerator is the defautl url generator
var DefaultGenerator = func() string {
	id, _ := uuid.NewV4()
	return id.String()
}

//工厂负责生成密钥（短网址）

// Factory is responsible to generate keys(short urls)
type Factory struct {
	store     Store
	generator Generator
}

// NewFactory接收一个生成器和一个存储，并返回一个新的URL Factory。

// NewFactory receives a generator and a store and returns a new url Factory.
func NewFactory(generator Generator, store Store) *Factory {
	return &Factory{
		store:     store,
		generator: generator,
	}
}

// Gen生成密钥

// Gen generates the key.
func (f *Factory) Gen(uri string) (key string, err error) {
	//我们不返回已解析的url，因为#hash已转换为uri兼容，并且我们不想一直进行编码/解码，
	// 因此不需要这样做，我们将URL保存为用户期望的值，如果 uri验证已通过

	// we don't return the parsed url because #hash are converted to uri-compatible
	// and we don't want to encode/decode all the time, there is no need for that,
	// we save the url as the user expects if the uri validation passed.
	_, err = url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}

	key = f.generator()
	//确保密钥是唯一的

	// Make sure that the key is unique
	for {
		if v := f.store.Get(key); v == "" {
			break
		}
		key = f.generator()
	}

	return key, nil
}
