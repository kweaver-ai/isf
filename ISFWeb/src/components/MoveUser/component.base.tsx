import { isFunction, noop } from 'lodash';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { moveUserToDepartment, editUserOss, getOrgDepartmentById, getDepParentPathById } from '@/core/thrift/sharemgnt/sharemgnt';
import { displayUserOssInfo } from '@/core/oss/oss';
import { getObjectStorageInfoById } from '@/core/apis/console/ossgateway'
import { Message } from '@/sweet-ui';
import WebComponent from '../webcomponent';
import { Range, getSeletedUsers } from '../helper';
import __ from './locale';

export enum Status {
    NORMAL,

    CURRENT_USER_INCLUDED, // 当前用户

    LOADING, // 加载中

    DESTDEPARTMENT_NOT_EXIST = 20211, // 目标部门不存在

    SRRDEP_NOT_EXIST = 20210,   // 源部门不存在
}
interface TreeNodeStatus {
    disabled: boolean;
}

interface DepInfo {
    id: string;
    name: string;
    parentPath: string;
}
interface Props {
    users?: Array<any>; // 选择的用户 * any 后续补充

    dep: DepInfo; // 选择的部门 * any 后续补充

    userid: string; // 当前登录的用户

    onComplete: (dep?: Core.ShareMgnt.ncTDepartmentInfo) => any; // 移动结束的事件

    onSuccess: (range: Range) => {}; // 移动成功
}

interface State {
    selected: Range; // 选择移动的对象

    status: Status; // 移动的状态

    selectedDep: any; // 选中(移入)的部门

    dep: DepInfo; // 选择的部门
}

export default class MoveUserBase extends WebComponent<Props, State> {
    static defaultProps = {
        users: [],

        dep: null,

        userid: '',

        onComplete: noop,

        onSuccess: noop,
    }

    state = {
        selected: Range.USERS,

        status: Status.NORMAL,

        selectedDep: null,

        dep: this.props.dep,
    }

    users = []

    // 是否仅选择单个用户
    singleUser = this.props.users && this.props.users.length === 1 && this.state.selected === Range.USERS

    dialogRef = null

    organizationTreeRef = null

    componentDidMount() {
        this.setState({
            selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        })

        // 对话框组件加载完毕，更新state，对话框宽高变化，导致再操作时对话框会重新居中一下
        if (this.dialogRef && !this.dialogRef.moved) {
            this.dialogRef.moved = true
        }
    }

    /**
     * 检查是否存在不能被删除的用户
     */
    checkUser(users: Array<any>): boolean {
        return !users.some((value, index) => {
            if (value.id === this.props.userid) {
                this.setState({
                    status: Status.CURRENT_USER_INCLUDED,
                })
                return true
            }
        })
    }

    /**
     * 选择目标部门
     * @param value  目标部门
     */
    protected async selectDep(value) {
        this.setState({
            selectedDep: {
                ...value,
            },
        })
    }

    /**
     * 移动用户
     * @param changeOss 是否要改变存储位置
     */
    private async moverUser(changeOss, users) {
        const { selectedDep } = this.state;
        let editUserOssError = false;
        this.setState({
            status: Status.LOADING,
        })

        for (let user of users) {
            let ids: ReadonlyArray<string> = []

            // 是否仅移动单个用户
            if (this.singleUser) {
                ids = [this.state.dep.id]
            } else {
                const { user: userInfo, directDeptInfo } = user
                // 判断有无直属部门
                if (directDeptInfo) {
                    ids = [directDeptInfo.departmentId]
                } else {
                    if (userInfo && Array.isArray(userInfo.departmentIds)) {
                        const { departmentIds, departmentNames } = userInfo
                        // 获取当前默认部门信息
                        const parentPathList = await getDepParentPathById([...departmentIds, this.props.dep.id])
                        // 当前(进入)部门的路径
                        const depPath = `${parentPathList[parentPathList.length - 1].parentPath ? parentPathList[parentPathList.length - 1].parentPath + '/' : ''}${this.props.dep.name}`

                        departmentIds.forEach((id, index) => {
                            const parent = parentPathList[index].parentPath
                            // 拼接用户所属部门路径
                            const path = `${parent ? parent + '/' : ''}${departmentNames[index]}`

                            if (path.startsWith(depPath)) {
                                ids = [...ids, id]
                            }
                        })
                    }
                }
            }

            const promises = ids.map(async (id) => {
                if (changeOss) {
                    let result = await moveUserToDepartment(
                        [user.id],
                        id,
                        selectedDep.id,
                    )
                    if (result && result.length === 0) {
                        await this.logUserMove(user)
                    }
                    await editUserOss(user.id, selectedDep.ossInfo.ossId ? selectedDep.ossInfo.ossId : '')
                    await this.logOssEdit(user);

                } else {
                    let result = await moveUserToDepartment(
                        [user.id],
                        id,
                        selectedDep.id,
                    )
                    if (result && result.length === 0) {
                        await this.logUserMove(user)
                    }
                }
            })

            await Promise.all(promises).catch((ex) => {
                if (ex.error.errID === 24405) {
                    if (!editUserOssError) {
                        editUserOssError = true;
                        Message.info({ message: __('替换存储位置失败。\n目标部门的存储位置已不可用。') })
                    }
                } else if (ex.error.errID !== 20110) {
                    this.setState({
                        status: ex.error.errID,
                    })
                    throw ex
                }
            })
        }
    }

