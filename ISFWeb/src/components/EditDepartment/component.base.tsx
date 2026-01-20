import * as React from 'react'
import { noop, trim } from 'lodash';
import { Message2, Toast } from '@/sweet-ui';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { editDepartment, editUserOss, editDepartOSS, getOrgDepartmentById } from '@/core/thrift/sharemgnt/sharemgnt';
import { Doclibs, ValidateStatus, getValidateInfo } from '@/core/doclibs/doclibs';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { displayUserOssInfo } from '@/core/oss/oss';
import { getObjectStorageInfoById } from '@/core/apis/console/ossgateway'
import { ValidateState, Dep, DepartmentInfo, UserInfoType } from '@/core/user';
import { mailAndLenth, isNormalName, isNormalCode } from '@/util/validators';
import { listUsersSince } from '../helper';
import WebComponent from '../webcomponent';
import __ from './locale';

export enum Status {
    /**
     * 默认状态
    */
    Normal,

    /**
     * 编辑中
    */
    Loading,
}

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
}

interface EditDepartmentProps extends React.Props<void> {
    /**
     * 选择的部门
     */
    dep: Dep;

    /**
     * 上级部门名称
     */
    parentName: string;

    /**
     * 当前登录的用户
     */
    userid: string;

    /**
     * 取消编辑部门
     */
    onRequestCancelEditDep: () => void;

    /**
     * 编辑部门成功
     */
    onEditDepSuccess: (dep: Dep) => void;

    /**
     * 部门不存在，删除部门
     */
    onRequestDelDep: (dep: Dep) => any;
}

interface EditDepartmentState {
    /**
     * 编辑部门信息
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
    * 是否更改子部门及用户存储位置
    */
    changeOss: boolean;

    /**
    * 更改子部门存储位置个数
    */
    totalCount: number;

    /**
    * 弹框加载状态
    */
    status: Status;

    /**
     * 是否显示选择弹框
     */
    showAddDepartmentLeaderDialog: boolean;
     /**
     * 是否编辑了部门
     */
    isEditDepartment: boolean;
}

export default class EditDepartmentBase extends WebComponent<EditDepartmentProps, EditDepartmentState> {
    static defaultProps = {
        dep: null,
        parentName: '',
        userid: '',
        onRequestCancelEditDep: noop,
        onEditDepSuccess: noop,
    }

    state = {
        departmentInfo: {
            ossInfo: (this.props.dep.ossInfo && this.props.dep.ossInfo.ossId) ? this.props.dep.ossInfo : { enabled: true, ossId: '', ossName: '' },
            departName: this.props.dep.name,
            email: this.props.dep.email,
            code: this.props.dep.code,
            remark: this.props.dep.remark,
            status: this.props.dep.status,
            parentName: this.props.parentName,
        },
        validateState: {
            departName: ValidateState.Normal,
            code: ValidateState.Normal,
            remark: ValidateState.Normal,
            email: ValidateState.Normal,
            ossInfo: getValidateInfo(ValidateStatus.Normal),
        },
        changeOss: false,
        totalCount: 0,
        status: Status.Normal,
        managerInfo: [],
        showAddDepartmentLeaderDialog: false,
        isEditDepartment: false,
    }

    isRequest: boolean // 是否在请求中

    async componentDidMount() {
        try {
            let { ossInfo, email, managerID, managerDisplayName, code, remark, status } = await getOrgDepartmentById(this.props.dep.id)
            const { departmentInfo, validateState } = this.state;

            if (ossInfo && ossInfo.ossId) {
                const ossData = await getObjectStorageInfoById(ossInfo.ossId)
                ossInfo = {
                    enabled: ossData.enabled,
                    ossId: ossData.id,
                    ossName: ossData.name,
                }
            }
            const invalidateOssInfo = ossInfo.ossId && !ossInfo.enabled

            this.setState({
                departmentInfo: {
                    ...departmentInfo,
                    ossInfo,
                    email,
                    code,
                    remark,
                    status,
                },
                managerInfo: managerID && managerDisplayName ? [{ id: managerID, name: managerDisplayName, type: 'user' }] :[],
                validateState: {
                    ...validateState,
                    ossInfo: !invalidateOssInfo ? getValidateInfo(ValidateStatus.Normal) : getValidateInfo(ErrorCode.OSSDisabled),
                },
                isEditDepartment: false,
            })
        } catch ({ error }) {
            error && await Message2.info({ message: error.errMsg })
        }
    }

