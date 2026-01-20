import * as React from 'react';
import { map } from 'lodash';
import { getConfidentialConfig } from '@/core/apis/eachttp/config/config'
import { LoginAuthType } from '@/core/logincertification/logincertification';
import { EndpointType } from '@/core/endpointtype/endpointtype';
import { ShareMgnt } from '@/core/thrift';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import { LoginWay } from '@/core/logincertification/logincertification';
import { getVcodeConfig, setVcodeConfig, getCustomConfigOfString, setCustomConfigOfString, getTriSystemStatus } from '@/core/thrift/sharemgnt/sharemgnt';
import { getPloicyInfo, setBatchOSTypeForbidLoginInfo } from '@/core/apis/console/loginsecuritypolicy';
import { Message } from '@/sweet-ui';
import WebComponent from '../../../webcomponent';
import { OstypeInfo, rederErrorMsg, ValidateState, LoginAuthTypeMapper, LoginWays } from './helper';
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';

interface LoginPolicyState {

    /**
     * 身份证账号策略状态
     */
    idCardStatus: boolean;

    /**
     * 禁止登录的客户端
     */
    forbidOstypeInfo: object;

    /**
     *  自动禁用天数
     */
    autoDisableDays: number;

    /**
     * 客户端锁定状态
     */
    clientLockStatus: boolean;

    /**
     * 自动禁用状态
     */
    autoDisableStatus: boolean;

    /**
     * 登录方式
     */
    loginWay: LoginWay;

    /**
     * 允许密码连续输错次数
     */
    passwdErrCnt: any;

    /**
     * 表单项是否改变
     */
    isChanged: boolean;

    /**
     * 密码输错次数限制合法性
     */
    passwdErrCntValidity;

    /**
    * 登录认证锁定状态
    */
    multiAutnLocked: boolean;

    /**
     * 禁止同一帐号多地同时登录
     */
    loginStrategyStatus: boolean;

    /**
     * 是否屏蔽禁止同一账号多地登录
     */
    loginStrategyDisabled: boolean;

    /**
     * 域用户免登录
     */
    domainUserExemptLogin: boolean;

    /**
     * 屏蔽的客户端登陆选项
     */
    disabledEndpointTypes: ReadonlyArray<EndpointType>;
    /**
     * 屏蔽的登录认证选项
     */
    disabledLoginAuth: ReadonlyArray<LoginAuthType>;
    /**
     * 是否正在加载
     */
    loading: boolean;
}

interface LoginPolicyType {
    /**
   * 身份证账号策略状态
   */
    idCardStatus: boolean;

    /**
     * 禁止登录的客户端
     */
    forbidOstypeInfo: object;

    /**
     *  自动禁用天数
     */
    autoDisableDays: number;

    /**
     * 自动禁用状态
     */
    autoDisableStatus: boolean;

    /**
     * 登录方式
     */
    loginWay: LoginWay;

    /**
     * 允许密码连续输错次数
     */
    passwdErrCnt: any;

    /**
     * 客户端锁定状态
     */
    clientLockStatus: boolean;

    /**
    * 登录认证锁定状态
    */
    multiAutnLocked: boolean;

    /**
     * 禁止同一帐号多地同时登录
     */
    loginStrategyStatus: boolean;

    /**
     * 域用户免登录
     */
    domainUserExemptLogin: boolean;
}

export default class LoginPolicyBase extends WebComponent<any, LoginPolicyState> {
    static contextType = AppConfigContext
    /**
     * 初始登录策略配置信息
     */
    initLoginConfig = {
        idCardStatus: false,
        forbidOstypeInfo: {},
        autoDisableDays: 0,
        autoDisableStatus: false,
        loginWay: LoginWay.Account,
        passwdErrCnt: 0,
        clientLockStatus: false,
        multiAutnLocked: false,
        loginStrategyStatus: false,
        domainUserExemptLogin: false,
    }

    state: LoginPolicyState = {
        ...this.initLoginConfig,
        isChanged: false,
        passwdErrCntValidity: ValidateState.Normal,
        disabledEndpointTypes: [],
        disabledLoginAuth: [],
        loading: true,
        loginStrategyDisabled: false,
    }

    originVcodeConfig: Core.ShareMgnt.ncTVcodeConfig;
    triSystemStatus: boolean = false;

