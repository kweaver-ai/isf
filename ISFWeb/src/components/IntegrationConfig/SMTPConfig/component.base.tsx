import * as React from 'react';
import { getSMTPConfig, testSMTPServer, setSMTPConfig } from '@/core/thrift/sharemgnt';
import { positiveInteger, mailAndLenth } from '@/util/validators';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import { encrypt } from '@/core/auth';
import WebComponent from '../../webcomponent';
import { ValidateState, SafeMode, Port, TestStatus } from './helper';
import __ from './locale';

export default class SMTPConfigBase extends WebComponent<Console.SMTPConfig.Props, Console.SMTPConfig.State> {

    initialConfigInfo: Console.SMTPConfig.ConfigInfo = {
        server: '',
        safeMode: SafeMode.Default,
        port: Port.Default,
        email: '',
        password: '',
        openRelay: false,
    }

    state = {
        configInfo: {
            ...this.initialConfigInfo,
        },
        isFormChanged: false,
        testStatus: TestStatus.NoStart,
        isTestSuccess: false,
        testError: null,
        saveError: null,
        validateState: {
            server: ValidateState.OK,
            port: ValidateState.OK,
            email: ValidateState.OK,
            password: ValidateState.OK,
        },
    }

    password: string = '';
    pwdInput: HTMLInputElement = null;

    async componentDidMount() {
        try {
            const { server, safeMode, port, email, password, openRelay } = await getSMTPConfig();
            const configInfo = {
                server: server ? server : '',
                safeMode: safeMode ? safeMode : SafeMode.Default,
                port: port ? port : Port.Default,
                email: email ? email : '',
                password: password ? password : '',
                openRelay: openRelay,
            }

            this.setState({
                configInfo: {
                    ...configInfo,
                },
            })
            this.password = configInfo.password;
            this.initialConfigInfo = { ...configInfo };

        } catch (error) {
            this.setState({
                saveError: error,
            })
        }
    }

    componentWillUnmount() {
        this.pwdInput = null;
        // 组件销毁后设置state，防止内存泄漏
        this.setState = (state, callback) => {
            return;
        };
    }

    /**
     * 输入值
     * @param key 对应的输入框
     * @param configInfo 输入的值
     */
    protected handleChange(key: string, configInfo: object = {}) {
        const { validateState, testStatus } = this.state;
        this.resetMessage(testStatus);
        this.setState({
            isFormChanged: true,
            testError: null,
            saveError: null,
            configInfo: { ...this.state.configInfo, ...configInfo },
            validateState: {
                server: key === 'server' ? ValidateState.OK : validateState.server,
                port: key === 'port' ? ValidateState.OK : validateState.port,
                email: key === 'email' ? ValidateState.OK : validateState.email,
                password: key === 'password' ? ValidateState.OK : validateState.password,
            },
        })
    }

    /**
     * 下拉列表选择值
     * @param val 选择列表的值
     */
    protected selectChangeHandler({detail: val}) {
        const { configInfo, testStatus, validateState } = this.state;
        switch (val) {
            case SafeMode.SslOrTsl:
                this.setState({
                    configInfo: {
                        ...configInfo,
                        safeMode: SafeMode.SslOrTsl,
                        port: Port.SslOrTsl,
                    },
                })
                break;
            case SafeMode.Starttls:
                this.setState({
                    configInfo: {
                        ...configInfo,
                        safeMode: SafeMode.Starttls,
                        port: Port.Starttls,
                    },
                })
                break;
            default:
                this.setState({
                    configInfo: {
                        ...configInfo,
                        safeMode: SafeMode.Default,
                        port: Port.Default,
                    },
                })
                break;
        }

        this.resetMessage(testStatus);
        this.setState({
            isFormChanged: true,
            testError: null,
            saveError: null,
            validateState: {
                ...validateState,
                port: ValidateState.OK,
            },
        })
    }

    /**
     * 验证输入框输入的合法性
     * @return  boolean 验证的结果
     */
    protected validateCheck(): boolean {
        const { server, port, email, password, openRelay } = this.state.configInfo;
        const serverValidity = server && /^[\w,\.,\-,_,@]{3,100}$/i.test(server);
        const portValidity = port && positiveInteger(port) && Number(port) < 65536 && String(port).indexOf('.') === -1;
        const emailValidity = email && mailAndLenth(email, 4, 101);
        const passwordValidity = openRelay ? true : password;

        if (serverValidity && portValidity && emailValidity && passwordValidity) {
            return true

        } else {
            this.setState({
                validateState: {
                    server: server ? (serverValidity ? ValidateState.OK : ValidateState.ServerError) : ValidateState.Empty,
                    port: port ? (portValidity ? ValidateState.OK : ValidateState.PortError) : ValidateState.Empty,
                    email: email ? (emailValidity ? ValidateState.OK : ValidateState.EmailError) : ValidateState.Empty,
                    password: passwordValidity ? ValidateState.OK : ValidateState.Empty,
                },
            })

            return false;
        }
    }

