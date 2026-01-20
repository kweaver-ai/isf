import { noop } from 'lodash';
import { ListTipStatus } from '../ListTipComponent/helper';
import WebComponent from '../webcomponent';
import { UserGroup } from './helper';

interface UserGroupState {
    /**
     * 选中的用户组
     */
    selectedGroup: UserGroup.GroupInfo;

    /**
     * 用户组列表状态
     */
    groupStatus: ListTipStatus;
}

export default class UserGroupBase extends WebComponent<any, UserGroupState> {
    static defaultProps = {

    }

    state = {
        selectedGroup: null,
        groupStatus: ListTipStatus.Loading,
    }

    /**
     * 用户组列表的ref
     */
    groupGrid = {
        handleGroupNotExist: noop,
    }

    /**
     * 选中的用户组
     */
    protected selectGroup = (selectedGroup: UserGroup.GroupInfo): void => {
        this.setState({
            selectedGroup,
        })
    }

    /**
     * 用户组列表状态变化
     */
    protected changeGroupStatus = (groupStatus: ListTipStatus): void => {
        this.setState({
            groupStatus,
        })
    }

    /**
     * 用户组不存在时刷新用户组
     */
    protected updateGroup = (): void => {
        this.groupGrid.handleGroupNotExist(this.state.selectedGroup)
    }
}