package auditloghelper

// NcTManagementType 管理日志操作类型
type NcTManagementType int64

const (
	NcTManagementType_NCT_MNT_ALL               NcTManagementType = 0
	NcTManagementType_NCT_MNT_CREATE            NcTManagementType = 1
	NcTManagementType_NCT_MNT_ADD               NcTManagementType = 2
	NcTManagementType_NCT_MNT_SET               NcTManagementType = 3
	NcTManagementType_NCT_MNT_DELETE            NcTManagementType = 4
	NcTManagementType_NCT_MNT_COPY              NcTManagementType = 5
	NcTManagementType_NCT_MNT_MOVE              NcTManagementType = 6
	NcTManagementType_NCT_MNT_REMOVE            NcTManagementType = 7
	NcTManagementType_NCT_MNT_IMPORT            NcTManagementType = 8
	NcTManagementType_NCT_MNT_EXPORT            NcTManagementType = 9
	NcTManagementType_NCT_MNT_AUDIT_MGM         NcTManagementType = 10
	NcTManagementType_NCT_MNT_QUARANTINE        NcTManagementType = 11
	NcTManagementType_NCT_MNT_UPLOAD            NcTManagementType = 12
	NcTManagementType_NCT_MNT_PREVIEW           NcTManagementType = 13
	NcTManagementType_NCT_MNT_DOWNLOAD          NcTManagementType = 14
	NcTManagementType_NCT_MNT_RESTORE           NcTManagementType = 15
	NcTManagementType_NCT_MNT_QUARANTINE_APPEAL NcTManagementType = 16
	NcTManagementType_NCT_MNT_RESTART           NcTManagementType = 17
	NcTManagementType_NCT_MNT_SEND_EMAIL        NcTManagementType = 18
	NcTManagementType_NCT_MNT_RECOVER           NcTManagementType = 19
	NcTManagementType_NCT_MNT_EDIT              NcTManagementType = 20
	NcTManagementType_NCT_MNT_OTHER             NcTManagementType = 127
)
