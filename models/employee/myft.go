package employee

type Myft struct {
	StaffID string `db:"staff_id"`
	MyftID  string `db:"myft_id"`
}

type FtcAccount struct {
	ID    string `json:"id" db:"user_id"`
	Email string `json:"email" db:"email"`
	IsVIP bool   `json:"isVip" db:"is_vip"`
}
