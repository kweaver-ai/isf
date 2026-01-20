import { SafeMode, TestStatus } from './helper';

declare namespace Console {
  namespace SMTPConfig {

    interface Props extends React.Props<any> {

    }
    interface State {

      /**
       * 服务器配置信息
       */
      configInfo: ConfigInfo | null,

      /**
       * 是否显示保存/取消按钮
       */
      isFormChanged: boolean,

      /**
       * 测试进行状态
       */
      testStatus: TestStatus,

      /**
       * 测试结果成功
       */
      isTestSuccess: boolean,

      /**
       * 数据合法性状态
       */
      validateState: ValidateState,

      /**
       * 测试错误
       */
      testError: string,

      /**
       * 保存错误
       */
      saveError: string,
    }

    /**
     * 服务器配置信息对象
     */
    interface ConfigInfo {
      /**
       * 邮件服务器
       */
      server: string,

      /**
       * 安全连接
       */
      safeMode: SafeMode,

      /**
       * 端口
       */
      port: number,

      /**
       * 邮件地址
       */
      emial: string,

      /**
       * 邮件密码
       */
      password: string,

      /**
       * open Relay开关
       */
      openRelay: boolean,
    }

    /**
     * 合法性对象
     */
    interface ValidateState {
      /**
       * 邮件服务器合法性
       */
      server: ValidateState,

      /**
       * 端口合法性
       */
      port: ValidateState,

      /**
       * 邮件地址合法性
       */
      email: ValidateState,

      /**
       * 邮件密码合法性
       */
      password: ValidateState
    }

  }
}