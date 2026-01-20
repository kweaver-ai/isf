// eslint-disable-next-line spaced-comment
/// <reference path="./user.d.ts" />

import { ShareMgnt } from '../thrift';

export function getDepartmentOfUsers(id: string, start: number, limit: number): PromiseLike<any> {
    return ShareMgnt('Usrm_GetDepartmentOfUsers', [id, start, limit]);
}

export function getSubDepartments(parentId: string): PromiseLike<any> {
    return ShareMgnt('Usrm_GetSubDepartments', [parentId])
}

export function delUser(userid: string): PromiseLike<void> {
    return ShareMgnt('Usrm_DelUser', [userid]);
}

export function getUserInfo(userid: string): PromiseLike<any> {
    return ShareMgnt('Usrm_GetUserInfo', [userid]);
}