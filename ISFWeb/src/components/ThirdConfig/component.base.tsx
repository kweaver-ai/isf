import { assign, trim, isEqual, filter, findIndex } from 'lodash';
import { getDefaultStorage, getObjectStorageInfoByApp } from '@/core/apis/console/ossgateway'
import * as WebUploader from '@/libs/webuploader';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { getHeaders } from '@/core/token'
import { getThirdPartyAppConfig, addThirdPartyAppConfig, setThirdPartyAppConfig } from '@/core/thrift/sharemgnt/sharemgnt';
import { Message } from '@/sweet-ui';
import WebComponent from '../webcomponent'
import { PluginInfo, checkConfig, ConfigItem, Terminal, transformJsonToArray, transformArrayToObject, ValidateStatus, NcTPluginType, Types, defaultClientConfigArray, defaultInternalConfigArray, defaultClientConfig, defaultInternalConfig, formatterError, ErrorCode, Stage } from './helper';
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';

interface ThirdConfigProps {
    /**
     * swf文件路径
     */
    swf: string;

    /**
     * 上传地址
     */
    server: string;
}

interface ThirdConfigState {
    /**
     * 第三方认证服务ID
     */
    thirdPartyId: string;

    /**
     * 第三方认证服务名称
     */
    thirdPartyName: string;

    /**
     * 第三方认证服务启用状态
     */
    enabled: boolean;

    /**
     * 客户端参数配置列表
     */
    clientConfig: ReadonlyArray<ConfigItem>;

    /**
     * 服务端参数配置列表
     */
    internalConfig: ReadonlyArray<ConfigItem>;

    /**
     * 插件
     */
    plugin: null | PluginInfo;

    /**
     * 验证状态
     */
    validateStatus: {
        /**
         * 第三方认证服务ID验证状态
         */
        thirdPartyIdValidateStatus: ValidateStatus;

        /**
         * 第三方认证服务名称验证状态
         */
        thirdPartyNameValidateStatus: ValidateStatus;
    };

    /**
     * 是否处于编辑状态
     */
    edited: boolean;

    /**
     * 显示客户端高级配置对话框
     */
    showClientAdvancedConfigDialog: boolean;

    /**
     * 显示服务端高级配置对话框
     */
    showInternalAdvancedConfigDialog: boolean;

    /**
     * 控制高级配置时，修改了默认项被还原提示对话框的显示
     */
    invalidConfig: ReadonlyArray<ConfigItem> | null;

    /**
     * 错误信息
     */
    error: null | object;

    /**
     * 是否正在上传
     */
    uploading: boolean;

    /**
     * toast 错误弹窗
     */
    errorToast: null | object;
}

export default class ThirdConfigBase extends WebComponent<ThirdConfigProps, ThirdConfigState> {
    static contextType = AppConfigContext;
    // WebUploader 实例
    uploader = null;

    // 存储成功配置的值，用于需要还原的操作
    originConfig = {
        thirdPartyId: '',
        thirdPartyName: '',
        enabled: false,
        clientConfig: defaultClientConfigArray,
        internalConfig: defaultInternalConfigArray,
        plugin: null,
    };

    /**
     * 第三方认证服务唯一索引
     */
    indexId = null

    clientConfigContainer = null

    internalConfigContainer = null

    select = null

    static defaultProps = {
        swf: '',
        server: '',
    }

    state = {
        thirdPartyId: '',
        thirdPartyName: '',
        enabled: false,
        clientConfig: defaultClientConfigArray,
        internalConfig: defaultInternalConfigArray,
        plugin: null,
        validateStatus: {
            thirdPartyIdValidateStatus: ValidateStatus.Normal,
            thirdPartyNameValidateStatus: ValidateStatus.Normal,
        },
        edited: false,
        showClientAdvancedConfigDialog: false,
        showInternalAdvancedConfigDialog: false,
        invalidConfig: null,
        error: null,
        uploading: false,
        errorToast: null,
    }

