package types

type LoginMethodDetailReply struct {
	Name string `json:"name"`
}

type VerificationCodeReq struct {
	Email string `json:"email" binding:"required"`
}

type LoginMethodReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoginMethods LoginMethodDetailReply `json:"loginMethod"`
		Countries    []CountriesObjDetail   `json:"countries"`
		Languages    []LanguagesObjDetail   `json:"languages"`
	} `json:"data"`
}

type SignUpVerificationCodeRely struct {
	Result
}

type ResetVerificationCodeRely struct {
	Result
}
