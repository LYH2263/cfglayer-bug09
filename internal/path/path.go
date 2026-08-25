package path

import (
    "strings"
)

const sep = "."

func Join(parts ...string) string {
    clean := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.Trim(p, sep)
        if p != "" {
            clean = append(clean, p)
        }
    }
    return strings.Join(clean, sep)
}

func Split(key string) []string {
    if key == "" {
        return nil
    }
    return strings.Split(key, sep)
}

func Parent(key string) string {
    i := strings.LastIndex(key, sep)
    if i <= 0 {
        return ""
    }
    return key[:i]
}

func Base(key string) string {
    i := strings.LastIndex(key, sep)
    if i < 0 {
        return key
    }
    return key[i+1:]
}

func HasPrefix(key, prefix string) bool {
    if prefix == "" {
        return true
    }
    if key == prefix {
        return true
    }
    return strings.HasPrefix(key, prefix+sep)
}

func Normalize(key string) string {
    parts := Split(key)
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p != "" {
            out = append(out, p)
        }
    }
    return strings.Join(out, sep)
}

func Descendants(root string, keys []string) []string {
    var out []string
    for _, k := range keys {
        if HasPrefix(k, root) {
            out = append(out, k)
        }
    }
    return out
}
