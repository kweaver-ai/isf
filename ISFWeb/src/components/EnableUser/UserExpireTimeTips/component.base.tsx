import WebComponent from '../../webcomponent';
import { noop } from 'lodash';

export default class UserExpireTimeTipsBase extends WebComponent<Console.UserExpireTimeTips.Props, Console.UserExpireTimeTips.State> {

    static defaultProps = {
        expirtTimeUsers: [],
        completeExpireTimeTips: noop,
        cancelExpireTimeTips: noop,
        userid: '',
    }
}