package core

import (
	"strconv"
	"testing"
)

func TestPutData(t *testing.T) {

}

func BenchmarkPutData(b *testing.B) {
	go func() {
		for i := 0; i < b.N; i++ {
			UserHub.PutData(strconv.Itoa(i))
		}
	}()
	go func() {
		for i := b.N + 1; i < b.N*2; i++ {
			UserHub.PutData(strconv.Itoa(i))
		}
	}()
	for i := b.N*2 + 1; i < b.N*3; i++ {
		UserHub.PutData(strconv.Itoa(i))
	}
}
