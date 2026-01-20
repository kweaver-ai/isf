import * as React from 'react';
import WebComponent from '../../webcomponent';
import { Selection } from '../helper';
import { AnonymousPickProps, AnonymousPickState, DefaultSelection } from './type';

export default class AnonymousPickBase extends WebComponent<AnonymousPickProps, AnonymousPickState> {
    state: AnonymousPickState = {
        checkStatus: false,
        selections: [],
    };

    /**
     * 勾选匿名用户复选框
     */
    protected checkAnonymous = (checkStatus: boolean): void => {
        this.setState({
            checkStatus,
            selections: checkStatus ? [DefaultSelection] : [],
        })
    }

    /**
     * 无复选框时，点击选择匿名用户
     */
    protected selectAnonymous = (): void => {
        this.setState({
            selections: [DefaultSelection],
        }, () => {
            this.props.onRequsetSelection(this.state.selections);
        })
    }

    /**
     * 有复选框时，获取已选项
     */
    public getSelections = (): ReadonlyArray<Selection> => {
        const { checkStatus, selections } = this.state;

        return checkStatus ? selections : [];
    }

    /**
     * 有复选框时，清除已选项勾选状态
     */
    public cancelSelections = () => {
        this.setState({
            checkStatus: false,
            selections: [],
        })
    }
}