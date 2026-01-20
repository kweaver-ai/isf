import { noop, trim } from 'lodash';
import * as PropTypes from 'prop-types';
import { getUserInfo, getPwdControl, getPwdConfig, setPwdControl } from '@/core/thrift/sharemgnt';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { encrypt } from '@/core/auth';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { generateRandomPwd } from '@/core/password';
import { getRandom as secureRandom } from '@/util/random'
import { Message, Message2 } from '@/sweet-ui';
import WebComponent from '../webcomponent';
import __ from './locale';

/**
 * 验证密码状态类型
 */
export enum ValidateState {
    /**
     * 正常
     */
    OK,

    /**
    * 空值
    */
    Empty,

    /**
     * 强密码状态下报错
     */
    StrongPwdError,

    /**
     * 弱密码状态下报错
     */
    WeakPwdError,

    /**
     * 输入为初始密码状态下报错
     */
    OrgPwdError,
}

/**
 * 用户类型
 */
export enum UserType {
    /**
     * 本地用户
     */
    LocalUser = 1,

    /**
     * 域用户
     */
    DomainUser = 2,

    /**
     * 第三方用户
     */
    ThirdPartyUser = 3,
}

/**
 * 用户原始的密码配置
 */
interface UserData {
    /**
    * 当前密码
    */
    password: string;

    /**
     * 当前是否启用密码管控
     */
    pwdControl: boolean;

    /**
     * 当前是否锁定
     */
    lockStatus: boolean;

    /**
     * 当前是否初始密码
     */
    originalPwd: boolean;

    /**
     * 当前是否为强密码状态
     */
    strongStatus: boolean;

    /**
     * 用户显示名
     */
    displayName: string;

    /**
     * 登录名
     */
    loginName: string;

    /**
     * 用户类型
     */
    userType: UserType;

    /**
    * 强密码的最小长度
    */
    strongPwdLength: number;
}

interface PwdManageProps {
    /**
     * 选中的用户的信息
     */
    selectedUser: {
        id: string;
        [key: string]: any;
    };

    /**
     * 登录用户Id
     */
    userid: string;

    /**
     * 点击确定按钮回调函数
     */
    onRequestConfirm: () => void;

    /**
    * 点击取消按钮回调函数
    */
    onRequestCancel: () => void;
}

interface PwdManageState {
    /**
     * 是否允许用户自主修改密码
     */
    pwdControl: boolean;

    /**
     * 密码
     */
    password: string;

    /**
     * 解锁状态
     */
    lockStatus: boolean;

    /**
     * 输入框内的密码错误提示
     */
    validateState: ValidateState;
}

/**
 * 初始密码值
 */
export const systemOriginalPwd = '123456'

/**
 * 密码不可见时显示
 */
export const hiddenPwd = '**********'

/**
 * 特殊字符
 */
const SpecialChars = '~!%#$@-_.'

export default class PwdManageBase extends WebComponent<PwdManageProps, PwdManageState> {

    static defaultProps = {
        selectedUser: {
            id: '',
        },
        userid: '',
        onRequestConfirm: noop,
        onRequestCancel: noop,
    }

    static contextTypes = {
        toast: PropTypes.func,
    }

    state = {
        pwdControl: false,
        password: '',
        lockStatus: false,
        validateState: ValidateState.OK,
    }

    userData: UserData = {
        password: '',
        pwdControl: false,
        lockStatus: false,
        originalPwd: false,
        strongStatus: false,
        displayName: '',
        loginName: '',
        userType: UserType.LocalUser,
        strongPwdLength: 10,
    }

    isClickUnlock = false; // 是否点击解锁按钮

    isClickConfirm = false; // 确定按钮可用flag

    isInputPwd = false; // 是否手动输入密码用于区分重置密码的12346和不允许用户管控密码时手动输入123466

    async componentDidMount() {
        const { selectedUser: { id }, userid } = this.props;

        if (id === userid) {
            await Message.info({ message: __('您无法管控自己的密码。') })
            this.props.onRequestCancel && this.props.onRequestCancel();
        } else {
            try {
                const [
                    { originalPwd, user: { userType, displayName, loginName } },
                    { strongStatus, strongPwdLength },
                    { pwdControl, lockStatus, password },
                ] = await Promise.all([
                    getUserInfo([id]),
                    getPwdConfig(),
                    getPwdControl([id]),
                ]);

                this.userData = {
                    password: password ? password : '',
                    pwdControl,
                    lockStatus,
                    originalPwd,
                    strongStatus,
                    userType,
                    displayName,
                    loginName,
                    strongPwdLength,
                }

                this.setState({
                    pwdControl,
                    lockStatus,
                    password: pwdControl ? password ? hiddenPwd : '' : (originalPwd ? systemOriginalPwd : hiddenPwd),
                })
            } catch (err) {
                const { error } = err
                if(error?.errMsg) {
                    await Message2.info({ message: error.errMsg })
                    if(error?.errID === ErrorCode.UserNotExist) {
                        this.props.onRequestCancel && this.props.onRequestCancel();
                    }
                }
            }
        }
    }

    /**
     * 勾选或取消允许用户自主修改密码
     */
    protected toggleCheck(): void {
        const { pwdControl } = this.state,
            { originalPwd, password: userPassword } = this.userData;

        this.isInputPwd = false
        this.setState({
            validateState: ValidateState.OK,
            pwdControl: !pwdControl,
            password: !pwdControl ? userPassword ? hiddenPwd : '' : (originalPwd ? systemOriginalPwd : hiddenPwd),
        })
    }

