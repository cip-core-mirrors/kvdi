package grants

type RoleGrant int

const (
	ReadUsers RoleGrant = 1 << iota
	WriteUsers
	ReadRoles
	WriteRoles
	ReadTemplates
	WriteTemplates
	LaunchTemplates
	ReadDesktopSessions
)

const (
	All RoleGrant = ReadUsers | WriteUsers |
		ReadRoles | WriteRoles |
		ReadTemplates | WriteTemplates |
		LaunchTemplates | ReadDesktopSessions
)

var grantNames = []string{
	"ReadUsers",
	"WriteUsers",
	"ReadRoles",
	"WriteRoles",
	"ReadTemplates",
	"WriteTemplates",
	"LaunchTemplates",
	"ReadDesktopSessions",
}

func (r RoleGrant) Has(grant RoleGrant) bool { return r&grant != 0 }

func (r RoleGrant) Names() []string {
	var result []string
	for i := 0; i < len(grantNames); i++ {
		if r&(1<<uint(i)) != 0 {
			result = append(result, grantNames[i])
		}
	}
	return result
}