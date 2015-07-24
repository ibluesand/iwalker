package protocol

type Protocol struct {
    Time int64 `json:"t"` //协议创建时间
    ProtocolType string `json:"pt"` //协议类型
    Payload *Payload `json:"p,omitempty"`
}

//struc tag 冒号后不要空格
type Payload struct {
    Uid string  `json:"u"` //用户id
    MessageType string `json:"mt"`   //消息类型
    Content *Content `json:"c,omitempty"`
}


type Content struct {
    From string `json:"fr"` //from
    To string `json:"to,omitempty"` //to
    Message string `json:"m"` //message
}
