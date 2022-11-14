package types

type Users struct {
	UserId         int    `json:"user_id"`
	MailLogin      string `json:"mail_login"`
	MailPassword   string `json:"mail_password"`
	MailService    string `json:"mail_service"`
	TotalMsgCount  int    `json:"total_msg_count"`
	UnseenMsgCount int    `json:"unseen_msg_count"`
}

type Response struct {
	UserId         int
	From           string
	Body           string
	TotalMsgCount  int
	UnseenMsgCount int
}
