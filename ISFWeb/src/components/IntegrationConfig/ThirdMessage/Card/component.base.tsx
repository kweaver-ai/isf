import { assign, trim, isEqual, filter, findIndex, noop, isEmpty, uniqBy, values } from 'lodash';
import * as WebUploader from '@/libs/webuploader';
import { getDefaultStorage, getObjectStorageInfoByApp } from '@/core/apis/console/ossgateway'
import { ErrorCode as ErrCode } from '@/core/apis/openapiconsole/errorcode'
import { addThirdMessage, editThirdMessage, deleteThirdMessage } from '@/core/apis/console/thirdMessage';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { getAccessPrefix } from '@/core/accessprefix';
import { Message } from '@/sweet-ui';
import WebComponent from '../../../webcomponent';
import {
    ConfigInfo,
    checkConfig,
    PluginInfo,
    ConfigItem,
    transformJsonToArray,
    transformArrayToObject,
    ValidateStatus,
    Types,
    ErrorCode,
    formatterError,
    Stage,
    MessageConfigItem,
    checkMessages,
    InterfaceParams,
} from '../helper';
import { DefaultKeys, defaultKeys } from './type';
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';

interface CardProps {
    /**
     * swf文件路径
     */
    swf: string;

    /**
     * 配置参数信息
     */
    configInfo: ConfigInfo;

    /**
     * 是否显示删除按钮，如果是第一项则不显示
     */
    showDeleteIcon: boolean;

    /**
     * 响应删除未保存的卡片
     */
    onRequestDeleteUnSavedCard: () => void;

    /**
     * 响应删除已保存的卡片
     */
    onRequestDeleteSavedCard: (indexId: number) => void;

    /**
     * 添加一组配置成功
     */
    onRequestAddConfigSuccess: ({ ...ConfigInfo }) => void;

    /**
     * 编辑一组配置成功
     */
    onRequestEditConfigSuccess: ({ ...ConfigInfo }) => void;

    /**
     * 处于编辑状态的卡片增加
     */
    onRequestEditedCardIncrease: (indexId: number) => void;

    /**
     * 处于编辑状态的卡片减少
     */
    onRequestEditedCardDecrease: (indexId: number) => void;

    /**
     * 保存时，验证消息服务名称是否重复
     */
    onRequestValidateState: (thirdPartyName: string) => ValidateStatus;
}

interface CardState {
    /**
     * 第三方消息服务名称
     */
    thirdPartyName: string;

    /**
     * 插件类名
     */
    pluginClassName: string;

    /**
     * 第三方认证服务启用状态
     */
    enabled: boolean;

    /**
     * 服务端参数配置列表（去除msg_type_list属性）
     */
    internalConfig: ReadonlyArray<ConfigItem>;

    /**
     * 消息列表（internalConfig中的msg_type_list属性值）
     */
    messages: ReadonlyArray<MessageConfigItem>;

    /**
     * 插件
     */
    plugin: null | PluginInfo;

    /**
     * 验证状态
     */
    validateStatus: {
        /**
         * 第三方消息服务名称验证状态
         */
        thirdPartyNameValidateStatus: ValidateStatus;

        /**
         * 插件类名
         */
        pluginClassNameValidateStatus: ValidateStatus;
    };

    /**
     * 是否处于编辑状态
     */
    edited: boolean;

    /**
     * 显示服务端高级配置对话框
     */
    showInternalAdvancedConfigDialog: boolean;

    /**
     * 控制高级配置时，修改了默认项被还原提示对话框的显示
     */
    invalidConfig: InterfaceParams | null;

    /**
     * 错误信息
     */
    error: null | object;

    /**
     * 是否正在上传
     */
    uploading: boolean;

    /**
     * 显示删除（已保存）卡片对话框
     */
    showDeleteCardDialog: boolean;

    /**
    * toast 错误弹窗
    */
    errorToast: null | object;
}

export default class CardBase extends WebComponent<CardProps, CardState> {

    // WebUploader 实例
    uploader = null;

    static contextType = AppConfigContext;

