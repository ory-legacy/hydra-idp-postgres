package account

type Account struct {
	ID       string `json:"id" gorethink:"id"`
	Username string `json:"username" valid:"required" gorethink:"username"`
	Password string `json:"-" gorethink:"password"`
	Data     string `json:"data,omitempty" valid:"json" gorethink:"data"`
}
