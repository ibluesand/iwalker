package codec
import (
    "github.com/ibluesand/iwalker/protocol"

    "encoding/json"
)


func Decoder(data []byte, p *protocol.Protocol) error {
    err := json.Unmarshal(data, &p)
    return err
}