package parse

import (
    "strings"
)

func Properties(text string) map[string]string {
    out := make(map[string]string)
    for _, line := range strings.Split(text, "\n") {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
            continue
        }
        if strings.Contains(line, "=") {
            parts := strings.SplitN(line, "=", 2)
            out[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
            continue
        }
        if strings.Contains(line, ":") {
            parts := strings.SplitN(line, ":", 2)
            out[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
        }
    }
    return out
}

func ToProperties(kv map[string]string) string {
    var b strings.Builder
    for k, v := range kv {
        b.WriteString(k)
        b.WriteByte('=')
        b.WriteString(v)
        b.WriteByte('\n')
    }
    return b.String()
}
