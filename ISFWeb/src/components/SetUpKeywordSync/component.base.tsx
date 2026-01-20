import * as React from 'react';
import { noop } from 'lodash';
import { getDomainKeyConfig, setDomainKeyConfig } from '@/core/thrift/sharemgnt';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { Message2 } from '@/sweet-ui';
import { SetUpKeywordSyncProps, SetUpKeywordSyncState, ValidateStatusEnum, KeywordInput } from './helper';
import __ from './locale';

export default class SetUpKeywordSyncBase extends React.PureComponent<SetUpKeywordSyncProps, SetUpKeywordSyncState> {
    static defaultProps: SetUpKeywordSyncProps = {
        domainInfo: {
            id: -1,
            name: '',
        },
        onRequestEditStatus: noop,
        onRequestDomainInvalid: noop,
    }

    state: SetUpKeywordSyncState = {
        keywordInput: {
            departNameKeys: '',
            departThirdIdKeys: '',
            loginNameKeys: '',
            displayNameKeys: '',
            emailKeys: '',
            userThirdIdKeys: '',
            groupKeys: '',
            subOuFilter: '',
            subUserFilter: '',
            baseFilter: '',
            statusKeys: '',
            idcardNumberKeys: '',
            telNumberKeys: '',
        },
        isEditStatus: false,
        validateStatus: {
            departNameKeysValidateStatus: ValidateStatusEnum.Normal,
            departThirdIdKeysValidateStatus: ValidateStatusEnum.Normal,
            loginNameKeysValidateStatus: ValidateStatusEnum.Normal,
            displayNameKeysValidateStatus: ValidateStatusEnum.Normal,
            emailKeysValidateStatus: ValidateStatusEnum.Normal,
            userThirdIdKeysValidateStatus: ValidateStatusEnum.Normal,
            subOuFilterValidateStatus: ValidateStatusEnum.Normal,
            subUserFilterValidateStatus: ValidateStatusEnum.Normal,
            baseFilterValidateStatus: ValidateStatusEnum.Normal,
        },
    }

    /**
     * 存储每一次保存前输入框字段值，用于恢复先前状态
     */
    originKeywordInput: KeywordInput = {
        /**
         * 部门名对应的域key字段
         */
        departNameKeys: '',

        /**
         * 部门ID对应的域key字段
         */
        departThirdIdKeys: '',

        /**
         * 登录名对应的域key字段
         */
        loginNameKeys: '',

        /**
         * 显示名对应的域key字段
         */
        displayNameKeys: '',

        /**
         * 用户邮箱对应的域key字段
         */
        emailKeys: '',

        /**
         * 用户Id对应的域key字段
         */
        userThirdIdKeys: '',

        /**
         * 安全组信息的key字段
         */
        groupKeys: '',

        /**
         * 搜索子部门的Filter
         */
        subOuFilter: '',

        /**
         * 搜索子用户的Filter
         */
        subUserFilter: '',

        /**
         * 具体某个部门或用户信息的Filter
         */
        baseFilter: '',

        /**
         * 用户状态对应的域key字段
         */
        statusKeys: '',

        /**
         * 用户身份证号对应的域key字段
         */
        idcardNumberKeys: '',

        /**
         * 用户电话号对应的域key字段
         */
        telNumberKeys: '',
    }

    async componentDidMount() {
        /**
         * 获取域关键字配置信息
         */
        const {
            departNameKeys,
            departThirdIdKeys,
            loginNameKeys,
            displayNameKeys,
            emailKeys,
            userThirdIdKeys,
            groupKeys,
            idcardNumberKeys,
            telNumberKeys,
            statusKeys,
            subOuFilter,
            subUserFilter,
            baseFilter,
        } = await getDomainKeyConfig([this.props.domainInfo.id])

        this.setState({
            keywordInput: {
                departNameKeys: departNameKeys.join(),
                departThirdIdKeys: departThirdIdKeys.join(),
                loginNameKeys: loginNameKeys.join(),
                displayNameKeys: displayNameKeys.join(),
                emailKeys: emailKeys.join(),
                userThirdIdKeys: userThirdIdKeys.join(),
                groupKeys: groupKeys.join(),
                idcardNumberKeys: idcardNumberKeys.join(),
                telNumberKeys: telNumberKeys.join(),
                statusKeys: statusKeys.join(),
                subOuFilter,
                subUserFilter,
                baseFilter,
            },
        }, () => {
            /**
             * 保存初始字段信息（源信息），用作后续恢复状态
             */
            this.originKeywordInput = this.state.keywordInput;
        })
    }

    /**
     * 编辑输入框
     */
    protected handleInputChange = (value: string, name: string): void => {
        this.setState({
            keywordInput: {
                ...this.state.keywordInput,
                [name]: value,
            },
            isEditStatus: true,
            validateStatus: {
                ...this.state.validateStatus,
                [`${name}ValidateStatus`]: ValidateStatusEnum.Normal,
            },
        })

        this.props.onRequestEditStatus(true);
    }

    /**
     * 失焦时验证失焦输入框输入内容合规性
     */
    protected handleValidate = (name: string): void => {
        this.setState({
            validateStatus: {
                ...this.state.validateStatus,
                [`${name}ValidateStatus`]: this.state.keywordInput[name].trim() ? ValidateStatusEnum.Normal : ValidateStatusEnum.Empty,
            },
        })
    }

