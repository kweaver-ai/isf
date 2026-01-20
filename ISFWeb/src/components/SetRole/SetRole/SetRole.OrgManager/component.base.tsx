import { noop, values } from 'lodash';
import { getRoleName, SystemRoleType } from '@/core/role/role';
import { NodeType } from '@/core/organization';
import { Toast } from '@/sweet-ui';
import WebComponent from '../../../webcomponent';
import __ from './locale';

export enum ValidateState {
    Normal,
    Empty,
    InvalidSpace,
}

export const ValidateMessages = {
    [ValidateState.Empty]: __('此项不允许为空。'),
    [ValidateState.InvalidSpace]: __('配额空间值为不超过 1000000 的正数，支持小数点后两位，请重新输入。'),
}

export default class SetOrgManagerBase extends WebComponent<Console.SetOrgManager.Props, Console.SetOrgManager.State> {
    static defaultProps = {
        userid: '',
        editRateInfo: null,
        roleInfo: null,
        userInfo: null,
        directDeptInfo: null,
        limitSpaceInfo: null,
        onConfirmSetRoleConfig: noop,
        onCancelSetRoleConfig: noop,
    }

    state = {
        validateState: {
            userSpace: ValidateState.Normal,
            docSpace: ValidateState.Normal,
        },
        limitCheckStatus: {
            limitUserCheckSatus: false,
            limitDocCheckSatus: false,
        },
        limitCheckDisable: {
            limitUserCheckDisable: false,
            limitDocCheckDisable: false,
        },
        spaceConfig: {
            userSpace: '',
            docSpace: '',
        },
        selectDeps: [],
        selectState: false,
    }

    componentDidMount() {
        const { directDeptInfo, limitSpaceInfo: { limitUserSpace, limitDocSpace }, roles } = this.props;

        // 如果是组织管理员，并且被限额了，则限额复选框勾选并灰化，显示限额信息
        if (roles.some((item) => item.id === SystemRoleType.OrgManager) && (limitDocSpace !== -1 || limitUserSpace !== -1)) {
            this.setState({
                spaceConfig: {
                    userSpace: limitUserSpace === -1 ? '' : (limitUserSpace / Math.pow(1024, 3)).toFixed(2),
                    docSpace: limitDocSpace === -1 ? '' : (limitDocSpace / Math.pow(1024, 3)).toFixed(2),
                },
                limitCheckStatus: {
                    limitUserCheckSatus: limitUserSpace === -1 ? false : true,
                    limitDocCheckSatus: limitDocSpace === -1 ? false : true,
                },
                limitCheckDisable: {
                    limitUserCheckDisable: limitUserSpace === -1 ? false : true,
                    limitDocCheckDisable: limitDocSpace === -1 ? false : true,
                },
            })
        }

        // 如果是编辑，优先显示被编辑管理员信息
        if (this.props.editRateInfo) {
            const { editRateInfo: { manageDeptInfo } } = this.props;
            if (manageDeptInfo) {
                const { editRateInfo: { manageDeptInfo: { limitUserSpaceSize, limitDocSpaceSize } } } = this.props;

                this.setState({
                    selectDeps: manageDeptInfo.departmentIds.length ? manageDeptInfo.departmentIds.map((cur, index) => (
                        {
                            objectId: cur,
                            objectName: manageDeptInfo.departmentNames[index],
                            objType: NodeType.DEPARTMENT,
                        }
                    )) : [],
                    spaceConfig: {
                        userSpace: limitUserSpaceSize === -1 ? '' : (limitUserSpaceSize / Math.pow(1024, 3)).toFixed(2),
                        docSpace: limitDocSpaceSize === -1 ? '' : (limitDocSpaceSize / Math.pow(1024, 3)).toFixed(2),
                    },
                    limitCheckStatus: {
                        limitUserCheckSatus: limitUserSpaceSize === -1 && limitUserSpace == -1 ? false : true,
                        limitDocCheckSatus: limitDocSpaceSize === -1 && limitDocSpace == -1 ? false : true,
                    },
                })
            }
        } else {
            if (directDeptInfo && directDeptInfo.departmentId !== '-1') {
                this.setState({
                    selectDeps: [
                        {
                            objectId: directDeptInfo.departmentId,
                            objectName: directDeptInfo.departmentName,
                            objType: NodeType.DEPARTMENT,
                        },
                    ],
                })
            }
        }
    }

    /**
     * 转入前先转换数据格式
     * @param data
     */
    protected convertData = (data) => {
        return {
            id: data.objectId,
            name: data.objectName,
            type: data.objType,
        }

    }

    /**
     * 转出数据时转换数据格式
     */
    protected convertDataOut(data) {
        return {
            objectId: data.id,
            objectName: data.name || data.displayName || data.departmentName || (data.user && data.user.displayName),
            objType: data.type,
        }
    }

    /**
     * 选择部门
     */
    protected selectDeparment(data) {
        this.setState({
            selectDeps: data,
        }, () => {
            if (this.state.selectDeps.length) {
                this.setState({
                    selectState: false,
                })
            }
        })
    }

