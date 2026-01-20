import * as React from 'react';
import { noop, keys, trim } from 'lodash';
import { Message2 } from '@/sweet-ui'
import { isLoginName } from '@/util/validators';
import { UserManagementErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { createUserGroup, editUserGroup } from '@/core/apis/console/usergroup';
import WebComponent from '../../../webcomponent';
import { UserGroup } from '../../helper';
import { UserGroupType, ValidateState } from './helper';
import __ from './locale';

interface SetGroupProps extends React.Props<void> {
    /**
     * 编辑的用户组信息
     */
    editGroup: UserGroup.GroupInfo;

    /**
     * 取消新建/编辑
     */
    onRequestCancel: (group?: UserGroup.GroupInfo) => void;

    /**
     * 新建/编辑成功
     */
    onRequestSuccess: (group: UserGroup.GroupInfo) => void;
}

interface SetGroupState {
    /**
     * 用户组名称
     */
    name: string;

    /**
     * 用户组名称状态
     */
    nameStatus: ValidateState;

    /**
     * 备注
     */
    notes: string;

    /**
     * 用户组源
     */
    userGroupSources: ReadonlyArray<UserGroupType>;

    /**
     * 是否打开用户组选择弹框
     */
    isShowSelectGroup: boolean;
}

export default class SetGroupBase extends WebComponent<SetGroupProps, SetGroupState> {
    static defaultProps = {
        editGroup: {
            id: '',
            name: '',
            notes: '',
        },
        onRequestCancel: noop,
        onRequestSuccess: noop,
    }

    state = {
        name: '',
        nameStatus: ValidateState.Normal,
        notes: '',
        userGroupSources: [],
        isShowSelectGroup: false,
    }

    /**
     * 是否正在请求中
     */
    isRequesting = false

    /**
     * 是否新建
     */
    isCreate = true

    componentDidMount() {
        const { editGroup } = this.props

        if (editGroup && editGroup.id) {
            this.isCreate = false
            const { name, notes = '' } = editGroup

            this.setState({
                name,
                notes,
            }, () => {
                this.verifyName()
            })
        }
    }

    /**
     * 修改名称
     */
    protected changeName = (name: string): void => {
        this.setState({
            name,
            nameStatus: ValidateState.Normal,
        })
    }

    /**
     * 修改备注
     */
    protected changeNote = (notes: string): void => {
        this.setState({
            notes,
        })
    }

    /**
     * 修改选择的用户组
     */
    protected changeSelectGroups = (value: ReadonlyArray<UserGroupType>): void => {
        this.setState({
            userGroupSources: value,
        })
    }

    /**
     * 修改选择用户组弹框状态
     */
    protected changeSelectGroupsState = (state: boolean): void => {
        this.setState({
            isShowSelectGroup: state,
        })
    }

    /**
     * 确定新建/编辑
     */
    protected confirm = async (): Promise<void> => {
        if (this.verifyName() && !this.isRequesting) {
            try {
                this.isRequesting = true

                const { name, notes } = this.state

                let param = {
                    name: trim(name),
                    notes,
                    ...(
                        this.isCreate
                            ? {
                                group_ids_of_members: this.state.userGroupSources.map(({ id }) => id),
                            }
                            : {}
                    ),
                }

                let id = this.props.editGroup && this.props.editGroup.id || ''

                // 编辑
                if (id) {
                    await editUserGroup({ id, fields: keys(param).join(','), ...param })
                }
                // 新建
                else {
                    id = (await createUserGroup(param)).id
                }

                this.isRequesting = false

                this.props.onRequestSuccess({ ...param, id })
            } catch (ex) {
                if (ex && ex.code) {
                    switch (ex.code) {
                        case UserManagementErrorCode.GroupConflict:
                            this.setState({
                                nameStatus: ValidateState.NameConfilct,
                            })
                            break

                        case UserManagementErrorCode.GroupNotFound:
                            this.props.onRequestCancel(this.props.editGroup)
                            break

                        default:
                            ex.description && Message2.info({ message: ex.description })
                    }
                }

                this.isRequesting = false
            }
        }
    }

    /**
     * 校验名称
     */
    protected verifyName = (): boolean => {
        const name = trim(this.state.name)

        if (name && isLoginName(name)) {
            this.setState({
                nameStatus: ValidateState.Normal,
            })

            return true
        } else {
            this.setState({
                nameStatus:
                    !name ?
                        ValidateState.Empty
                        : !isLoginName(name) ?
                            ValidateState.InvalidName
                            : ValidateState.Normal,
            })

            return false
        }
    }
}