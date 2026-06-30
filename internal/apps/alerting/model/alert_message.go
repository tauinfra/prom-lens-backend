package model

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type AlertManagerMessage struct {
	Status       string            `json:"status"`
	Receiver     string            `json:"receiver"`
	Alerts       []Alert           `json:"alerts"`
	CommonLabels map[string]string `json:"commonLabels"`
}

type LarkMsg struct {
	MsgType string    `json:"msg_type"`
	Content LarkPost  `json:"content,omitempty"`
	Card    *LarkCard `json:"card,omitempty"`
}

type LarkPost struct {
	Post LarkLang `json:"post"`
}

type LarkLang struct {
	ZhCn LarkBody `json:"zh_cn"`
}

type LarkBody struct {
	Title   string           `json:"title"`
	Content [][]LarkRichText `json:"content"`
}

type LarkRichText struct {
	Tag  string `json:"tag"`
	Text string `json:"text,omitempty"`
	Href string `json:"href,omitempty"`
}

type LarkCard struct {
	Header   LarkCardHeader    `json:"header"`
	Elements []LarkCardElement `json:"elements"`
}

type LarkCardHeader struct {
	Title    LarkCardText `json:"title"`
	Template string       `json:"template,omitempty"`
}

type LarkCardElement struct {
	Tag     string           `json:"tag"`
	Text    *LarkCardText    `json:"text,omitempty"`
	Actions []LarkCardAction `json:"actions,omitempty"`
}

type LarkCardText struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type LarkCardAction struct {
	Tag   string        `json:"tag"`
	Text  *LarkCardText `json:"text,omitempty"`
	Type  string        `json:"type,omitempty"`
	Value interface{}   `json:"value,omitempty"`
}