    /**
     * 点击确定
     */
    protected async confirmMoveUsers() {
        this.users = await getSeletedUsers(this.state.selected, this.props.dep, this.props.users);
        if (this.checkUser(this.users)) {
            // 获取目标部门存储信息
            try {
                const { selectedDep: { id } } = this.state
                let { ossInfo } = await getOrgDepartmentById(id)

                if (ossInfo && ossInfo.ossId) {
                    const ossData = await getObjectStorageInfoById(ossInfo.ossId)
                    ossInfo = {
                        enabled: ossData.enabled,
                        ossId: ossData.id,
                        ossName: ossData.name,
                    }
                }
                this.setState({
                    selectedDep: {
                        ...this.state.selectedDep,
                        ossInfo,
                    },
                })

                await this.moverUser(false, this.users)
                this.props.onSuccess(this.state.selected)
                
            } catch (error) {
                // 点击确定，获取部门/组织的存储信息，若返回20215不存在，则弹窗提示（notExist 将该部门标记为不存咋）
                if (error && error.error && error.error.errID === 20215) {
                    this.setState({
                        selectedDep: {
                            ...this.state.selectedDep,
                            notExist: true,
                        },
                    })

                }
            }
        }
    }

    /**
     * 部门/组织不存在
     */
    protected depNotExist = (): void => {
        // 获取部门/组织存储信息时  部门/组织不存在，则将已选部门置空，并刷新组织结构树
        this.organizationTreeRef && isFunction(this.organizationTreeRef.getOrganization) && this.organizationTreeRef.getOrganization()

        this.setState({ selectedDep: null })
    }

    /**
     * 确定改变对象存储
     */
    protected async confirmChangeOss() {
        try {
            await this.moverUser(true, this.users)
            this.props.onSuccess(this.state.selected)

        } catch { }
    }

    /**
     * 取消改变对象存储
     */
    protected async cancelChangeOss() {
        await this.moverUser(false, this.users)
        this.props.onSuccess(this.state.selected)
    }

    /**
     * 禁用当前部门
     * @param node 当前部门
     */
    getDepartmentStatus(node): TreeNodeStatus {
        if (this.state.dep.id === node.id) {
            return {
                disabled: true,
            }
        } else {
            return {
                disabled: false,
            }
        }
    }

    /**
     * 记录移动用户日志
     * @param user  当前移动的日志
     */
    logUserMove(user) {
        return manageLog(ManagementOps.MOVE,
            __('移动用户“${username}(${loginName})”至部门“${orgname}”成功', {
                username: user.user.displayName,
                loginName: user.user.loginName,
                orgname: this.state.selectedDep.name,
            }),
            __('原部门：${originDepName}，新部门：${newDepName}', {
                originDepName:user.directDeptInfo.departmentName,
                newDepName: this.state.selectedDep.name,
            }),
            Level.INFO,
        )
    }

    logOssEdit(user) {
        return manageLog(
            ManagementOps.SET,
            __('编辑用户 "${displayName}(${loginName})" 成功', { displayName: user.user.displayName, loginName: user.user.loginName }),
            __('存储位置 “${ossName}”', { ossName: displayUserOssInfo(this.state.selectedDep.ossInfo) }),
            Level.INFO)
    }

    /**
    * 选择要删除的对象
    * @param value 部门、部门及其子部门、所选中的
    */
    onSelectedType(value) {
        this.singleUser = value === Range.USERS && this.props.users && this.props.users.length === 1

        this.setState({
            selected: value,
            dep: this.props.dep,
        })
    }

    /**
     * 操作单个用户时，选择用户的部门
     */
    protected onDepChange = (dep): void => {
        this.setState({ dep })
    }
}