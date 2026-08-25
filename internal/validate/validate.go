package validate

import (
    "errors"
    "strings"
)

var ErrBad = errors.New("validate: bad input")

func LayerID(id string) error {
    if id == "" || len(id) > 128 {
        return ErrBad
    }
    if strings.ContainsAny(id, "/\\ \t\n") {
        return ErrBad
    }
    return nil
}

func Key(key string) error {
    if key == "" || len(key) > 512 {
        return ErrBad
    }
    if strings.ContainsAny(key, "\n\r\x00") {
        return ErrBad
    }
    return nil
}

type LayerInput struct {
    ID     string
    Values map[string]string
}

func Layer(layer LayerInput) error {
    if err := LayerID(layer.ID); err != nil {
        return err
    }
    for k := range layer.Values {
        if err := Key(k); err != nil {
            return err
        }
    }
    return nil
}