    /**
     * 切换open relay开关
     */
    protected switchOpenRelay() {
        const { configInfo, testStatus } = this.state;
        this.resetMessage(testStatus);
        this.setState({
            configInfo: {
                ...configInfo,
                openRelay: !configInfo.openRelay,
            },
            isFormChanged: true,
            testError: null,
            saveError: null,
        })
    }

    /**
     * 清空密码框
     */
    protected cleanPassword() {
        const { configInfo, testStatus, isFormChanged } = this.state;
        this.resetMessage(testStatus);
        this.setState({
            configInfo: {
                ...configInfo,
                password: '',
            },
            isFormChanged: isFormChanged ? true : (configInfo.password ? true : false),
            testError: null,
            saveError: null,
        })
    }

    /**
     * 测试服务器
     */
    protected async testHandler() {
        let { configInfo, testStatus } = this.state;
        if (this.validateCheck()) {

            try {
                if (testStatus === TestStatus.Testing) {
                    this.setState({
                        testStatus: TestStatus.NoStart,
                    }, () => {
                        setTimeout(() => {
                            this.setState({
                                testStatus: TestStatus.Testing,
                            })
                        }, 20)
                    })
                } else {
                    this.setState({
                        testStatus: TestStatus.Testing,
                    })
                }
                this.setState({
                    testError: null,
                    saveError: null,
                })

                let smtpConfig = {
                    ...configInfo,
                    port: Number(configInfo.port),
                    password: configInfo.password === this.password ?
                        this.password
                        : encrypt(configInfo.password),
                }

                if (smtpConfig.openRelay) {
                    delete smtpConfig.password;
                }

                const ncTSmtpSrvConf = smtpConfig;
                await testSMTPServer([ncTSmtpSrvConf])
                this.setState({
                    isTestSuccess: true,
                })

            } catch (error) {
                this.setState({
                    isTestSuccess: false,
                    testError: error,
                })
            }

            this.setState({
                testStatus: TestStatus.Tested,
            })
        }
    }

    /**
     * 保存服务器信息
     */
    protected async saveHandler() {
        if (this.validateCheck()) {
            try {
                let { configInfo, testStatus } = this.state;
                this.resetMessage(testStatus);
                this.setState({
                    testError: null,
                    saveError: null,
                })

                const encryptPwd = configInfo.password === this.password ?
                    this.password
                    : encrypt(configInfo.password);
                let smtpConfig = {
                    ...configInfo,
                    port: Number(configInfo.port),
                    password: encryptPwd,
                }

                if (configInfo.openRelay) {
                    delete smtpConfig.password;
                }
                const ncTSmtpSrvConf = smtpConfig;
                await setSMTPConfig([ncTSmtpSrvConf]);

                if (configInfo.openRelay) {
                    this.setState({
                        configInfo: {
                            ...configInfo,
                            password: '',
                        },
                    })
                }
                this.initialConfigInfo = {
                    ...configInfo,
                    password: configInfo.openRelay ? '' : configInfo.password,
                }
                this.password = configInfo.openRelay ? '' : this.password;
                this.setState({
                    isFormChanged: false,
                })

                // 记录日志
                manageLog(
                    ManagementOps.SET,
                    __('设置 SMTP服务器 成功'),
                    __(
                        configInfo.openRelay ?
                            '邮件服务器地址“${server}”；安全连接“${safeMode}”；端口“${port}”；Open Relay“开启”；邮箱地址“${email}”'
                            : '邮件服务器地址“${server}”；安全连接“${safeMode}”；端口“${port}”；Open Relay“关闭”；邮箱地址“${email}”；邮箱密码“******”',
                        {
                            server: configInfo.server,
                            safeMode: this.formateSafeMode(configInfo.safeMode),
                            port: configInfo.port,
                            email: configInfo.email,
                        }),
                    Level.INFO,
                );
            } catch (error) {
                this.setState({
                    saveError: error,
                })
            }
            this.initValidateState()
        }
    }

    /**
     * 格式化安全连接
     */
    protected formateSafeMode(safeMode) {
        switch (safeMode) {
            case SafeMode.SslOrTsl:
                return 'SSL/TLS'
            case SafeMode.Starttls:
                return 'STARTTLS'
            default:
                return __('无')
        }
    }

    /**
     * 取消修改
     */
    protected cancalHandler() {
        this.setState({
            isFormChanged: false,
            isTestSuccess: false,
            testError: null,
            saveError: null,
            configInfo: {
                ...this.initialConfigInfo,
            },
        })
        this.initValidateState();
    }

    /**
     * 初始化验证信息
     */
    protected initValidateState() {
        this.setState({
            validateState: {
                server: ValidateState.OK,
                port: ValidateState.OK,
                email: ValidateState.OK,
                password: ValidateState.OK,
            },
        })
    }

    /**
     * 重置错误重置信息
     * @param testStatus 测试状态
     */
    protected resetMessage(testStatus: number) {
        if (testStatus === TestStatus.Tested) {
            this.setState({
                testStatus: TestStatus.NoStart,
            })
        }
    }
}