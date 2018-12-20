package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

//go:generate protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. config.proto

func json(body []byte) (*Configuration, error) {
	cfg := &Configuration{}
	if err := jsonpb.Unmarshal(bytes.NewReader(body), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func protobin(body []byte) (*Configuration, error) {
	cfg := &Configuration{}
	if err := proto.Unmarshal(body, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func prototxt(body []byte) (*Configuration, error) {
	cfg := &Configuration{}
	if err := proto.UnmarshalText(string(body), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func yml(body []byte) (*Configuration, error) {
	jsn, err := yaml.YAMLToJSON(body)
	if err != nil {
		return nil, err
	}
	return json(jsn)
}

type Parser = func([]byte) (*Configuration, error)

func defaultParserIndex() map[string]Parser {
	parserIndex := make(map[string]Parser)
	parserIndex[".json"] = json
	parserIndex[".yaml"] = yml
	parserIndex[".yml"] = yml
	parserIndex[".protobin"] = protobin
	parserIndex[".bin"] = protobin
	parserIndex[".prototxt"] = prototxt
	parserIndex[".txt"] = prototxt
	return parserIndex
}

func Load(url string) (*Configuration, error) {
	idx := defaultParserIndex()

	ext := path.Ext(url)
	parser, ok := idx[ext]
	if !ok {
		return nil, fmt.Errorf("")
	}

	body, err := ioutil.ReadFile(url)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	config, err := parser(body)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return config, nil
}
