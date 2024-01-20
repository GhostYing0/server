package role

const (
	RoleNumber = 3
)

var RoleTable [RoleNumber]string

func InitRoleTable() {
	RoleTable = [RoleNumber]string{
		"cms",
		"student",
		"teacher",
	}
}
