import { noop } from 'lodash';
import WebComponent from '../webcomponent';

interface Props {
    users?: Array<any>; // 选择的用户 * any 后续补充

    dep: any; // 选择的部门 * any 后续补充

    userid: string; // 当前登录的用户

    onComplete: () => any; // 设置角色结束的事件
}

export default class SetRoleBase extends WebComponent<Props, any> {
    static defaultProps ={
        users: [],
        dep: null,
        userid: '',
        onComplete: noop,
    }
}