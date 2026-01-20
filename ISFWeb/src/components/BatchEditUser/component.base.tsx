import * as React from 'react';
import { noop } from 'lodash';
import { editUser } from '@/core/thrift/sharemgnt/sharemgnt';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import session from '@/util/session/index';
import { SystemRoleType } from '@/core/role/role';
import { ValidateState } from '@/core/user';
import { Toast, Message2 } from '@/sweet-ui';
import { Range, getSeletedUsers } from '../helper';
import WebComponent from '../webcomponent';
import { Status, BatchEditUserProps, BatchEditUserState } from './type'
import __ from './locale';
import { getLevelConfig } from '@/core/apis/console/usermanagement';

export default class BatchEditUserBase extends WebComponent<BatchEditUserProps, BatchEditUserState> {
    static defaultProps: BatchEditUserProps = {
        users: [],
        dep: null,
        onRequestCancel: noop,
        onRequestSuccess: noop,
    }

    state: BatchEditUserState = {
        selected: (this.props.users && this.props.users.length) ? Range.USERS : Range.DEPARTMENT,
        expireTime: -1,
        status: Status.Config,
        csfLevel: null,
        csfLevel2: null,
        csfOptions: [],
        csfOptions2: [],
        show_csf_level2: false,
        csfIsChecked: true,
        csf2IsChecked: true,
        expIsChecked: true,
        progress: 0,
        currentUserName: {
            loginName: '',
            displayName: '',
        },
        csfValidateState: ValidateState.Normal,
        csf2ValidateState: ValidateState.Normal,
    }

    isAdmin: boolean // 是否系统管理员
    isSecurit: boolean // 是否安全管理员
    currentUserLength: number // 当前编辑所有用户的长度
    currentLength: number = 0 // 计算当前编辑个数

    componentDidMount() {
        this.isAdmin = session.get('isf.userInfo').user.roles.some((role) => [SystemRoleType.Admin].includes(role.id))
        this.isSecurit = session.get('isf.userInfo').user.roles.some((role) => [SystemRoleType.Securit].includes(role.id))
        this.getUserCsfInfos();
    }

    /**
     * 选择要被设置有效期限的对象
     * @param value 部门、部门及其子部门、所选中的用户
     */
    protected selectedType(selected: Range) {
        this.setState({
            selected,
        })
    }

    /**
     * 更改日期组件的时间时触发事件
     * @param expireTime 从 1970-01-01 开始算，到截止日期之间的时间，单位：毫秒
     */
    protected changeExpireTime = (expireTime: number) => {
        this.setState({
            expireTime,
        })
    }

    /**
     * 设置用户有效期限点击【确定】按钮时触发事件
     */
    protected confirmBatchEditUser = async () => {
        let errors: ReadonlyArray<any> = []

        const { csfIsChecked, csfLevel, csfLevel2, expireTime, expIsChecked, csf2IsChecked, selected, progress } = this.state
        // 判断必选框是否为空
        if ((csfIsChecked && csfLevel !== null) || !csfIsChecked) {

            this.setState({
                status: Status.Progress,
            })

            const users = await getSeletedUsers(this.state.selected, this.props.dep, this.props.users)

            this.currentUserLength = users.length

            // 不可为当前登录用户设置
            if (!users.some((user) => user.id === session.get('isf.userid'))) {
                for (const user of users) {

                    // 如果取消批量编辑，中断设置
                    if(this.state.status === Status.Close){
                        break
                    }
                    // 累加已编辑用户个数
                    this.currentLength ++

                    // this.state.expireTime 单位为微妙，editUser 接口需要参数单位为秒
                    const expireTimeTemp: number = expireTime === -1 ? -1 : expireTime / 1000 / 1000

                    let data: { [key: string]: string | number } = { id: user.id }
                    if(csfIsChecked) {
                        data = {
                            ...data,
                            csfLevel,
                        }
                    }
                    if(csf2IsChecked) {
                        data = {
                            ...data,
                            csfLevel2,
                        }
                    }
                    if(expIsChecked) {
                        data = {
                            ...data,
                            expireTime: expireTimeTemp,
                        }
                    }

                    const { user: { loginName, displayName } } = user
                    //  更新当前正在设置的用户名
                    this.setState({
                        currentUserName: {
                            loginName,
                            displayName,
                        },
                        progress: this.state.progress + 1,
                    })

                    try {
                        // 设置用户
                        await editUser([{ ncTEditUserParam: data }, user.id])
                        await this.logBatchEditUser(user)
                    } catch ({ error }) {
                        if (error) {
                            errors = [
                                ...errors,
                                {
                                    loginName: loginName,
                                    displayName: displayName,
                                    ...error,
                                },
                            ]
                        }
                    }
                }

                // 当有设置失败的时候，弹框提示
                if (errors.length) {
                    await Message2.info({
                        message: __('无法对以下用户进行设置：'),
                        detail: errors.map(({ loginName, displayName, errCode, errMsg }, index) => (
                            <div
                                key={index}
                            >
                                {`${displayName}(${loginName})： ${errCode === ErrorCode.UserNotExist ? __('用户不存在。') : (errMsg || '')}`}
                            </div>
                        )),
                    })
                    Toast.open(__('操作完成，本次成功编辑${length}个用户', { length: this.currentLength - errors.length }))
                } else {
                    Toast.open(__('操作完成，本次成功编辑${length}个用户', { length: this.currentLength }))
                }
            } else {
                Message2.info({
                    message: __('您无法编辑自身账号。'),
                })

                // 如果包含自身账号，关闭弹窗
                this.props.onRequestCancel()
            }

            this.props.onRequestSuccess(selected)
        } else {
            this.setState({
                csfValidateState: ValidateState.Empty,
            })
        }
    }

