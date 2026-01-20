import * as React from 'react'
import { noop } from 'lodash';
import { variousIdCard } from '@/util/validators'
import WebComponent from '../webcomponent'
import __ from './locale';

enum IDNumberStatusItems {
    /**
     * 只读
     */
    OnlyRead = 0,

    /**
     * 编辑
     */
    AllowEidt = 1,
}

export enum ValidateState {
    /**
     * 正常
     */
    Normal = 0,

    /**
     * 身份证错误
     */
    Error = 1,
}

export default class IDNumberBoxBase extends WebComponent<Console.IDnumberBox.Props, Console.IDNumberBox.State> {

    static defaultProps = {
        width: 200,
        idcardNumber: '',
        isClickBtn: true,
        isIDNumEdit: false,
        defaultidcardNumber: '',
        onChange: noop,
    }

    state = {
        status: IDNumberStatusItems.AllowEidt,
        IDCardStatus: ValidateState.Normal,
        IDNumber: '',
        showIDNumber: '',
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (nextProps.isClickBtn) {
            if (prevState.status === IDNumberStatusItems.AllowEidt) {
                return {
                    IDCardStatus: ValidateState.Normal,
                    IDNumber: nextProps.defaultidcardNumber,
                    showIDNumber: nextProps.defaultidcardNumber,
                }
            } else {
                if (prevState.showIDNumber && nextProps.isIDNumEdit) {
                    return {
                        showIDNumber: nextProps.idcardNumber,
                        IDCardStatus: nextProps.defaultidcardNumber === prevState.showIDNumber ? ValidateState.Normal : variousIdCard(prevState.showIDNumber) ? ValidateState.Normal : ValidateState.Error,
                    }
                } else {
                    return {
                        showIDNumber: nextProps.idcardNumber,
                        IDCardStatus: ValidateState.Normal,
                    }
                }
            }
        }
        return null
    }

    componentDidMount() {
        this.setState({
            IDNumber: this.props.defaultidcardNumber,
            showIDNumber: this.props.idcardNumber,
        })
    }

    /**
     * 转换编辑状态
     * @param status 状态
     */
    protected toggleEditStatus(status) {
        if (status === IDNumberStatusItems.AllowEidt) {
            this.setState({
                status: IDNumberStatusItems.OnlyRead,
                showIDNumber: '',
                IDCardStatus: ValidateState.Normal,
            })
            this.props.onChange('')
        } else {
            this.setState({
                status: IDNumberStatusItems.AllowEidt,
                showIDNumber: this.state.IDNumber,
                IDCardStatus: ValidateState.Normal,
            })
            this.props.onChange(null)
        }
    }

    /**
     * 输入框值改变时出发
     * @param idcard 新值
     */
    protected handleChange(idcard: string) {
        this.setState({
            showIDNumber: idcard,
        })
        this.props.onChange(idcard)
    }
}
