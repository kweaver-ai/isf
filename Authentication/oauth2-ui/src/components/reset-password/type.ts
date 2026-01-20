export enum VerificationType {
  NONE = "",
  PHONE = "telephone",
  EMAIL = "email",
}

export interface IForgetPasswordProps {
  /**
   * 重置密码成功后的重定向地址
   */
  redirect: string;

  /**
   * oem 配置
   */
  oemconfig?: { [key: string]: any };

  /**
   * 认证配置
   */
  authconfig?: { [key: string]: any };

  /**
   * 错误信息
   */
  error?: any;

  /**
   * 客户端类型
   */
  client_type?: string;

  /**
   * IE浏览器
   */
  isIE?: boolean;
  /**
   * 前缀
   */
  urlPrefix?: string;
}

export enum UserverificationStatusEnum {
  //无效账户
  invalid_account = "invalid_account",
  //用户被禁用
  disable_user = "disable_user",
  //密码找回功能未开启
  unable_pwd_retrieval = "unable_pwd_retrieval",
  //非本地用户
  non_local_user = "non_local_user",
  //密码管控开启
  enable_pwd_control = "enable_pwd_control",
  //密码找回功能可用
  available = "available",
}

export interface verificationValueType {
  telephone?: string;
  email?: string;
}

export interface IForgetPasswordState {
  /**
   * 是否处于获取验证码界面
   */
  isVerifying: boolean;

  /**
   * 验证唯一标识
   */
  verificationId: string;

  /**
   * 验证手机号/邮箱/错误码的状态
   */
  verificationValue: verificationValueType | undefined;

  /**
   * 使用的验证类型
   */
  verificationType: VerificationType;

  /**
   * 是否开启强密码
   */
  strongPasswordStatus: boolean;

  /**
   * 强密码最小长度
   */
  strongPasswordLength: number;

  /**
   * 支持的重置类型
   */
  sendVcodeType: SendVcodeType;

  /**
   * 是否处于账号验证界面
   */
  isUserVerification: boolean;

  /**
   * 用户账号
   */
  account: string;

  /**
   * 更新验证码值
   */
  updateVerificationValue: (value: verificationValueType) => void;

  /**
   * 更新验证唯一标识
   */
  updateVerificationId: (id: string) => void;

  /**
   * 发送验证码成功回调
   */
  sendVcodeSuccess: (uuid: string) => void;

  /**
   * 账号检查成功回调
   */
  checkAccountSuccess: () => void;

  /**
   * 返回账号验证界面
   */
  returnUserVerification: () => void;

  /**
   * 更新账号值
   */
  updateAccount: (account: string) => void;

  /**
   * 更新验证方式
   */
  updateVerificationType: (type: VerificationType) => void;

  /**
   * 返回发送验证码界面
   */
  returnSendVcode: () => void;

  /**
   * 更新密码配置（强密码、重置密码）
   */
  updatePasswordConfig: () => void;

  [key: string]: any;
}

export interface SendVcodeType {
  sendVcodeBySMS: boolean;
  sendVcodeByEmail: boolean;
}

export interface ISendVerificationProps {
  /**
   * 全球化函数
   */
  t: any;
}

export interface IUserVerificationProps {
  /**
   * 返回按钮的重定向地址
   */
  redirect: string;
  /**
   * 全球化函数
   */
  t: any;
}

export interface ISendVerificationState {
  /**
   * 错误状态码
   */
  errorStatus: number;

  /**
   * 发送验证码
   */
  sendVcode: () => void;

  [key: string]: any;
}

export interface IResetPasswordProps {
  /**
   * 重置密码成功后的重定向地址
   */
  redirect: string;

  /**
   * 全球化函数
   */
  t: any;
}

export interface IResetPasswordState {
  /**
   * 输入验证码值
   */
  captcha: string;

  /**
   * 新密码
   */
  newPassword: string;

  /**
   * 确认新密码
   */
  confirmPassword: string;

  /**
   * 倒计时计数
   */
  count: number;

  /**
   * 错误状态码
   */
  errorStatus: number;

  /**
   * 错误信息
   */
  errorInfo: any;

  /**
   * 校验输入
   */
  checkInput: () => boolean;

  /**
   * 重置密码
   */
  reset: () => void;

  /**
   * 发送验证码
   */
  sendVcodeAgain: () => void;

  /**
   * 校验重置密码功能状态
   */
  checkSendVcodeStatus: () => boolean;

  [key: string]: any;
}
