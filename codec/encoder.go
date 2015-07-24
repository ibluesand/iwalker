package codec
import (
    "encoding/json"
    "github.com/ibluesand/iwalker/protocol"
)


func Eecoder(p protocol.Protocol) ([]byte, error) {
    data, err := json.Marshal(p)
    return data,err
}