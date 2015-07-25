package protocol

type Protocol struct {
    Time int64 `json:"t"` //协议创建时间
    Type string `json:"pt"` //协议类型
    Payload *Payload `json:"p,omitempty"`
}

//struc tag 冒号后不要空格
type Payload struct {   //谁做什么
    Uid string  `json:"u"` //用户id
    MessageType string `json:"mt"`   //消息类型
    Content *Content `json:"c,omitempty"`
}


type Content struct {
    From string `json:"fr"` //from
    To string `json:"to,omitempty"` //to
    Message string `json:"m"` //message
}


//type Request struct {
//    Time int64 `json:"t"` //客户端请求时间
//    Type string `json:"ty"` //请求类型
//    Payload *Payload `json:"p,omitempty"`
//}
//
//type Response struct {
//    Time int64 `json:"t"` //服务器响应时间
//    Type string `json:"ty"` //请求类型
//    Payload *Payload `json:"p,omitempty"`
//}
