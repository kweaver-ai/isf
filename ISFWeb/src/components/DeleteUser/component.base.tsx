import { noop } from 'lodash';
import { getDepartmentOfUsers, delUser } from '@/core/thrift/user/user';
import { manageLog, Level, ManagementOps } from '@/core/log';
import WebComponent from '../webcomponent';
import { listUsersSince } from '../helper';
import __ from './locale';

export enum Range {
    DEPARTMENT_DEEP, // 部门及其子部门

    DEPARTMENT, // 当前部门

    USERS, // 当前选中的用户
}

export enum Status {
    NORMAL,

    CURRENT_USER_INCLUDED, // 当前用户

    HAS_UNCLOSED_DOC, // 个人文档未关闭

    LOADING, // 加载中
}

interface Props {
    users?: Array<any>; // 选择的用户 * any 后续补充

    onComplete: () => any; // 删除结束的事件

    dep: any; // 选择的部门 * any 后续补充

    userid: string; // 当前登录的用户

    onSuccess: (range: Range) => {}; // 删除成功
}

interface State {
    selected: Range; // 选择删除的对象

    status: Status; // 删除的状态
}

export default class DeleteUserBase extends WebComponent<any, any> {
    static defaultProps = {
        users: [],
        onComplete: noop,
        dep: null,
        userid: '',
        onSuccess: noop,
    }

    state = {
        selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        status: Status.NORMAL,
    }

    componentDidMount() {
        this.setState({
            selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        })
    }

    /**
     * 获取选中的用户
     */
    async getSeletedUsers() {
        this.setState({
            status: Status.LOADING,
        })
        if (this.state.selected === Range.USERS) {
            return [...this.props.users]
        } else if (this.state.selected === Range.DEPARTMENT) {
            return await getDepartmentOfUsers(this.props.dep.id, 0, -1)
        } else if (this.state.selected === Range.DEPARTMENT_DEEP) {
            return await listUsersSince(this.props.dep);
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
     * 删除用户
     */
    async deleteUsers(users: Array<any>) {
        for (let user of users) {
            try {
                await delUser(user.id)
            }
            catch (ex) {
                if (ex.error.errID !== 20110) {
                    this.setState({
                        status: ex.error.errID,
                    })
                    throw (ex);
                }
            }

            await this.setLog(user)
        }
    }

    /**
     * 确定事件（不存在无法删除的用户，删除即可）
     */
    async confirmDeleteUsers() {
        const users = await this.getSeletedUsers();
        if (this.checkUser(users)) {
            await this.deleteUsers(users);
            this.props.onSuccess(this.state.selected);
        }
    }

    async setLog(user) {
        await manageLog(ManagementOps.DELETE,
            __('删除用户 “${displayName}(${loginName})” 成功', {
                displayName: user.user.displayName,
                loginName: user.user.loginName,
            }),
            null,
            Level.WARN,
        )
    }

    /**
     * 选择要删除的对象
     * @param value 部门、部门及其子部门、所选中的
     */
    onSelectedType(value) {
        this.setState({
            selected: value,
        })
    }
}