    // 存储成功配置的值，用于需要还原的操作
    originConfig = {
        thirdPartyName: '',
        pluginClassName: '',
        enabled: false,
        internalConfig: [],
        plugin: null,
        messages: [],
    };

    // 第三方认证服务唯一索引
    indexId = null;

    internalConfigContainer = null;

    // 消息类型父级盒子ref
    messageConfigItemContainer = null;

    // 消息类型ref
    messagesItemRef = null;

    select = null

    // 不能重复点击保存按钮
    loading = false;

    static defaultProps = {
        swf: '',
        configInfo: {
            indexId: null,
            thirdPartyName: '',
            enabled: false,
            internalConfig: '',
            pluginClassName: '',
            messages: [],
            plugin: null,
        },
        showDeleteIcon: false,
        onRequestDeleteUnSavedCard: noop,
        onRequestDeleteSavedCard: noop,
        onRequestAddConfigSuccess: noop,
        onRequestEditConfigSuccess: noop,
        onRequestEditedCardIncrease: noop,
        onRequestEditedCardDecrease: noop,
    }

    state = {
        thirdPartyName: '',
        pluginClassName: '',
        enabled: false,
        internalConfig: [],
        messages: [],
        plugin: null,
        validateStatus: {
            thirdPartyNameValidateStatus: ValidateStatus.Normal,
            pluginClassNameValidateStatus: ValidateStatus.Normal,
        },
        edited: false,
        showInternalAdvancedConfigDialog: false,
        invalidConfig: null,
        error: null,
        uploading: false,
        showDeleteCardDialog: false,
        errorToast: null,
    }

    // 将 props 中传递的配置信息存入state和originConfig
    async componentDidMount() {
        const { indexId, thirdPartyName, enabled, internalConfig, plugin, pluginClassName, messages } = this.props.configInfo
        if (indexId !== null) {
            this.indexId = indexId;

            this.setState({
                thirdPartyName,
                pluginClassName,
                enabled,
                internalConfig: transformJsonToArray(internalConfig, Stage.Initialize),
                messages,
                plugin,
            }, () => {
                // 成功后更新originConfig
                const { thirdPartyName, pluginClassName, enabled, internalConfig, plugin, messages } = this.state
                this.originConfig = {
                    thirdPartyName,
                    pluginClassName,
                    enabled,
                    internalConfig,
                    plugin,
                    messages,
                }
            })
        }
        this.initWebUpload()
    }