    /**
    * 存储位置的设置应用到子部门及用户成员
    */
    protected setHomeOssStatus(changeOss: boolean) {
        this.setState({
            changeOss,
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
        const { departmentInfo: { departName, email, code, remark }, validateState } = this.state;
        const validateName = isNormalName(trim(departName));
        const validateCode = isNormalCode(trim(code));
        const validateRemark = isNormalName(trim(remark));
        const validateEmail = mailAndLenth(email, 4, 101);

        if (validateName && (validateEmail || !email) && (validateCode || !code) && (validateRemark || !remark)) {
            return true;
        } else {
            this.setState({
                validateState: {
                    ...validateState,
                    departName: validateName ? ValidateState.Normal : trim(departName) ? ValidateState.DepartmentInvalid : ValidateState.Empty,
                    code: validateCode ? ValidateState.Normal : trim(code) ? ValidateState.CodeInvalid : ValidateState.Normal,
                    remark: validateRemark ? ValidateState.Normal : trim(remark) ? ValidateState.RemarksInvalid : ValidateState.Normal,
                    email: validateEmail ? ValidateState.Normal : email ? ValidateState.EamilInvalid : ValidateState.Normal,
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
            isEditDepartment: true,
        })
    }

    /**
     * 保存编辑部门
     */
    protected editDepartment = async () => {
        if (this.checkForm()) {
            const { departmentInfo, managerInfo, departmentInfo: { code, remark, status, email, ossInfo, ossInfo: { ossId }, parentName }, changeOss, validateState, isEditDepartment } = this.state;
            const { dep: { id, name } } = this.props;
            const departName = trim(departmentInfo.departName);
            const editParma = {
                ncTEditDepartParam: {
                    departName,
                    code: trim(code),
                    remark: trim(remark),
                    managerID: managerInfo.length ? (managerInfo as UserInfoType[])[0].id : '',
                    status,
                    departId: id,
                    ossId: ossId || '',
                    email,
                },
            }
            let dep = this.props.dep;

            if (!this.isRequest) {
                try {
                    this.isRequest = true
                    await editDepartment([editParma]);
                    dep = { ...dep, name: departName, code: trim(code), remark: trim(remark), status, managerInfo }
                    if (email !== dep.email) {
                        dep = { ...dep, email }
                    }

                    if (dep.ossInfo && dep.ossInfo.ossId !== ossId) {
                        dep = { ...dep, ossInfo }
                    }
                    if(isEditDepartment){
                        let nameText = name;
                        if(nameText !== departName.replace(/\.+$/, '')){
                            nameText =  __('由 ${oldText} 改为 ${newText}', { oldText: nameText, newText: departName.replace(/\.+$/, '') });
                        }

                        manageLog(
                            ManagementOps.SET,
                            __('编辑部门 “${depName}” 成功', { depName: departName.replace(/\.+$/, '') }),
                            __('部门名 “${oldName}”；部门编码 “${code}”； 部门负责人“${managerDisplayName}”；备注 “${remark}”；邮箱地址 “${emailAddress}”；存储位置 “${ossName}”；状态 “${status}”；', {
                                oldName: nameText,
                                emailAddress: email,
                                code,
                                managerDisplayName: managerInfo.length ? (managerInfo as UserInfoType[])[0].name : '',
                                remark,
                                status: status ? __('启用') : __('禁用'),
                                ossName: displayUserOssInfo(ossInfo),
                            }),
                            Level.INFO,
                        )
                    }
                    if (changeOss) {
                        this.setState({
                            status: Status.Loading,
                        })
                        this.changeOss(dep)
                    } else {
                        if(isEditDepartment){
                            Toast.open(__('编辑成功'))
                            this.props.onEditDepSuccess(dep);
                        }else{
                            this.props.onRequestCancelEditDep()
                        }
                    }
                    this.setState({
                        isEditDepartment: false,
                    })
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

                            case ErrorCode.ParentDepartmentNotExist:
                            case ErrorCode.DepNameNotExist:
                                this.props.onRequestCancelEditDep()
                                await Message2.info({ message: __('部门不存在。') })
                                this.props.onRequestDelDep(dep)
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

    /**
    * 更改存储位置
    */
    private async changeOss(orgNode: any) {
        try {
            let { users: elements, deps: deparments } = await listUsersSince(orgNode, true);
            let that = this;
            const { departmentInfo: { ossInfo } } = this.state;

            deparments = deparments.slice(1);

            elements = elements.concat(deparments);
            if (elements.length === 0) {
                that.setState({
                    totalCount: 100,
                })
                Toast.open(__('编辑成功'))
                setTimeout(() => {
                    that.props.onEditDepSuccess({ ...orgNode, changeChild: false });
                }, 300)
            } else {
                for (let i = 0; i < elements.length; i++) {
                    that.setState({
                        totalCount: Math.round((i + 1) / elements.length * 100),
                    })

                    if (elements.length === i + 1) {
                        Toast.open(__('编辑成功'))
                        setTimeout(() => {
                            that.props.onEditDepSuccess({ ...orgNode, changeChild: true });
                        }, 300)
                    }

                    const element = elements[i]

                    if (element.user) {
                        await editUserOss(element.id, ossInfo.ossId || '')
                        await manageLog(ManagementOps.SET,
                            __('编辑用户 "${displayName}(${loginName})" 成功', { displayName: element.user.displayName, loginName: element.user.loginName }),
                            __('存储位置 “${ossName}”', { ossName: displayUserOssInfo(ossInfo) }),
                        )
                    } else {
                        await editDepartOSS([element.id, ossInfo.ossId || ''])
                        await manageLog(ManagementOps.SET,
                            __('编辑部门 “${depName}” 成功', { depName: element.name }),
                            __('存储位置 “${ossName}”', { ossName: displayUserOssInfo(ossInfo) }),
                        );
                    }
                }
            }
        } catch ({ error }) {
            error && await Message2.info({ message: error.errMsg })
        }
    }
}