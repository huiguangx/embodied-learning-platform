package auth

type Role string

const (
	Engineer Role = "engineer"
	DataEngineer Role = "data_engineer"
	Operator Role = "operator"
	BusinessAdmin Role = "business_admin"
)

func (r Role) CanWrite() bool { return r == Operator || r == BusinessAdmin }