    /*
     * 获取用户密级枚举
     */
    private async getUserCsfInfos() {
        const {csf_level_enum, csf_level2_enum, show_csf_level2} = await getLevelConfig({fields: 'csf_level_enum,csf_level2_enum,show_csf_level2'})
        this.setState({
            csfOptions: csf_level_enum,
            csfOptions2: csf_level2_enum,
            show_csf_level2,
        })
    }

    /**
     * 密级切换
     */
    protected updateCsfLevel(type: 'csfLevel' | 'csfLevel2', csfLevel: number): void {
        const validateStateType = type === 'csfLevel' ? 'csfValidateState' : 'csf2ValidateState'
        this.setState({
            [type]: csfLevel,
            [validateStateType]: this.state[validateStateType] === ValidateState.Empty ? ValidateState.Normal : ValidateState.Normal,
        })
    }

    /**
     * 修改用户密级是否必选
     */
    protected updateCsfIsChecked = (type: 'csfIsChecked' | 'csf2IsChecked', csfIsChecked: boolean): void => {
        const validateStateType = type === 'csfIsChecked' ? 'csfValidateState' : 'csf2ValidateState'
        this.setState({
            [type]: csfIsChecked,
            [validateStateType]: ValidateState.Normal,
        })
    }

    /**
     * 修改用户有效期限是否必选
     */
    protected updateExpIsChecked = (expIsChecked: boolean): void => {
        this.setState({
            expIsChecked,
        })
    }

    /**
     * 失焦事件
     */
    protected handleOnBlur = (type: 'csfLevel' | 'csfLevel2') => {
        // 避免失焦判断还未完成时，去勾密级选择，下拉框出现气泡提示
        setTimeout(() => {
            const { csfIsChecked, csfLevel, csf2IsChecked, csfLevel2 } = this.state
            if(type === 'csfLevel') {
                this.setState({
                    csfValidateState: (csfIsChecked && csfLevel === null) ? ValidateState.Empty : ValidateState.Normal,
                })
            } else {
                this.setState({
                    csf2ValidateState: (csf2IsChecked && csfLevel2 === null) ? ValidateState.Empty : ValidateState.Normal,
                })
            }
        }, 100)
    }

    /**
     * 切换设置进度状态
     */
    protected changeStatus = () => {
        this.setState({
            status: this.state.status === Status.Progress ? Status.Confirm : Status.Progress,
        })
    }

    /**
     * 确定取消设置
     */
    protected confirmCancel = () => {
        this.setState({
            status: Status.Close,
        })
    }

    /**
     * 记录批量编辑用户密级、有效期限日志
     * @param user 批量编辑的用户
     */
    private logBatchEditUser(user) {
        const { expireTime, csfLevel, csfLevel2, show_csf_level2, csfOptions, csfOptions2, expIsChecked, csfIsChecked, csf2IsChecked } = this.state
        let newDate = new Date()
        let expireTimeDate = __('永久有效')
        let csfLevelText = ''
        let csfLevel2Text = ''
        const originalCsfLevel = user?.user?.csfLevel
        const originalCsfLevel2 = user?.user?.csfLevel2
        let originalCsfLevelText = ''
        let originalCsfLevel2Text = ''

        if (expireTime !== -1) {
            newDate.setTime(expireTime / 1000)
            expireTimeDate = newDate.toLocaleDateString()
        }

        if (csfLevel !== null) {
            csfOptions.map((csfLevels) => {
                if (csfLevels.value === csfLevel) {
                    csfLevelText = csfLevels.name
                }
                if (csfLevels.value === originalCsfLevel) {
                    originalCsfLevelText = csfLevels.name
                }
            })
        }
        if (csfLevel2 !== null) {
            csfOptions2.map((csfLevels) => {
                if (csfLevels.value === csfLevel2) {
                    csfLevel2Text = csfLevels.name
                }
                if (csfLevels.value === originalCsfLevel2) {
                    originalCsfLevel2Text = csfLevels.name
                }
            })
        }

        const isChangeCsfLevel = originalCsfLevelText !== csfLevelText || originalCsfLevel2Text !== csfLevel2Text;
        const isChangeExpireTime = user.user.expireTime !== expireTime;

        // 辅助函数：生成变更文本
        const getChangeText = (oldText: string, newText: string): string => {
            return oldText === newText ? newText : __('由 ${oldText} 改为 ${newText}', { oldText, newText });
        };

        // 构建变更描述数组
        const changeDescriptions: string[] = [];
        
        // 添加密级信息
        if (csfIsChecked && csfLevel !== null) {
            changeDescriptions.push(__('密级“ ${csfLevelText} ”', { csfLevelText: getChangeText(originalCsfLevelText, csfLevelText) }));
        }
        
        // 添加密级2信息
        if (show_csf_level2 && csf2IsChecked && csfLevel2 !== null) {
            changeDescriptions.push(__('密级2“ ${csfLevel2Text} ”', { csfLevel2Text: getChangeText(originalCsfLevel2Text, csfLevel2Text) }));
        }
        
        // 添加有效期信息
        if (expIsChecked) {
            changeDescriptions.push(__('有效期“ ${expireTime} ”', { expireTime: expireTimeDate }));
        }
        
        // 将变更描述数组拼接为最终文本
        const changeDescription = changeDescriptions.join('，');

        return isChangeCsfLevel || isChangeExpireTime ? 
            manageLog(ManagementOps.SET,
                __('编辑 用户“ ${displayName}(${loginName}) ”成功', {
                    displayName: user.user.displayName,
                    loginName: user.user.loginName,
                }),
                changeDescription,
                Level.INFO,
            ) : null
    }
}
