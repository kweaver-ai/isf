import * as React from 'react';
import { noop } from 'lodash';
import { Message2 } from '@/sweet-ui';
import { PublicErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { getConfig } from '@/core/apis/eachttp/auth1';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { range, number } from '@/util/validators';
import { ShareMgnt } from '@/core/thrift';
import { getConfidentialConfig } from '@/core/apis/eachttp/config/config';
import { ProtLevelPassExpire } from '@/core/systemprotectionlevel';
import { setPwdStrengthMeter, setSystemProtectionLevels } from '@/core/apis/console/loginsecuritypolicy';
import { ProtLevel } from '@/core/systemprotectionlevel';
import WebComponent from '../webcomponent';
import { Status, PwdValidity, PwdPolicy, UserCsfLevelType } from './helper';
import __ from './locale';
import { getLevelConfig, setLevelConfig } from '@/core/apis/console/usermanagement';

interface Props {
    // 更新用户信息列表中的密级信息
    onUpdateCsfLevelText: () => void;
}

interface State {
    // 是否初始部署状态
    init: boolean;
    // 是否开启强密码
    strongPassword: boolean;

    // 密码错误次数
    passwdErrCnt: number;

    // 打开自定义密级Dialog
    openCustomSecurity: boolean;

    // 自定义密级列表是否为空
    noCustomedSecurity: boolean;

    // 最长有效期天数
    timeCtrl: number;

    /**
     * 密码策略配置
     */
    pwdPolicy: PwdPolicy;

    // 用户密级是否已初始化
    csfLevelsInit: boolean;

    // 用户密级2是否已初始化
    csfLevels2Init: boolean;

    // 用户密级默认值
    defaultCsfLevel: Array<{value: number; name: string}>;

    // 用户密级2默认值
    defaultCsfLevel2: Array<{value: number; name: string}>;

    // 密级临时值
    tempCsfLevel: Array<any>;

    // 用户密级
    userCsfLevel: Array<any>;
    
    // 密级是否已自定义
    isCustomed: boolean;

    
    // 是否支持密级
    supportCsf: boolean;

    // 是否确认自定义密级
    isConfirmCustSecu: boolean;

    
    // 是否触发自定义密级
    isClickCustomSecur: boolean;

    // 是否开启密码错误锁定
    lockStatus: boolean;

    // 新增的记录值
    addedRecord: string;

    // 密码有效期
    expireTime: number;

        // 是否显示系统保护等级配置
    isShowProtLevel: boolean;

    // 系统保护等级
    protLevel: ProtLevel;

    // 密级信息验证状态
    validateSecuStatus: number;

    // 密码错误锁定次数验证状态
    validatePwdStatus: number;

    // 密码强度错误状态
    validatePwdLenStatus: number;

    // 强密码初始长度
    passwordLength: number;

    // 密码锁定时间
    passwdLockTime: number;

    // 密码锁定时间状态
    validateLockTimeStatus: number;
}

export default class SecurityIntegrationBase extends WebComponent<Props, any> {

    static defaultProps = {
        onUpdateCsfLevelText: noop,
    }

    props: Props;

    state: State = {
        init: false,
        strongPassword: false,
        passwordLength: 8,
        passwdErrCnt: 5,
        passwdLockTime: 60,
        openCustomSecurity: false,
        noCustomedSecurity: true,
        csfLevelsInit: false,
        csfLevels2Init: false,
        defaultCsfLevel: [{ value: 5, name: '非密' }, { value: 6, name: '内部' }, { value: 7, name: '秘密' }, { value: 8, name: '机密' }],
        defaultCsfLevel2: [{ value: 51, name: '公开' }],
        tempCsfLevel: [],
        isCustomed: false,
        supportCsf: false,
        isConfirmCustSecu: false,
        isClickCustomSecur: false,
        lockStatus: false,
        addedRecord: '',
        expireTime: 3,
        isShowProtLevel: false,
        protLevel: ProtLevel.Common,
        validateSecuStatus: Status.OK,
        validatePwdStatus: Status.OK,
        validatePwdLenStatus: Status.OK,
        validateLockTimeStatus: Status.OK,
        pwdPolicy: {
            weak_pwd_disabled: false,
            min_strong_pwd_length: 8,
            max_err_count: 99,
        },
        timeCtrl: -1,
    };

    componentDidMount() {
        getLevelConfig({ fields: 'csf_level_enum,csf_level2_enum,show_csf_level2' }).then(({csf_level_enum, csf_level2_enum, show_csf_level2}) => {
            // 判断是否初始化部署状态
            if (!csf_level_enum.length) {
                Promise.all([
                    ShareMgnt('Usrm_GetPasswordConfig'),
                    getConfidentialConfig('pwd_policy'),
                    getConfidentialConfig('protection_level_init_disabled'),
                ]).then(
                    async ([
                        { expireTime, strongStatus, lockStatus, passwdErrCnt, passwdLockTime, strongPwdLength },
                        pwdPolicy,
                        protection_level_init_disabled,
                    ]) => {
                        let protLevel = ProtLevel.Common

                        // 是否需要配置系统保护等级
                        if (!protection_level_init_disabled) {
                            protLevel = ProtLevel.Classified

                            // 若保护等级为Classified及以上，需要默认选中强密码 且
                            // 只能有强密码一个选项 -> (由涉密字段：pwd_policy控制)
                            strongStatus = true

                            // 强密码最小长度
                            strongPwdLength = pwdPolicy.min_strong_pwd_length

                            // 启用密码错误锁定
                            lockStatus = true

                            // 设置密码有效期
                            if (expireTime === -1 || expireTime > ProtLevelPassExpire[protLevel]) {
                                expireTime = ProtLevelPassExpire[protLevel]
                            }
                        }

                        // 依据当前保护等级，设置密码有效期可设定的范围
                        const timeCtrl = ProtLevelPassExpire[protLevel]

                        this.setState({
                            csfLevelsInit: csf_level_enum.length,
                            csfLevels2Init: csf_level2_enum.length,
                            init: true,
                            expireTime: expireTime,
                            strongPassword: strongStatus,
                            lockStatus: lockStatus,
                            passwordLength: strongPwdLength,
                            passwdErrCnt: passwdErrCnt,
                            passwdLockTime,
                            supportCsf: true,
                            timeCtrl,
                            pwdPolicy,
                            isShowProtLevel: !protection_level_init_disabled,
                            protLevel,
                        }, () => {
                            let userCsfLevel = [{
                                    type: UserCsfLevelType.UserLevel,
                                    label: __('用户密级'),
                                    value: csf_level_enum.length ? csf_level_enum : this.state.defaultCsfLevel,
                                }]

                                userCsfLevel = [...userCsfLevel, ...(show_csf_level2 ? [{type: UserCsfLevelType.UserLevel2, label: __('用户密级2'), value: csf_level2_enum.length ? csf_level2_enum : this.state.defaultCsfLevel2 }]: []) ]

                                this.setState({
                                    tempCsfLevel: userCsfLevel,
                                    userCsfLevel: userCsfLevel,
                                })
                        });
                    })
            }

            this.setState({ init: false })
        })
    }

    /**
 * 选择密码有效期
 */
    selectPasswordExpiration({ detail: expiration }) {
        this.setState({ expireTime: expiration })

    }

    /**
 * 选择密码强度
 */
    selectPasswordStrength({ detail: strength }) {
        strength === 1 ? this.setState({ strongPassword: true }) : this.setState({ strongPassword: false })

    }

    /**
     * 初始化配置
     */
    async setSecuInit() {
        const {
            validatePwdStatus,
            strongPassword,
            passwordLength,
            expireTime,
            lockStatus,
            passwdErrCnt,
            passwdLockTime,
            userCsfLevel,
            csfLevelsInit,
            csfLevels2Init,
            defaultCsfLevel2,
            protLevel,
        } = this.state

        if (validatePwdStatus === Status.OK) {
            try {
                await setPwdStrengthMeter(
                    {
                        name: 'password_strength_meter',
                        value: {
                            enable: strongPassword,
                            length: strongPassword ? Number(passwordLength) : 8,
                        },
                    },
                )
            } catch (ex) {
                if (ex.code && ex.code === PublicErrorCode.BadRequest && ex.detail.policys.length) {
                    await Message2.info({
                        message: __('因受父域登录策略管控，无法执行此操作。请在父域的【多文域管理-策略同步】页面关闭密码强度设置才可正常使用。'),
                    })
                    return
                }
            }

            ShareMgnt('Usrm_SetPasswordConfig', [{
                ncTUsrmPasswordConfig: {
                    strongStatus: strongPassword,
                    expireTime: expireTime,
                    lockStatus: lockStatus,
                    passwdErrCnt: Number(passwdErrCnt),
                    strongPwdLength: Number(passwordLength),
                    passwdLockTime: Number(passwdLockTime),
                },
            }]).then(() => {
                manageLog(ManagementOps.SET, __('设置 密码策略 成功'), strongPassword ? __('强密码') : __('弱密码'), Level.WARN);
                if (strongPassword) {
                    manageLog(ManagementOps.SET, __('设置 密码策略 成功'), __('强密码，密码长度至少为${passwordLength}个字符', { passwordLength: passwordLength }), Level.WARN);
                }
                if (lockStatus) {
                    manageLog(
                        ManagementOps.SET,
                        __('设置 密码策略 成功'),
                        __('启用 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
                            { passwdErrCnt: passwdErrCnt, passwdLockTime: passwdLockTime },
                        ),
                        Level.WARN,
                    );
                } else {
                    manageLog(
                        ManagementOps.SET,
                        __('设置 密码策略 成功'),
                        __('关闭 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
                            { passwdErrCnt: passwdErrCnt, passwdLockTime: passwdLockTime },
                        ),
                        Level.WARN,
                    );
                }
            })

            // 设置系统保护等级
            if (protLevel !== ProtLevel.Common) {
                // 更改保护等级
                await setSystemProtectionLevels(
                    {
                        name: 'system_protection_levels',
                        value: { level: protLevel },
                    },
                )
            }

            if (!csfLevelsInit || !csfLevels2Init) {
                // 查询用户密级
                const levelConfig = await getLevelConfig({ fields: 'csf_level_enum,csf_level2_enum' })
                const {csf_level_enum, csf_level2_enum} = levelConfig

                const csfLevel = userCsfLevel[0]?.value
                const csfLevel2 = userCsfLevel[1]?.value || defaultCsfLevel2
            
                let fields = !csf_level_enum.length ? ['csf_level_enum'] : []
                fields =  !csf_level2_enum.length ? [...fields, 'csf_level2_enum'] : fields
                let param: {fields?: string, csf_level_enum ?: Array<{name: string, value: number}>, csf_level2_enum?: Array<{name: string, value: number}>} = { fields: fields.join(',') }
                param = !csf_level_enum.length ? {...param, csf_level_enum: csfLevel} : param
                param = !csf_level2_enum.length ? {...param, csf_level2_enum: csfLevel2 } : param
                
                await setLevelConfig(param)
                this.setState({ csfLevelsInit: true, csfLevels2Init: true })
            }

            getConfig(null, { useCache: 0 })
            this.props.onUpdateCsfLevelText()
        }
    }

    /**
     * 确认初始化对话框
     */

    async triggerConfirmInit() {
        if (this.state.validatePwdStatus === Status.OK) {
            await this.setSecuInit()
        }
    }

    /**
 * 触发确认自定义密级对话框
 */

    triggerConfirmCust(tempCsfLevel) {
        const tempUserCsfLevel = tempCsfLevel.map((cur) => {
            if(cur.type === UserCsfLevelType.UserLevel) {
               const value = cur.value.map((item, index) => ({ ...item, value: index + 5 }))
               return {...cur, value: value }
            }else {
               const value = cur.value.map((item, index) => ({ ...item, value: index + 51 }))
               return {...cur, value: value }
            }
        })
        this.setState({ tempCsfLevel: tempUserCsfLevel, openCustomSecurity: false, isConfirmCustSecu: true })
    }

    /**
     * 设置自定义密级
     */
    setCustomedSecurity() {
        this.setState({
            userCsfLevel: this.state.tempCsfLevel, 
            isConfirmCustSecu: false,
            isCustomed: true,
            csfLevel: 5,
            validateSecuStatus: Status.OK,
        })
    }

    cancelCustom() {
        this.setState({ isConfirmCustSecu: false, openCustomSecurity: true, validateSecuStatus: Status.OK, isClickCustomSecur: false  })
        
    }

    errCntAdd() {
        if (this.state.lockStatus && number(this.state.passwdErrCnt)) {
            this.setState({ passwdErrCnt: Number(this.state.passwdErrCnt) + 1 }, () => {
                if (!this.passwordLockCnt(this.state.passwdErrCnt)) {
                    this.setState({ validatePwdStatus: Status.COUNT_RANGE_ERROR })
                } else {
                    this.setState({ validatePwdStatus: Status.OK });
                }
            })
        }
    }

    errCntSub() {
        if (this.state.lockStatus && number(this.state.passwdErrCnt)) {
            this.setState({ passwdErrCnt: Number(this.state.passwdErrCnt) - 1 }, () => {
                if (!this.passwordLockCnt(this.state.passwdErrCnt)) {
                    this.setState({ validatePwdStatus: Status.COUNT_RANGE_ERROR })
                } else {
                    this.setState({ validatePwdStatus: Status.OK });
                }
            })
        }
    }

    passwordLengthAdd() {
        if (this.state.strongPassword && number(this.state.passwordLength)) {
            this.setState({ passwordLength: Number(this.state.passwordLength) + 1 }, () => {
                if (!this.passwordLengthRange(this.state.passwordLength)) {
                    if (this.state.passwordLength > 99) {
                        this.setState({ passwordLength: Number(99) })
                    } else {
                        this.setState({ validatePwdLenStatus: Status.COUNT_PWD_RANGE_ERROR })
                    }
                } else {
                    this.setState({ validatePwdLenStatus: Status.OK });
                }
            })
        }
    }

    passwordLengthSub() {
        if (this.state.strongPassword && number(this.state.passwordLength)) {
            this.setState({ passwordLength: Number(this.state.passwordLength) - 1 }, () => {
                if (!this.passwordLengthRange(this.state.passwordLength)) {
                    this.setState({ validatePwdLenStatus: Status.COUNT_PWD_RANGE_ERROR })
                } else {
                    this.setState({ validatePwdLenStatus: Status.OK });
                }
            })
        }
    }

    handleErrCntChange(cnt) {
        this.setState({
            passwdErrCnt: cnt,
        }, () => {
            if (!this.passwordLockCnt(cnt)) {
                this.setState({ validatePwdStatus: Status.COUNT_RANGE_ERROR });
            } else {
                this.setState({ validatePwdStatus: Status.OK });
            }
        })
    }

    handlePwdCntChange(pass) {
        this.setState({
            passwordLength: pass,
        }, () => {
            if (!this.passwordLengthRange(pass)) {
                this.setState({ validatePwdLenStatus: Status.COUNT_PWD_RANGE_ERROR });
            } else {
                this.setState({ validatePwdLenStatus: Status.OK });
            }
        })
    }

    /**
     * 输入框修改密码锁定时间
     */
    protected changePasswdLockTime = (passwdLockTime) => {
        this.setState({
            passwdLockTime,
        }, () => {
            this.setState({
                validateLockTimeStatus: this.passwdLockTime(this.state.passwdLockTime) ? Status.OK : Status.PASSWD_LOCK_TIME_ERROR,
            })
        })
    }

    /**
     * 点击上箭头增加密码锁定时间
     */
    protected addPasswdLockTime = () => {
        if (this.state.lockStatus && number(this.state.passwdLockTime)) {
            this.setState({
                passwdLockTime: Number(this.state.passwdLockTime) + 1,
            }, () => {
                this.setState({
                    validateLockTimeStatus: this.passwdLockTime(this.state.passwdLockTime) ? Status.OK : Status.PASSWD_LOCK_TIME_ERROR,
                })
            })
        }
    }

    /**
     * 点击下箭头减少密码锁定时间
     */
    protected subPasswdLockTime = () => {
        if (this.state.lockStatus && number(this.state.passwdLockTime)) {
            this.setState({
                passwdLockTime: Number(this.state.passwdLockTime) - 1,
            }, () => {
                this.setState({
                    validateLockTimeStatus: this.passwdLockTime(this.state.passwdLockTime) ? Status.OK : Status.PASSWD_LOCK_TIME_ERROR,
                })
            })
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

    /**
     * 选择系统保护等级
     */
    protected selectProLevel = (protLevel: ProtLevel): void => {
        const { expireTime: preExpireTime } = this.state
        const timeCtrl = ProtLevelPassExpire[protLevel]
        let expireTime = preExpireTime

        if (preExpireTime === PwdValidity.Permanent || preExpireTime > timeCtrl) {
            expireTime = timeCtrl
        }

        this.setState({
            protLevel,
            expireTime,
            timeCtrl: ProtLevelPassExpire[protLevel],
        })
    }

    /**
 * 验证密码锁定次数
*/
    passwordLockCnt(input) {
        const { pwdPolicy: { max_err_count } } = this.state
        return range(1, max_err_count)(input)
    }

    /**
 * 验证密码锁定时间
 */
    passwdLockTime(input) {
        return range(10, 180)(input)
    }

    /**
 * 验证密码长度
*/
    passwordLengthRange(input) {
        const { pwdPolicy: { min_strong_pwd_length } } = this.state

        return range(min_strong_pwd_length, 99)(input)
    }

    customSecurityClassification() {
        this.setState({ openCustomSecurity: true, validateSecuStatus: Status.OK, isClickCustomSecur: true })
    }

    closeCustomSecurity() {
        this.setState({ openCustomSecurity: false, noCustomedSecurity: true, isClickCustomSecur: false })
    }

    /**
 * 启用密码错误锁定
 */
    setPasswordLock() {
        this.setState({
            lockStatus: !this.state.lockStatus,
            passwdErrCnt: 5,
            passwdLockTime: 60,
        }, () => {
            if (!this.state.lockStatus) {
                this.setState({
                    validatePwdStatus: Status.OK,
                    validateLockTimeStatus: Status.OK,
                });
            }
        });
    }
}