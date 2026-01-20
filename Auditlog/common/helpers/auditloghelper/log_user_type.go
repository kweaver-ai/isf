package auditloghelper

type NcTLogUserType string

const (
	NcTLogUserType_NCT_LUT_AUTHUSER        NcTLogUserType = "authenticated_user"
	NcTLogUserType_NCT_LUT_ANONYUSER       NcTLogUserType = "anonymous_user"
	NcTLogUserType_NCT_LUT_APP             NcTLogUserType = "app"
	NcTLogUserType_NCT_LUT_INTERNALSERVICE NcTLogUserType = "internal_service"
)
