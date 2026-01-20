
import { noop } from 'lodash';
import { setUserExpireTime, setUserStatus } from '@/core/thrift/sharemgnt/sharemgnt';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import session from '@/util/session/index';
import WebComponent from '../webcomponent';
import { Range, getSeletedUsers } from '../helper';
import __ from './locale';

export enum Status {
    /**
     * 没有任何弹窗和提示状态
     */
    None,

    /**
     * 有设置有效期限弹窗状态
     */
    Normal,

    /**
     * 转圈圈组件出现，正在启用用户
     */
    Loading,

    /**
     * 错误弹窗提示
     */
    Error,

    /**
     * 当前用户
     */
    CurrentUser,
}

export default class SetUserExpireTimeBase extends WebComponent<Console.SetUserExpireTime.Props, Console.SetUserExpireTime.State> {
    static defaultProps = {
        users: [],
        dep: [],
        userid: '',
        shouldEnableUsers: false,
        onCancel: noop,
        onSuccess: noop,
    }

    state = {
        selected: (this.props.users && this.props.users.length) ? Range.USERS : Range.DEPARTMENT,
        expireTime: (this.props.users && this.props.users.length === 1 && this.props.users[0].user.expireTime !== -1) ? (this.props.users[0].user.expireTime * 1000 * 1000) : -1,
        errors: [],
        status: Status.Normal,
        invalidExpireTime: false,
    }

    /**
     * 选择要被设置有效期限的对象
     * @param value 部门、部门及其子部门、所选中的用户
     */
    protected onSelectedType(value) {
        if(value !== Range.USERS){
            this.setState({
                expireTime: -1,
            })
        }else{
            const { users } = this.props;
            this.setState({
                expireTime: (users && users.length === 1 && users[0].user.expireTime !== -1) ? (users[0].user.expireTime * 1000 * 1000) : -1,
            })
        }
        this.setState({
            selected: value,
        })
    }

    /**
     * 更改日期组件的时间时触发事件
     * @param expireTime 从 1970-01-01 开始算，到截止日期之间的时间，单位：毫秒
     */
    protected changeExpireTime(expireTime: number): void {
        this.setState({
            expireTime: expireTime,
        })
    }

    /**
     * 设置用户有效期限点击【确定】按钮时触发事件
     */
    protected async confirmSetUserExpireTime() {
        let errors: Array<any> = []

        this.setState({
            status: Status.Loading,
        })

        const users = await getSeletedUsers(this.state.selected, this.props.dep, this.props.users)

        let invalidExpireTime = false;

        // 不可为当前登录用户设置有效期
        if (!users.some((user) => user.id === session.get('isf.userid'))) {
            for (const user of users) {
                // this.state.expireTime 单位为微妙，setUserExpireTime 接口需要参数单位为秒
                const expireTimeTemp: number = this.state.expireTime === -1 ? -1 : this.state.expireTime / 1000 / 1000

                // 0代表当前是启用状态 newStatus为将要设置的用户状态 false代表禁用
                const newStatus: boolean = user.status === 0 ? false : true

                try {
                    // 用户只是到期被自动禁用，在设置有效期之后需要重新启用用户
                    if (this.props.shouldEnableUsers) {
                        // 在调用 setUserExpireTime 设置用户有效期限之后，才能调用 setUserStatus 接口启用用户
                        await setUserExpireTime([user.id, expireTimeTemp])
                        await setUserStatus(user.id, newStatus)

                        await this.logEnabled(user)
                        await this.logSetExpireTime(user)
                    } else {
                        await setUserExpireTime([user.id, expireTimeTemp])
                        await this.logSetExpireTime(user)
                    }

                } catch ({ error }) {
                    if (error) {
                        // 若用户不存在，在批量设置有效期时不做任何处理
                        if (error.errID !== 20110) {
                            errors = [
                                ...errors,
                                error,
                            ]
                            if (error && error.errID === 20005) {
                                // 如果有效期不合法，则批量中所有用户的日期都不合法，所以直接 break
                                invalidExpireTime = true
                                break
                            }
                        }
                    }
                }
            }

            if (!errors.length) {
                this.setState({
                    status: Status.None,
                })

                this.props.onSuccess(this.state.selected)
            } else {
                this.setState({
                    status: invalidExpireTime ? Status.Normal : Status.Error,
                    errors,
                    invalidExpireTime,
                })
            }
        } else {
            this.setState({
                status: Status.CurrentUser,
            })
        }
    }

    /**
     * 记录设置用户有效期限日志
     * @param user 设置有效期限的用户
     */
    private logSetExpireTime(user) {
        let newDate = new Date()
        let expireTimeDate = __('永久有效')

        if (this.state.expireTime !== -1) {
            newDate.setTime(this.state.expireTime / 1000)
            expireTimeDate = newDate.toLocaleDateString()
        }

        return manageLog(ManagementOps.SET,
            __('设置用户账号 “${displayName}(${loginName})” 有效期限成功', {
                displayName: user.user.displayName,
                loginName: user.user.loginName,
            }),
            __('有效期限 “${expireTime}”', {
                expireTime: expireTimeDate,
            }),
            Level.INFO,
        )
    }

    /**
     * 启用用户日志
     * @param user 用户
     */
    private logEnabled(user) {
        return manageLog(ManagementOps.SET,
            __('启用 用户 “${displayName}(${loginName})” 成功', {
                displayName: user.user.displayName,
                loginName: user.user.loginName,
            }),
            null,
            Level.INFO,
        )
    }
}