    // 进入界面，获取数据，并且将获取的数据用副本originConfig保存
    async componentDidMount() {
        try {
            const configList = await getThirdPartyAppConfig(NcTPluginType.Authentication);
            // 如果配置成功过,则根据返回来的结果填充state
            if (configList && configList.length > 0) {
                const { indexId, thirdPartyId, thirdPartyName, enabled, config, internalConfig, plugin } = configList[0]
                this.indexId = indexId
                this.setState({
                    thirdPartyId,
                    thirdPartyName,
                    enabled,
                    clientConfig: transformJsonToArray(config, Terminal.Client, Stage.Initialize),
                    internalConfig: transformJsonToArray(internalConfig, Terminal.Internal, Stage.Initialize),
                    plugin,
                }, () => {
                    // 成功后更新originConfig
                    const { thirdPartyId, thirdPartyName, enabled, clientConfig, internalConfig, plugin } = this.state
                    this.originConfig = {
                        thirdPartyId: thirdPartyId,
                        thirdPartyName: thirdPartyName,
                        enabled: enabled,
                        clientConfig: clientConfig,
                        internalConfig: internalConfig,
                        plugin: plugin,
                    }
                })
            }
        } catch (err) {
            Message.alert({ message: err.error.errMsg });
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
            server: self.props.server,
            auto: false,
            threads: 1,
            duplicate: true,
            fileVal: 'package',
            pick: {
                id: self.select,
                multiple: false,
            },
            accept: {
                title: '*.tar.gz',
                extensions: 'tar.gz',
                mimeTypes: 'application/x-gzip',
            },
            // 当文件被加入队列之前触发
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

                        if (isHasOss) {
                            self.setState({
                                uploading: true,
                            })
                            self.uploader.upload()
                        } else {
                            self.uploader.cancelFile(file)
                            self.uploader.removeFile(file, true)
                            await Message.alert({ message: formatterError(ErrorCode.NoStorage) });
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
                assign(headers, getHeaders(self.context?.getToken?.()).headers);
                data.filename = object.file.name;
                data.indexId = self.indexId;
                data.thirdPartyId = self.state.thirdPartyId;
                data.type = NcTPluginType.Authentication
            },

            // 当文件上传成功时触发。
            onUploadSuccess: async function (file, response) {
                self.setState({
                    plugin: {
                        filename: file.name,
                    },
                });
                // 更新原始数据
                self.originConfig = assign({}, self.originConfig, { plugin: self.state.plugin });

                try {
                    await manageLog(
                        ManagementOps.SET,
                        __('设置 第三方认证 "${name}" 成功', { name: self.state.thirdPartyName }),
                        __('认证服务ID：${id}；认证服务名称：${name}；认证模块插件：${filename}；客户端参数配置：***；服务端参数配置：***',
                            {
                                id: self.state.thirdPartyId,
                                name: self.state.thirdPartyName,
                                filename: self.state.plugin.filename,
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
            },
            // 当文件上传出错时触发。
            onUploadError: function () {
                self.setState({
                    uploading: false,
                });
                Message.alert({ message: __('上传失败。') });
            },
            // 当某个文件上传到服务端响应后，询问服务端响应是否有效
            onUploadAccept: function (object, ret) {
                if (ret.error) {
                    self.setState({
                        uploading: false,
                    });
                    Message.alert({ message: formatterError(ret.error.errID, ret.error.errMsg) });
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
            ...this.originConfig,
            enabled: checked,
            edited: true,
            validateStatus: {
                thirdPartyIdValidateStatus: ValidateStatus.Normal,
                thirdPartyNameValidateStatus: ValidateStatus.Normal,
            },
            error: null,
        })
    }

    /**
     * 认证服务ID输入框改变 - 改变值，并且取消验证状态,显示保存取消按钮
     * @param {*} value 输入的值
     */
    protected handleThirdPartyIdChange(value: string): void {
        this.setState({
            thirdPartyId: value,
            validateStatus: {
                ...this.state.validateStatus,
                thirdPartyIdValidateStatus: ValidateStatus.Normal,
            },
            edited: true,
        })
    }

    /**
     * 认证服务名称输入框值改变
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
        })
    }

    /**
     * 点击添加，向clientConfig中插入一项，显示保存和取消按钮
     */
    protected addClientConfig(): void {
        this.setState({
            clientConfig: [
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
                ...this.state.clientConfig,
            ],
            edited: true,
        })
        // 新增时滚动到最上方
        this.clientConfigContainer.scrollTop = 0
    }

    /**
    * 点击添加，向=internalConfig中插入一项，显示保存和取消按钮
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
        })
        // 新增时滚动到最上方
        this.internalConfigContainer.scrollTop = 0
    }

    /**
     * 点击高级配置，打开高级配置对话框
     */
    protected openClientAdvancedConfigDialog(): void {
        this.setState({
            showClientAdvancedConfigDialog: true,
        })
    }

    /**
    * 点击高级配置，打开高级配置对话框
    */
    protected openInternalAdvancedConfigDialog(): void {
        this.setState({
            showInternalAdvancedConfigDialog: true,
        })
    }

    /**
     * 保存配置
     */
    protected async save() {
        const { enabled } = this.state

        try {
            if (enabled) {
                // 开关打开，打开状态检验各个输入的合法性
                if (await this.validate()) {
                    const { thirdPartyId, thirdPartyName, clientConfig, internalConfig, plugin } = this.state
                    const stringifyClientConfig = JSON.stringify(transformArrayToObject(clientConfig))
                    const stringifyInternalConfig = JSON.stringify(transformArrayToObject(internalConfig))

                    const info = {
                        thirdPartyId,
                        thirdPartyName,
                        config: stringifyClientConfig,
                        internalConfig: stringifyInternalConfig,
                        enabled,
                    }
                    if (!this.indexId) { // 添加
                        const indexId = await addThirdPartyAppConfig({
                            ncTThirdPartyConfig: {
                                ...info,
                                plugin: {
                                    ncTThirdPartyPluginInfo: {
                                        type: NcTPluginType.Authentication,
                                    },
                                },
                            },
                        });
                        this.indexId = indexId
                    } else { // 修改
                        await setThirdPartyAppConfig({
                            ncTThirdPartyConfig: {
                                ...info,
                                indexId: this.indexId,
                                plugin: {
                                    ncTThirdPartyPluginInfo: {
                                        type: NcTPluginType.Authentication,
                                    },
                                },
                            },
                        });
                    }
                    // 更新原始数据
                    this.originConfig = assign({}, {
                        thirdPartyId,
                        thirdPartyName,
                        clientConfig: [...clientConfig],
                        internalConfig: [...internalConfig],
                        enabled,
                    }, { plugin })

                    this.setState({
                        edited: false,
                        error: null,
                    });

                    manageLog(
                        ManagementOps.SET,
                        __('设置 第三方认证 "${name}" 成功', { name: thirdPartyName }),
                        plugin && plugin.filename ?
                            __(
                                '认证服务ID：${id}；认证服务名称：${name}；认证模块插件：${filename}；客户端参数配置：***；服务端参数配置：***',
                                {
                                    id: thirdPartyId,
                                    name: thirdPartyName,
                                    filename: plugin.filename,
                                },
                            )
                            :
                            __(
                                '认证服务ID：${id}；认证服务名称：${name}；认证模块插件：未上传；客户端参数配置：***；服务端参数配置：***',
                                {
                                    id: thirdPartyId,
                                    name: thirdPartyName,
                                },
                            )
                        ,
                        Level.WARN,
                    );
                }
            } else { // 开关关闭
                if (!this.indexId) {
                    // 没有配置过，保存时，直接修改enable状态（不用去改originConfig，因为他已经是false了）
                    this.setState({
                        edited: false,
                        error: null,
                    })
                } else {
                    // 修改
                    const { thirdPartyId, thirdPartyName, clientConfig, internalConfig, plugin } = this.originConfig
                    const stringifyClientConfig = JSON.stringify(transformArrayToObject(clientConfig))
                    const stringifyInternalConfig = JSON.stringify(transformArrayToObject(internalConfig))
                    const info = {
                        thirdPartyId,
                        thirdPartyName,
                        config: stringifyClientConfig,
                        internalConfig: stringifyInternalConfig,
                        enabled,
                    }
                    await setThirdPartyAppConfig({
                        ncTThirdPartyConfig: {
                            ...info,
                            indexId: this.indexId,
                            plugin: {
                                ncTThirdPartyPluginInfo: {
                                    type: NcTPluginType.Authentication,
                                },
                            },
                        },
                    });

                    this.originConfig = {
                        ...this.originConfig,
                        enabled: false,
                    }

                    this.setState({
                        edited: false,
                        error: null,
                    })

                    manageLog(
                        ManagementOps.SET,
                        __('禁用 第三方认证 "${name}" 成功', { name: thirdPartyName }),
                        plugin && plugin.filename ?
                            __(
                                '认证服务ID：${id}；认证服务名称：${name}；认证模块插件：${filename}；客户端参数配置：***；服务端参数配置：***',
                                {
                                    id: thirdPartyId,
                                    name: thirdPartyName,
                                    filename: plugin.filename,
                                },
                            )
                            :
                            __(
                                '认证服务ID：${id}；认证服务名称：${name}；认证模块插件：未上传；客户端参数配置：***；服务端参数配置：***',
                                {
                                    id: thirdPartyId,
                                    name: thirdPartyName,
                                },
                            )
                        ,
                        Level.WARN,
                    );
                }
            }
        } catch (error) {
            this.setState({
                error,
            })
        }
    }

    /**
     * 保存时验证表单输入的合法性
     */
    private async validate() {
        const { thirdPartyId, thirdPartyName, clientConfig, internalConfig } = this.state
        // 处理过名称和值，包含有验证状态的clientConfig
        const newClientConfig = checkConfig(clientConfig)
        const newInternalConfig = checkConfig(internalConfig)
        const trimThirdPartyId = trim(thirdPartyId)
        const trimThirdPartyName = trim(thirdPartyName)

        this.setState({
            validateStatus: {
                thirdPartyIdValidateStatus: trimThirdPartyId ? trimThirdPartyId.length > 128 ? ValidateStatus.IDLengthExceedsLimit : ValidateStatus.Normal : ValidateStatus.Empty,
                thirdPartyNameValidateStatus: trimThirdPartyName ? trimThirdPartyName.length > 128 ? ValidateStatus.NameLengthExceedsLimit : ValidateStatus.Normal : ValidateStatus.Empty,
            },
            thirdPartyId: trimThirdPartyId,
            thirdPartyName: trimThirdPartyName,
            clientConfig: [...newClientConfig],
            internalConfig: [...newInternalConfig],
        })

        // 查找客户端参数中，验证状态不为Normal的位置
        const clientIndex = findIndex(newClientConfig, (item) => {
            return item.nameValidateStatus !== ValidateStatus.Normal || item.valueValidateStatus !== ValidateStatus.Normal
        })

        // 如果客户端参数中，找到验证状态不为Normal的项则滚动最上方
        if (clientIndex !== -1) {
            this.clientConfigContainer.scrollTop = clientIndex * 40
        }

        const internalIndex = findIndex(newInternalConfig, (item) => {
            return item.nameValidateStatus !== ValidateStatus.Normal || item.valueValidateStatus !== ValidateStatus.Normal
        })

        if (internalIndex !== -1) {
            this.internalConfigContainer.scrollTop = internalIndex * 40
        }

        return !!(trimThirdPartyId && trimThirdPartyId.length <= 128 &&
            trimThirdPartyName && trimThirdPartyName.length <= 128 &&
            clientIndex === -1 &&
            internalIndex === -1)
    }

    /**
     * 取消保存 - 配置的值还原到上一次成功保存的状态，隐藏保存和取消按钮，所有错误信息复位
     */
    protected cancel(): void {
        this.setState({
            ...this.originConfig,
            edited: false,
            validateStatus: {
                thirdPartyIdValidateStatus: ValidateStatus.Normal,
                thirdPartyNameValidateStatus: ValidateStatus.Normal,
            },
            error: null,
        })
    }

    /**
     * 处理子组件ConfigItem名称改变事件（客户端） - 改变相应项名称的值和验证状态，并且显示保存和取消按钮
     * @param {string} value 输入的值
     * @param {number} index 被修改项的数组下标
     * @memberof ThirdConfigBase
     */
    protected clientConfigNameChange(value: string, index: number): void {
        const { clientConfig } = this.state

        this.setState({
            clientConfig: clientConfig.map((item, i) => {
                return i === index ? {
                    ...item,
                    configName: value,
                    nameValidateStatus: ValidateStatus.Normal,
                } : item
            }),
            edited: true,
        })
    }

    /**
     * 处理子组件ConfigItem名称改变事件（服务端） - 改变相应项名称的值和验证状态，并且显示保存和取消按钮
     * @param {string} value 输入的值
     * @param {number} index 被修改项的数组下标
     * @memberof ThirdConfigBase
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
        })
    }

    /**
     * 处理子组件ConfigItem类型改变事件（客户端）
     * @param {*} type 选择的类型
     * @param {*} index 被修改项的数组下标
     */
    protected clientConfigTypeChange(type: Types, index: number): void {
        const { clientConfig } = this.state

        this.setState({
            clientConfig: clientConfig.map((item, i) => {
                return i === index ? {
                    ...item,
                    configType: type,
                    configValue: type === Types.BooleanType ? true : '',
                } : item
            }),
            edited: true,
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
        })
    }

    /**
     * 处理子组件ConfigItem值改变事件（客户端） - 当为数组输入框的时候要实时验证输入值的正确性
     * @param {*} value 输入的值
     * @param {*} type 选择的类型
     * @param {*} index 被修改项的数组下标
     */
    protected clientConfigValueChange(value: any, type: Types, index: number): void {
        const { clientConfig } = this.state

        this.setState({
            clientConfig: clientConfig.map((item, i) => {
                return i === index ? {
                    ...item,
                    configValue: (type === Types.NumberType && value === null) ? '' : value,
                    valueValidateStatus: ValidateStatus.Normal,
                } : item
            }),
            edited: true,
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
        })
    }

    /**
     * 处理子组件ConfigItem删除一条配置项（客户端）
     * @param {*} index 被删除项的数组下标
     */
    protected deleteClientConfigItem(index: number): void {
        const { clientConfig } = this.state

        this.setState({
            clientConfig: clientConfig.filter((item, i) => {
                return i !== index
            }),
            edited: true,
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
        })
    }

    /**
     * 处理子组件AdvancedConfig传出的关闭、取消高级配置对话框事件
     */
    protected closeClientAdvancedConfigDialog(): void {
        // 只需要关闭弹窗，无需重置数据，也无需改变保存、取消按钮的状态
        this.setState({
            showClientAdvancedConfigDialog: false,
        })
    }

    /**
     * 处理子组件AdvancedConfig传出的关闭、取消高级配置对话框事件
     */
    protected closeInternalAdvancedConfigDialog(): void {
        this.setState({
            showInternalAdvancedConfigDialog: false,
        })
    }

    /**
     * 处理子组件AdvancedConfig传出的确认高级配置对话框事件(客户端)
     * @param {*} config 已经将非法参数还原的配置项
     */
    protected updateClientConfig(config: string): void {
        // 先将传递过来的字符串"前后的空格处理掉之后在进行转换
        const configArr = transformJsonToArray(config.replace(/\s*\"\s*/g, '"'), Terminal.Client, Stage.Advanced)
        const { clientConfig } = this.state
        // 比较传递过来的config 和 state上的是否一致
        if (!isEqual(configArr, clientConfig)) {
            // 从高级配置返回的数据中的默认项
            const defaultConfig = filter(configArr, (item) => {
                return !!(defaultClientConfig[item.configName])
            })

            // 用于还原的项所有默认项
            const revertConfig = filter(clientConfig, (item) => {
                return !!(defaultClientConfig[item.configName] && item.defaultConfig)
            })

            const { invalidConfig, validConfig } = this.updateConfig(configArr, revertConfig, defaultConfig, Terminal.Client)

            this.setState({
                clientConfig: validConfig.sort(this.compare('index')),
                invalidConfig: invalidConfig.length === 0 ? null : [...invalidConfig],
                edited: true,
                showClientAdvancedConfigDialog: false,
            })
        } else {
            // 一致，则不更新数据，也不改变保存取消按钮状态，仅关闭对话框
            this.setState({
                showClientAdvancedConfigDialog: false,
            })
        }
    }

    /**
     * 处理子组件AdvancedConfig传出的确认高级配置对话框事件（服务端）
     * @param {*} config 已经将非法参数还原的配置项
     */
    protected updateInternalConfig(config: string): void {
        const configArr = transformJsonToArray(config.replace(/\s*\"\s*/g, '"'), Terminal.Internal, Stage.Advanced)
        const { internalConfig } = this.state
        // 比较传递过来的config 和 state上的是否一致
        if (!isEqual(configArr, internalConfig)) {
            // 从高级配置返回的数据中的默认项
            const defaultConfig = filter(configArr, (item) => {
                return !!(defaultInternalConfig[item.configName])
            })

            // 用于还原的项所有默认项
            const revertConfig = filter(internalConfig, (item) => {
                return !!(defaultInternalConfig[item.configName] && item.defaultConfig)
            })

            const { invalidConfig, validConfig } = this.updateConfig(configArr, revertConfig, defaultConfig, Terminal.Internal)

            this.setState({
                internalConfig: validConfig.sort(this.compare('index')),
                invalidConfig: invalidConfig.length === 0 ? null : [...invalidConfig],
                edited: true,
                showInternalAdvancedConfigDialog: false,
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
    private updateConfig(configArr: ReadonlyArray<ConfigItem>, revertConfig: ReadonlyArray<ConfigItem>, defaultConfig: ReadonlyArray<ConfigItem>, terminal: Terminal) {
        // 记录非法修改的默认参数数组
        let invalidConfig: ReadonlyArray<ConfigItem> = []
        // 还原非法修改后的合法项
        let validConfig: ReadonlyArray<ConfigItem> = [...configArr]

        revertConfig.forEach((revert) => {
            // 从高级配置返回的数据中查找是否包含某一具体的默认配置项revert
            const result = defaultConfig.filter((item) => item.configName === revert.configName)

            if (result.length === 0) {
                // 默认配置项被删除
                invalidConfig = [...invalidConfig, { ...revert }]
                validConfig = [...validConfig, { ...revert }]
            } else {
                // 从高级配置项中找到了默认项，表明没有被删除，检查类型
                if (revert.configType !== result[0].configType) {
                    // 类型被修改，还原数据
                    // 如果是服务端配置的syncModule的值不为BaseSyncer时（因为该值不能被修改），还原数据
                    invalidConfig = [...invalidConfig, { ...revert }]
                    validConfig = validConfig.map((item) => {
                        return item.configName === revert.configName ? revert : item
                    })
                }
            }
        })

        return { invalidConfig, validConfig }
    }

    /**
     * 数组排序
     */
    private compare(p) {
        return (m, n) => {
            return m[p] - n[p]
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
}