    /**
     * 保存操作
     */
    protected handleRequestSaveKeyword = async (): Promise<void> => {
        try {
            /**
             * 判断输入内容是否合法
             */
            const isValidate: boolean = Object.keys(this.state.validateStatus).every((key) => {
                return this.state.validateStatus[key] === ValidateStatusEnum.Normal;
            })

            if (isValidate) {
                const {
                    keywordInput: {
                        departNameKeys,
                        departThirdIdKeys,
                        loginNameKeys,
                        displayNameKeys,
                        emailKeys,
                        userThirdIdKeys,
                        subOuFilter,
                        subUserFilter,
                        baseFilter,
                        groupKeys,
                        statusKeys,
                        idcardNumberKeys,
                        telNumberKeys,
                    },
                } = this.state;

                /**
                 * 设置域关键字配置信息
                 */
                await setDomainKeyConfig([
                    this.props.domainInfo.id,
                    {
                        ncTUsrmDomainKeyConfig: {
                            departNameKeys: departNameKeys.trim().split(','),
                            departThirdIdKeys: departThirdIdKeys.trim().split(','),
                            loginNameKeys: loginNameKeys.trim().split(','),
                            displayNameKeys: displayNameKeys.trim().split(','),
                            emailKeys: emailKeys.trim().split(','),
                            userThirdIdKeys: userThirdIdKeys.trim().split(','),
                            subOuFilter: subOuFilter.trim(),
                            subUserFilter: subUserFilter.trim(),
                            baseFilter: baseFilter.trim(),
                            groupKeys: groupKeys.trim().split(','),
                            statusKeys: statusKeys.trim().split(','),
                            idcardNumberKeys: idcardNumberKeys.trim().split(','),
                            telNumberKeys: telNumberKeys.trim().split(','),
                        },
                    },
                ])

                /**
                 * 打印日志
                 */
                manageLog(
                    ManagementOps.SET,
                    __('编辑 同步关键字设置 成功'),
                    __('子部门搜索Filter为“${subOuFilter}”; 子用户搜索Filter为“${subUserFilter}”; 部门和用户Filter为“${baseFilter}”; 部门名关键字为“${departNameKeys}”; 部门ID关键字为“${departThirdIdKeys}”; 登录名关键字为“${loginNameKeys}”; 显示名关键字为“${displayNameKeys}”; 邮箱关键字为“${emailKeys}”; 用户ID关键字为“${userThirdIdKeys}”; 身份证号关键字为“${idcardNumberKeys}”; 手机号码关键字为“${telNumberKeys}”; 安全组关键字为“${groupKeys}”; 禁用关键字为“${statusKeys}”',
                        {
                            departNameKeys: departNameKeys.trim(),
                            departThirdIdKeys: departThirdIdKeys.trim(),
                            loginNameKeys: loginNameKeys.trim(),
                            displayNameKeys: displayNameKeys.trim(),
                            emailKeys: emailKeys.trim(),
                            userThirdIdKeys: userThirdIdKeys.trim(),
                            subOuFilter: subOuFilter.trim(),
                            subUserFilter: subUserFilter.trim(),
                            baseFilter: baseFilter.trim(),
                            groupKeys: groupKeys.trim(),
                            statusKeys: statusKeys.trim(),
                            idcardNumberKeys: idcardNumberKeys.trim(),
                            telNumberKeys: telNumberKeys.trim(),
                        },
                    ),
                    Level.INFO,
                )

                this.setState({
                    keywordInput: {
                        departNameKeys: departNameKeys.trim(),
                        departThirdIdKeys: departThirdIdKeys.trim(),
                        loginNameKeys: loginNameKeys.trim(),
                        displayNameKeys: displayNameKeys.trim(),
                        emailKeys: emailKeys.trim(),
                        userThirdIdKeys: userThirdIdKeys.trim(),
                        subOuFilter: subOuFilter.trim(),
                        subUserFilter: subUserFilter.trim(),
                        baseFilter: baseFilter.trim(),
                        groupKeys: groupKeys.trim(),
                        statusKeys: statusKeys.trim(),
                        idcardNumberKeys: idcardNumberKeys.trim(),
                        telNumberKeys: telNumberKeys.trim(),
                    },
                    isEditStatus: false,
                }, () => {
                    /**
                     * 保存成功后更新源信息
                     */
                    this.originKeywordInput = {
                        ...this.state.keywordInput,
                    }
                })

                this.props.onRequestEditStatus(false);
            }
        } catch (ret) {
            if (ret.error && ret.error.errID && ret.error.errID === ErrorCode.SetDomainKeyConfigDomainInvalid) {
                const domainName = this.props.domainInfo.name.trim();
                this.props.onRequestDomainInvalid(domainName)
            } else if (ret.error && ret.error.errID && ret.error.errID === ErrorCode.DomainUnavailable) {
                Message2.info({ message: __('连接LDAP服务器失败，请检查域控制器地址是否正确，或者域控制器是否开启。'), zIndex: 10000 })
            } else {
                ret.error && ret.error.errMsg && Message2.alert({ message: ret.error.errMsg });
            }
        }
    }

    /**
     * 取消操作
     */
    protected handleCancelEdit = (): void => {
        /**
         * 还原编辑前的状态
         */
        this.setState({
            keywordInput: {
                ...this.originKeywordInput,
            },
            isEditStatus: false,
            validateStatus: {
                departNameKeysValidateStatus: ValidateStatusEnum.Normal,
                departThirdIdKeysValidateStatus: ValidateStatusEnum.Normal,
                loginNameKeysValidateStatus: ValidateStatusEnum.Normal,
                displayNameKeysValidateStatus: ValidateStatusEnum.Normal,
                emailKeysValidateStatus: ValidateStatusEnum.Normal,
                userThirdIdKeysValidateStatus: ValidateStatusEnum.Normal,
                subOuFilterValidateStatus: ValidateStatusEnum.Normal,
                subUserFilterValidateStatus: ValidateStatusEnum.Normal,
                baseFilterValidateStatus: ValidateStatusEnum.Normal,
            },
        })

        this.props.onRequestEditStatus(false);
    }
}