    /**
     *修改密码
     */
    protected updatePwd(value: string): void {
        this.isInputPwd = true
        this.setState({
            password: value,
            validateState: ValidateState.OK,
        })
    }

    /**
     * 点击重置按钮
     */
    protected resetPwd(): void {
        this.isInputPwd = false
        this.setState({
            password: systemOriginalPwd,
            validateState: ValidateState.OK,
        })
    }

    /**
     * 随机产生密码
     */
    protected handlePwdRandomly(): void {
        const { strongStatus, strongPwdLength } = this.userData;

        const pwdStr = generateRandomPwd(strongStatus ? strongPwdLength - 1 : 9) + SpecialChars[Math.floor(secureRandom(0, 1, 2) * SpecialChars.length)]

        this.setState({
            password: pwdStr.split('').sort(() => secureRandom(0, 1, 1) > 0.5 ? 1 : -1).join(''),
            validateState: ValidateState.OK,
        }, () => {
            this.userData = {
                ...this.userData,
                password: this.state.password,
            }
        })
    }

    /**
     * 解锁账户
     */
    protected unlockUser(): void {
        this.setState({
            lockStatus: false,
        })
        this.isClickUnlock = true;
    }

    /**
     * 验证密码函数
     */
    private validatePwd(): boolean {
        const { strongPwdLength } = this.userData,
            strongPwdReg = new RegExp('^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])[\x20-\x7E]{' + strongPwdLength + ',100}$'),
            weakPwdReg = /^[\x20-\x7E]{6,100}$/,
            { password, pwdControl } = this.state,
            { strongStatus, password: originalPassword } = this.userData,
            trimPassword = trim(password === hiddenPwd ? originalPassword : password);
        if (pwdControl) {
            if (trimPassword === '') {
                this.setState({
                    validateState: ValidateState.Empty,
                })
                return false;
            } else {
                if (strongStatus) {
                    if (!strongPwdReg.test(trimPassword)) {
                        this.setState({
                            validateState: ValidateState.StrongPwdError,
                        })
                        return false;
                    } else {
                        return true
                    }
                } else {
                    if (!weakPwdReg.test(trimPassword)) {
                        this.setState({
                            validateState: ValidateState.WeakPwdError,
                        })
                        return false;
                    } else {
                        return true
                    }
                }
            }
        } else {
            return true
        }
    }

    /**
     * 点击确定按钮
     */
    protected async confirm() {
        if (!this.isClickConfirm) {
            try {
                this.isClickConfirm = true;
                const { selectedUser: { id } } = this.props,
                    { lockStatus, password, pwdControl } = this.state,
                    { pwdControl: originalPwdControl, displayName, loginName, password: originalPassword, lockStatus: originalLockStatus, originalPwd } = this.userData,
                    { lockStatus: currentLockStatus } = await getPwdControl([this.props.selectedUser.id]),
                    trimPassword = trim(password === hiddenPwd ? originalPassword : password);

                if (this.validatePwd()) {
                    const ncTUsrmPwdControlConfig = { lockStatus: this.isClickUnlock ? lockStatus : currentLockStatus, password: originalPwd && trimPassword === '123456' && !this.isInputPwd ? null : encrypt(trimPassword), pwdControl };

                    try {
                        const res = await setPwdControl([id, ncTUsrmPwdControlConfig]);
                        if (res === null) {
                            if (!pwdControl) {
                                manageLog(
                                    ManagementOps.SET,
                                    __('允许 用户“${userName}”自主修改密码，设置密码成功',
                                        { userName: displayName + '(' + loginName + ')' },
                                    ),
                                    null,
                                    Level.WARN,
                                )
                            } else if ((!originalPwdControl && ncTUsrmPwdControlConfig) || ncTUsrmPwdControlConfig) {
                                manageLog(
                                    ManagementOps.SET,
                                    __('不允许 用户“${userName}”自主修改密码，设置密码成功',
                                        { userName: displayName + '(' + loginName + ')' },
                                    ),
                                    null,
                                    Level.WARN,
                                )
                            }

                            if (!originalPwd && password === systemOriginalPwd) {
                                manageLog(
                                    ManagementOps.SET,
                                    __('将用户“${userName}”的密码重置为初始密码 成功',
                                        { userName: displayName + '(' + loginName + ')' },
                                    ),
                                    null,
                                    Level.WARN,
                                )
                            }

                            if (originalLockStatus && !lockStatus) {
                                manageLog(
                                    ManagementOps.SET,
                                    __('解锁 用户“${userName}” 成功',
                                        { userName: displayName + '(' + loginName + ')' },
                                    ),
                                    null,
                                    Level.WARN,
                                )
                            }
                            this.props.onRequestConfirm && this.props.onRequestConfirm();
                        }
                    } catch ({ error }) {

                        // 密码是否是原密码校验改由后端校验
                        if (error.errID === ErrorCode.CannotUseInitPwd) {
                            this.setState({
                                validateState: ValidateState.OrgPwdError,
                            })
                            this.isClickConfirm = false;
                            return
                        }

                        this.isClickConfirm = false;
                        Message.error({ message: __('无效的密码') });
                    }
                } else {
                    this.isClickConfirm = false;
                }
            }catch({ error }) {
                this.isClickConfirm = false;
                if(error?.errMsg) {
                    await Message2.info({ message: error.errMsg })
                    if(error?.errID === ErrorCode.UserNotExist) {
                        this.props.onRequestCancel && this.props.onRequestCancel();
                    }
                }
            }
        }
    }

    /**
     * 点击取消按钮
     */
    protected cancel(): void {
        this.props.onRequestCancel && this.props.onRequestCancel();
    }
}