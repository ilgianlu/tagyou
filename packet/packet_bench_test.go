package packet

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func BenchmarkStartConnect(b *testing.B) {
	buffer := bytes.NewReader([]byte{16, 59, 0, 4, 77, 81, 84, 84, 4, 6, 0, 60, 0, 15, 109, 113, 116, 116, 106, 115, 95, 53, 51, 48, 48, 102, 100, 54, 51, 0, 8, 108, 97, 115, 116, 119, 105, 108, 108, 0, 20, 97, 32, 118, 101, 114, 121, 32, 115, 104, 111, 114, 116, 32, 109, 101, 115, 115, 97, 103, 101})
	session := model.RunningSession{}
	for n := 0; n < b.N; n++ {
		p := Packet{}
		reader := bufio.NewReader(buffer)
		p.Parse(reader, &session)
	}
}

func BenchmarkStartSubscribe(b *testing.B) {
	buffer := []byte{130, 13, 149, 223, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 0}
	session := model.RunningSession{}
	for n := 0; n < b.N; n++ {
		p := Packet{}
		reader := bufio.NewReader(bytes.NewReader(buffer))
		p.Parse(reader, &session)
	}
}

func BenchmarkStartPublish(b *testing.B) {
	buffer := []byte{48, 20, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 72, 101, 108, 108, 111, 32, 109, 113, 116, 116}
	session := model.RunningSession{}
	for n := 0; n < b.N; n++ {
		p := Packet{}
		reader := bufio.NewReader(bytes.NewReader(buffer))
		p.Parse(reader, &session)
	}
}

func BenchmarkParseConnect(b *testing.B) {
	buffer := []byte{16, 59, 0, 4, 77, 81, 84, 84, 4, 6, 0, 60, 0, 15, 109, 113, 116, 116, 106, 115, 95, 53, 51, 48, 48, 102, 100, 54, 51, 0, 8, 108, 97, 115, 116, 119, 105, 108, 108, 0, 20, 97, 32, 118, 101, 114, 121, 32, 115, 104, 111, 114, 116, 32, 109, 101, 115, 115, 97, 103, 101}
	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		LastConnect:    time.Now().Unix(),
	}

	for n := 0; n < b.N; n++ {
		p := Packet{}
		reader := bufio.NewReader(bytes.NewReader(buffer))
		p.Parse(reader, &session)
	}
}

func BenchmarkParseSubscribe(b *testing.B) {
	buffer := []byte{130, 13, 149, 223, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 0}
	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		LastConnect:    time.Now().Unix(),
	}

	for n := 0; n < b.N; n++ {
		p := Packet{}
		reader := bufio.NewReader(bytes.NewReader(buffer))
		p.Parse(reader, &session)
	}
}

func BenchmarkParsePublish(b *testing.B) {
	buffer := []byte{48, 20, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 72, 101, 108, 108, 111, 32, 109, 113, 116, 116}
	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		LastConnect:    time.Now().Unix(),
	}

	for n := 0; n < b.N; n++ {
		p := Packet{}
		reader := bufio.NewReader(bytes.NewReader(buffer))
		p.Parse(reader, &session)
	}
}
