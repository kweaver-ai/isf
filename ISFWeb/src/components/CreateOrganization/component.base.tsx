import * as React from 'react'
import { noop, trim } from 'lodash';
import { Message2, Toast } from '@/sweet-ui';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { Doclibs, ValidateStatus, getValidateInfo } from '@/core/doclibs/doclibs';
import { createOrganization } from '@/core/thrift/sharemgnt/sharemgnt';
import { displayUserOssInfo } from '@/core/oss/oss';
import { ValidateState, Dep, OrganizeInfo, NodeInfo, UserInfoType } from '@/core/user';
import { manageLog, Level, ManagementOps } from '@/core/log';
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
    orgName: ValidateState;

    /**
     * 组织编码
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

interface CreateOrganizationProps extends React.Props<void> {
    /**
     * 选择的部门
     */
    dep: Dep;

    /**
     * 当前登录的用户
     */
    userid: string;

    /**
     * 取消新建组织
     */
    onRequestCancelCreateOrg: () => void;

    /**
     * 新建组织成功
     */
    onCreateOrgSuccess: (nodeInfo: NodeInfo) => void;
}

interface CreateOrganizationState {
    /**
     * 新建组织信息
     */
    organizeInfo: OrganizeInfo;

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
}
export default class CreateOrganizationBase extends WebComponent<CreateOrganizationProps, CreateOrganizationState> {
    static defaultProps = {
        dep: null,
        userid: '',
        onRequestCancelCreateOrg: noop,
        onCreateOrgSuccess: noop,
    }

    state = {
        organizeInfo: {
            ossInfo: { enabled: true, ossId: '', ossName: '' },
            orgName: '',
            code: '',
            remark: '',
            status: true,
            email: '',
        },
        managerInfo: [],
        validateState: {
            orgName: ValidateState.Normal,
            code: ValidateState.Normal,
            remark: ValidateState.Normal,
            email: ValidateState.Normal,
            ossInfo: getValidateInfo(ValidateStatus.Normal),
        },
        showAddDepartmentLeaderDialog: false,
    }

    isRequest: boolean // 是否在请求中

    /**
     * 修改存储位置
     */
    protected updateSelectedOss = (ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo): void => {
        const { validateState, organizeInfo } = this.state;

        this.setState({
            organizeInfo: {
                ...organizeInfo,
                ossInfo,
            },
            validateState: {
                ...validateState,
                ossInfo: getValidateInfo(ValidateStatus.Normal),
            },
        })
    }

    /**
     * 组织名称失焦事件
     */
    protected handleOnBlurOrg(): void {
        const { organizeInfo: { orgName }, validateState } = this.state;
        const validateOrgName = isNormalName(trim(orgName));

        this.setState({
            validateState: {
                ...validateState,
                orgName: validateOrgName ? ValidateState.Normal : trim(orgName) ? ValidateState.OrgInValid : ValidateState.Empty,
            },
        })
    }

    /**
     * 组织编码失焦事件
     */
    protected handleOnBlurCode(): void {
        const { organizeInfo: { code }, validateState } = this.state;
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
        const { organizeInfo: { remark }, validateState } = this.state;
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
        const { organizeInfo: { email }, validateState } = this.state;
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
        const { organizeInfo: { orgName, code, remark, email }, validateState } = this.state;
        const validateOrgName = isNormalName(trim(orgName));
        const validateCode = isNormalCode(trim(code))
        const validateEmail = mailAndLenth(email, 4, 101);
        const validateRemark = isNormalName(trim(remark));

        if (validateOrgName && (validateEmail || !email)&& (validateCode || !code) && (validateRemark || !remark)) {
            return true;
        } else {
            this.setState({
                validateState: {
                    ...validateState,
                    orgName: validateOrgName ? ValidateState.Normal : trim(orgName) ? ValidateState.OrgInValid : ValidateState.Empty,
                    email: validateEmail ? ValidateState.Normal : email ? ValidateState.EamilInvalid : ValidateState.Normal,
                    code: validateCode ? ValidateState.Normal : trim(code) ? ValidateState.CodeInvalid : ValidateState.Normal,
                    remark: validateRemark ? ValidateState.Normal : trim(remark) ? ValidateState.RemarksInvalid : ValidateState.Normal,
                },
            })
            return false;
        }
    }

    /**
     * 设置文本框变更事件
     */
    protected handleValueChange(inputInfo: { orgName?: string; code?: string; remark?: string; status?: boolean; email?: string }) {
        const { organizeInfo, validateState } = this.state;

        this.setState({
            organizeInfo: { ...organizeInfo, ...inputInfo },
            validateState: {
                ...validateState,
                orgName: 'orgName' in inputInfo ? ValidateState.Normal : validateState.orgName,
                code: 'code' in inputInfo ? ValidateState.Normal : validateState.code,
                remark: 'remark' in inputInfo ? ValidateState.Normal : validateState.remark,
                email: 'email' in inputInfo ? ValidateState.Normal : validateState.email,
            },
        })
    }

    /**
     * 保存新建组织
     */
    protected createOrganization = async () => {
        if (this.checkForm()) {
            const { organizeInfo, managerInfo, organizeInfo: { code, remark, status, email, ossInfo, ossInfo: { ossId } }, validateState } = this.state;
            const orgName = trim(organizeInfo.orgName);
            const data = {
                ncTAddOrgParam: {
                    orgName,
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
                    const organizationId = await createOrganization([data]);
                    const nodeInfo = {
                        id: organizationId,
                        name: orgName.replace(/\.+$/, ''),
                        organizationId,
                        organizationName: orgName,
                        is_root: true,
                        ossInfo,
                        email,
                        depart_existed: false,
                        code: trim(code),
                        remark: trim(remark),
                        status,
                        managerInfo,
                    }

                    manageLog(
                        ManagementOps.CREATE,
                        __('新建 组织 “${name}” 成功', { name: orgName.replace(/\.+$/, '') }),
                        __('组织编码 “${code}”；组织负责人 “${managerDisplayName}”；备注 “${remark}”；邮箱地址 “${emailAddress}”；存储位置 “${ossName}”；状态 “${status}”；', {
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
                    this.props.onCreateOrgSuccess(nodeInfo);

                } catch ({ error }) {
                    this.isRequest = false;
                    switch (error && error.errID) {
                        case ErrorCode.InvalidDpCode:
                            this.setState({
                                validateState: { ...validateState, code: ValidateState.CodeInvalid },
                            })
                            break;

                        case ErrorCode.DpCodeExit:
                            this.setState({
                                validateState: { ...validateState, code: ValidateState.OrgCodeExit },
                            })
                            break;
                        case ErrorCode.OrgNameExist:
                            this.setState({
                                validateState: {
                                    ...validateState,
                                    orgName: ValidateState.OrgExist,
                                },
                            })
                            break;

                        case ErrorCode.InvalidDepartName:
                            this.setState({
                                validateState: {
                                    ...validateState,
                                    orgName: ValidateState.OrgInValid,
                                },
                            })
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
                            Message2.info({ message: error.errMsg })
                            break;
                    }
                }
            }
        }
    }
}