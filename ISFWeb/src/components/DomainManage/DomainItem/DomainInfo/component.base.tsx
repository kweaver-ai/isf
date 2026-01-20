import * as React from 'react'
import { trim, noop } from 'lodash'
import { addDomain, editDomain } from '@/core/thrift/sharemgnt/sharemgnt'
import { Message2, Toast } from '@/sweet-ui'
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode'
import { isDomainName, IP, IPV6 } from '@/util/validators';
import { encrypt } from '@/core/auth'
import { manageLog, Level, ManagementOps } from '@/core/log'
import WebComponent from '../../../webcomponent'
import { ActionType, DomainType, ValidateStatus } from '../../helper'
import __ from './locale'

/**
 * 验证类型
 */
export enum VerifyType {
    /**
     * 域名
     */
    DomainName,

    /**
     * ip
     */
    DomainIP,

    /**
     * 端口
     */
    DomainPort,

    /**
     * 账号
     */
    Account,

    /**
     * 密码
     */
    Password,
}

interface DomainInfoProps {
    /**
     * 渲染类型
     */
    actionType: ActionType;

    /**
     * 选中项
     */
    selection: Core.ShareMgnt.ncTUsrmDomainInfo;

    /**
     * 编辑项
     */
    editDomain: Core.ShareMgnt.ncTUsrmDomainInfo;

    /**
     * 传递ref
     */
    onRef: (child: object) => void;

    /**
     * 新建或编辑域信息成功回调
     */
    onSetDomainInfoSuccess: (DomainInfo) => void;

    /**
     * 向父组件传递编辑状态
     */
    onRequestEditStatus: (status: boolean) => void;

    /**
     * 域错误回调
     */
    onRequestDomainInvalid: (id) => void;
}

interface DomainInfoState {
    /**
     * 域控信息
     */
    domainInfo: {
        id: number;
        domainType: DomainType;
        domainName: string;
        domainIP: string;
        useSSL: boolean;
        domainPort: number | string;
        account: string;
        password: string;
    };

    /**
     * 账号密码错误提示
     */
    accountError: boolean;

    /**
     * 域控信息保存取消按钮是否显示
     */
    isDomainInfoEditStatus: boolean;

    /**
     * 验证状态
     */
    validateStatus: {
        domainNameValidateStatus: ValidateStatus;
        domainIPValidateStatus: ValidateStatus;
        domainPortValidateStatus: ValidateStatus;
        accountValidateStatus: ValidateStatus;
        passwordValidateStatus: ValidateStatus;
    };

    /**
     * 是否显示域名输入框
     */
    showDomainNameInput: boolean;
}

export default class DomainInfoBase extends WebComponent<DomainInfoProps, DomainInfoState> {
    static defaultProps = {
        actionType: ActionType.Add,
        selection: {},
        editDomain: {},
        onRef: noop,
        onSetDomainInfoSuccess: noop,
        onRequestEditStatus: noop,
    }

    state: DomainInfoState = {
        /**
        * 域控信息
        */
        domainInfo: {
            id: this.props.actionType === ActionType.Edit ? this.props.editDomain.id : this.props.selection.id ? this.props.selection.id : - 1,
            domainType: this.props.actionType === ActionType.Edit ? this.props.editDomain.type : this.props.selection.id ? DomainType.Sub : DomainType.Primary,
            domainName: this.props.actionType === ActionType.Edit ? this.props.editDomain.name : '',
            domainIP: this.props.actionType === ActionType.Edit ? this.props.editDomain.ipAddress : '',
            useSSL: this.props.actionType === ActionType.Edit ? this.props.editDomain.useSSL : false,
            domainPort: this.props.actionType === ActionType.Edit ? this.props.editDomain.port : 389,
            account: this.props.actionType === ActionType.Edit ? this.props.editDomain.adminName : '',
            password: this.props.actionType === ActionType.Edit ? this.props.editDomain.password : '',
        },

        /**
         * 账号密码错误提示
         */
        accountError: false,

        /**
         * 域控信息保存取消按钮是否显示
         */
        isDomainInfoEditStatus: this.props.actionType === ActionType.Edit ? false : true,

        /**
         * 验证状态
         */
        validateStatus: {
            domainNameValidateStatus: ValidateStatus.Normal,
            domainIPValidateStatus: ValidateStatus.Normal,
            domainPortValidateStatus: ValidateStatus.Normal,
            accountValidateStatus: ValidateStatus.Normal,
            passwordValidateStatus: ValidateStatus.Normal,
        },

        /**
         * 域名输入框显示
         */
        showDomainNameInput: this.props.actionType === ActionType.Edit ? false : true,
    }

    /**
    * 编辑时存储原始密码
    */
    password = this.props.actionType === ActionType.Edit ? this.props.editDomain.password : ''

    /**
     * 初始域控信息
     */
    originDomainInfo = this.state.domainInfo

    componentDidMount() {
        this.props.onRef(this)
    }

