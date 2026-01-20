import * as React from 'react'
import { noop, trim } from 'lodash';
import { Message2, Toast } from '@/sweet-ui';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { editOrganization, editUserOss, editDepartOSS, getOrgDepartmentById } from '@/core/thrift/sharemgnt/sharemgnt';
import { Doclibs, ValidateStatus, getValidateInfo } from '@/core/doclibs/doclibs';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { displayUserOssInfo } from '@/core/oss/oss';
import { getObjectStorageInfoById } from '@/core/apis/console/ossgateway'
import { ValidateState, Dep, OrganizeInfo, UserInfoType } from '@/core/user';
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

interface EditOrganizationProps extends React.Props<void> {
    /**
     * 选择的部门
     */
    dep: Dep;

    /**
     * 当前登录的用户
     */
    userid: string;

    /**
     * 取消编辑组织
     */
    onRequestCancelEditOrg: () => void;

    /**
     * 编辑组织成功
     */
    onEditOrgSuccess: (dep: Dep) => void;

    /**
     * 组织不存在，删除该组织
     */
    onRequestDelOrg: (dep: Dep) => any;
}

interface EditOrganizationState {
    /**
     * 编辑组织信息
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
     * 是否修改了组织信息
     */
    isEditOrgInfo: boolean;

    /**
     * 是否显示选择弹框
     */
    showAddDepartmentLeaderDialog: boolean;
}
export default class EditOrganizationBase extends WebComponent<EditOrganizationProps, EditOrganizationState> {
    static defaultProps = {
        dep: null,
        userid: '',
        onRequestCancelEditOrg: noop,
        onEditOrgSuccess: noop,
    }

    state = {
        organizeInfo: {
            ossInfo: { enabled: true, ossId: '', ossName: '' },
            orgName: this.props.dep.name,
            email: '',
            code: '',
            remark: '',
            status: true,
        },
        managerInfo: [],
        validateState: {
            orgName: ValidateState.Normal,
            code: ValidateState.Normal,
            remark: ValidateState.Normal,
            email: ValidateState.Normal,
            ossInfo: getValidateInfo(ValidateStatus.Normal),
        },
        changeOss: false,
        totalCount: 0,
        status: Status.Normal,
        showAddDepartmentLeaderDialog: false,
        isEditOrgInfo: false,
    }

    isRequest: boolean // 是否在请求中

    async componentDidMount() {
        try {
            const { organizeInfo, validateState } = this.state
            let { ossInfo, email, code, managerID, managerDisplayName, remark, status } = await getOrgDepartmentById(this.props.dep.id)
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
                organizeInfo: {
                    ...organizeInfo,
                    ossInfo,
                    email,
                    code,
                    remark,
                    status,
                },
                managerInfo: managerID && managerDisplayName ? [{ id: managerID, name: managerDisplayName, type: 'user' }] : [],
                validateState: {
                    ...validateState,
                    ossInfo: !invalidateOssInfo ? getValidateInfo(ValidateStatus.Normal) : getValidateInfo(ErrorCode.OSSDisabled),
                },
                isEditOrgInfo: false,
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
            isEditOrgInfo: true,
        })
    }

    /**
     * 保存编辑组织
     */
    protected editOrganization = async () => {
        if (this.checkForm()) {
            const { changeOss, managerInfo, organizeInfo, organizeInfo: { code, remark, status, email, ossInfo, ossInfo: { ossId } }, validateState, isEditOrgInfo  } = this.state;
            const { dep: { id, name } } = this.props;
            const orgName = trim(organizeInfo.orgName);
            const editParma = {
                ncTEditDepartParam: {
                    departId: id,
                    departName: orgName,
                    managerID: managerInfo.length ? (managerInfo as UserInfoType[])[0].id : '',
                    code: trim(code),
                    remark: trim(remark),
                    status,
                    ossId: ossId || '',
                    email,
                },
            };
            let dep = this.props.dep;

            if (!this.isRequest) {
                try {
                    this.isRequest = true;
                    await editOrganization([editParma]);
                    dep = { ...dep, name: orgName, email, ossInfo, managerInfo, status, code: trim(code), remark: trim(remark) }
                    if(isEditOrgInfo) {
                        let nameText = name;
                        if(nameText !== orgName.replace(/\.+$/, '')) {
                            nameText = __('由 ${oldText} 改为 ${newText}', { oldText: nameText, newText: orgName.replace(/\.+$/, '') })
                        }
                        manageLog(
                            ManagementOps.SET,
                            __('编辑 组织 “${name}” 成功', { name: orgName.replace(/\.+$/, '') }),
                            __('组织名 “${oldName}”；组织编码 “${code}”；组织负责人 “${managerDisplayName}”；备注 “${remark}”；邮箱地址 “${emailAddress}”；存储位置 “${ossName}”；状态 “${status}”；', {
                                oldName: nameText,
                                emailAddress: email,
                                ossName: displayUserOssInfo(ossInfo),
                                code,
                                managerDisplayName: managerInfo.length ? (managerInfo as UserInfoType[])[0].name : '',
                                remark,
                                status: status ? __('启用') : __('禁用'),
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
                        if(isEditOrgInfo) {
                            Toast.open(__('编辑成功'))
                            this.props.onEditOrgSuccess(dep);
                        }else{
                            this.props.onRequestCancelEditOrg()
                        }
                    }
                    this.setState({
                        isEditOrgInfo: false,
                    })
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

                        case ErrorCode.OrgNameNotExist:
                            this.props.onRequestCancelEditOrg()
                            await Message2.alert({ message: __('组织 “${orgName}” 不存在。', { orgName: orgName }) })
                            this.props.onRequestDelOrg(this.props.dep)
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
                            Message2.alert({ message: error.errMsg })
                            break;
                    }
                }
            }
        }
    }

    /**
     * 更改存储位置
    */
    private async changeOss(orgNode: Dep) {
        try {
            let { users: elements, deps: deparments } = await listUsersSince(orgNode, true);

            let that = this;
            const { organizeInfo: { ossInfo } } = this.state;

            deparments = deparments.slice(1);

            elements = elements.concat(deparments);
            if (elements.length === 0) {
                that.setState({
                    totalCount: 100,
                })
                Toast.open(__('编辑成功'))
                setTimeout(() => {
                    that.props.onEditOrgSuccess({ ...orgNode, changeChild: false });
                }, 300)
            } else {
                for (let i = 0; i < elements.length; i++) {
                    that.setState({
                        totalCount: Math.round((i + 1) / elements.length * 100),
                    })

                    if (elements.length === i + 1) {
                        Toast.open(__('编辑成功'))
                        setTimeout(() => {
                            that.props.onEditOrgSuccess({ ...orgNode, changeChild: true });
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
                            __('编辑组织 “${orgName}” 成功', { orgName: element.name }),
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