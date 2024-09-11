//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"
import "unsafe"

type Blowfish struct { //comment
	key C.BF_KEY
}

func (b Blowfish) Encrypt(out, in []byte) {
	C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&in[0])), (*C.uchar)(unsafe.Pointer(&out[0])), &b.key, C.BF_ENCRYPT)
}

func (b Blowfish) Decrypt(out, in []byte) {
	C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&in[0])), (*C.uchar)(unsafe.Pointer(&out[0])), &b.key, C.BF_DECRYPT)
}

func New(key []byte) *Blowfish {
	newbf := &Blowfish{}
	C.BF_set_key(&newbf.key, (C.int)(len(key)), (*C.uchar)(unsafe.Pointer(&key[0])))
	return newbf
}

func (b Blowfish) BlockSize() int {
	return 8
}
