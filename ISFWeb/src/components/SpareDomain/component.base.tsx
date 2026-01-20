import * as React from 'react'
import { noop, trim } from 'lodash'
import { Message2 as Message } from '@/sweet-ui'
import { encrypt } from '@/core/auth/auth'
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode'
import { manageLog, Level, ManagementOps } from '@/core/log'
import { usrmCheckFailoverDomainAvailable, usrmEditFailoverDomains, usrmGetFailoverDomains } from '@/core/thrift/sharemgnt/sharemgnt'
import { isDomainName, IP, IPV6 } from '@/util/validators';
import { ValidateStatus } from '../DomainManage/helper'
import { VerifyType } from '../DomainManage/DomainItem/DomainInfo/component.base'
import WebComponent from '../webcomponent'
import __ from './locale'

let self

interface DomainInfo {
    /**
     * 域控制器地址
     */
    address: string;

    /**
     * 使用SSL 勾选状态
     */
    useSSL: boolean;

    /**
     * 域控制器端口
     */
    port: string;

    /**
     * 域管理员账号
     */
    account: string;

    /**
     * 域管理员密码
     */
    password: string;

    /**
     * 编辑中
     */
    editable?: boolean;

    /**
     * id
     */
    id: number;

    /**
     * 账号密码是否正确
     */
    accountError: boolean;
}

interface Props {
    /**
     * 主域信息
     */
    mainDomain: {
        /**
         * 域控制器地址
         */
        ipAddress: string;

        /**
         * 主域的id
         */
        id: number;
    };

    /**
     * 传递是否编辑状态
     */
    onRequestEditStatus: (status: boolean) => void;

    /**
     * 保存完成，返回上一步
     */
    onBackLastStep: () => void;

    /**
     * 保存完成，进入下一步
     */
    onJumpNextStep: () => void;

    /**
     * 将上一步、下一步暴露出来
     */
    onExposeJumpNext: (fn: (next: boolean) => any) => any;

    /**
     * 域不存在时回调
     */
    onRequestDomainInvalid: (name: string) => void;
}

interface State {
    /**
     * 输入框的验证状态
     */
    validateState: {
        /**
         * 域控制器地址
         */
        address: number;

        /**
         * 域控制器端口
         */
        port: number;

        /**
         * 域管理员账号
         */
        account: number;

        /**
         * 域管理员密码
         */
        password: number;
    };

    /**
     * 保存的域控信息
     */
    domains: Array<DomainInfo>;

    /**
     * 当前正在编辑的是第几条数据
     */
    editingIndex: number;
}

export default class SpareDomainBase extends WebComponent<Props, State> {

    static defaultProps = {
        mainDomain: null,
        onBackLastStep: noop,
        onJumpNextStep: noop,
        onExposeJumpNext: noop,
        onRequestEditStatus: noop,
        onRequestDomainInvalid: noop,
    }

    state = {
        validateState: {
            address: ValidateStatus.Normal,
            port: ValidateStatus.Normal,
            account: ValidateStatus.Normal,
            password: ValidateStatus.Normal,
        },
        domains: [],
        editingIndex: -1,
    }

    templateDomain: DomainInfo = { // 新建域的模板 -- 根据主域获取
        address: '',
        useSSL: false,
        port: '389',
        account: '',
        password: '',
        editable: true,
        id: -1,
        accountError: false,
    }

    originalDomains: ReadonlyArray<DomainInfo> = []

    originalEncryptPasswords: ReadonlyArray<string> = [] // 获取的初始的密码（已加密的，包括主域、备用域）

    async componentDidMount() {
        // 获取主域的使用SSL、端口、账号、密码
        const { mainDomain: { adminName: account, password, useSSL, port, id } } = this.props

        this.templateDomain = {
            address: '',
            useSSL,
            port,
            account,
            password,
            editable: true,
            id,
            accountError: false,
        }

        // 根据主域id获取备用域信息
        const domains = await usrmGetFailoverDomains([id])
        const formatDomains = domains.map(({ parentId: id, address, port, adminName: account, password, useSSL }) => (
            { id, address, port, account, password, useSSL, editable: false }
        ))
        // 设置originalDomains
        this.originalDomains = formatDomains

        this.setState({
            domains: formatDomains,
        })

        this.props.onExposeJumpNext(this.jumpStep)

        this.originalEncryptPasswords = [password, ...(domains.map(({ password }) => password))]

        self = this
    }