    /**
    * 初始化上传组件
    */
    private initWebUpload() {
        const self = this;
        self.uploader = new WebUploader.Uploader({
            swf: self.props.swf,
            server: `${getAccessPrefix()}/api/thirdparty-message-plugin/v1/plugins/${this.indexId}`,
            auto: false,
            threads: 1,
            duplicate: true,
            sendAsBinary: true,
            method: 'PUT',
            pick: {
                id: self.select,
                multiple: false,
            },
            accept: {
                title: '*.tar.gz',
                extensions: 'tar.gz',
                mimeTypes: 'application/x-gzip',
            },
            onBeforeFileQueued: async function (file) {
                self.uploader.reset();
                if (!/\.tar\.gz$/.test(file.name) || file.size > 200 * 1024 * 1024) {
                    self.uploader.cancelFile(file)
                    self.uploader.removeFile(file, true)
                    self.setState({
                        errorToast: { errCode: ErrorCode.WrongFormat },
                    })
                    setTimeout(() => {
                        self.setState({
                            errorToast: null,
                        })
                    }, 3000)
                } else if ((/[\\\/\@\#\$\%\^\&\*\(\)\[\]]/.test(file.name)) || (/^[\.\+\-]/.test(file.name)) || (file.name.length > 255)) {
                    self.uploader.cancelFile(file)
                    self.uploader.removeFile(file, true)
                    self.setState({
                        errorToast: { errCode: ErrorCode.IllegalFileName },
                    })
                    setTimeout(() => {
                        self.setState({
                            errorToast: null,
                        })
                    }, 3000)
                } else {
                    let storage_id: string
                    try {
                        ({ storage_id } = await getDefaultStorage())
                    } catch {
                        storage_id = ''
                    }

                    if (!storage_id) {
                        // 获取当前站点的存储信息
                        const ossInfoDatas = await getObjectStorageInfoByApp({ app: 'as' })
                        // 当前站点是否有开启可用存储
                        const isHasOss = ossInfoDatas.some((item) => (item.enabled))

                        if (!isHasOss) {
                            self.uploader.cancelFile(file)
                            self.uploader.removeFile(file, true)
                            await Message.alert({ message: formatterError(ErrorCode.NoStorage) });
                        } else {
                            self.setState({
                                uploading: true,
                            })
                            self.uploader.upload()
                        }

                    } else {
                        // 如果有默认存储则使用默认存储
                        self.setState({
                            uploading: true,
                        })
                        self.uploader.upload()
                    }
                }
            },

            onUploadBeforeSend: function (object, data, headers) {
                assign(headers, {
                    Authorization: 'Bearer ' + self.context?.getToken?.(),
                    'x-filename': object.file.name,
                });
            },

            // 当文件上传成功时触发。
            onUploadSuccess: function (file, response) {
                self.setState({
                    plugin: {
                        filename: file.name,
                    },
                }, async () => {
                    // 更新原始数据
                    self.originConfig = assign({}, self.originConfig, { plugin: self.state.plugin });
                    // 通知父组件信息更新成功
                    const { thirdPartyName, pluginClassName, internalConfig, messages, enabled, plugin } = self.state

                    self.props.onRequestEditConfigSuccess({
                        thirdPartyName,
                        internalConfig: JSON.stringify({ ...transformArrayToObject(internalConfig) }),
                        enabled,
                        indexId: self.indexId,
                        plugin,
                        pluginClassName,
                        messages,
                    })

                    try {
                        await manageLog(
                            ManagementOps.SET,
                            __('设置 第三方消息集成 “${name}” 成功', { name: thirdPartyName }),
                            __('消息服务名称：${name}；消息模块插件：${filename}；插件参数配置：***',
                                {
                                    name: thirdPartyName,
                                    filename: plugin.filename,
                                }),
                            Level.WARN,
                        );

                        self.setState({
                            uploading: false,
                        });
                    } catch (ex) {
                        self.setState({
                            uploading: false,
                        });
                    }
                });
            },
            // 当文件上传出错时触发。
            onUploadError: function (object, code) {
                self.setState({
                    uploading: false,
                });

                if (self.state.errorToast && self.state.errorToast.errCode !== ErrorCode.ServerError) {
                    self.setState({
                        errorToast: {
                            errCode: ErrorCode.UnknownError,
                        },
                    })

                    setTimeout(() => {
                        self.setState({
                            errorToast: null,
                        })
                    }, 3000)
                }
            },
            // 当某个文件上传到服务端响应后，询问服务端响应是否有效
            onUploadAccept: function (object, ex) {
                self.setState({
                    uploading: false,
                });

                const { code, message } = ex;

                if (ex._raw && /(502)|(503)/.test(ex._raw)) {
                    self.setState({
                        errorToast: {
                            errCode: ErrorCode.ServerError,
                        },
                    })

                    setTimeout(() => {
                        self.setState({
                            errorToast: null,
                        })
                    }, 3000)

                    return
                }

                if (code) {
                    if (code === ErrCode.TokenExpire) {

                        return
                    }

                    message && Message.alert({ message });
                }
            },
        })
    }

    /**
     * 状态按钮状态改变 - 改变装填按钮的选中状态，配置的值还原到上一次成功保存的状态，显示保存和取消按钮，所有错误信息复位
     * @param {*} checked 选中状态
     */
    protected statusChange(checked: boolean): void {
        this.setState({
            enabled: checked,
            edited: true,
            validateStatus: {
                thirdPartyNameValidateStatus: ValidateStatus.Normal,
                pluginClassNameValidateStatus: ValidateStatus.Normal,
            },
            error: null,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 消息服务名称输入框值改变
     * @param {*} value 输入的值
     */
    protected handleThirdPartyNameChange(value: string): void {
        this.setState({
            thirdPartyName: value,
            validateStatus: {
                ...this.state.validateStatus,
                thirdPartyNameValidateStatus: ValidateStatus.Normal,
            },
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 插件类名输入框值改变
     * @param {*} value 输入的值
     */
    protected handlePluginClassNameChange(value: string): void {
        this.setState({
            pluginClassName: value,
            validateStatus: {
                ...this.state.validateStatus,
                pluginClassNameValidateStatus: ValidateStatus.Normal,
            },
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
    * 点击添加，向internalConfig中插入一项，显示保存和取消按钮
    */
    protected addInternalConfig(): void {
        this.setState({
            internalConfig: [
                {
                    index: null,
                    configName: '',
                    configType: Types.StringType,
                    configValue: '',
                    // 新添加的项都不是默认项，不是默认项的名称都是必填
                    defaultConfig: false,
                    valueRequired: false,
                    nameValidateStatus: ValidateStatus.Normal,
                    valueValidateStatus: ValidateStatus.Normal,
                },
                ...this.state.internalConfig,
            ],
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
        // 新增时滚动到最上方
        this.internalConfigContainer.scrollTop = 0
    }

    /**
    * 点击高级配置，打开高级配置对话框(点击高级配置按钮的时候，不把showSaveButton设置为true，当关闭的时候再去比较配置参数是否改变，变了再去显示保存、取消按钮)
    */
    protected openInternalAdvancedConfigDialog(): void {
        this.setState({
            showInternalAdvancedConfigDialog: true,
        })
    }

    /**
     * 点击添加消息类型配置
     */
    protected addMessageConfig = () => {
        this.setState({
            messages: [
                {
                    configValue: '',
                    valueValidateStatus: ValidateStatus.Normal,
                },
                ...this.state.messages,
            ],
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
        // 新增时滚动到最上方
        this.messageConfigItemContainer.scrollTop = 0
    }

    /**
     * 删除消息类型
     */
    protected deleteMessageConfigItem = (index: number): void => {
        this.setState({
            messages: this.state.messages.filter((item, i) => {
                return i !== index
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 保存配置
     */
    protected async save() {
        // 检验各个输入的合法性
        if (!this.loading && await this.validate()) {
            this.loading = true;

            try {
                const { thirdPartyName, pluginClassName, internalConfig, plugin, messages, enabled } = this.state

                // 去重
                const newMessages = uniqBy(messages, 'configValue');
                const allConfig = transformArrayToObject(internalConfig);
                const params = {
                    thirdPartyName,
                    enabled,
                    internalConfig: JSON.stringify({ ...allConfig }),
                    pluginClassName,
                    messages: newMessages,
                }

                this.setState({
                    messages: newMessages,
                })

                // 通过验证，向后端发送请求
                const info = {
                    thirdparty_name: thirdPartyName,
                    class_name: pluginClassName,
                    channels: newMessages.map(({ configValue }) => configValue),
                    enabled,
                    config: { ...allConfig },
                }

                if (!this.indexId) {
                    // 添加
                    const { id: indexId } = await addThirdMessage(info);
                    this.indexId = indexId
                    this.props.onRequestAddConfigSuccess({
                        ...params,
                        indexId,
                    })
                } else {
                    // 修改
                    await editThirdMessage({ id: this.indexId, ...info });
                    this.props.onRequestEditConfigSuccess({ ...params, indexId: this.indexId, plugin })
                }

                // 更新原始数据
                this.originConfig = assign({}, {
                    thirdPartyName,
                    pluginClassName,
                    internalConfig: [...internalConfig],
                    messages: newMessages,
                    enabled,
                }, { plugin })

                this.setState({
                    edited: false,
                }, () => {
                    if (this.indexId) {
                        this.props.onRequestEditedCardDecrease(this.indexId)
                    }
                });

                enabled ?
                    manageLog(
                        ManagementOps.SET,
                        __('设置 第三方消息集成 “${name}” 成功', { name: thirdPartyName }),
                        plugin && plugin.filename ?
                            __(
                                '消息服务名称：${name}；消息模块插件：${filename}；插件参数配置：***',
                                {
                                    name: thirdPartyName,
                                    filename: plugin.filename,
                                },
                            )
                            :
                            __(
                                '消息服务名称：${name}；消息模块插件：未上传；插件参数配置：***',
                                {
                                    name: thirdPartyName,
                                },
                            )
                        ,
                        Level.WARN,
                    )
                    :
                    manageLog(
                        ManagementOps.SET,
                        __('禁用 第三方消息集成 “${name}” 成功', { name: thirdPartyName }),
                        plugin && plugin.filename ?
                            __(
                                '消息服务名称：${name}；消息模块插件：${filename}；插件参数配置：***',
                                {
                                    name: thirdPartyName,
                                    filename: plugin.filename,
                                },
                            )
                            :
                            __(
                                '消息服务名称：${name}；消息模块插件：未上传；插件参数配置：***',
                                {
                                    name: thirdPartyName,
                                },
                            )
                        ,
                        Level.WARN,
                    );
            } catch (ex) {
                const { code, message } = ex;

                if (code) {
                    if (code === ErrCode.AppNotExist) {
                        this.setState({
                            errorToast: {
                                errCode: ErrorCode.AppNotExist,
                            },
                        })

                        setTimeout(() => {
                            this.setState({
                                errorToast: null,
                            })
                            this.props.onRequestDeleteSavedCard(this.indexId);
                        }, 3000)

                        return/*  */
                    }

                    message && Message.alert({ message });
                }
            }

            this.loading = false;
        }
    }

    /**
    * 保存时验证表单输入的合法性
    */
    private async validate() {
        const { thirdPartyName, pluginClassName, internalConfig, messages } = this.state
        // 处理过名称和值，包含有验证状态的internalConfig
        const newInternalConfig = checkConfig(internalConfig)
        const trimThirdPartyName = trim(thirdPartyName)
        const trimPluginClassName = trim(pluginClassName)
        const { newMessages, validateIndex } = checkMessages(messages)
        const thirdPartyNameValidateStatus = trimThirdPartyName ?
            trimThirdPartyName.length > 128 || /[\\\/\:\*\?\"\<\>\|]/.test(trimThirdPartyName) ?
                ValidateStatus.NameLengthExceedsLimitOrHasInvalid
                : trimThirdPartyName === this.originConfig.thirdPartyName ?
                    ValidateStatus.Normal
                    : this.props.onRequestValidateState(trimThirdPartyName)
            : ValidateStatus.Empty;

        const pluginClassNameValidateStatus = trimPluginClassName ?
            trimPluginClassName.length > 128 || /[\\\/\:\*\?\"\<\>\|]/.test(trimPluginClassName) ?
                ValidateStatus.PluginClassNameValueToLongOrHasInvalid
                : ValidateStatus.Normal
            : ValidateStatus.Empty;

        this.setState({
            thirdPartyName: trimThirdPartyName,
            pluginClassName: trimPluginClassName,
            validateStatus: {
                thirdPartyNameValidateStatus,
                pluginClassNameValidateStatus,
            },
            internalConfig: [...newInternalConfig],
            messages: newMessages,
            error: messages.length === 0 ? { error: { errID: ErrorCode.MessageEmpty } } : null,
        })

        if (validateIndex !== -1) {
            // 如果有验证不通过的项，把第一个不通过的项定位到父级盒子顶部
            this.messageConfigItemContainer.scrollTop = this.messagesItemRef.clientHeight * validateIndex;
        }

        const internalIndex = findIndex(newInternalConfig, (item) => {
            return item.nameValidateStatus !== ValidateStatus.Normal || item.valueValidateStatus !== ValidateStatus.Normal
        })

        if (internalIndex !== -1) {
            this.internalConfigContainer.scrollTop = internalIndex * 40
        }

        return !!(
            thirdPartyNameValidateStatus === ValidateStatus.Normal &&
            pluginClassNameValidateStatus === ValidateStatus.Normal &&
            internalIndex === -1 &&
            validateIndex === -1 &&
            messages.length !== 0)
    }

    /**
     * 取消保存 - 配置的值还原到上一次成功保存的状态，隐藏保存和取消按钮，所有错误信息复位
     */
    protected cancel(): void {
        this.setState({
            ...this.originConfig,
            edited: false,
            validateStatus: {
                thirdPartyNameValidateStatus: ValidateStatus.Normal,
                pluginClassNameValidateStatus: ValidateStatus.Normal,
            },
            error: null,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardDecrease(this.indexId)
            }
        })
    }

    /**
     * 处理子组件ConfigItem名称改变事件（服务端） - 改变相应项名称的值和验证状态，并且显示保存和取消按钮
     * @param {string} value 输入的值
     * @param {number} index 被修改项的数组下标
     */
    protected internalConfigNameChange(value: string, index: number): void {
        const { internalConfig } = this.state

        this.setState({
            internalConfig: internalConfig.map((item, i) => {
                return i === index ? {
                    ...item,
                    configName: value,
                    nameValidateStatus: ValidateStatus.Normal,
                } : item
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 处理子组件ConfigItem类型改变事件（服务端）
     * @param {*} type 选择的类型
     * @param {*} index 被修改项的数组下标
     */
    protected internalConfigTypeChange(type: Types, index: number): void {
        const { internalConfig } = this.state

        this.setState({
            internalConfig: internalConfig.map((item, i) => {
                return i === index ? {
                    ...item,
                    configType: type,
                    configValue: type === Types.BooleanType ? true : '',
                } : item
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
    * 处理子组件ConfigItem值改变事件（服务端） - 当为数组输入框的时候要实时验证输入值的正确性
    * @param {*} value 输入的值
    * @param {*} type 选择的类型
    * @param {*} index 被修改项的数组下标
    */
    protected internalConfigValueChange(value: any, type: Types, index: number): void {
        const { internalConfig } = this.state

        this.setState({
            internalConfig: internalConfig.map((item, i) => {
                return i === index ? {
                    ...item,
                    configValue: (type === Types.NumberType && value === null) ? '' : value,
                    valueValidateStatus: ValidateStatus.Normal,
                } : item
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 处理子组件ConfigItem删除一条配置项（服务端）
     * @param {*} index 被删除项的数组下标
     */
    protected deleteInternalConfigItem(index: number): void {
        const { internalConfig } = this.state

        this.setState({
            internalConfig: internalConfig.filter((item, i) => {
                return i !== index
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 处理子组件AdvancedConfig传出的关闭、取消高级配置对话框事件
     */
    protected closeInternalAdvancedConfigDialog(): void {
        // 只需要关闭弹窗，无需重置数据，也无需改变保存、取消按钮的状态
        this.setState({
            showInternalAdvancedConfigDialog: false,
        })
    }

    /**
     * 处理子组件AdvancedConfig传出的确认高级配置对话框事件（服务端）
     * @param {*} config 已经将非法参数还原的配置项
     */
    protected updateConfig(config: string): void {
        const allConfig = JSON.parse(config);

        const { thirdPartyName, pluginClassName, enabled, messages, internalConfig } = this.state

        const oldAllConfig = {
            thirdparty_name: thirdPartyName,
            class_name: pluginClassName,
            enabled,
            channels: messages.map(({ configValue }) => configValue),
            config: transformArrayToObject(internalConfig),
        }

        // 比较传递过来的config 和 state上的是否一致
        if (!isEqual(oldAllConfig, allConfig)) {
            const { invalidConfig, validConfig } = this.updateInternalConfig(allConfig, oldAllConfig)

            const { thirdparty_name, class_name, enabled, channels, config } = validConfig;

            // 先将传递过来的字符串"前后的空格处理掉之后在进行转换
            const configArr = transformJsonToArray(JSON.stringify(config).replace(/\s*\"\s*/g, '"'), Stage.Advanced)

            this.setState({
                thirdPartyName: thirdparty_name,
                pluginClassName: class_name,
                enabled,
                messages: channels.map((configValue) => {
                    return {
                        configValue,
                        valueValidateStatus: ValidateStatus.Normal,
                    }
                }),
                internalConfig: [...configArr],
                invalidConfig: isEmpty(invalidConfig) ? null : { ...invalidConfig },
                edited: true,
                showInternalAdvancedConfigDialog: false,
                validateStatus: {
                    thirdPartyNameValidateStatus: ValidateStatus.Normal,
                    pluginClassNameValidateStatus: ValidateStatus.Normal,
                },
            }, () => {
                if (this.indexId) {
                    this.props.onRequestEditedCardIncrease(this.indexId)
                }
            })
        } else {
            // 一致，则不更新数据，也不改变保存取消按钮状态，仅关闭对话框
            this.setState({
                showInternalAdvancedConfigDialog: false,
            })
        }
    }

    /**
    * 更新数据
    */
    private updateInternalConfig(allConfig: InterfaceParams, oldAllConfig: InterfaceParams): {
        invalidConfig: InterfaceParams;
        validConfig: InterfaceParams;
    } {
        // 记录非法修改的默认参数数组
        let invalidConfig = null;
        // 还原非法修改后的合法项
        let validConfig = null;

        defaultKeys.forEach((key) => {
            const value = allConfig[key];

            if (value !== undefined && this.check(key, value)) {
                validConfig = {
                    [key]: value,
                    ...validConfig,
                }
            } else {
                validConfig = {
                    [key]: oldAllConfig[key],
                    ...validConfig,
                }

                invalidConfig = {
                    [key]: oldAllConfig[key],
                    ...invalidConfig,
                }
            }
        })

        return { invalidConfig, validConfig }
    }

    /**
     * 检查config每一项是否合规
     */
    private check = (key: DefaultKeys, value: any): boolean => {
        switch (key) {
            case DefaultKeys.Enabled:

                return typeof value === 'boolean';

            case DefaultKeys.ThirdpartyName:
            case DefaultKeys.ClassName:

                return typeof value === 'string';

            case DefaultKeys.Config:

                return typeof value === 'object' && !Array.isArray(value) && value !== null && values(value).every((val) => val !== null);

            case DefaultKeys.Channels:

                return Array.isArray(value) && value.every((item) => typeof item === 'string');

            default:

                return false;
        }
    }

    /**
     * 处理子组件ResetInvalidConfig传出的关闭、确认高级配置还原默认参数配置事件
     */
    protected resetInvalidConfig() {
        this.setState({
            invalidConfig: null,
        })
    }

    /**
     * 响应消息类型的全选事件
     * @param checked
     */
    protected handleConfigValueChange(configValue: string, index: number): void {
        this.setState({
            messages: this.state.messages.map((item, i) => {
                return i === index ? {
                    configValue,
                    valueValidateStatus: ValidateStatus.Normal,
                } : item;
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 获取消息列表单项盒子ref
     */
    protected handleMessagesItemRef = (ref) => {
        this.messagesItemRef = ref;
    }

    /**
     * 响应选中或者取消选中一条消息
     * @param checked
     * @param value
     */
    protected handleSelectOne(checked: boolean, value: number): void {
        const { messages } = this.state
        // 选中则将该项添加到messages列表中，取消选中从列表中去除
        this.setState({
            messages: checked ? [...messages, value] : filter(messages, (message) => {
                return message !== value
            }),
            edited: true,
        }, () => {
            if (this.indexId) {
                this.props.onRequestEditedCardIncrease(this.indexId)
            }
        })
    }

    /**
     * 点击删除卡片
     */
    protected handleDelete() {
        // 如果尚未保存则直接删除卡片
        if (this.indexId === null) {
            this.props.onRequestDeleteUnSavedCard()
        } else {
            this.setState({ showDeleteCardDialog: true })
        }
    }

    /**
     * 确认删除一个已经保存的卡片
     */
    protected async handleDeleteSavedCard() {
        this.setState({
            showDeleteCardDialog: false,
        })
        try {
            await deleteThirdMessage({ id: this.indexId })
            this.props.onRequestDeleteSavedCard(this.indexId)
            manageLog(
                ManagementOps.DELETE,
                __('删除 第三方消息集成 “${name}” 成功', { name: this.state.thirdPartyName }),
                undefined,
                Level.WARN,
            );
        } catch ({ message }) {
            if (message) {
                Message.alert({ message });
            }
        }
    }

    /**
     * 取消删除一个已经保存的卡片
     */
    protected handleCancleDelete() {
        this.setState({ showDeleteCardDialog: false })
    }
}