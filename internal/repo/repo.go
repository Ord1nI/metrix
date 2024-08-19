package repo

import (
    "encoding/json"
)

type Repo interface {
    GetAdder
    Closer
    json.Marshaler
}

type Adder interface {
    Add(name string, val interface{}) (error)
}

type GetAdder interface {
    Adder
    Getter
}

type Getter  interface {
    Get(name string, val interface{}) (error)
}

type Closer interface {
    Close() error
}

