import * as React from 'react';
import { includes } from 'lodash';
import { getConfidentialConfig } from '@/core/apis/eachttp/config/config';
import { getSystemProtectionLevel, ProtLevelPassExpire } from '@/core/systemprotectionlevel'
import { ShareMgnt } from '@/core/thrift';
import { getThirdMessage } from '@/core/apis/console/thirdMessage';
import { updateDefaultConfigs, checkDefaultPwd } from '@/core/apis/console/usermanagement'
import { Configs } from '@/core/apis/console/usermanagement/types'
import { jsencrypt2048 } from '@/core/auth'
import { ErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import { Message } from '@/sweet-ui';
import { getPloicyInfo, setPwdStrengthMeter } from '@/core/apis/console/loginsecuritypolicy';
import WebComponent from '../../../webcomponent';
import { PwdValidity, PwdStrength, ValidateState, rederErrorMsg, MessageTypes, pluginParamConf, hiddenPwd } from './helper'
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';

interface PasswordPolicyState {
    /**
     * 密码策略配置信息
     */
    configInfo: ConfigInfo;

    /**
     * 密码策略配置
     */
    pwdPolicy: PwdPolicy;

    /**
     * 是否屏蔽密码错误锁定
     */
    lockStatusDisabled: boolean;

    /**
     * 是否选择通过短信验证进行密码重置
     */
    isSMSReset: boolean;

    /**
     * 是否选择通过邮箱验证进行密码重置
     */
    isEmailReset: boolean;

    /**
     * 短信验证是否配置
     */
    isSMSConfig: boolean;

    /**
     * 邮箱验证是否配置
     */
    isEmailConfig: boolean;

    /**
     * 表单是否改变
     */
    isChanged: boolean;

    /**
     * 是否开启三权分立
     */
    isTriSystem: boolean;

    /**
     * 最长有效期天数
     */
    timeCtrl: number;

    /**
     * 是否屏蔽【忘记密码重置】-【通过短信验证】选项
     */
    resetViaSMSDisabled: boolean;

    /**
     * 数字输入框验证
     */
    validateState: {
        /**
         * 强密码长度
         */
        strongPwdLength: ValidateState;

        /**
         * 密码错误次数
         */
        passwdErrCnt: ValidateState;

        /**
         * 锁定时间
         */
        passwdLockTime: ValidateState;

        /**
         * 初始化密码
         */
        initPwdState: ValidateState;
    };

    /**
     * 是否正在加载
     */
    loading: boolean;
}

interface ConfigInfo {
    /**
     * 初始密码
     */
    initPwd: string;

    /**
     * 密码有效期
     */
    expireTime: number;

    /**
     * 密码错误锁定启用状态
     */
    lockStatus: boolean;

    /**
     * 密码强度策略是否被锁定（文档域策略同步）
     */
    pwdLockStatus: boolean;

    /**
     * 密码连续输错次数
     */
    passwdErrCnt: number;

    /**
     * 密码解锁时间
     */
    passwdLockTime: number;

    /**
     * 强密码长度
     */
    strongPwdLength: number;

    /**
     * 密码强度
     */
    strongStatus: PwdStrength;
}

interface PwdPolicy {
    /**
     * 是否屏蔽弱密码选项
     */
    weak_pwd_disabled: boolean;

    /**
     * 强密码最小长度
     */
    min_strong_pwd_length: number;

    /**
     * 密码错误最大次数
     */
    max_err_count: number;
}

export default class PasswordPolicyBase extends WebComponent<any, PasswordPolicyState> {
    static contextType = AppConfigContext
    initPwdConfig = {
        configInfo: {
            expireTime: PwdValidity.None,
            lockStatus: false,
            pwdLockStatus: false,
            passwdErrCnt: 5,
            passwdLockTime: 60,
            strongPwdLength: 8,
            strongStatus: PwdStrength.None,
            initPwd: hiddenPwd,
        },
        isSMSReset: false,
        isEmailReset: false,
    }

    state: PasswordPolicyState = {
        ...this.initPwdConfig,
        lockStatusDisabled: false,
        isSMSConfig: false,
        isEmailConfig: false,
        isChanged: false,
        isTriSystem: false,
        timeCtrl: -1,
        resetViaSMSDisabled: true,
        validateState: {
            strongPwdLength: ValidateState.Normal,
            passwdErrCnt: ValidateState.Normal,
            passwdLockTime: ValidateState.Normal,
            initPwdState: ValidateState.Normal,
        },
        loading: true,
        pwdPolicy: {
            max_err_count: 99,
            min_strong_pwd_length: 8,
            weak_pwd_disabled: false,
        },
    }

    /**
     * 错误信息
     */
    errMsg = {}

    async componentDidMount() {
        try {
            const resultInfo = await this.getInterfaceData();
            const resetViaSMSDisabled = await getConfidentialConfig('reset_via_SMS_disabled');
            const lockStatusDisabled = await getConfidentialConfig('lockdown_criteria_disabled');

            this.setState({
                isChanged: false,
                resetViaSMSDisabled,
                ...resultInfo,
                loading: false,
                lockStatusDisabled,
            })
            this.getServerPluginStatus();
            this.initPwdConfig = {
                configInfo: {
                    ...resultInfo.configInfo,
                },
                isSMSReset: resultInfo.isSMSReset,
                isEmailReset: resultInfo.isEmailReset,
            }
        } catch (error) {
            rederErrorMsg(error);
        }
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
     * 获取接口返回的数据
     * @return 结果类型为对象的Promise
     */
    protected async getInterfaceData(): Promise<object> {
        const configInfo = await ShareMgnt('Usrm_GetPasswordConfig');
        const pwdPolicy = await getConfidentialConfig('pwd_policy');

        let locked, value

        const passwordPolicyInfo = await getPloicyInfo({ mode: 'current', name: 'password_strength_meter' })
        const level = await getSystemProtectionLevel()

        try {
            ({ locked = false, value = { enable: false, length: 10 } } = passwordPolicyInfo.data[0]);
        } catch (ex) {
            throw null
        }
        const vcodeServerStatus = await ShareMgnt('GetCustomConfigOfString', ['vcode_server_status']);
        const isTriSystem = await ShareMgnt('Usrm_GetTriSystemStatus');

        return {
            configInfo: {
                ...configInfo,
                pwdLockStatus: locked,
                strongStatus: value.enable ? PwdStrength.StrongPwd : PwdStrength.WeakPwd,
                strongPwdLength: value.length,
                initPwd: hiddenPwd,
            },
            pwdPolicy,
            timeCtrl: ProtLevelPassExpire[level],
            isTriSystem,
            isSMSReset: JSON.parse(vcodeServerStatus).send_vcode_by_sms,
            isEmailReset: JSON.parse(vcodeServerStatus).send_vcode_by_email,
        }
    }

    /**
     * 判断短信/邮箱是否配置服务器
     */
    protected async getServerPluginStatus() {
        const { isSMSReset, isEmailReset } = this.state;
        const { server } = await ShareMgnt('SMTP_GetConfig');

        const pluginParamConfItem = pluginParamConf[MessageTypes.ResetPWDVerificationCode];

        let isSMSPluginExist = false;

        if (isSMSReset) {
            isSMSPluginExist = (await getThirdMessage()).some(({ config, channels }) => {
                return includes(channels, MessageTypes.ResetPWDVerificationCode) && config[pluginParamConfItem.key] === pluginParamConfItem.value;
            })
        }

        this.setState({
            isSMSConfig: isSMSReset && !isSMSPluginExist,
            isEmailConfig: !server && isEmailReset,
        })
    }

    /**
     * 更改初始密码
     */
    protected updateInitPwd(initPwd: string) {
        this.setState({
            configInfo: {
                ...this.state.configInfo,
                initPwd,
            },
            validateState: {
                ...this.state.validateState,
                initPwdState: ValidateState.Normal,
            },
            isChanged: true,
        })
    }

    /**
     * 选择密码有效期
     * @param detail number 选项值
     */
    protected selectPwdValidity({ detail }: { detail: number }) {
        const { configInfo } = this.state;
        this.setState({
            configInfo: {
                ...configInfo,
                expireTime: detail,
            },
            isChanged: true,
        })
    }

    /**
     * 选择密码强度
     * @param detail number 密码强度选项
     */
    protected selectPwdStrength({ detail }: { detail: number }) {
        const { configInfo, validateState } = this.state;
        this.setState({
            configInfo: {
                ...configInfo,
                strongStatus: detail,
            },
            isChanged: true,
        })
        if (detail === PwdStrength.WeakPwd) {
            this.setState({
                validateState: {
                    ...validateState,
                    strongPwdLength: ValidateState.Normal,
                },
            })
        }
    }

    /**
     * 输入框输入值
     * @param changedConfigInfo 输入框输入标识和值组成的对象
     */
    protected handleChange(changedConfigInfo: object = {}) {
        const { configInfo } = this.state;

        this.setState({
            isChanged: true,
            configInfo: {
                ...configInfo,
                ...changedConfigInfo,
            },
        })
    }

    /**
     * 检测初始密码格式
     */
    protected checkInitPwd = async (): Promise<boolean> => {
        const { initPwd } = this.state.configInfo

        if (!initPwd) {
            this.setState({
                validateState: {
                    ...this.state.validateState,
                    initPwdState: ValidateState.Empty,
                },
            })
            return false
        }

        if (initPwd !== hiddenPwd) {
            try {
                const { result, err_msg } = await checkDefaultPwd({ password: jsencrypt2048(initPwd) })
                if (!result) {
                    this.errMsg = { [ValidateState.InvalidInitPwd]: err_msg }
                    this.setState({
                        validateState: {
                            ...this.state.validateState,
                            initPwdState: ValidateState.InvalidInitPwd,
                        },
                    })
                    return false
                }
            } catch (ex) {
                return false
            }
        }

        return true
    }

    /**
     * 数字输入框失焦校验
     * @param key 对应的输入框键值
     * @param value 填写的值
     */
    protected handleBlurValidate(key, value) {
        const { validateState, pwdPolicy: { min_strong_pwd_length, max_err_count } } = this.state;
        const numberValue = Number(value);
        switch (key) {
            case 'strongPwdLength':
                this.setState({
                    validateState: {
                        ...validateState,
                        strongPwdLength: !value ? ValidateState.Empty
                            : (
                                numberValue < min_strong_pwd_length ? ValidateState.InvalidStrongPwdLengthMin
                                    : (
                                        numberValue > 99 ? ValidateState.InvalidStrongPwdLengthMax :
                                            ValidateState.Normal
                                    )
                            ),
                    },
                })
                break;
            case 'passwdErrCnt':
                this.setState({
                    validateState: {
                        ...validateState,
                        passwdErrCnt: !value ? ValidateState.Empty
                            : (numberValue >= 1 && numberValue <= max_err_count ? ValidateState.Normal
                                : ValidateState.InvalidPasswdErrCnt
                            ),
                    },
                })
                break;
            case 'passwdLockTime':
                this.setState({
                    validateState: {
                        ...validateState,
                        passwdLockTime: !value ? ValidateState.Empty
                            : (numberValue >= 10 && numberValue <= 180 ? ValidateState.Normal
                                : ValidateState.InvalidPasswdLockTime
                            ),
                    },
                })
                break;
            default:
                break;
        }
    }

    /**
     * 开启密码错误锁定
     */
    protected openPwdLock() {
        const { configInfo, validateState } = this.state;

        this.setState({
            validateState: {
                ...validateState,
                passwdErrCnt: ValidateState.Normal,
                passwdLockTime: ValidateState.Normal,
            },
            configInfo: {
                ...configInfo,
                lockStatus: !configInfo.lockStatus,
                passwdErrCnt: configInfo.lockStatus ? this.initPwdConfig.configInfo.passwdErrCnt : configInfo.passwdErrCnt,
                passwdLockTime: configInfo.lockStatus ? this.initPwdConfig.configInfo.passwdLockTime : configInfo.passwdLockTime,
            },
            isChanged: true,
        })
    }

    /**
     * 选择短信验证
     */
    protected selectSMSReset() {
        const { isSMSReset } = this.state;
        this.setState({
            isSMSReset: !isSMSReset,
            isChanged: true,
        })
    }

    /**
     * 选择邮箱验证
     */
    protected selectEmailReset() {
        const { isEmailReset } = this.state;
        this.setState({
            isEmailReset: !isEmailReset,
            isChanged: true,
        })
    }

    /**
     * 跳转到第三方消息集成
     */
    protected openThirdMessage() {
        if(this.context.history && this.context.history.navigateToMicroWidget) {
            this.context.history.navigateToMicroWidget({ name: 'third-party-messaging-plugin' })
        }
    }

    /**
     * 跳转到第三方服务器配置
     */
    protected openThirdPartyServer() {
        if(this.context.history && this.context.history.navigateToMicroWidget) {
            this.context.history.navigateToMicroWidget({ name: 'mailconfig' })
        }
    }

    /**
     * 检查保存时的合法性
     */
    protected async checkSaveValidation(): Promise<boolean> {
        const { pwdPolicy: { min_strong_pwd_length, max_err_count } } = this.state
        const {
            strongPwdLength,
            passwdErrCnt,
            passwdLockTime,
            strongStatus,
        } = this.state.configInfo;

        const initPwdValidity = await this.checkInitPwd()

        const strongPwdLengthValidity = strongPwdLength && strongPwdLength >= min_strong_pwd_length && strongPwdLength <= 99;

        const passwdErrCntValidity = passwdErrCnt && passwdErrCnt >= 1 && passwdErrCnt <= max_err_count;

        const passwdLockTimeValidity = passwdLockTime && passwdLockTime >= 10 && passwdLockTime <= 180;

        if ((strongStatus === PwdStrength.WeakPwd || strongPwdLengthValidity) && passwdErrCntValidity && passwdLockTimeValidity && initPwdValidity) {
            return true;
        } else {
            this.setState({
                validateState: {
                    ...this.state.validateState,
                    strongPwdLength: !strongPwdLength && strongPwdLength !== 0 ?
                        ValidateState.Empty
                        :
                        (
                            strongPwdLength < min_strong_pwd_length ?
                                ValidateState.InvalidStrongPwdLengthMin :
                                (strongPwdLength > 99 ? ValidateState.InvalidStrongPwdLengthMax : ValidateState.Normal)),
                    passwdErrCnt: !passwdErrCnt && passwdErrCnt !== 0 ? ValidateState.Empty
                        : (passwdErrCnt >= 1 && passwdErrCnt <= max_err_count ? ValidateState.Normal
                            : ValidateState.InvalidPasswdErrCnt),
                    passwdLockTime: !passwdLockTime && passwdLockTime !== 0 ? ValidateState.Empty
                        : (passwdLockTime >= 10 && passwdLockTime <= 180 ? ValidateState.Normal
                            : ValidateState.InvalidPasswdLockTime
                        ),
                },
            })
            return false;
        }

    }

    /**
     * 保存设置
     */
    protected async saveConfigInfo() {
        try {
            const { configInfo, isSMSReset, isEmailReset } = this.state;
            const {
                pwdLockStatus,
                strongStatus,
                expireTime,
                lockStatus,
                passwdErrCnt,
                passwdLockTime,
                strongPwdLength,
                initPwd,
            } = this.state.configInfo;

            if (!configInfo.lockStatus) {
                this.setState({
                    configInfo: {
                        ...configInfo,
                        passwdErrCnt: this.initPwdConfig.configInfo.passwdErrCnt,
                    },
                })
            }

            if (await this.checkSaveValidation()) {
                if (initPwd !== hiddenPwd) {
                    await updateDefaultConfigs({ [Configs.DefaultUserPwd]: jsencrypt2048(initPwd) })
                    this.initPwdConfig = {
                        ...this.initPwdConfig,
                        configInfo: {
                            ...configInfo,
                            initPwd: hiddenPwd,
                        },
                    }

                    this.setState({
                        configInfo: {
                            ...configInfo,
                            initPwd: hiddenPwd,
                        },
                    })
                }

                if (!pwdLockStatus) {
                    await setPwdStrengthMeter(
                        {
                            name: 'password_strength_meter',
                            value: {
                                enable: strongStatus ? true : false,
                                length: strongStatus ? strongPwdLength : 8,
                            },
                        },
                    )
                }

                await ShareMgnt('Usrm_SetPasswordConfig', [{
                    ncTUsrmPasswordConfig: {
                        strongStatus: strongStatus,
                        expireTime: Number(expireTime),
                        lockStatus: Number(lockStatus),
                        passwdErrCnt: Number(passwdErrCnt),
                        passwdLockTime: Number(passwdLockTime),
                        strongPwdLength: Number(strongPwdLength),
                    },
                }]);

                await ShareMgnt('SetCustomConfigOfString', [
                    'vcode_server_status',
                    JSON.stringify({
                        send_vcode_by_sms: isSMSReset,
                        send_vcode_by_email: isEmailReset,
                    }),
                ]);
                this.setState({
                    isChanged: false,
                }, () => {
                    // 重新获取短信/邮箱服务器配置状态
                    this.getServerPluginStatus();
                })

                // 记录日志
                this.recordManageLog()

                this.initPwdConfig = {
                    configInfo: {
                        ...configInfo,
                        passwdErrCnt,
                        passwdLockTime,
                        initPwd: hiddenPwd,
                    },
                    isSMSReset: isSMSReset,
                    isEmailReset: isEmailReset,
                }
                Message.info({ message: __('保存成功') });
            }

        } catch (ex) {
            if (ex.error && ex.error.errID === ErrorCode.InvalidSecretPasswdErr) {
                this.setState({
                    validateState: {
                        ...this.state.validateState,
                        passwdErrCnt: ValidateState.InvalidPasswdErrCnt,
                    },
                })
            } else {
                rederErrorMsg(ex);
            }
        }
    }

    /**
     * 取消设置
     */
    protected cancalConfigInfo() {
        this.setState({
            isChanged: false,
            ...this.initPwdConfig,
            validateState: {
                strongPwdLength: ValidateState.Normal,
                passwdErrCnt: ValidateState.Normal,
                passwdLockTime: ValidateState.Normal,
                initPwdState: ValidateState.Normal,
            },
        })
    }

    /**
     * 记录日志
     */
    protected recordManageLog() {
        const { configInfo, isSMSReset, isEmailReset } = this.state
        const { configInfo: initConfigInfo } = this.initPwdConfig;
        if (configInfo.expireTime !== initConfigInfo.expireTime) {
            manageLog(
                ManagementOps.SET,
                (
                    configInfo.expireTime === -1 ?
                        __('设置 密码时效为 永久有效 成功') :
                        __('设置 密码时效为 ${time}天 成功',
                            {
                                time: configInfo.expireTime,
                            },
                        )
                ),
                null,
                Level.WARN,
            )
        }
        if (configInfo.lockStatus !== initConfigInfo.lockStatus ||
            String(configInfo.passwdErrCnt) !== String(initConfigInfo.passwdErrCnt) ||
            String(configInfo.passwdLockTime) !== String(initConfigInfo.passwdLockTime)) {
            manageLog(
                ManagementOps.SET,
                __('设置 密码策略 成功'),
                (
                    configInfo.lockStatus ?
                        __('启用 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
                            { passwdErrCnt: configInfo.passwdErrCnt, passwdLockTime: configInfo.passwdLockTime })
                        : __('关闭 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
                            { passwdErrCnt: configInfo.passwdErrCnt, passwdLockTime: configInfo.passwdLockTime })
                ),
                Level.WARN,
            )

        }
        // 记录弱密码：strongStatus由true --> false
        if (!configInfo.strongStatus && initConfigInfo.strongStatus) {
            manageLog(
                ManagementOps.SET,
                __('设置 密码策略 成功'),
                __('弱密码'),
                Level.WARN,
            )
        }
        // 记录强密码: (1)strongStatus由false-->true (2)strongStatus不变（true）,strongPwdLength发生变化
        if (configInfo.strongStatus && !initConfigInfo.strongStatus || configInfo.strongPwdLength !== initConfigInfo.strongPwdLength) {
            manageLog(
                ManagementOps.SET,
                __('设置 密码策略 成功'),
                __('强密码，密码长度至少为${strongPwdLength}个字符', { strongPwdLength: configInfo.strongPwdLength }),
                Level.WARN,
            )
        }

        if (this.initPwdConfig.isSMSReset !== isSMSReset) {
            manageLog(
                ManagementOps.SET,
                isSMSReset ? __('启用 忘记密码重置通过短信验证 成功') : __('关闭 忘记密码重置通过短信验证 成功'),
                null,
                Level.WARN,
            )
        }

        if (this.initPwdConfig.isEmailReset !== isEmailReset) {
            manageLog(
                ManagementOps.SET,
                isEmailReset ? __('启用 忘记密码重置通过邮箱验证 成功') : __('关闭 忘记密码重置通过邮箱验证 成功'),
                null,
                Level.WARN,
            )
        }
    }

    /**
     * 获取密码有效期下拉框加载的选项
     * @timeCtrl 有效期
     * @return 选项数组
     */
    protected expireTimeSelected(timeCtrl: number) {
        const selects: Array<number> = [
            PwdValidity.OneDay,
            PwdValidity.TreeDays,
            PwdValidity.SevenDays,
            PwdValidity.OneMonth,
            PwdValidity.TreeMonths,
            PwdValidity.SixMonths,
            PwdValidity.TwelveMonths,
            PwdValidity.Permanent,
        ]

        if (timeCtrl === -1) { return selects }

        return selects.filter((item) => item > 0 && item <= timeCtrl)
    }
}