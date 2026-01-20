import { map, flattenDeep, isArray } from 'lodash'
import { getDepartmentOfUsers, getSubDepartments } from '@/core/thrift/user/user';

export enum Range {
    DEPARTMENT_DEEP, // 部门及其子部门

    DEPARTMENT, // 当前部门

    USERS, // 当前选中的用户
}

/**
 * 系统类型
 */
export enum SystemType {
    /**
     * 控制台
     */
    Console,
}

/**
 * 递归获取所有部门及子部门
 * @param deps {object|object[]}
 * @param result {undefined || object[]}
 * @return {promise}
 */
export async function listDepsSince(deps) {
    const department = isArray(deps) ? deps : [deps];
    let result = department

    for (const dep of department) {
        const resultDep = await getSubDepartments(dep.id);

        if (resultDep.length) {
            const newReuslt = await listDepsSince(resultDep);

            result = [...result, ...newReuslt]
        }
    }

    return result
}

/**
 * 列举部门下所有用户（含子部门）
 * @param dep {object}
 * @param includeDep {boolean}是否返回部门
 * @return {promise}
 */
export async function listUsersSince(dep, includeDep?) {
    const department = await listDepsSince(dep)
    let results = [];

    for (const dep of department) {
        const result = await getDepartmentOfUsers(dep.id, 0, -1);

        results = results.concat(flattenDeep(result));
    }

    return includeDep ? { users: results, deps: department } : results
}

/**
  * 获取选中的用户
  */
export function getSeletedUsers(range, dep, users?) {

    if (range === Range.USERS) {
        return users
    } else if (range === Range.DEPARTMENT) {
        return getDepartmentOfUsers(dep.id, 0, -1)
    } else if (range === Range.DEPARTMENT_DEEP) {
        return listUsersSince(dep);
    }
}
