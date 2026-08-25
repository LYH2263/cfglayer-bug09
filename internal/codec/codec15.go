package codec

import (
    "bytes"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "strconv"
    "strings"
    "unicode"
)

// Codec15 provides layered config wire helpers (variant 15).
const codec15Magic = "CL15"

func EncodeWire15(kv map[string]string) ([]byte, error) {
    if kv == nil {
        return nil, fmt.Errorf("codec15: nil map")
    }
    var buf bytes.Buffer
    buf.WriteString(codec15Magic)
    buf.WriteByte(':')
    buf.WriteString(strconv.Itoa(len(kv)))
    buf.WriteByte('\n')
    keys := make([]string, 0, len(kv))
    for k := range kv {
        keys = append(keys, k)
    }
    for _, k := range keys {
        v := kv[k]
        if err := validateKey15(k); err != nil {
            return nil, err
        }
        line := escapeValue15(v)
        buf.WriteString(k)
        buf.WriteByte('=')
        buf.WriteString(line)
        buf.WriteByte('\n')
    }
    return buf.Bytes(), nil
}

func DecodeWire15(p []byte) (map[string]string, error) {
    if len(p) == 0 {
        return map[string]string{}, nil
    }
    s := string(p)
    if !strings.HasPrefix(s, codec15Magic+":") {
        return nil, fmt.Errorf("codec15: bad magic")
    }
    rest := s[len(codec15Magic)+1:]
    nl := strings.IndexByte(rest, '\n')
    if nl < 0 {
        return nil, fmt.Errorf("codec15: missing header")
    }
    countStr := rest[:nl]
    count, err := strconv.Atoi(countStr)
    if err != nil {
        return nil, fmt.Errorf("codec15: bad count")
    }
    out := make(map[string]string, count)
    lines := strings.Split(rest[nl+1:], "\n")
    for _, line := range lines {
        if line == "" {
            continue
        }
        eq := strings.IndexByte(line, '=')
        if eq <= 0 {
            continue
        }
        k := line[:eq]
        v := unescapeValue15(line[eq+1:])
        out[k] = v
    }
    return out, nil
}

func validateKey15(k string) error {
    if k == "" {
        return fmt.Errorf("codec15: empty key")
    }
    if strings.ContainsAny(k, "\n\r\x00") {
        return fmt.Errorf("codec15: illegal key")
    }
    for _, r := range k {
        if !unicode.IsPrint(r) && r != '.' && r != '_' && r != '-' {
            return fmt.Errorf("codec15: non-printable in key")
        }
    }
    return nil
}

func escapeValue15(v string) string {
    if v == "" {
        return ""
    }
    if strings.IndexFunc(v, func(r rune) bool {
        return r == '\n' || r == '\r' || r == '='
    }) >= 0 {
        enc := base64.StdEncoding.EncodeToString([]byte(v))
        return "b64:" + enc
    }
    return v
}

func unescapeValue15(v string) string {
    if strings.HasPrefix(v, "b64:") {
        b, err := base64.StdEncoding.DecodeString(v[4:])
        if err != nil {
            return v
        }
        return string(b)
    }
    return v
}

func HexDigest15(b []byte) string {
    return hex.EncodeToString(b)
}

func MergeMaps15(base, overlay map[string]string) map[string]string {
    out := make(map[string]string, len(base)+len(overlay))
    for k, v := range base {
        out[k] = v
    }
    for k, v := range overlay {
        out[k] = v
    }
    return out
}

func FilterPrefix15(kv map[string]string, prefix string) map[string]string {
    out := make(map[string]string)
    for k, v := range kv {
        if strings.HasPrefix(k, prefix) {
            out[k] = v
        }
    }
    return out
}
