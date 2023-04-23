package aes

import "testing"

const (
	key = "051487b9f9b666e1"
	iv  = "9d487f08e3674ae2"
)

func TestEncrypt(t *testing.T) {
	t.Log(NewAes(key, iv).Encrypt("123456"))
}

func TestDecrypt(t *testing.T) {
	t.Log(NewAes(key, iv).Decrypt("zc7LHIpY/RLi/z/cFTpnyw=="))
}

func BenchmarkEncAndDec(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewAes(key, iv).Encrypt("123456")
		NewAes(key, iv).Decrypt("zc7LHIpY/RLi/z/cFTpnyw==")
	}
	b.StopTimer()
}
