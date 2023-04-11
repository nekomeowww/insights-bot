package telegram

type MemberStatus string

const (
	MemberStatusCreator       MemberStatus = "creator"
	MemberStatusAdministrator MemberStatus = "administrator"
	MemberStatusMember        MemberStatus = "member"
	MemberStatusRestricted    MemberStatus = "restricted"
	MemberStatusLeft          MemberStatus = "left"
	MemberStatusKicked        MemberStatus = "kicked"
)
