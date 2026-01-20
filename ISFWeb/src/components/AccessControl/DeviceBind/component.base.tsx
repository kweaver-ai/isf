import * as React from 'react'
import { trim, assign } from 'lodash'
import { ShareMgnt, EACP } from '@/core/thrift'
import { getHeaders } from '@/core/token'
import { manageLog, Level, ManagementOps } from '@/core/log'
import * as WebUploader from '@/libs/webuploader'
import { isMac } from '@/util/validators';
import { Message } from '@/sweet-ui';
import WebComponent from '../../webcomponent'
import __ from './locale'
import AppConfigContext from '@/core/context/AppConfigContext'

export enum SearchField {
    AllUser = 1,
    BindUser,
    NoBindUser,
}

export enum OsType {
    UnKnown,
    IOS,
    Android,
    WindowsPhone,
    Windows,
    MacOS,
    Web,
}

export enum ValidateStates {
    Normal,
    InvalidMac,
}

export default class DeviceBindBase extends WebComponent<any, any> {
    static  contextType = AppConfigContext;
    state = {
        currentSelectUser: null, // 当前选中用户
        scope: SearchField.AllUser, // select搜索范围
        searchKey: '', // 搜索用户名关键字
        searchudidKey: '',
        searchResults: [], // 搜索结果
        deviceInfos: [], // 用户绑定设备信息
        disabled: true, // 禁用用户设备面板
        addingDevice: false,  // 添加Mac的Dialog显示隐藏标识
        deviceIsExist: false, // 错误提示框显示隐藏标识
        validateState: ValidateStates.Normal,
        addBox: {
            osType: 4, // 设备类型
            mac: '', // mac地址
        },
        uploading: false, // 批量导入用户设备
        loadingUsers: true, // 懒加载->加载用户列表
        loadingDevices: false, // 懒加载->加载设备列表
        errorlist: null,
        clearing: false,
    }

    UserListLimit = 50;

    DeviceListLimit = 50;

    uploader = null;

