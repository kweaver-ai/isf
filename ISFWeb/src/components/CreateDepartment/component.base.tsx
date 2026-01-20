import * as React from 'react'
import { noop, trim } from 'lodash';
import { Message2, Toast } from '@/sweet-ui';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { addDepartment } from '@/core/thrift/sharemgnt/sharemgnt';
import { Doclibs, ValidateStatus, getValidateInfo } from '@/core/doclibs/doclibs';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { displayUserOssInfo } from '@/core/oss/oss';
import { getObjectStorageInfoById } from '@/core/apis/console/ossgateway'
import { ValidateState, Dep, DepartmentInfo, NodeInfo, UserInfoType } from '@/core/user';
import { mailAndLenth, isNormalName, isNormalCode } from '@/util/validators';
import WebComponent from '../webcomponent';
import __ from './locale';

/**
 * 错误提示索引
 */
interface Validate {
    /**
     * 组织名称错误提示
     */
    departName: ValidateState;

    /**
     * 部门编码
     */
    code: ValidateState;

    /**
     * 备注
     */
    remark: ValidateState;

    /**
     * 邮箱错误提示
     */
    email: ValidateState;

    /**
     * 存储位置提示
     */
    ossInfo: Doclibs.ValidateInfo;

    /**
     * 上级部门
     */
    parentId: ValidateState;
}

interface CreateDepartmentProps extends React.Props<void> {
    /**
     * 选择的部门信息
     */
    dep?: Dep;

    /**
     * 来源
     */
    sourcePage?: string;

    /**
     * 当前登录的用户
     */
    userid: string;

    /**
     * 取消新建部门
     */
    onRequestCancelCreateDep: () => void;

    /**
     * 新建部门成功
     */
    onCreateDepSuccess: (nodeInfo: NodeInfo) => void;

    /**
     * 上级部门不存在，删除部门
     */
    onRequestDelDep: (dep: Dep) => any;
}

interface CreateDepartmentState {
    /**
     * 新建部门信息
     */
    departmentInfo: DepartmentInfo;

    /**
     * 部门负责人
     */
    managerInfo: UserInfoType[];

    /**
     * 错误提示信息索引
     */
    validateState: Validate;

    /**
     * 是否显示选择弹框
     */
    showAddDepartmentLeaderDialog: boolean;

    /**
     * 是否显示选择部门弹框
     */
    showAddDepartmentDialog: boolean;
}
export default class CreateDepartmentBase extends WebComponent<CreateDepartmentProps, CreateDepartmentState> {
    static defaultProps = {
        dep: null,
        sourcePage: 'default',
        userid: '',
        onRequestCancelCreateDep: noop,
        onCreateDepSuccess: noop,
    }

    state = {
        departmentInfo: {
            ossInfo: (this.props.dep && this.props.dep.ossInfo && this.props.dep.ossInfo.ossId) ? this.props.dep.ossInfo : { enabled: true, ossId: '', ossName: '' },
            departName: '',
            code: '',
            remark: '',
            status: true,
            email: '',
            parentName: this.props.dep && this.props.dep.name || '',
            parentId: this.props.dep && this.props.dep.id || '',
            parentType: '',
        },
        managerInfo: [],
        validateState: {
            departName: ValidateState.Normal,
            code: ValidateState.Normal,
            remark: ValidateState.Normal,
            email: ValidateState.Normal,
            ossInfo: getValidateInfo(ValidateStatus.Normal),
            parentId: ValidateState.Normal,
        },
        showAddDepartmentLeaderDialog: false,
        showAddDepartmentDialog: false,
    }

    isRequest: boolean // 是否在请求中

    async componentDidMount() {
        try {
            if(this.props.sourcePage !== 'default') {
                return
            }
            const { departmentInfo, validateState } = this.state;
            let ossInfo = departmentInfo.ossInfo;
            if (ossInfo && ossInfo.ossId) {
                const ossData = await getObjectStorageInfoById(ossInfo.ossId)
                ossInfo = {
                    enabled: ossData.enabled,
                    ossId: ossData.id,
                    ossName: ossData.name,
                }
            }
            const invalidateOssInfo = ossInfo.ossId && !ossInfo.enabled;

            this.setState({
                departmentInfo: {
                    ...departmentInfo,
                    ossInfo,
                },
                validateState: {
                    ...validateState,
                    ossInfo: !invalidateOssInfo ? getValidateInfo(ValidateStatus.Normal) : getValidateInfo(ErrorCode.OSSDisabled),
                },
            })
        } catch ({ error }) {
            error && await Message2.info({ message: error.errMsg })
        }
    }