    async componentDidMount() {
        try {
            // 获取是否开启三权分立
            this.triSystemStatus = await getTriSystemStatus();
            const loginConfigInfo = await this.getInterfaceData();
            const disabledEndpointTypes = await getConfidentialConfig('disabled_endpoint_types');
            const disabledLoginAuth = await getConfidentialConfig('disabled_login_auth');
            const loginStrategyDisabled = await getConfidentialConfig('different_devices_login_disabled');
            const transferDisabledLoginAuth = disabledLoginAuth.map((auth) => LoginAuthTypeMapper[auth]);
            const allowedAuth = LoginWays.filter((way) => !transferDisabledLoginAuth.includes(way));
            const loginWay = transferDisabledLoginAuth.includes(loginConfigInfo.loginWay) ? allowedAuth[0] : loginConfigInfo.loginWay;
            this.setState({
                ...loginConfigInfo,
                loginWay,
                disabledEndpointTypes,
                disabledLoginAuth,
                loading: false,
                loginStrategyDisabled,
            })
            this.initLoginConfig = { ...loginConfigInfo };

        } catch (e) {
            rederErrorMsg(e);
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
    * @return 结果类型为LoginPolicyType对象的Promise
    */
    protected async getInterfaceData(): Promise<LoginPolicyType> {
        const loginStrategyStatus = await ShareMgnt('GetLoginStrategyStatus');
        const idCardStatus = await ShareMgnt('GetCustomConfigOfBool', ['id_card_login_status']);
        const { isEnabled, days } = await ShareMgnt('Usrm_GetAutoDisable');
        // 获取是否为图形验证码
        this.originVcodeConfig = { ...(await getVcodeConfig()) };

        // 获取所有设备类型禁止登录状态
        const { locked, value } = (await getPloicyInfo({ mode: 'current', name: 'client_restriction' })).data[0];
        // 获取登录认证禁用状态
        const multiAutnLocked = (await getPloicyInfo({ mode: 'current', name: 'multi_factor_auth' })).data[0].locked;

        // 获取域用户免登录状态、
        const domainUserExemptLogin = await ShareMgnt('Usrm_GetADSSOStatus')

        let forbidOstypeInfo = {
            [OstypeInfo.Web]: value.pc_web,
            [OstypeInfo.Windows]: value.windows,
            [OstypeInfo.iOS]: value.ios,
            [OstypeInfo.Mobileweb]: value.mobile_web,
            [OstypeInfo.Mac]: value.mac,
            [OstypeInfo.Android]: value.android,
            [OstypeInfo.Linux]: value.linux,
        };

        let loginWay = LoginWay.Account;
        let passwdErrCnt = 0;

        if (this.originVcodeConfig.isEnable) {
            loginWay = LoginWay.AccountAndImgCaptcha;
            passwdErrCnt = this.originVcodeConfig.passwdErrCnt;
        } else {
            // 获取是否为短信验证码/动态密码
            const authServerStatus = JSON.parse(await getCustomConfigOfString(['dualfactor_auth_server_status']))
            if (authServerStatus.auth_by_sms) {
                loginWay = LoginWay.AccountAndSmsCaptcha;
            } else if (authServerStatus.auth_by_OTP) {
                loginWay = LoginWay.AccountAndDynamicPassword;
            }
            else {
                // 若非短信验证码/动态密码/图形验证码，则判断登录方式为：账号密码
                loginWay = LoginWay.Account;
            }
        }
        return {
            idCardStatus,
            forbidOstypeInfo,
            autoDisableStatus: isEnabled,
            autoDisableDays: days,
            loginWay,
            passwdErrCnt,
            clientLockStatus: locked,
            multiAutnLocked,
            loginStrategyStatus,
            domainUserExemptLogin,
        }
    }

    /**
     * 自动禁用超过时间未登录的用户进行登录
     */
    protected autoDisableLogin() {
        const { autoDisableStatus } = this.state;
        this.setState({
            autoDisableStatus: !autoDisableStatus,
            isChanged: true,
        })
    }

    /**
     * 自动禁用的限制时间
     * @param detail 限制时间
     */
    protected selectAutoDisableDays({ detail }) {
        this.setState({
            autoDisableDays: detail,
            isChanged: true,
        })
    }

    /**
     * 禁止登录选项
     * @param plafform 禁止登录的平台
     */
    protected forbidTerminalLogin(osType: number) {
        const { forbidOstypeInfo } = this.state;
        this.setState({
            forbidOstypeInfo: {
                ...forbidOstypeInfo,
                [osType]: !forbidOstypeInfo[osType],
            },
            isChanged: true,
        })
    }

    /**
     * 禁止同一帐号多地同时登录
     */
    protected setloginStrategyStatus() {
        const { loginStrategyStatus } = this.state;

        this.setState({
            loginStrategyStatus: !loginStrategyStatus,
            isChanged: true,
        })
    }

    /**
     * 设置域用户免登录
     */
    protected setDomainUserExemptLogin = (): void => {
        const { domainUserExemptLogin } = this.state;

        this.setState({
            domainUserExemptLogin: !domainUserExemptLogin,
            isChanged: true,
        })
    }

    /**
     * 允许身份证号作为账号登录
     */
    protected allowIdCardAccount() {
        const { idCardStatus } = this.state;
        this.setState({
            idCardStatus: !idCardStatus,
            isChanged: true,
        })
    }

    /**
     * 选择登录方式
     * @param detal 登录方式
     */
    protected selectLoginWay({ detail: loginWay }) {
        if (loginWay === LoginWay.AccountAndImgCaptcha) {
            this.setState({
                loginWay,
                isChanged: true,
                passwdErrCnt: this.originVcodeConfig.passwdErrCnt,
            })
        } else {
            this.setState({
                loginWay,
                isChanged: true,
                passwdErrCntValidity: ValidateState.Normal,
            })
        }
    }

    /**
     * 密码输错次数失焦校验
     * @param value 填写的值
     */
    protected handleBlurPwdErrCntValidate(value) {
        this.setState({
            passwdErrCntValidity: value > 99 ? ValidateState.InvalidPasswdErrCnt : !value && value !== 0 ? ValidateState.Empty : ValidateState.Normal,
        })
    }

    /**
     * 修改输入密码错误次数
     * @param passwdErrCnt 密码错误次数
     */
    changePwdErrCnt(passwdErrCnt: number | string) {
        this.setState({
            passwdErrCnt,
            isChanged: true,
        })
    }

    /**
     * 跳转到第三方认证
     */
    protected handleJump() {
        if(this.context.history && this.context.history.navigateToMicroWidget) {
            this.context.history.navigateToMicroWidget({ name: 'cert-manage', path: "?tab=thrid-auth" })
        }
    }

    /**
     * 保存登录策略配置信息
     */
    protected async saveConfigInfo() {
        try {
            const {
                idCardStatus,
                forbidOstypeInfo,
                autoDisableDays,
                autoDisableStatus,
                loginWay,
                passwdErrCnt,
                clientLockStatus,
                loginStrategyStatus,
                domainUserExemptLogin,
            } = this.state;
            if (this.checkPasswdErrCntValidity()) {
                await this.saveIdCardStatus();
                if (!clientLockStatus) {
                    await this.saveForbidOstypeInfo();
                }
                await this.saveAutoDisable();
                await this.saveLoginWayConfig();
                await this.saveLoginStrategyStatus();
                await this.saveDomainUserExemptLogin();
                this.setState({
                    isChanged: false,
                })
                this.initLoginConfig = {
                    idCardStatus,
                    forbidOstypeInfo,
                    autoDisableDays,
                    autoDisableStatus,
                    loginWay,
                    passwdErrCnt,
                    clientLockStatus,
                    loginStrategyStatus,
                    domainUserExemptLogin,
                }
                Message.info({ message: __('保存成功') });
            }

        } catch (e) {
            rederErrorMsg(e);
        }
    }

    /**
     * 保存时校验密码输错次数合法性
     */
    protected checkPasswdErrCntValidity(): boolean {
        const { passwdErrCnt, loginWay } = this.state;
        if (loginWay === LoginWay.AccountAndImgCaptcha && !passwdErrCnt && passwdErrCnt !== 0 || passwdErrCnt > 99) {
            this.setState({
                passwdErrCntValidity: passwdErrCnt > 99 ? ValidateState.InvalidPasswdErrCnt : ValidateState.Empty,
            })
            return false;
        } else {
            return true;
        }
    }

    /**
     * 保存身份证账号策略状态
     */
    protected async saveIdCardStatus() {
        const { idCardStatus } = this.state;
        const status = await ShareMgnt('GetCustomConfigOfBool', ['id_card_login_status']);

        if (idCardStatus !== status) {
            await ShareMgnt('SetCustomConfigOfBool', ['id_card_login_status', idCardStatus]);
            manageLog(
                ManagementOps.SET,
                idCardStatus ?
                    __('设置 允许身份证号登录 成功')
                    : __('取消 允许身份证号登录 成功'),
                __(''),
                Level.WARN,
            );
        }
    }

    /**
     * 禁止同一帐号多地同时登录
     */
    protected async saveLoginStrategyStatus() {
        const { loginStrategyStatus } = this.state;
        const status = await ShareMgnt('GetLoginStrategyStatus');

        if (loginStrategyStatus !== status) {
            await ShareMgnt('SetLoginStrategyStatus', [loginStrategyStatus]);
            manageLog(
                ManagementOps.SET,
                loginStrategyStatus ?
                    __('设置 禁止同一账号多地同时登录 成功，（仅对使用Windows客户端登录进行限制，使用其他客户端登录不受限制）')
                    : __('取消设置 禁止同一账号多地同时登录 成功，（仅对使用Windows客户端登录进行限制，使用其他客户端登录不受限制）'),
                __(''),
                Level.WARN,
            );
        }
    }

    /**
     * 保存域用户免登录设置
     */
    protected saveDomainUserExemptLogin = async (): Promise<void> => {
        const { domainUserExemptLogin } = this.state;
        const status = await ShareMgnt('Usrm_GetADSSOStatus');

        if (domainUserExemptLogin !== status) {
            await ShareMgnt('Usrm_SetADSSOStatus', [domainUserExemptLogin]);
            manageLog(
                ManagementOps.SET,
                domainUserExemptLogin ?
                    __('设置 域用户免登录 成功')
                    : __('取消设置 域用户免登录 成功'),
                __(''),
                Level.WARN,
            );
        }
    }

    /**
     * 保存自动禁用设置
     */
    protected async saveAutoDisable() {
        const { autoDisableStatus, autoDisableDays } = this.state;
        if (autoDisableStatus !== this.initLoginConfig.autoDisableStatus || autoDisableDays !== this.initLoginConfig.autoDisableDays) {
            await ShareMgnt('Usrm_SetAutoDisable', [{
                ncTUserAutoDisableConfig: {
                    isEnabled: autoDisableStatus,
                    days: autoDisableDays,
                },
            }]);
            if (autoDisableStatus) {
                manageLog(
                    ManagementOps.SET,
                    __('设置 禁用长期未登录的账号 成功'),
                    __('登录周期${time}个月 自动禁用', { time: parseInt(String(autoDisableDays / 30)) }),
                    Level.INFO,
                );
            } else {
                manageLog(
                    ManagementOps.SET,
                    __('取消 禁用长期未登录的账号 成功'),
                    __(''),
                    Level.INFO,
                );
            }
        }
    }

    /**
     * 批量设置指定设备类型禁止登录状态
     */
    protected async saveForbidOstypeInfo() {
        const { forbidOstypeInfo } = this.state;
        await setBatchOSTypeForbidLoginInfo(
            {
                name: 'client_restriction',
                value: {
                    pc_web: forbidOstypeInfo[OstypeInfo.Web],
                    mobile_web: forbidOstypeInfo[OstypeInfo.Mobileweb],
                    windows: forbidOstypeInfo[OstypeInfo.Windows],
                    ios: forbidOstypeInfo[OstypeInfo.iOS],
                    mac: forbidOstypeInfo[OstypeInfo.Mac],
                    android: forbidOstypeInfo[OstypeInfo.Android],
                    linux: forbidOstypeInfo[OstypeInfo.Linux],
                },
            },
        );
        map(forbidOstypeInfo, (info, key) => {
            // 设置或取消禁止成功,记录一条管理日志
            const index = parseInt(key);
            if (this.initLoginConfig.forbidOstypeInfo[key] !== info) {
                if (info) {
                    if (index === 7) {
                        manageLog(
                            ManagementOps.SET,
                            __('设置 禁止移动Web客户端登录 成功'),
                            '',
                            Level.WARN,
                        );
                    } else {
                        manageLog(
                            ManagementOps.SET,
                            __('设置 禁止${client}登录 成功', { client: OstypeInfo[key] === 'Web' ? __('网页端') : OstypeInfo[key] + __('客户端') }),
                            '',
                            Level.WARN);
                    }
                } else {
                    if (index === 7) {
                        manageLog(
                            ManagementOps.SET,
                            __('取消 禁止移动Web客户端登录 成功'),
                            '',
                            Level.WARN,
                        );
                    } else {
                        manageLog(
                            ManagementOps.SET,
                            __('取消 禁止${client}登录 成功', { client: OstypeInfo[key] === 'Web' ? __('网页端') : OstypeInfo[key] + __('客户端') }),
                            '',
                            Level.WARN,
                        );
                    }
                }
            }
        })

    }

    /**
     * 保存登录方式配置
     */
    protected async saveLoginWayConfig() {
        const { loginWay, passwdErrCnt } = this.state;
        if (loginWay !== this.initLoginConfig.loginWay || passwdErrCnt !== this.initLoginConfig.passwdErrCnt) {
            if (loginWay === LoginWay.AccountAndImgCaptcha && passwdErrCnt !== '') {
                await setVcodeConfig([{
                    ncTVcodeConfig: {
                        isEnable: true,
                        passwdErrCnt: Number(passwdErrCnt),
                    },
                }]);
                await setCustomConfigOfString([
                    'dualfactor_auth_server_status',
                    JSON.stringify({
                        auth_by_sms: false,
                        auth_by_email: false,
                        auth_by_OTP: false,
                        auth_by_Ukey: false,
                    }),
                ])
                this.recordSaveLoginWaylog();
                this.originVcodeConfig = {
                    isEnable: true,
                    passwdErrCnt,
                };
            } else if (loginWay !== LoginWay.AccountAndImgCaptcha) {
                if (this.originVcodeConfig.isEnable) {
                    /**
                     * 若原始选项为图形验证码，需先将图形验证码关闭，再设置对应认证方式
                     */
                    await setVcodeConfig([{
                        ncTVcodeConfig: {
                            isEnable: false,
                            passwdErrCnt: Number(this.originVcodeConfig.passwdErrCnt),
                        },
                    }]);
                    this.originVcodeConfig.isEnable = false;

                }
                // 设置对应认证方式
                await setCustomConfigOfString([
                    'dualfactor_auth_server_status',
                    JSON.stringify({
                        auth_by_sms: loginWay === LoginWay.AccountAndSmsCaptcha,
                        auth_by_email: false,
                        auth_by_OTP: loginWay === LoginWay.AccountAndDynamicPassword,
                        auth_by_Ukey: false,
                    }),
                ]);
                this.recordSaveLoginWaylog();
            }
        }

    }

    /**
     * 记录保存登录方式配置日志
     */
    private recordSaveLoginWaylog() {
        const { loginWay, passwdErrCnt } = this.state;
        switch (loginWay) {
            case LoginWay.Account:
                manageLog(
                    ManagementOps.SET,
                    __('设置 用户必须通过“账号密码”登录 成功'),
                    null,
                    Level.WARN,
                );
                break;
            case LoginWay.AccountAndImgCaptcha:
                manageLog(
                    ManagementOps.SET,
                    __('设置 用户必须通过“账号密码+图形验证码”登录 成功'),
                    __('连续输错密码次数为 ${passwdErrCnt} 次出现登录验证码', { passwdErrCnt }),
                    Level.WARN,
                );
                break;
            case LoginWay.AccountAndSmsCaptcha:
                manageLog(
                    ManagementOps.SET,
                    __('设置 用户必须通过“账号密码+短信验证码”登录 成功'),
                    null,
                    Level.WARN,
                );
                break;
            case LoginWay.AccountAndDynamicPassword:
                manageLog(
                    ManagementOps.SET,
                    __('设置 用户必须通过“账号密码+动态密码”登录 成功'),
                    null,
                    Level.WARN,
                );
                break;
        }
    }

    /**
     * 取消修改
     */
    protected cancalConfigInfo() {
        this.setState({
            isChanged: false,
            ...this.initLoginConfig,
            passwdErrCntValidity: ValidateState.Normal,
        })
    }
}