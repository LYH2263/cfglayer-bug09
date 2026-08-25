package parse

import (
    "bufio"
    "strings"

    "github.com/LYH2263/go-cfglayer"
)

func LinesKV(text string) (map[string]string, error) {
    out := make(map[string]string)
    sc := bufio.NewScanner(strings.NewReader(text))
    lineNo := 0
    for sc.Scan() {
        lineNo++
        line := strings.TrimSpace(sc.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        eq := strings.IndexByte(line, '=')
        if eq <= 0 {
            return nil, cfglayer.ErrBadInput
        }
        k := strings.TrimSpace(line[:eq])
        v := strings.TrimSpace(line[eq+1:])
        if k == "" {
            return nil, cfglayer.ErrBadInput
        }
        out[k] = v
    }
    if err := sc.Err(); err != nil {
        return nil, err
    }
    return out, nil
}

func LayerFromText(id, source, text string, priority int) (cfglayer.Layer, error) {
    vals, err := LinesKV(text)
    if err != nil {
        return cfglayer.Layer{}, err
    }
    return cfglayer.Layer{
        ID:       id,
        Priority: priority,
        Source:   source,
        Values:   vals,
    }, nil
}