    /**
     * 改变域类型
     */
    protected changeDomainType = (type: DomainType): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                domainType: type,
            },
            isDomainInfoEditStatus: true,
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 改变域名
     */
    protected changeDomainName = ({ detail }: { detail: string }): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                domainName: detail,
            },
            isDomainInfoEditStatus: true,
            accountError: false,
            validateStatus: {
                ...this.state.validateStatus,
                domainNameValidateStatus: ValidateStatus.Normal,
            },
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 改变域IP
     */
    protected changeDomainIP = ({ detail }: { detail: string }): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                domainIP: detail,
            },
            isDomainInfoEditStatus: true,
            accountError: false,
            validateStatus: {
                ...this.state.validateStatus,
                domainIPValidateStatus: ValidateStatus.Normal,
            },
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 改变域端口
     */
    protected changeDomainPort = ({ detail }: { detail: number | string }): void => {
        if (detail !== this.state.domainInfo.domainPort) {
            this.setState({
                domainInfo: {
                    ...this.state.domainInfo,
                    domainPort: detail,
                },
                isDomainInfoEditStatus: true,
                accountError: false,
                validateStatus: {
                    ...this.state.validateStatus,
                    domainPortValidateStatus: ValidateStatus.Normal,
                },
            })
            this.props.onRequestEditStatus(true)
        }
    }

    /**
     * 改变ssl
     */
    protected changeSSL = ({ detail }: { detail: boolean }): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                useSSL: detail,
                domainPort: detail ? 636 : 389,
            },
            accountError: false,
            validateStatus: {
                ...this.state.validateStatus,
                domainPortValidateStatus: ValidateStatus.Normal,
            },
            isDomainInfoEditStatus: true,
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 聚焦清空密码
     */
    protected clearPassword = (): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                password: '',
            },
            accountError: false,
            isDomainInfoEditStatus: true,
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 改变账号
     */
    protected changeAccount = ({ detail }: { detail: string }): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                account: detail,
            },
            isDomainInfoEditStatus: true,
            accountError: false,
            validateStatus: {
                ...this.state.validateStatus,
                accountValidateStatus: ValidateStatus.Normal,
            },
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 改变密码
     */
    protected changePassword = ({ detail }: { detail: string }): void => {
        this.setState({
            domainInfo: {
                ...this.state.domainInfo,
                password: detail,
            },
            isDomainInfoEditStatus: true,
            accountError: false,
            validateStatus: {
                ...this.state.validateStatus,
                passwordValidateStatus: ValidateStatus.Normal,
            },
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 校验域控信息
     */
    protected verifyDomainInfo = (verifyType?: VerifyType): boolean => {
        const { domainInfo: { domainName, domainIP, domainPort, account, password } } = this.state,
            trimDomainName = trim(domainName),
            trimDomainIP = trim(domainIP),
            trimDomainPort = trim(domainPort),
            domainNameCheckResult = trimDomainName && isDomainName(String(trimDomainName)),
            domainIPString = String(trimDomainIP),
            isDomain = /^.*[a-zA-Z]+.*$/.test(domainIPString) && !domainIPString.includes(':'),
            domainIPCheckResult = trimDomainIP && (isDomain ? isDomainName(String(trimDomainIP)) : (IP(trimDomainIP) || IPV6(trimDomainIP))),
            domainPortCheckResult = trimDomainPort && trimDomainPort >= 1 && trimDomainPort <= 65535,
            accountCheckResult = trim(account) !== '',
            passwordCheckResult = password !== '';

        switch (verifyType) {
            case VerifyType.DomainName:
                if (domainNameCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateStatus: {
                            ...this.state.validateStatus,
                            domainNameValidateStatus: trimDomainName ? ValidateStatus.InvalidDomainName : ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.DomainIP:
                if (domainIPCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateStatus: {
                            ...this.state.validateStatus,
                            domainIPValidateStatus: trimDomainIP ? isDomain ? ValidateStatus.InvalidDomainName : ValidateStatus.InvalidDomainIP : ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.DomainPort:
                if (domainPortCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateStatus: {
                            ...this.state.validateStatus,
                            domainPortValidateStatus: trimDomainPort ? ValidateStatus.InvalidDomainPort : ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.Account:
                if (accountCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateStatus: {
                            ...this.state.validateStatus,
                            accountValidateStatus: ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.Password:
                if (passwordCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateStatus: {
                            ...this.state.validateStatus,
                            passwordValidateStatus: ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            default:
                if (domainNameCheckResult && domainIPCheckResult && domainPortCheckResult && accountCheckResult && passwordCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateStatus: {
                            ...this.state.validateStatus,
                            domainNameValidateStatus: trimDomainName ? isDomainName(String(trimDomainName)) ? ValidateStatus.Normal : ValidateStatus.InvalidDomainName : ValidateStatus.Empty,
                            domainIPValidateStatus: domainIPCheckResult ? ValidateStatus.Normal : trimDomainIP ? /^.*[a-zA-Z]+.*$/.test(String(trimDomainIP)) ? ValidateStatus.InvalidDomainName : ValidateStatus.InvalidDomainIP : ValidateStatus.Empty,
                            domainPortValidateStatus: trimDomainPort !== '' ? domainPort >= 1 && domainPort <= 65535 ? ValidateStatus.Normal : ValidateStatus.InvalidDomainPort : ValidateStatus.Empty,
                            accountValidateStatus: accountCheckResult ? ValidateStatus.Normal : ValidateStatus.Empty,
                            passwordValidateStatus: passwordCheckResult ? ValidateStatus.Normal : ValidateStatus.Empty,
                        },
                    })
                    return false
                }
        }
    }

    /**
     * 域控信息保存
     */
    protected saveDomainInfo = async (): Promise<void> => {
        if (this.verifyDomainInfo()) {
            const { actionType, selection, editDomain: { parentId } } = this.props;
            const { domainInfo: { id, domainType, domainName, domainIP, useSSL, domainPort, account, password } } = this.state,
                name = trim(domainName),
                ipAddress = trim(domainIP),
                adminName = trim(account),
                setParams = {
                    ncTUsrmDomainInfo: {
                        id: actionType === ActionType.Add && (id === -1 || selection.id) ? -1 : id,
                        parentId: actionType === ActionType.Add ? this.props.selection.id ? this.props.selection.id : - 1 : parentId,
                        name,
                        ipAddress,
                        useSSL,
                        port: domainPort,
                        adminName,
                        password: password === this.password ? password : encrypt(password),
                        status: true,
                        type: domainType,
                    },
                }

            try {
                const newId = actionType === ActionType.Add && (id === -1 || selection.id) ? await addDomain([setParams]) : await editDomain([setParams]);

                actionType === ActionType.Add && (id === -1 || selection.id) ?
                    manageLog(
                        ManagementOps.CREATE,
                        __('新建 域控 “${name}” 成功', { name }),
                        '',
                        Level.INFO,
                    ) :
                    manageLog(
                        ManagementOps.SET,
                        __('编辑 域控 “${name}” 成功', { name }),
                        '',
                        Level.INFO,
                    )

                this.setState({
                    domainInfo: {
                        ...this.state.domainInfo,
                        domainName: name,
                        domainIP: ipAddress,
                        account: adminName,
                        password,
                        id: newId || id,
                    },
                    isDomainInfoEditStatus: false,
                    showDomainNameInput: false,
                })

                this.originDomainInfo = this.state.domainInfo

                this.password = setParams.ncTUsrmDomainInfo.password;

                actionType === ActionType.Add && id === -1 ? Toast.open(__('新建成功，可进行以下配置')) : null

                this.props.onRequestEditStatus(false)

                this.props.onSetDomainInfoSuccess({ id: newId || id, parentId, name, adminName, password: actionType === ActionType.Add ? encrypt(password) : password, useSSL, ipAddress, domainPort, isDomainInfoEditStatus: false })
            } catch (ex) {
                this.handleError(ex)
            }
        }
    }

    /**
     * 域控信息取消
     */
    protected cancelDomainInfo = (): void => {
        this.setState({
            domainInfo: {
                ...this.originDomainInfo,
            },
            accountError: false,
            isDomainInfoEditStatus: this.props.actionType === ActionType.Edit || this.state.domainInfo.id !== -1 ? false : true,
            validateStatus: {
                ...this.state.validateStatus,
                domainNameValidateStatus: ValidateStatus.Normal,
                domainIPValidateStatus: ValidateStatus.Normal,
                domainPortValidateStatus: ValidateStatus.Normal,
                accountValidateStatus: ValidateStatus.Normal,
                passwordValidateStatus: ValidateStatus.Normal,
            },
        })
        this.props.onRequestEditStatus(false)
    }

    /**
     * 错误处理函数
     */
    private handleError = async (ex): Promise<void> => {
        switch (ex.error.errID) {
            case ErrorCode.DomainUnavailable:
                Message2.info({ message: __('连接LDAP服务器失败，请检查域控制器地址是否正确，或者域控制器是否开启。'), zIndex: 10000 })
                break

            case ErrorCode.AccountOrPwdInError:
                this.setState({
                    accountError: true,
                })
                break

            case ErrorCode.DomainAlreadyExists:
                this.setState({
                    validateStatus: {
                        ...this.state.validateStatus,
                        domainNameValidateStatus: ValidateStatus.DomainNameExist,
                    },
                })
                break

            case ErrorCode.DomainNotExists: {
                const { actionType, editDomain, selection } = this.props;
                const { domainName } = this.state.domainInfo;
                this.props.onRequestDomainInvalid(actionType === ActionType.Add ? selection.id ? selection.name : domainName : editDomain.name)
                break
            }

            case ErrorCode.SpareAddressDuplicateWithMainDomain:
                this.setState({
                    validateStatus: {
                        ...this.state.validateStatus,
                        domainIPValidateStatus: ValidateStatus.DuplicateWithSpareDomain,
                    },
                })
                break

            case ErrorCode.DomainNameIncorrect:
                this.setState({
                    validateStatus: {
                        ...this.state.validateStatus,
                        domainNameValidateStatus: ValidateStatus.AddressNotMatchIP,
                    },
                })
                break
        }
    }
}