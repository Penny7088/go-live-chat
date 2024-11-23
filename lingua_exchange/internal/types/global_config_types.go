package types

type LoginMethodDetailReply struct {
	Name string `json:"name"`
}

type LoginMethodReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoginMethods []LoginMethodDetailReply `json:"loginMethods"`
	} `json:"data"`
}