    /**
     * 勾选/取消复选框
     * @param status 点击状态
     */
    protected handleCheckStateChange(status = {}) {
        const { limitCheckStatus, validateState, spaceConfig } = this.state
        this.setState({
            validateState: {
                ...validateState,
                userSpace: ValidateState.Normal,
                docSpace: ValidateState.Normal,
            },
            limitCheckStatus: {
                ...limitCheckStatus,
                limitUserCheckSatus: 'userManageSpace' in status ? status.userManageSpace : limitCheckStatus.limitUserCheckSatus,
                limitDocCheckSatus: 'docManageSpace' in status ? status.docManageSpace : limitCheckStatus.limitDocCheckSatus,
            },
            spaceConfig: {
                userSpace: 'userManageSpace' in status ? '' : spaceConfig.userSpace,
                docSpace: 'docManageSpace' in status ? '' : spaceConfig.docSpace,
            },
        })
    }

    /**
     * 验证输入值以及是否选择部门
     */
    protected validateRole = () => {
        const { validateState, limitCheckStatus, spaceConfig, selectDeps } = this.state;
        this.setState({
            validateState: {
                ...validateState,
                userSpace: limitCheckStatus.limitUserCheckSatus ?
                    spaceConfig.userSpace ?
                        Number(spaceConfig.userSpace) <= 1000000 && Number(spaceConfig.userSpace) > 0 ?
                            ValidateState.Normal : ValidateState.InvalidSpace
                        : ValidateState.Empty
                    : validateState.userSpace,
                docSpace: limitCheckStatus.limitDocCheckSatus ?
                    spaceConfig.docSpace ?
                        Number(spaceConfig.docSpace) <= 1000000 && Number(spaceConfig.docSpace) > 0 ?
                            ValidateState.Normal : ValidateState.InvalidSpace
                        : ValidateState.Empty
                    : validateState.docSpace,
            },
            selectState: !selectDeps.length ? true : false,
        }, () => {
            // 如果输入框为空，不进行toast提示
            if (
                this.state.validateState.userSpace !== ValidateState.Empty
                && this.state.validateState.userSpace !== ValidateState.InvalidSpace
                && this.state.validateState.docSpace !== ValidateState.Empty
                && this.state.validateState.docSpace !== ValidateState.InvalidSpace
            ) {
                // 判断编辑的组织管理员限额是否超过当前登录的组织管理员
                const { limitSpaceInfo: { limitUserSpace, limitDocSpace }, roles } = this.props

                if (roles.some((item) => item.id === SystemRoleType.OrgManager) && (limitDocSpace !== -1 || limitUserSpace !== -1)) {
                    if (this.props.editRateInfo) {
                        const { editRateInfo: { manageDeptInfo } } = this.props;
                        const { limitCheckStatus: { limitUserCheckSatus, limitDocCheckSatus }, spaceConfig: { userSpace, docSpace } } = this.state
                        const userSize = (limitUserSpace / Math.pow(1024, 3)).toFixed(2)
                        const docSize = (limitDocSpace / Math.pow(1024, 3)).toFixed(2)

                        if (manageDeptInfo) {
                            if (
                                (limitUserCheckSatus && limitUserSpace !== -1 && (Number(userSpace) > Number(userSize) || userSpace === ''))
                                && (limitDocCheckSatus && limitDocSpace !== -1 && (Number(docSpace) > Number(docSize) || docSpace === ''))
                            ) {
                                Toast.open(
                                    __('配额空间不足，用户管理可配空间最大限额为${userSize}、文档管理可配空间最大限额为${docSize}。',
                                        { userSize, docSize }),
                                )
                                return
                            } else if (limitUserCheckSatus && limitUserSpace !== -1 && (Number(userSpace) > Number((limitUserSpace / Math.pow(1024, 3)).toFixed(2)) || userSpace === '')) {
                                Toast.open(
                                    __('配额空间不足，用户管理可配空间最大限额为${userSize}。',
                                        { userSize }),
                                )
                                return
                            } else if (limitDocCheckSatus && limitDocSpace !== -1 && (Number(docSpace) > Number((limitDocSpace / Math.pow(1024, 3)).toFixed(2)) || docSpace === '')) {
                                Toast.open(
                                    __('配额空间不足，文档管理可配空间最大限额为${docSize}。',
                                        { docSize }),
                                )
                                return
                            }
                        }
                    }
                }
            }

            if (!(values(this.state.validateState).some((state) => state !== ValidateState.Normal)) && !this.state.selectState) {
                this.confirmRoleRateConfig()
            }
        })
    }

    /**
     * 将数据传出去
     */
    private confirmRoleRateConfig() {
        if (this.state.selectDeps.length) {
            let depInfo = this.state.selectDeps.reduce((pre, cur) => (
                {
                    depIds: [...pre.depIds, cur.objectId],
                    depNames: [...pre.depNames, cur.objectName],
                }
            ), { depIds: [], depNames: [] })
            let manageRange = {
                ncTManageDeptInfo: {
                    departmentIds: depInfo.depIds,
                    departmentNames: depInfo.depNames,
                    limitUserSpaceSize: -1,
                    limitDocSpaceSize: -1,
                },
            }
            this.props.onConfirmSetRoleConfig({
                name: getRoleName(this.props.roleInfo),
                id: this.props.roleInfo.id,
                manageRange,
            })
        }
    }

    /**
     * 取消本次操作
     */
    protected cancelSetRoleConfig = () => {
        this.setState({
            selectDeps: [],
        }, () => {
            this.props.onCancelSetRoleConfig();
        })
    }
}