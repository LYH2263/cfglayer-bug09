package cfglayer

import "errors"

var (
    ErrClosed         = errors.New("cfglayer: merger closed")
    ErrNotFound       = errors.New("cfglayer: layer or key not found")
    ErrDuplicateLayer = errors.New("cfglayer: duplicate layer id")
    ErrBadInput       = errors.New("cfglayer: bad input")
    ErrEmptyStack     = errors.New("cfglayer: empty layer stack")
)
