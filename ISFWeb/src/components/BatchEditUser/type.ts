import { ValidateState } from '@/core/user';
import { Range } from '../helper';

/**
 * 组件进度
 */
export enum Status {
  /**
   * 创建设置
   */
  Config,

  /**
    * 创建进度
    */
  Progress,

  /**
    * 取消批量编辑确认
    */
  Confirm,

  /**
    * 编辑关闭
    */
  Close,
}

/**
* 密级数组内对象
*/
interface CsfOption {
  /**
  * 密级级别
  */
  value: number;

  /**
  * 密级文本
  */
  name: string;
}

export interface BatchEditUserProps {
  /**
   * 当前选择的用户
   */
  users?: ReadonlyArray<any>;

  /**
   * 选中的部门
   */
  dep?: any;

  /**
   * 批量编辑取消事件
   */
  onRequestCancel: () => void;

  /**
   * 批量编辑成功事件
   */
  onRequestSuccess: ( range: Range ) => void;
}

/**
 * 正在设置的用户名
 */
interface UserName {
  /**
   * 登录名
   */
  loginName: string;

  /**
   * 显示名
   */
  displayName: string;
}

export interface BatchEditUserState {
  /**
   * 选择的用户
   */
  selected: Range;

  /**
   * 有效期限
   */
  expireTime?: number;

  /**
   * 页面加载进度
   */
  status: number;

  /**
   * 密级
   */
  csfLevel?: number;

  /**
   * 密级
   */
  csfLevel2?: number;

  /**
   * 密级列表
   */
  csfOptions: Array<CsfOption>;
  /**
   * 密级2列表
   */
  csfOptions2: Array<CsfOption>;
  /**
   * 是否显示密级2
   */
  show_csf_level2: boolean;

  /**
   * 密级必选框状态
   */
  csfIsChecked: boolean;
  /**
   * 密级2必选框状态
   */
  csf2IsChecked: boolean;

  /**
   * 有效期必选框
   */
  expIsChecked: boolean;

  /**
   * 设置进度条进度
   */
  progress: number;

  /**
   * 当前正在设置的用户名
   */
  currentUserName: UserName;

  /**
   * 错误提示信息索引
   */
  csfValidateState: ValidateState;
  /**
   * 密级2错误提示信息索引
   */
  csf2ValidateState: ValidateState;
}