    /**
     * 修改存储位置
     */
    protected updateSelectedOss = (ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo): void => {
        const { validateState, departmentInfo } = this.state;

        this.setState({
            departmentInfo: {
                ...departmentInfo,
                ossInfo,
            },
            validateState: {
                ...validateState,
                ossInfo: getValidateInfo(ValidateStatus.Normal),
            },
        })
    }

    /**
     * 部门名称失焦事件
     */
    protected handleOnBlurDep(): void {
        const { departmentInfo: { departName }, validateState } = this.state;
        const validateName = isNormalName(trim(departName));

        this.setState({
            validateState: {
                ...validateState,
                departName: validateName ? ValidateState.Normal : trim(departName) ? ValidateState.DepartmentInvalid : ValidateState.Empty,
            },
        })
    }

    /**
     * 部门编码失焦事件
     */
    protected handleOnBlurCode(): void {
        const { departmentInfo: { code }, validateState } = this.state;
        const validateCode = isNormalCode(trim(code));

        this.setState({
            validateState: {
                ...validateState,
                code: validateCode ? ValidateState.Normal : trim(code) ? ValidateState.CodeInvalid : ValidateState.Normal,
            },
        })
    }

    /**
     * 备注失焦事件
     */
    protected handleOnBlurRemark(): void {
        const { departmentInfo: { remark }, validateState } = this.state;
        const validateRemark = isNormalName(trim(remark));

        this.setState({
            validateState: {
                ...validateState,
                remark: validateRemark ? ValidateState.Normal : trim(remark) ? ValidateState.RemarksInvalid : ValidateState.Normal,
            },
        })
    }

    /**
     * 邮箱失焦事件
     */
    protected handleOnBlurEmail(): void {
        const { departmentInfo: { email }, validateState } = this.state;
        const validateEmail = mailAndLenth(email, 4, 101);

        this.setState({
            validateState: {
                ...validateState,
                email: validateEmail ? ValidateState.Normal : email ? ValidateState.EamilInvalid : ValidateState.Normal,
            },
        })
    }

    /**
    * 检查表单合法性
    */
    private checkForm = (): boolean => {
        const { departmentInfo: { departName, email, code, remark, parentId }, validateState } = this.state;
        const validateName = isNormalName(trim(departName));
        const validateCode = isNormalCode(trim(code));
        const validateRemark = isNormalName(trim(remark));
        const validateEmail = mailAndLenth(email, 4, 101);

        if (validateName && (validateEmail || !email)&& (validateCode || !code) && (validateRemark || !remark) && !(this.props.sourcePage !== 'default' && !parentId)) {
            return true;
        } else {
            this.setState({
                validateState: {
                    ...validateState,
                    departName: validateName ? ValidateState.Normal : trim(departName) ? ValidateState.DepartmentInvalid : ValidateState.Empty,
                    email: validateEmail ? ValidateState.Normal : email ? ValidateState.EamilInvalid : ValidateState.Normal,
                    code: validateCode ? ValidateState.Normal : trim(code) ? ValidateState.CodeInvalid : ValidateState.Normal,
                    remark: validateRemark ? ValidateState.Normal : trim(remark) ? ValidateState.RemarksInvalid : ValidateState.Normal,
                    parentId: this.props.sourcePage !== 'default' && parentId ? ValidateState.Normal : ValidateState.Empty,
                },
            })
            return false;
        }
    }

    /**
     * 设置文本框变更事件
     */
    protected handleValueChange(inputInfo: { departName?: string; code?: string; remark?: string; status?: boolean; email?: string }) {
        const { departmentInfo, validateState } = this.state;

        this.setState({
            departmentInfo: { ...departmentInfo, ...inputInfo },
            validateState: {
                ...validateState,
                departName: 'departName' in inputInfo ? ValidateState.Normal : validateState.departName,
                code: 'code' in inputInfo ? ValidateState.Normal : validateState.code,
                remark: 'remark' in inputInfo ? ValidateState.Normal : validateState.remark,
                email: 'email' in inputInfo ? ValidateState.Normal : validateState.email,
            },
        })
    }

