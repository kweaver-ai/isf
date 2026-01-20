import * as React from 'react';
import { noop } from 'lodash'
import session from '@/util/session';
import { getDepParentPathById, moveDepartment, getSubDepartments, editDepartOSS, getDepartmentOfUsersCount, getDepartmentUser, editUserOss, getOrgDepartmentById } from '@/core/thrift/sharemgnt/sharemgnt';
import { getObjectStorageInfoById } from '@/core/apis/console/ossgateway'
import { manageLog, Level, ManagementOps } from '@/core/log';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { displayUserOssInfo } from '@/core/oss/oss';
import { Message2, Toast } from '@/sweet-ui';
import WebComponent from '../webcomponent';
import { Status } from './helper';
import __ from './locale';

interface MoveDepartmentProps {
    /**
     * 要移动的部门
     */
    srcDep: Core.ShareMgnt.ncTDepartmentInfo;

    /**
     * 取消移动部门
     */
    onRequestCancelMoveDep: () => void;

    /**
     * 移动部门完成
     */
    onRequestMoveDepFinished: (srcDep: Core.ShareMgnt.ncTDepartmentInfo, targetDep: Core.ShareMgnt.ncTDepartmentInfo, ossInfo: { ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo } | null) => void;

    /**
     * 源部门不存在，移除源部门
     */
    onRequestRemoveSrcDep: (srcDep: Core.ShareMgnt.ncTDepartmentInfo) => void;
}

interface MoveDepartmentState {
    /**
     * 目标部门
     */
    targetDep: Core.ShareMgnt.ncTDepartmentInfo;

    /**
     * 显示界面状态
     */
    status: Status;
}

export default class MoveDepartmentBase extends WebComponent<MoveDepartmentProps, MoveDepartmentState> {
    static defaultProps = {
        onRequestCancelMoveDep: noop,
        onRequestMoveDepFinished: noop,
        onRequestRemoveSrcDep: noop,
    }

    state: MoveDepartmentState = {
        targetDep: null,
        status: Status.Config,
    }

    userid = session.get('isf.userid');

    ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo;

    /**
     * 选中目标部门
     */
    protected selectDep = (targetDep: Core.ShareMgnt.ncTDepartmentInfo): void => {
        this.setState({
            targetDep,
        })
    }

    /**
     * 确定
     */
    protected confirmMoveDep = async (): Promise<void> => {
        try {
            this.setState({
                status: Status.None,
            })

            let isChangeOss = false
            let ossInfo = this.state.targetDep.ossInfo

            if (!ossInfo) {
                ({ ossInfo } = await getOrgDepartmentById(this.state.targetDep.id))
            }

            if (ossInfo.ossId) {
                const ossData = await getObjectStorageInfoById(ossInfo.ossId)
                this.ossInfo = {
                    enabled: ossData.enabled,
                    ossId: ossData.id,
                    ossName: ossData.name,
                }
            } else {
                this.ossInfo = ossInfo
            }

            this.moveDepRequst(false);
        }catch(ex) {
            if (ex && ex.error && ex.error.errID) {
                if (ex.error.errID === ErrorCode.DepNameNotExist) {
                    this.setState({
                        status: Status.None,
                    })

                    await Message2.info({ message: this.getErrorMsg(ex,this.props.srcDep, this.state.targetDep) })

                    this.props.onRequestRemoveSrcDep(this.state.targetDep)
                } else {
                    await Message2.info({ message: this.getErrorMsg(ex, this.props.srcDep, this.state.targetDep) })

                    this.setState({
                        status: Status.Config,
                        targetDep: null,
                    })
                }
            }
        }
    }

    /**
     * 发起移动部门
     */
    private async moveDepRequst(isChangeOss: boolean): Promise<void> {
        const { srcDep } = this.props
        const { targetDep } = this.state

        let ossChanged = isChangeOss

        try {
            const srcDepPath = await this.getDepPath(srcDep)

            await moveDepartment([srcDep.id, targetDep.id])

            const newDepPath = `${(await this.getDepPath(targetDep))}/${srcDep.name}`

            await manageLog(
                ManagementOps.MOVE,
                __(
                    targetDep.is_root ?
                        '移动部门 “${name}” 至组织 “${targetDepName}” 成功' :
                        '移动部门 “${name}” 至部门 “${targetDepName}” 成功',
                    {
                        name: srcDep.name,
                        targetDepName: targetDep.name,
                    },
                ),
                __('源组织路径 “${srcDepPath}”；新组织路径 “${newDepPath}”', { srcDepPath, newDepPath }),
                Level.INFO,
            )

            if (isChangeOss) {
                this.setState({
                    status: Status.Loading,
                })

                ossChanged = await this.changeStorage(srcDep, this.ossInfo)
            }

            Toast.open(__('移动部门成功'))

            this.props.onRequestMoveDepFinished(srcDep, targetDep, ossChanged ? { ossInfo: this.ossInfo } : null)
        } catch (ex) {
            if (ex && ex.error && ex.error.errID) {
                if (ex.error.errID === ErrorCode.DepNotExist) {
                    this.setState({
                        status: Status.None,
                    })

                    await Message2.info({ message: this.getErrorMsg(ex, srcDep, targetDep) })

                    this.props.onRequestRemoveSrcDep(this.props.srcDep)
                } else {
                    await Message2.info({ message: this.getErrorMsg(ex, srcDep, targetDep) })

                    this.setState({
                        status: Status.Config,
                        targetDep: null,
                    })
                }
            }
        }
    }

