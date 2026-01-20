import { SystemRoleType } from '../role/role'
enum SystemRole {
    /*
     * 超级管理员
     */
    Super = 'super_admin',
    /**
     * 系统管理员
     */
    Admin = 'sys_admin',
    /**
     * 安全管理员
     */
    Securit = 'security',
    /**
     * 审计管理员
     */
    Audit = 'audit_admin',
    /**
     * 组织管理员
     */
    OrgManager = 'org_admin',
    /**
     * 组织审计员
     */
    OrgAudit = 'org_audit',
}
/**
 * 涉密模式下将后端返回的字段转换为对应角色的id
 */
export const transformAdmins = (admins: ReadonlyArray<SystemRole>): ReadonlyArray<SystemRoleType | undefined> => {
    return admins.map((admin) => {
        switch (admin) {
            case SystemRole.Super:
                return SystemRoleType.Supper
            case SystemRole.Admin:
                return SystemRoleType.Admin
            case SystemRole.Securit:
                return SystemRoleType.Securit
            case SystemRole.OrgManager:
                return SystemRoleType.OrgManager
            case SystemRole.Audit:
                return SystemRoleType.Audit
            case SystemRole.OrgAudit:
                return SystemRoleType.OrgAudit
            default:
                return
        }
    })
}