    /**
     * 点击“添加备用域控制器”触发
     */
    protected async addDomain() {
        // 编辑的是第几条数据
        const { editingIndex } = this.state

        if (editingIndex !== -1) {
            try {
                // 保存编辑中的这条域
                await this.testDomain(editingIndex)
                // 添加一条新的默认域
                this.addEmptyDomain()
            } catch (ex) {
            }

        } else {
            // 没有正在编辑的域控，添加一条新的默认域
            this.addEmptyDomain()
        }

        this.props.onRequestEditStatus(true)
    }

    /**
     * 添加一条默认的域
     */
    private addEmptyDomain() {
        // 没有正在编辑的域控
        const newDomains = [
            ...this.state.domains,
            this.templateDomain,
        ]
        this.setState({
            domains: newDomains,
            editingIndex: newDomains.length - 1,
            validateState: {
                address: ValidateStatus.Normal,
                port: ValidateStatus.Normal,
                account: ValidateStatus.Normal,
                password: ValidateStatus.Normal,
            },
        })

        // 添加一条备用域，更新originalDomains
        this.originalDomains = newDomains
    }

    /**
     * 编辑域控制器地址输入框触发
     */
    protected editDomain(index: number, area: 'address' | 'port' | 'account' | 'password', value: string) {
        // 检查一下除了当前正在编辑的这条数据外，还有没有其他正在编辑中的数据
        const { validateState, domains } = this.state

        if (value !== domains[index][area]) {
            const newDomains = domains.map((domain, key) => {
                return key === index ? { ...domain, [area]: value, editable: true, accountError: area === 'account' || area === 'password' ? false : domain.accountError } : domain
            })

            this.setState({
                validateState: {
                    ...validateState,
                    [area]: ValidateStatus.Normal,
                },
                domains: newDomains,
                editingIndex: index,
            })

            this.props.onRequestEditStatus(true)
        }
    }

    /**
     * 删除域控
     */
    protected deleteDomain(index: number) {
        const { domains } = this.state
        const newDomains = domains.filter((item, key) => index !== key)
        // 删除一条备用域，更新originalDomains
        this.originalDomains = newDomains

        this.setState({
            domains: newDomains,
            editingIndex: -1,
        }, async () => {
            const formatDomains = this.state.domains.map(({ id: parentId, address, account: adminName, password, port, useSSL }) => (
                {
                    ncTUsrmFailoverDomainInfo: {
                        id: -1,
                        parentId,
                        address: trim(address),
                        port: Number(port),
                        adminName: trim(adminName),
                        password:  // 如果密码和originalEncryptPasswords中的某个密码一样，则不用加密；反之，要加密
                            this.originalEncryptPasswords.some((pass) => pass === password) ? password : encrypt(password),
                        useSSL,
                    },
                }
            ))

            try {
                await usrmEditFailoverDomains([formatDomains, this.props.mainDomain.id])
                manageLog(
                    ManagementOps.DELETE,
                    __('删除 备用域控 “${name}” 成功', { name: domains[index].address }),
                    '',
                    Level.WARN,
                )
                this.props.onRequestEditStatus(false)
            } catch (ex) {
                this.handleError(ex)
            }
        })
    }

    /**
     * 取消编辑
     */
    protected cancelEdit(index: number) {
        this.setState({
            domains: this.originalDomains,
            editingIndex: this.originalDomains[index].editable ? index : -1,
            validateState: {
                address: ValidateStatus.Normal,
                port: ValidateStatus.Normal,
                account: ValidateStatus.Normal,
                password: ValidateStatus.Normal,
            },
        }, () => {

            this.state.editingIndex === -1 ? this.props.onRequestEditStatus(false) : null
        })
    }

