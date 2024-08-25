package signature

import (
	"bytes"
	"crypto/hmac"
    "hash"
	"encoding/gob"
)

type Signer struct {
    hash.Hash
    Key []byte
}

func New(h func() hash.Hash, key []byte) *Signer{
    return &Signer{
        Hash: hmac.New(h, key),
        Key: key,
    }
}

func (s *Signer) Sign(data any) ([]byte, error) {
    if v,ok := data.([]byte); ok {
        _, err := s.Write(v)
        if err != nil {
            return nil, err
        }
        return s.Sum(nil), nil
    }

    buf := bytes.NewBuffer(nil)

    err := gob.NewEncoder(buf).Encode(data)
    if err != nil {
        return nil, err
    }

    _, err = s.Write(buf.Bytes())
    if err != nil {
        return nil, err
    }
    return s.Sum(nil), nil

}
