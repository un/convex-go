package sync

import (
    "fmt"
    "net/url"
)

func DeploymentURLToWebSocketURL(input string) (string, error) {
    parsed, err := url.Parse(input)
    if err != nil {
        return "", err
    }
    switch parsed.Scheme {
    case "http":
        parsed.Scheme = "ws"
    case "https":
        parsed.Scheme = "wss"
    case "ws", "wss":
    default:
        return "", fmt.Errorf("unsupported scheme %q", parsed.Scheme)
    }
    parsed.Path = "/api/sync"
    parsed.RawQuery = ""
    parsed.Fragment = ""
    return parsed.String(), nil
}
