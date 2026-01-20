import { getDepartmentRoots as getOrgDepartmentRoots, deleteDepartment } from '../apis/console/usermanagement';
import { DepartmentEntry } from '../apis/console/usermanagement/types'
import { PublicErrorCode } from '../apis/openapiconsole/errorcode';
import { ROOT_DEPARTMENT_ID } from '../constants';
import { UserRole, getRoleType } from '../role/role';

const LIMIT = 1000

/**
 * 循环获取所有组织信息
 */
export async function getOrganizations(isRequestNormal?: boolean): Promise<ReadonlyArray<DepartmentEntry>> {
    let offset = 0
    let allEntries: ReadonlyArray<any> = []
    const role = getRoleType()

    while (true) {
        const { departments: { entries } } = await getOrgDepartmentRoots({ departmentId: ROOT_DEPARTMENT_ID, fields: ['departments'], role: isRequestNormal ? UserRole.NormalUser : role, offset, limit: LIMIT })

        allEntries = [...allEntries, ...entries]

        if (entries.length < LIMIT) {
            break
        }

        offset += LIMIT
    }

    return allEntries
}

/**
 * 删除部门、组织
 */
export const deleteDep = async (id: string): Promise<void> => {
    try {
        await deleteDepartment(id)
    } catch(ex) {
        // 组织、部门不存在，直接删除成功
        if (ex.code !== PublicErrorCode.NotFound) {
            throw ex
        }
    }
}
