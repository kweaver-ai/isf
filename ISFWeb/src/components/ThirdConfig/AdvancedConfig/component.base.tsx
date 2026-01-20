import { noop } from 'lodash';
import WebComponent from '../../webcomponent'
import { ConfigItem, transformArrayToObject, isJsonObject } from '../helper'
import __ from './locale';

interface AdvancedConfigProps {
    /**
     * 提示信息（客户端还是服务端参数配置）
     */
    title: string;

    /**
     * 传递过来的初始配置数组
     */
    originalConfig: ReadonlyArray<ConfigItem>;

    /**
     * 关闭或者取消高级配置对话框
     */
    onRequestClose: () => void;

    /**
     * 确认高级配置对话框
     * @param {*} config 已经将非法参数还原的配置项
     */
    onRequestConfirm: (config: string) => void;
}

interface AdvancedConfigState {
    /**
     * string类型的高级配置
     */
    config: string;

    /**
     * 参数格式是否错误，格式错误展示错误提示
     */
    isInvalidFormat: boolean;
}

export default class AdvancedConfigBase extends WebComponent<AdvancedConfigProps, AdvancedConfigState> {
    static defaultProps = {
        title: __('客户端参数配置：'),
        originalConfig: [],
        onRequestClose: noop,
        onRequestConfirm: noop,
    }

    state = {
        config: '',
        isInvalidFormat: false,
    }

    componentDidMount() {
        this.setState({
            config: JSON.stringify(transformArrayToObject(this.props.originalConfig), null, 4),
        })
    }

    /**
     * 高级配置输入框值发生改变
     * @param {*} value 输入的值
     */
    protected handleConfigChange(value: string): void {
        this.setState({
            config: value,
        })
    }

    /**
     * 点击确认的时候触发 - 把格式合法的Json字符串传出去
     */
    protected confirm(): void {
        const { onRequestConfirm } = this.props
        const { config } = this.state

        if (isJsonObject(config)) {
            onRequestConfirm(config)
        } else {
            this.setState({
                isInvalidFormat: true,
            })
        }
    }
}