    /**
     * 点击“确定” -- 测试域控信息合法性、是否能连接上
     * @param index 测试第几条数据
     */
    protected testDomain(index: number): Promise<void> {
        return new Promise(async (resolve, reject) => {
            const { domains } = this.state

            // 检测各个输入框的合法性
            if (this.verifyDomainInfo(index)) {
                try {
                    // 检测domian是否能使用
                    const { id: parentId, address, account: adminName, password, port, useSSL } = domains[index]

                    const domainInfo = {
                        id: -1,
                        parentId,
                        address: trim(address),
                        port: Number(port),
                        adminName: trim(adminName),
                        password: this.originalEncryptPasswords.some((pass) => pass === password) ? password : encrypt(password),
                        useSSL,
                    }

                    await usrmCheckFailoverDomainAvailable([[{
                        ncTUsrmFailoverDomainInfo: domainInfo,
                    }]])

                    resolve()

                } catch ({ error: { errID, errMsg } }) {
                    switch (errID) {
                        case ErrorCode.DomainUnavailable:
                            Message.info({ message: __('连接LDAP服务器失败，请检查域控制器地址是否正确，或者域控制器是否开启。'), zIndex: 10000 })
                            break

                        case ErrorCode.AccountOrPwdInError:
                            this.setState({
                                domains: this.state.domains.map((item, id) => id === index ? { ...item, accountError: true } : item),
                            })
                            break

                        case ErrorCode.SpareAddressDuplicateWithMainDomain:
                            this.setState({
                                validateState: {
                                    ...this.state.validateState,
                                    address: ValidateStatus.SpareAddressDuplicateWithMainDomain,
                                },
                            })
                            break

                        case ErrorCode.DomainNotExists:
                            this.props.onRequestDomainInvalid(this.props.mainDomain.name)
                            break

                        case ErrorCode.DomainsNotInOneDomain:
                            this.setState({
                                validateState: {
                                    ...this.state.validateState,
                                    address: ValidateStatus.DomainsNotInOneDomain,
                                },
                            })
                            break

                        default:
                            Message.info({ message: errMsg, zIndex: 10000 })
                    }

                    reject()
                }
            } else {
                reject()
            }
        })
    }