    /**
     * 获取部门路径
     */
    private async getDepPath(dep: Core.ShareMgnt.ncTDepartmentInfo): Promise<string> {
        if (dep.is_root) {
            return dep.name
        } else {
            const [{ parentPath }] = await getDepParentPathById([dep.id])

            return parentPath ? `${parentPath}/${dep.name}` : dep.name
        }
    }

    /**
     * 更改存储位置
     */
    private async changeStorage(dep: Core.ShareMgnt.ncTDepartmentInfo, ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo): Promise<boolean> {
        let storageEnable = true

        storageEnable = await this.editDepOSS(dep, ossInfo)

        const subUserCount = await getDepartmentOfUsersCount([dep.id])

        if (subUserCount && storageEnable) {
            const limit = 200
            let start = 0

            while (start < subUserCount && storageEnable) {
                const subUsers = await getDepartmentUser([dep.id, start, limit])

                for (let user of subUsers) {
                    storageEnable = await this.editUserOssInfo(user, ossInfo)

                    if (!storageEnable) {
                        break
                    }
                }

                start += limit
            }
        }

        const subDeps = await getSubDepartments([dep.id])

        if (subDeps.length && storageEnable) {

            for (let sdep of subDeps) {
                if (sdep.subUserCount || sdep.subDepartmentCount) {
                    await this.changeStorage(sdep, ossInfo)
                } else {
                    storageEnable = await this.editDepOSS(sdep, ossInfo)

                    if (!storageEnable) {
                        break
                    }
                }
            }
        }

        if (!storageEnable) {
            await Message2.info({ message: __('目标部门的存储位置已不可用，无法替换，仍保留部门当前的存储位置。') })

            return false
        }

        return true
    }

    /**
     * 编辑用户的存储位置
     */
    private async editUserOssInfo(userInfo: Core.ShareMgnt.ncTUsrmGetUserInfo, ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo): Promise<boolean> {
        try {
            await editUserOss(userInfo.id, ossInfo.ossId || '')

            await manageLog(
                ManagementOps.SET,
                __('编辑用户 "${displayName}(${loginName})"的存储位置 成功', { displayName: userInfo.user.displayName, loginName: userInfo.user.loginName }),
                __('存储位置 “${storage}”', { storage: displayUserOssInfo(ossInfo) }),
                Level.INFO,
            )

            return true
        } catch (ex) {
            if (
                ex &&
                ex.error &&
                ex.error.errID &&
                (ex.error.errID === ErrorCode.OSSNotExist || ex.error.errID === ErrorCode.OSSDisabled)
            ) {
                return false
            }

            return true
        }
    }

    /**
     * 编辑部门的存储位置
     */
    private async editDepOSS(dep: Core.ShareMgnt.ncTDepartmentInfo, ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo): Promise<boolean> {
        try {
            await editDepartOSS([dep.id, ossInfo.ossId || ''])

            await manageLog(
                ManagementOps.SET,
                __('编辑部门 “${depName}”的存储位置 成功', { depName: dep.name }),
                __('存储位置 “${storage}”', { storage: displayUserOssInfo(ossInfo) }),
                Level.INFO,
            )

            return true
        } catch (ex) {
            if (
                ex &&
                ex.error &&
                ex.error.errID &&
                (ex.error.errID === ErrorCode.OSSNotExist || ex.error.errID === ErrorCode.OSSDisabled)
            ) {
                return false
            }

            return true
        }
    }

    /**
     * 获取错误提示
     */
    private getErrorMsg(ex: any, srcDep: Core.ShareMgnt.ncTDepartmentInfo, targetDep: Core.ShareMgnt.ncTDepartmentInfo): string {
        if (ex && ex.error && ex.error.errID) {
            switch (ex.error.errID) {
                case ErrorCode.DepNotExist:
                    return __('无法移动“${depName}”，此部门已不存在。', { depName: srcDep.name })

                case ErrorCode.TargetDepNotExist:
                    return __('无法移动“${depName}”，您选中的目标部门“${targetDepName}”已不存在，请重新选择。', { depName: srcDep.name, targetDepName: targetDep.name })

                case ErrorCode.TargetDepIncludeSameNameDep:
                    return __('无法移动“${depName}”，您选中的目标部门下存在与待移动部门同名的子部门，请重新选择或修改部门名称。', { depName: srcDep.name })

                default:
                    return __('移动部门“${depName}”失败。错误原因：${messge}', { depName: srcDep.name, messge: ex.error.errMsg })
            }
        }
    }
}