import { noop } from 'lodash';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { setUserStatus } from '@/core/thrift/sharemgnt/sharemgnt';
import WebComponent from '../webcomponent';
import { Range, getSeletedUsers } from '../helper';
import { Status } from './helper';
import __ from './locale';

interface Props {
    users?: Array<any>; // 选择的用户 * any 后续补充

    dep: any; // 选择的部门 * any 后续补充

    userid: string; // 当前登录的用户

    onComplete: () => any; // 启用结束的事件

    onSuccess: (range: Range) => {}; // 启用成功
}

interface State {
    selected: Range; // 选择启用的对象

    status: Status; // 启用的状态

    /**
     * 过期用户
     */
    expireTimeUsers: Array<any>;
}

export default class RemoveUserBase extends WebComponent<Props, any> {
    static defaultProps = {
        users: [],
        dep: null,
        userid: '',
        onComplete: noop,
        onSuccess: noop,
    }

    state = {
        selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        status: Status.NORMAL,
        errorUser: null,
        expireTimeUsers: [],
    }

    componentDidMount() {
        this.setState({
            selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        })
    }

    /**
     * 检查是否存在不能被移除的用户
     */
    checkUser(users: Array<any>): boolean {
        return !users.some((value, index) => {
            return value.id === this.props.userid
        })
    }

    /**
     * 启用用户
     */
    async enableUsers(users: Array<any>) {
        let expireTimeUsersTemp: ReadonlyArray<any> = []
        for (let user of users) {
            try {
                await setUserStatus(user.id, true)
                await this.logEnabled(user);
            }
            catch (ex) {
                if (ex.error.errID === 20157) {
                    expireTimeUsersTemp = [...expireTimeUsersTemp, user]
                }
                else if (ex.error.errID !== 20110 && ex.error.errID !== 20157) {
                    this.setState({
                        status: ex.error.errID,
                        errorUser: user,
                    })
                    throw (ex);
                }
            }

        }

        expireTimeUsersTemp.length ?
            this.setState({
                expireTimeUsers: expireTimeUsersTemp,
            })
            :
            this.props.onSuccess(this.props.users)
    }

    /**
     * 确定事件
     */
    async confirmEnableUsers() {
        this.setState({
            status: Status.LOADING,
        })
        try {
            const users = await getSeletedUsers(this.state.selected, this.props.dep, this.props.users);
            if (this.checkUser(users)) {
                await this.enableUsers(users);
            } else {
                this.setState({
                    status: Status.CURRENT_USER_INCLUDED,
                })
            }
        } catch (ex) {
            this.setState({
                status: ex.error.errID,
            })
        }

    }

    /**
     * 启用用户日志
     * @param user 用户
     */
    logEnabled(user) {
        return manageLog(ManagementOps.SET,
            __('启用 用户 “${displayName}(${loginName})” 成功', {
                displayName: user.user.displayName,
                loginName: user.user.loginName,
            }),
            null,
            Level.INFO,
        )
    }

    /**
     * 选择要启用的对象
     * @param value 部门、部门及其子部门、所选中的
     */
    onSelectedType(value) {
        this.setState({
            selected: value,
        })
    }

    /**
     * 【提示】弹窗点击确定后，关闭原【提示】弹窗&&出现【设置用户有效期限】弹窗
     */
    protected completeExpireTimeTips() {
        this.setState({
            status: Status.SET_EXPIRE_TIME,
        })
    }
}