    /**
     * 校验域控信息
     */
    protected verifyDomainInfo = (index: number, verifyType?: VerifyType): boolean => {
        const { address, port, account, password } = this.state.domains[index],
            trimDomainIP = trim(address),
            trimDomainPort = trim(port),
            domainIPString = String(trimDomainIP),
            isDomain = /^.*[a-zA-Z]+.*$/.test(domainIPString) && !domainIPString.includes(':'),
            domainIPCheckResult = trimDomainIP && (isDomain ? isDomainName(String(trimDomainIP)) : (IP(trimDomainIP) || IPV6(trimDomainIP))),
            domainPortCheckResult = trimDomainPort && trimDomainPort >= 1 && trimDomainPort <= 65535,
            accountCheckResult = trim(account) !== '',
            passwordCheckResult = password !== '';

        switch (verifyType) {
            case VerifyType.DomainIP:
                if (domainIPCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateState: {
                            ...this.state.validateState,
                            address: trimDomainIP ? isDomain ? ValidateStatus.InvalidDomainName : ValidateStatus.InvalidDomainIP : ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.DomainPort:
                if (domainPortCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateState: {
                            ...this.state.validateState,
                            port: trimDomainPort ? ValidateStatus.InvalidDomainPort : ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.Account:
                if (accountCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateState: {
                            ...this.state.validateState,
                            account: ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            case VerifyType.Password:
                if (passwordCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateState: {
                            ...this.state.validateState,
                            password: ValidateStatus.Empty,
                        },
                    })
                    return false
                }

            default:
                if (domainIPCheckResult && domainPortCheckResult && accountCheckResult && passwordCheckResult) {
                    return true
                } else {
                    this.setState({
                        validateState: {
                            ...this.state.validateState,
                            address: domainIPCheckResult ? ValidateStatus.Normal : trimDomainIP ? /^.*[a-zA-Z]+.*$/.test(String(trimDomainIP)) ? ValidateStatus.InvalidDomainName : ValidateStatus.InvalidDomainIP : ValidateStatus.Empty,
                            port: trimDomainPort !== '' ? trimDomainPort >= 1 && trimDomainPort <= 65535 ? ValidateStatus.Normal : ValidateStatus.InvalidDomainPort : ValidateStatus.Empty,
                            account: accountCheckResult ? ValidateStatus.Normal : ValidateStatus.Empty,
                            password: passwordCheckResult ? ValidateStatus.Normal : ValidateStatus.Empty,
                        },
                    })
                    return false
                }
        }
    }

    /**
     * 清空密码输入框的内容
     */
    protected clearPassword(index: number) {
        const { domains } = this.state

        const newDomains = domains.map((domain, key) => {
            return key === index ? { ...domain, password: '', accountError: false, editable: true } : domain
        })

        this.setState({
            domains: newDomains,
            editingIndex: index,
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 勾选SSL 发生变化
     */
    protected changeSSL(checked: boolean, index: number) {
        // 勾选 -- 636；取消勾选 -- 389
        let { domains, validateState } = this.state
        const newDomains = domains.map((domain, key) => {
            return (index === key) ? { ...domain, port: checked ? '636' : '389', useSSL: checked, editable: true, accountError: false } : domain
        })

        this.setState({
            domains: newDomains,
            editingIndex: index,
            validateState: {  // 当勾选SSL发生变化，需要清空port的气泡提示
                ...validateState,
                port: ValidateStatus.Normal,
            },
        })
        this.props.onRequestEditStatus(true)
    }

    /**
     * 跳转到上一步或下一步
     * @param nextStep true -- 下一步; false -- 上一步
     */
    public async jumpStep(nextStep: boolean) {
        // 检测是否有正在编辑中的
        const { editingIndex } = self.state

        try {
            if (editingIndex !== -1) {
                // 有正在编辑中的域
                // 保存编辑中的域
                try {
                    await self.testDomain(editingIndex)
                    // 添加备用域
                    await self.addRequestDomains()
                    nextStep ? self.props.onJumpNextStep() : self.props.onBackLastStep()
                } catch (ex) {
                    if (ex) {
                        const { errID, errMsg } = ex.error

                        switch (errID) {
                            case ErrorCode.DomainUnavailable:
                                Message.info({ message: __('连接LDAP服务器失败，请检查域控制器地址是否正确，或者域控制器是否开启。'), zIndex: 10000 })
                                break

                            default:
                                Message.info({ message: errMsg, zIndex: 10000 })
                        }
                    }
                }

            } else {
                // 没有正在编辑中的
                // 添加备用域
                await self.addRequestDomains()
                nextStep ? self.props.onJumpNextStep() : self.props.onBackLastStep()
            }
        } catch (ex) {
            if (ex) {
                const { errID, errMsg } = ex.error

                switch (errID) {
                    case ErrorCode.DomainUnavailable:
                        Message.info({ message: __('连接LDAP服务器失败，请检查域控制器地址是否正确，或者域控制器是否开启。'), zIndex: 10000 })
                        break

                    case ErrorCode.SpareAddressDuplicateWithMainDomain:
                        this.setState({
                            validateState: {
                                ...this.state.validateState,
                                address: ValidateStatus.SpareAddressDuplicateWithMainDomain,
                            },
                        })
                        break

                    case ErrorCode.DomainNotExists:
                        this.props.onRequestDomainInvalid(this.props.mainDomain.name)
                        break

                    default:
                        Message.info({ message: errMsg, zIndex: 10000 })
                }
            }
        }
    }

    /**
     * 点击“确定”按钮，测试备用域是否有用
     */
    protected async clickTestDomain(editingIndex: number) {
        try {
            await this.testDomain(editingIndex)
            await this.addRequestDomains(editingIndex)

            this.props.onRequestEditStatus(false)
        } catch (ex) {
            // this.handleError(ex)
        }
    }

    /**
     * 添加备用域，向后端发送请求
     */
    private addRequestDomains(editingIndex) {
        return new Promise(async (resolve, reject) => {
            try {
                const { domains } = this.state
                const formatDomains = domains.map(({ id: parentId, address, account: adminName, password, port, useSSL }) => (
                    {
                        ncTUsrmFailoverDomainInfo: {
                            id: -1,
                            parentId,
                            address: trim(address),
                            port: Number(port),
                            adminName: trim(adminName),
                            password:  // 如果密码和originalEncryptPasswords中的某个密码一样，则不用加密；反之，要加密
                                this.originalEncryptPasswords.some((pass) => pass === password) ? password : encrypt(password),
                            useSSL,
                        },
                    }
                ))

                const { address, account: adminName } = domains[editingIndex]

                await usrmEditFailoverDomains([formatDomains, this.props.mainDomain.id])

                manageLog(
                    ManagementOps.ADD,
                    __('添加 备用域控 “${name}” 成功', { name: domains[editingIndex].address }),
                    '',
                    Level.INFO,
                )

                // 检测成功，修改domian的editable为false
                const newDomains = domains.map((domain, key) => {
                    return key === editingIndex ? { ...domain, address: trim(address), account: trim(adminName), editable: false, accountError: false } : domain
                })
                // 更新originalDomains
                this.originalDomains = newDomains

                this.setState({
                    domains: newDomains,
                    editingIndex: -1,
                }, () => {
                    resolve(editingIndex)
                })

            } catch (error) {
                this.handleError(error)
                reject(error)
            }
        })
    }

    /**
     * 错误处理函数
     */
    private handleError = async (ex): Promise<void> => {
        switch (ex.error.errID) {
            case ErrorCode.DomainNotExists:
                this.props.onRequestDomainInvalid(this.props.mainDomain.name)
                break

            case ErrorCode.SameSpareDomain:
                this.setState({
                    validateState: {
                        ...this.state.validateState,
                        address: ValidateStatus.SpareDomainExist,
                    },
                })
                break
        }
    }
}