    /**
     * 提交新建部门信息
     */
    protected createDepartment = async () => {
        if (this.checkForm()) {
            const { managerInfo, departmentInfo, departmentInfo: { code, remark, status, email, ossInfo, ossInfo: { ossId }, parentName, parentId }, validateState } = this.state;
            const departName = trim(departmentInfo.departName);
            const addParmas = {
                ncTAddDepartParam: {
                    parentId,
                    departName,
                    managerID: managerInfo.length ? (managerInfo as UserInfoType[])[0].id : null,
                    code: trim(code),
                    remark: trim(remark),
                    status,
                    email,
                    ossId: ossId || '',
                },
            }

            if (!this.isRequest) {
                try {
                    this.isRequest = true;
                    const createId = await addDepartment([addParmas]);
                    const nodeInfo = {
                        id: createId,
                        name: departName.replace(/\.+$/, ''),
                        departmentId: createId,
                        departmentName: departName.replace(/\.+$/, ''),
                        managerID: managerInfo.length ? (managerInfo as UserInfoType[])[0].id : null,
                        code,
                        remark,
                        status,
                        ossInfo,
                        email,
                        responsiblePersons: [],
                        subDepartmentCount: 0,
                        parentId,
                        managerInfo,
                    }

                    manageLog(
                        ManagementOps.CREATE,
                        __('新建部门 “${name}” 成功', { name: departName.replace(/\.+$/, '') }),
                        __('部门编码 “${code}”； 部门负责人 “${managerDisplayName}”；备注 “${remark}”；邮箱地址 “${emailAddress}”；存储位置 “${ossName}”；状态 “${status}”；', {
                            emailAddress: email,
                            ossName: displayUserOssInfo(ossInfo),
                            code,
                            managerDisplayName: managerInfo.length ? (managerInfo as UserInfoType[])[0].name : '',
                            remark,
                            status: status ? __('启用') : __('禁用'),
                        }),
                        Level.INFO,
                    )
                    Toast.open(__('新建成功'))
                    this.props.onCreateDepSuccess(nodeInfo);

                } catch ({ error }) {
                    this.isRequest = false;
                    if (error) {
                        switch (error.errID) {
                            case ErrorCode.InvalidDpCode:
                                this.setState({
                                    validateState: { ...validateState, code: ValidateState.CodeInvalid },
                                })
                                break;

                            case ErrorCode.DpCodeExit:
                                this.setState({
                                    validateState: { ...validateState, code: ValidateState.DpCodeExit },
                                })
                                break;

                            case ErrorCode.DepNameExist:
                                this.setState({
                                    validateState: {
                                        ...validateState,
                                        departName: ValidateState.DepartmentExist,
                                    },
                                })
                                break;

                            case ErrorCode.InvalidDepName:
                                this.setState({
                                    validateState: {
                                        ...validateState,
                                        departName: ValidateState.DepartmentInvalid,
                                    },
                                })
                                break;

                            case ErrorCode.DepNameNotExist:
                                this.props.onRequestCancelCreateDep()
                                await Message2.info({ message: __('新建失败，上级部门 “${parentName}” 不存在，请重新选择。', { parentName }) })
                                this.props.sourcePage === 'default' && this.props.onRequestDelDep(this.props.dep)
                                break;

                            case ErrorCode.EmailExist:
                                this.setState({
                                    validateState: {
                                        ...validateState,
                                        email: ValidateState.EmailExist,
                                    },
                                })
                                break;

                            case ErrorCode.OSSNotExist:
                            case ErrorCode.OSSDisabled:
                            case ErrorCode.OSSInvalid:
                            case ErrorCode.OSSUnabled:
                                this.setState({
                                    validateState: {
                                        ...validateState,
                                        ossInfo: getValidateInfo(error.errID),
                                    },
                                })
                                break;

                            default:
                                await Message2.info({ message: error.errMsg })
                                break;
                        }
                    }
                }
            }
        }
    }
}