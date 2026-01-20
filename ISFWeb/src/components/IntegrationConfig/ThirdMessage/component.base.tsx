import { union } from 'lodash';
import { getThirdMessage } from '@/core/apis/console/thirdMessage';
import { Message } from '@/sweet-ui';
import WebComponent from '../../webcomponent';
import { ConfigInfo, ValidateStatus } from './helper';

interface ThirdMessageProps {
    /**
     * swf文件路径
     */
    swf: string;
}

interface ThirdMessageState {
    /**
     * 第三方消息信息
     */
    thirdPartyConfig: ReadonlyArray<ConfigInfo>;

    /**
     * 处于编辑状态的卡片（未保存过的卡片不计入其中），当有处于编辑中的卡片时，不允许添加新的卡片
     */
    editedCards: ReadonlyArray<number>;
}

export default class ThirdMessageBase extends WebComponent<ThirdMessageProps, ThirdMessageState> {
    static defaultProps = {
        swf: '',
    }

    state = {
        thirdPartyConfig: [],
        editedCards: [],
    }

    async componentDidMount() {
        try {
            const thirdPartyConfig = await getThirdMessage();

            if (thirdPartyConfig.length === 0) {
                this.setState({
                    thirdPartyConfig: [
                        {
                            indexId: null,
                            thirdPartyName: '',
                            internalConfig: '',
                            pluginClassName: '',
                            messages: [],
                            enabled: false,
                            plugin: null,
                        },
                    ],
                })
            } else {
                this.setState({
                    thirdPartyConfig: thirdPartyConfig.map(({ id, thirdparty_name, class_name, channels = [], config, enabled, filename }) => {

                        return {
                            indexId: id,
                            thirdPartyName: thirdparty_name,
                            internalConfig: JSON.stringify({
                                ...config,
                            }),
                            pluginClassName: class_name,
                            messages: channels.map((configValue) => {
                                return {
                                    configValue,
                                    valueValidateStatus: ValidateStatus.Normal,
                                }
                            }),
                            enabled,
                            plugin: {
                                filename,
                            },
                        }
                    }),
                })
            }
        } catch ({ message }) {
            if (message) {
                Message.alert({ message });
            }
        }
    }

    /**
     * 响应子组件删除卡片 -- 配置尚未保存过 indexId===null
     */
    protected handleDeleteUnSavedCard(): void {
        const { thirdPartyConfig } = this.state
        this.setState({
            thirdPartyConfig: thirdPartyConfig.filter((item) => item.indexId !== null),
        })
    }

    /**
     * 响应子组件删除卡片 -- 配置已经保存过 indexId!==null
     */
    protected handleDeleteSavedCard(indexId: number): void {
        const { thirdPartyConfig, editedCards } = this.state
        this.setState({
            thirdPartyConfig: thirdPartyConfig.filter((item) => item.indexId !== indexId),
            editedCards: editedCards.filter((item) => item !== indexId),
        })
    }

    /**
     * 响应子组件添加配置成功
     */
    protected handleAddConfigSuccess({ indexId, thirdPartyName, internalConfig, enabled, pluginClassName, messages }) {
        const { thirdPartyConfig } = this.state
        this.setState({
            thirdPartyConfig: [
                ...thirdPartyConfig.filter((item) => item.indexId !== null),
                {
                    indexId,
                    thirdPartyName,
                    internalConfig,
                    enabled,
                    pluginClassName,
                    messages,
                    plugin: {},
                },
            ],
        })
    }

    /**
     * 响应子组件编辑配置成功
     */
    protected handleEditConfigSuccess({ indexId, thirdPartyName, internalConfig, enabled, plugin, pluginClassName, messages }) {
        const { thirdPartyConfig, editedCards } = this.state
        this.setState({
            thirdPartyConfig: thirdPartyConfig.map((item) => {
                return item.indexId === indexId ?
                    {
                        ...item,
                        thirdPartyName,
                        internalConfig,
                        enabled,
                        pluginClassName,
                        messages,
                        plugin,
                    }
                    : item
            }),
            editedCards: editedCards.filter((item) => item !== indexId),
        })
    }

    /**
     * 子组件处于编辑状态时，向editedCards写入处于编辑状态的indexId
     */
    protected handleEditedCardIncrease(indexId: number): void {
        const { editedCards } = this.state

        if (indexId !== null) {
            this.setState({
                editedCards: union(editedCards, [indexId]),
            })
        }
    }

    /**
     * 子组件取消编辑状态，从editedCards移除取消编辑状态的indexId
     */
    protected handleEditedCardDecrease(indexId: number): void {
        const { editedCards } = this.state

        if (indexId !== null) {
            this.setState({
                editedCards: editedCards.filter((item) => item !== indexId),
            })
        }
    }

    /**
     * 点击“添加第三方应用按钮”，增加新建卡片
     */
    protected handleAddCard() {
        this.setState({
            thirdPartyConfig: [
                ...this.state.thirdPartyConfig,
                {
                    indexId: null,
                    thirdPartyName: '',
                    internalConfig: '',
                    pluginClassName: '',
                    messages: [],
                    enabled: false,
                    plugin: null,
                },
            ],
        })
    }

    /**
     * 验证消息服务名称是否重复
     */
    protected validate = (currentThirdPartyName: string): ValidateStatus => {
        return this.state.thirdPartyConfig.some(({ thirdPartyName }) => thirdPartyName === currentThirdPartyName) ?
            ValidateStatus.ThirdPartyNameRepeat : ValidateStatus.Normal;
    }
}