    async componentDidMount() {
        // 初始化查询
        const { scope, searchKey } = this.state
        const searchResults = await ShareMgnt('Devicem_SearchUsersBindStatus', [scope, searchKey, 0, this.UserListLimit])

        this.setState({
            searchResults,
            loadingUsers: false,
        })

        // 初始化上传组件
        this.initWebUpload()
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return;
        };
    }

    /**
     * 初始化上传组件
     */
    initWebUpload = () => {
        const self = this;
        self.uploader = new WebUploader.Uploader({
            swf: '/res/libs/webuploader/Uploader.swf',
            server: `${self.context?.prefix || ''}/isfweb/api/device/batchimport/`,
            auto: true,
            threads: 1,
            duplicate: true,
            pick: {
                id: self.refs.batchImport,
                multiple: false,
            },
            accept: {
                extensions: 'csv',
            },
            timeout: 0,
            onBeforeFileQueued: function (_file) {
                self.uploader.reset();
            },
            onUploadBeforeSend: function (object, data, headers) {
                assign(headers, getHeaders(self.context?.getToken?.()).headers);
                data.filename = object.file.name;
                data.userId = self.state.currentSelectUser.id;
                data.osType = 0
            },
            onStartUpload: function (_file) {
                self.setState({
                    uploading: true,
                })
            },
            onUploadSuccess: async function (file, response) {
                const { currentSelectUser, searchudidKey } = self.state;
                self.setState({
                    uploading: false,
                    deviceInfos: searchudidKey ?
                        await EACP('EACP_SearchDevicesByUserIdAndUdid', [currentSelectUser.id, searchudidKey, 0, self.DeviceListLimit]) :
                        await EACP('EACP_GetDevicesByUserId', [currentSelectUser.id, 0, self.DeviceListLimit]),
                    errorlist: response.length ? response : null,
                })

                self.updateCurrentUser()
            },
            onUploadError: function () {
                self.setState({
                    uploading: false,
                });
                Message.alert({ message: __('导入失败。') });
            },
            onUploadAccept: function (object, ret) {
                if (ret.error) {
                    self.setState({
                        uploading: false,
                    });
                    Message.alert({ message: ret.error.errMsg });
                }
            },
            onError: function (type) {
                if (type === 'Q_TYPE_DENIED') {
                    Message.alert({ message: __('请选择扩展名为 csv 的文件。') });
                }
            },
        })
    }

    /**
     * 根据用户名和用户id
     * @param params
     * @param params.name 要查找的用户名
     * @param params.id 要查找的用户id
     */
    private findUser = async ({ name, id }: { name: string; id: string }) => {
        let users = []

        // 如果是所有用户，则只通过搜索接口得到第1条所有用户的数据
        if (id === '-2') {
            users = await ShareMgnt('Devicem_SearchUsersBindStatus', [1, '', 0, 0])
        } else {
            users = await ShareMgnt('Devicem_SearchUsersBindStatus', [1, name, 0, -1])
        }

        return users.find((user) => user.id === id)
    }

    /**
     * 更新当前选中的用户状态
     */
    private updateCurrentUser = async () => {
        const { currentSelectUser, searchResults } = this.state

        if (currentSelectUser) {
            const { loginName, id } = currentSelectUser;
            const nextUserInfo = await this.findUser({ name: loginName, id })

            if (nextUserInfo) {
                this.setState({
                    searchResults: searchResults.map((user) => user.id === id ? nextUserInfo : user),
                    currentSelectUser: nextUserInfo,
                })
            }
        }
    }

    /**
     * 搜索用户
     * 执行搜索并处理搜索结果
     * @memberof DeviceBindBase
     */
    updateUsers = async ({ searchKey = this.state.searchKey, scope = this.state.scope }) => {
        const { currentSelectUser } = this.state
        const searchResults = await ShareMgnt('Devicem_SearchUsersBindStatus', [scope, searchKey, 0, this.UserListLimit])
        let selectnew;
        if (currentSelectUser) {
            selectnew = searchResults.find((searchResult) => searchResult.id === currentSelectUser.id)
        }
        this.setState({
            searchResults,
            currentSelectUser: selectnew,
        })
    }

    /**
     * 更新User搜索关键字
     */
    protected updateUserSearchKey = (searchKey: string) => {
        this.setState({ searchKey })
    }

    /**
     * 搜索关键字改变触发搜索
     * @param {string} searchKey
     * @memberof DeviceBindBase
     */
    handleSearchKeyChange(searchKey: string) {
        this.updateUsers({ searchKey })
    }

    /**
     * 更新UDID
     */
    protected updateUDIDSearchKey = (searchudidKey: string) => {
        this.setState({ searchudidKey })
    }

    /**
     * 搜索用户设备
     * @param udid
     */
    async handleSearchUdidChange(udid) {
        const { currentSelectUser } = this.state;

        if (currentSelectUser) {
            this.setState({
                deviceInfos: await EACP('EACP_SearchDevicesByUserIdAndUdid', [currentSelectUser.id, trim(udid), 0, this.DeviceListLimit]),
            });
        }

    }

    /**
     * 搜索范围发生改变触发搜索
     * @param {number} scope
     * @memberof DeviceBindBase
     */
    handleChangeSearchField({ detail: scope }: { detail: number }) {
        this.setState({
            scope,
        })

        this.updateUsers({ scope, searchKey: this.state.searchKey })
    }

    /**
     * 选中用户后解除设备绑定面板禁用，请求所选中用户设备绑定信息并进行渲染
     * @param {any} UserInfoItem
     * @memberof DeviceBindBase
     */
    handleSelectUser({ detail: nextSelectedUser }) {
        const { currentSelectUser } = this.state;

        this.setState({
            currentSelectUser: nextSelectedUser,
        })

        if (!nextSelectedUser) {
            this.setState({
                disabled: true,
                deviceInfos: [],
                searchudidKey: '',
            })
        }
        else {
            this.setState({
                disabled: false,
            })

            if (!currentSelectUser) {
                this.updateDevices(nextSelectedUser.id);
            } else {
                // 切换了选择的用户
                // 清空udid，并更新用户下的设备
                if (nextSelectedUser.id !== currentSelectUser.id) {
                    this.setState({
                        searchudidKey: '',
                    })

                    this.updateDevices(nextSelectedUser.id);
                }
            }
        }
    }

    /**
     * 获取选中用户设备信息
     * @memberof DeviceBindBase
     */
    async updateDevices(id: string) {
        this.setState({
            deviceInfos: await EACP('EACP_GetDevicesByUserId', [id, 0, this.DeviceListLimit]),
        });
    }

    /**
     * 切换绑定状态
     * @param {any} data
     * @memberof DeviceBindBase
     */
    async switchDeviceBind(data) {
        const { deviceInfos, currentSelectUser } = this.state
        const { bindFlag, baseInfo: { udid: mac } } = data;
        const index = deviceInfos.indexOf(data)
        const nextDeviceInfos = [...deviceInfos.slice(0, index), { ...data, bindFlag: bindFlag ^ 1 }, ...deviceInfos.slice(index + 1)]

        await EACP(bindFlag === 1 ? 'EACP_UnbindDevice' : 'EACP_BindDevice', [currentSelectUser.id, mac])

        this.setState({
            deviceInfos: nextDeviceInfos,
        })

        await manageLog(
            ManagementOps.SET,
            __(bindFlag === 1 ? '用户${userName} 绑定设备 ${device} 成功' : '用户${userName} 解除设备 ${device}的绑定 成功', { userName: currentSelectUser.id, device: mac }),
            '',
            Level.WARN,
        )

        this.updateCurrentUser()
    }

    /**
     * 删除设备绑定并重新获取用户列表
     * @param {any} data
     * @memberof DeviceBindBase
     */
    async deleteDeviceBind(data) {
        const { deviceInfos, currentSelectUser } = this.state
        const nextDeviceInfos = deviceInfos.filter((item) => item.baseInfo.udid !== data.baseInfo.udid)

        await EACP('EACP_DeleteDevices', [currentSelectUser.id, [data.baseInfo.udid]])

        this.setState({
            deviceInfos: nextDeviceInfos,
        })

        this.updateCurrentUser()
    }

    /**
     * 点击添加显示Mac绑定的Dialog
     * @memberof DeviceBindBase
     */
    handleAddDevice() {
        this.setState({
            addingDevice: true,
            addBox: {
                osType: 4,
                mac: '',
            },
        })
    }

    /**
     * 点击取消按钮或者X隐藏Mac绑定的Dialog
     * @memberof DeviceBindBase
     */
    handleCancelAddDevice() {
        this.setState({
            addingDevice: false,
            validateState: ValidateStates.Normal,
            addBox: {
                osType: 4,
                mac: '',
            },
        })
    }

    /**
     * 处理设备类型改变和Mac地址输入改变
     * @param {any} [value={}]
     * @memberof DeviceBindBase
     */
    handleAddBoxChange(value = {}) {
        this.setState({
            validateState: ValidateStates.Normal,
            addBox: { ...this.state.addBox, ...value },
        })
    }

    /**
     * 处理添加设备绑定
     * @memberof DeviceBindBase
     */
    async handleSubmitAddDevice() {
        const { currentSelectUser, addBox } = this.state
        if (isMac(addBox.mac)) {
            try {
                await EACP('EACP_AddDevice', [currentSelectUser.id, addBox.mac, addBox.osType]);
                const deviceInfonew = await EACP('EACP_GetDevicesByUserId', [currentSelectUser.id, 0, this.DeviceListLimit])
                this.setState({
                    deviceInfos: deviceInfonew,
                    addingDevice: false,
                    validateState: ValidateStates.Normal,
                    addBox: {
                        osType: 4,
                        mac: '',
                    },
                });

                this.updateCurrentUser()
            } catch ({ error }) {
                if (error.errID === 4197) {
                    this.setState({
                        deviceIsExist: true,
                    })
                }
            }
        } else {
            this.setState({
                validateState: ValidateStates.InvalidMac,
            })
        }
    }

    /**
     * 关闭错误提示框
     * @memberof DeviceBindBase
     */
    handleCancelTip() {
        this.setState({
            deviceIsExist: false,
        })
    }

    /**
     * 用户列表懒加载回调
     * @param detail 分页信息
     */
    handleRequestLoadUsers({ detail }) {
        const { scope, searchKey, searchResults } = this.state;
        const { limit, start } = detail
        this.setState({
            loadingUsers: true,
        }, async () => {
            const more = await ShareMgnt('Devicem_SearchUsersBindStatus', [scope, searchKey, start, limit])
            this.setState({
                searchResults: [...searchResults, ...more],
                loadingUsers: false,
            })
        })
    }

    /**
     * 设备列表懒加载回调
     * @param detail 分页信息
     */
    handleRequestLoadDevice({ detail }) {
        const { deviceInfos, currentSelectUser, searchudidKey } = this.state;

        /**
         * 当未选中用户时，不加载数据
         */
        if (!currentSelectUser.id) {
            return
        }

        const { limit, start } = detail
        this.setState({
            loadingDevices: true,
        }, async () => {
            let result = [];
            if (searchudidKey) {
                result = await EACP('EACP_SearchDevicesByUserIdAndUdid', [currentSelectUser.id, searchudidKey, start, limit])
            } else {
                result = await EACP('EACP_GetDevicesByUserId', [currentSelectUser.id, start, limit])
            }
            this.setState({
                deviceInfos: [...deviceInfos, ...result],
                loadingDevices: false,
            })
        })
    }

    async handleClearDevices() {
        const { currentSelectUser } = this.state

        this.setState({
            clearing: false,
        })

        await EACP('EACP_DeleteDevices', [currentSelectUser.id, []]);

        this.setState({
            deviceInfos: [],
        })

        this.updateCurrentUser()
    }
}