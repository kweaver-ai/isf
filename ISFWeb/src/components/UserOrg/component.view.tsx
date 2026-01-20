import * as React from 'react'
import { SystemRoleType } from '@/core/role/role'
import { ToolBar, UIIcon } from '@/ui/ui.desktop'
import { PopMenu, Text } from '@/sweet-ui'
import SecurityIntegration from '../SecurityIntegration/component.view'
import DisplayManager from '../DisplayManager/component.view'
import CreateOrganization from '../CreateOrganization/component.view'
import EditOrganization from '../EditOrganization/component.view'
import DeleteOrganization from '../DeleteOrganization/component.view'
import CreateDepartment from '../CreateDepartment/component.view'
import EditDepartment from '../EditDepartment/component.view'
import DeleteDepartment from '../DeleteDepartment/component.view'
import MoveDepartment from '../MoveDepartment/component.view'
import AddUsersToDep from '../AddUsersToDep/component.view'
import CreateUser from '../CreateUser/component.view'
import EditUser from '../EditUser/component.view'
import DeleteUser from '../DeleteUser/component.view'
import SetUserExpireTime from '../SetUserExpireTime/component.view'
import MoveUser from '../MoveUser/component.view'
import RemoveUser from '../RemoveUser/component.view'
import EnableUser from '../EnableUser/component.view'
import DisableUser from '../DisableUser/component.view'
import SetRole from '../SetRole/component.view'
import PwdManage from '../PwdManage/component.view'
import ExportImportOrganize from '../ExportImportOrganize/component.view'
import ImportDomainUser from '../ImportDomainUser/component.view'
import ImportOrganization from '../ImportOrganization/component.view'
import SetUsersFreezeStatus from '../SetUsersFreezeStatus/component.view'
import BatchEditUser from '../BatchEditUser/component.view'
import { Range, TabEnum } from './helper'
import OrgTree from './OrgTree/component.view'
import UserGrid from './UserGrid/component.view'
import UserOrgMgntBase from './component.base'
import { Action, EnableStatus } from './helper'
import __ from './locale'
import styles from './styles.view';
import DepartmentGrid from './DepartmentGrid'
import UserGroup from '../UserGroup/component.view'
import { UseAccountMgnt } from '../UseAccountMgnt/index'
import { Button, Popover, Tabs } from 'antd'
import { ProductLicense } from '../ProductLicense'
import ProductLicenseIcon from '../../icons/product-license.svg'
import { ProductLicenseOverview } from '../ProductLicense/Overview'
import intl from 'react-intl-universal'

export default class UserOrgMgnt extends UserOrgMgntBase {
    render() {
        const {
            menus,
            selectedDep,
            isShowSetRole,
            isShowInit,
            urlParams
        } = this.state

        const isAdmin = this.isAdmin()

        let tabs = [{id: TabEnum.User, label: __('用户')}, {id: TabEnum.Department, label: __('部门')}]

        const isSuperOrAdmin = this.isSuperOrAdmin()

        tabs =[...tabs, ...(isSuperOrAdmin? [{id: TabEnum.UserGroup, label: __('用户组')}, {id: TabEnum.AppAccount, label: __('应用账户')}] : [])]

        return (
            <div className={styles['container']}>
                <Tabs activeKey={urlParams?.get("tab") || TabEnum.User} onChange={this.onChangeTab} destroyOnHidden={true}>
                    {
                        tabs.map((tab) => {
                            return (
                                <Tabs.TabPane key={tab.id} tab={tab.label} className={styles['tab']}>
                                    {
                                        tab.id === TabEnum.User ? 
                                            <div className={styles['user-container']}>
                                                {
                                                    this.isKjzDisabled ?
                                                        <div className={styles['head-bar']}>
                                                            <ToolBar role={'ui-toolbar'}>
                                                                {
                                                                    menus.map((menu) => {
                                                                        return menu ?
                                                                            (
                                                                                <PopMenu
                                                                                    role={'sweetui-popmenu'}
                                                                                    key={menu.name}
                                                                                    alignOrigin={[-20, 0]}
                                                                                    triggerEvent={'hover'}
                                                                                    onRequestCloseWhenClick={(close) => close()}
                                                                                    freeze={false}
                                                                                    trigger={({ setPopupVisibleOnMouseEnter, setPopupVisibleOnMouseLeave }) =>
                                                                                        <div
                                                                                            key={menu.name}
                                                                                            className={styles['menu']}
                                                                                            onMouseEnter={setPopupVisibleOnMouseEnter}
                                                                                            onMouseLeave={setPopupVisibleOnMouseLeave}
                                                                                        >
                                                                                            <UIIcon
                                                                                                role={'ui-uiicon'}
                                                                                                code={'\u0000'}
                                                                                                fallback={menu.fallback}
                                                                                                size={16}
                                                                                            />
                                                                                            <div className={styles['menu-name']}><Text role={'sweetui-text'}>{menu.text}</Text></div>
                                                                                            <UIIcon
                                                                                                role={'ui-uiicon'}
                                                                                                code={'\uf00b'}
                                                                                                size={12}
                                                                                            />
                                                                                        </div>
                                                                                    }
                                                                                >
                                                                                    {
                                                                                        menu.actions && menu.actions.map((action, i) =>
                                                                                            action ?
                                                                                                action.type === Action.Line ?
                                                                                                    <span
                                                                                                        className={styles['line']}
                                                                                                        key={i + 'line'}
                                                                                                    >
                                                                                                    </span>
                                                                                                    :
                                                                                                    <PopMenu.Item
                                                                                                        role={'sweetui-popmenu.item'}
                                                                                                        className={styles['action']}
                                                                                                        key={action.type}
                                                                                                        label={action.text}
                                                                                                        disabled={action.disabled}
                                                                                                        onClick={(e) => action.disabled ? e.stopPropagation() : this.changeAction(action.type)}
                                                                                                    />
                                                                                                : null,
                                                                                        )
                                                                                    }
                                                                                </PopMenu>
                                                                            )
                                                                            : null
                                                                    })
                                                                }
                                                                {
                                                                    isSuperOrAdmin &&
                                                                    <div className={styles['authorized-overview']}>
                                                                        <Popover
                                                                            trigger={'click'}
                                                                            placement={'bottom'}
                                                                            destroyOnHidden={true}
                                                                            title={intl.get('authorized.user.count.overview')}
                                                                            content={<ProductLicenseOverview />}
                                                                        >
                                                                            <Button
                                                                                color="default"
                                                                                variant={'link'}
                                                                                icon={<ProductLicenseIcon style={{ width: 14, height: 14 }} />}
                                                                            >
                                                                                {intl.get('authorized.overview')}
                                                                            </Button>
                                                                        </Popover>
                                                                    </div>
                                                                }
                                                            </ToolBar>
                                                        </div>
                                                        : null
                                                }
                                                <div className={styles['main-content']}>
                                                    <div className={styles['org']}>
                                                        <div className={styles['org-tree']}>
                                                            <OrgTree
                                                                ref={(orgTree) => this.orgTree = orgTree}
                                                                userid={this.userid}
                                                                isKjzDisabled={this.isKjzDisabled}
                                                                onRequestSelectDep={this.selectDep}
                                                            />
                                                        </div>
                                                        {
                                                            !!selectedDep && !isAdmin && (isShowSetRole === SystemRoleType.OrgManager) && this.isKjzDisabled
                                                                ? (
                                                                    <div className={styles['dep-manager-display']}>
                                                                        <DisplayManager
                                                                            userid={this.userid}
                                                                            departmentId={selectedDep.data.id}
                                                                            departmentName={selectedDep.data.name}
                                                                            hasPermission={!!this.userInfo.user.roles.find((role) => role.id === SystemRoleType.Securit || SystemRoleType.Supper)}
                                                                            onComplete={() => this.userGrid.updateCurrentPage()}
                                                                        />
                                                                    </div>
                                                                )
                                                                : null
                                                        }
                                                    </div>
                                                    {
                                                        menus.length ?
                                                            <div className={styles['grid-container']}>
                                                                <UserGrid
                                                                    ref={(userGrid) => this.userGrid = userGrid}
                                                                    userid={this.userid}
                                                                    selectedDep={selectedDep && selectedDep.data}
                                                                    freezeStatus={this.freezeStatus}
                                                                    isShowSetRole={!!isShowSetRole}
                                                                    isShowEnableAndDisableUser={this.isShowEnableAndDisableUser}
                                                                    isKjzDisabled={this.isKjzDisabled}
                                                                    onRequestSelectUsers={this.selectUsers}
                                                                    onRequestDelDepNode={(dep) => { this.orgTree.deleteNode(dep) }}
                                                                />
                                                            </div>
                                                            : null
                                                    }
                                                </div>
                                            </div>
                                            : tab.id === TabEnum.Department ? 
                                                <DepartmentGrid />
                                                : tab.id === TabEnum.UserGroup ?
                                                    <UserGroup /> :
                                                    tab.id === TabEnum.AppAccount ?
                                                        <UseAccountMgnt /> : null
                                    }
                                </Tabs.TabPane>
                            )
                        })
                    }
                </Tabs>
                {
                    isShowInit ?
                        <SecurityIntegration
                            onUpdateCsfLevelText={this.closeInit}
                        />
                        : null
                }
                {
                    this.rebderAction()
                }
            </div>
        )
    }

    /**
     * 根据Action渲染弹窗
     */
    protected rebderAction() {
        const { action, selectedDep, selectedUsers } = this.state

        if (selectedDep) {
            switch (action) {
                case Action.CreateOrg:
                    return <CreateOrganization
                        dep={selectedDep.data}
                        userid={this.userid}
                        onRequestCancelCreateOrg={() => this.changeAction(Action.None)}
                        onCreateOrgSuccess={(orgInfo) => this.addTreeNode(orgInfo, Action.CreateOrg)}
                    />

                case Action.EditOrg:
                    return <EditOrganization
                        dep={selectedDep.data}
                        userid={this.userid}
                        onRequestCancelEditOrg={() => this.changeAction(Action.None)}
                        onEditOrgSuccess={this.updateTreeNode}
                        onRequestDelOrg={(dep) => { this.changeAction(Action.None); this.orgTree.deleteNode(dep) }}
                    />

                case Action.DelOrg:
                    return <DeleteOrganization
                        dep={selectedDep.data}
                        userid={this.userid}
                        onRequestCancelDeleteOrg={() => this.changeAction(Action.None)}
                        onDeleteOrgSuccess={this.deleteTreeNode}
                    />

                case Action.CreateDep:
                    return <CreateDepartment
                        dep={selectedDep.data}
                        userid={this.userid}
                        onRequestCancelCreateDep={() => this.changeAction(Action.None)}
                        onCreateDepSuccess={(depInfo) => this.addTreeNode(depInfo, Action.CreateDep)}
                        onRequestDelDep={(dep) => { this.changeAction(Action.None); this.orgTree.deleteNode(dep) }}
                    />

                case Action.EditDep:
                    return <EditDepartment
                        dep={selectedDep.data}
                        parentName={selectedDep.parent ? selectedDep.parent.data.name : '' }
                        userid={this.userid}
                        onRequestCancelEditDep={() => this.changeAction(Action.None)}
                        onEditDepSuccess={this.updateTreeNode}
                        onRequestDelDep={(dep) => { this.changeAction(Action.None); this.orgTree.deleteNode(dep) }}
                    />

                case Action.DelDep:
                    return <DeleteDepartment
                        dep={selectedDep.data}
                        userid={this.userid}
                        onRequestCancelDeleteDep={() => this.changeAction(Action.None)}
                        onDeleteDepSuccess={this.deleteTreeNode}
                    />

                case Action.MoveDep:
                    return <MoveDepartment
                        srcDep={selectedDep.data}
                        onRequestCancelMoveDep={() => this.changeAction(Action.None)}
                        onRequestMoveDepFinished={this.moveDep}
                        onRequestRemoveSrcDep={(srcDep) => { this.changeAction(Action.None); this.orgTree.deleteNode(srcDep) }}
                    />

                case Action.AddUsersToDep:
                    return <AddUsersToDep
                        targetDep={selectedDep.data}
                        onRequestCancel={() => this.changeAction(Action.None)}
                        onRequestSuccess={() => { this.changeAction(Action.None); this.userGrid.backToFirst() }}
                        onRequestRemoveDep={(srcDep) => { this.changeAction(Action.None); this.orgTree.deleteNode(srcDep) }}
                    />

                case Action.CreateUser:
                    return <CreateUser
                        dep={selectedDep.data}
                        userid={this.userid}
                        onRequestCancel={() => this.changeAction(Action.None)}
                        onRequestCreateUserSuccess={(userInfo) => { this.changeAction(Action.None); this.userGrid.addUser(userInfo) }}
                        onRequestRemoveDep={(dep) => { this.changeAction(Action.None); this.orgTree.deleteNode(dep) }}
                    />

                case Action.EditUser:
                    if (selectedUsers && selectedUsers.length == 1) {
                        return <EditUser
                            selectUserId={selectedUsers[0].id}
                            dep={selectedDep.data}
                            userid={this.userid}
                            onRequestCancel={() => this.changeAction(Action.None)}
                            onRequestEditUserSuccess={(userInfo) => { this.changeAction(Action.None); this.userGrid.updateUserInfo(userInfo) }}
                        />
                    } else {
                        return <BatchEditUser
                            dep={selectedDep.data}
                            users={selectedUsers}
                            onRequestCancel={() => this.changeAction(Action.None)}
                            onRequestSuccess={(range) => this.updateGrid(range !== Range.USERS)}
                        />
                    }

                case Action.DelUser:
                    return <DeleteUser
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        shouldEnableUsers={false}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={(range) => this.updateGrid(range !== Range.USERS)}
                    />

                case Action.SetExpiration:
                    return <SetUserExpireTime
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        shouldEnableUsers={false}
                        onCancel={() => this.changeAction(Action.None)}
                        onSuccess={(range) => this.updateGrid(range !== Range.USERS)}
                    />

                case Action.MoveUser:
                    return <MoveUser
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        onComplete={(dep) => { this.changeAction(Action.None); dep && this.orgTree.deleteNode(dep) }}
                        onSuccess={(range) => this.updateGrid(range !== Range.USERS)}
                    />

                case Action.RemoveUser:
                    return <RemoveUser
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={(range) => this.updateGrid(range !== Range.USERS)}
                    />

                case Action.EnableUser:
                    return <EnableUser
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={(users) => this.enableUser(users.length > 1, { status: EnableStatus.Enabled })}
                    />

                case Action.DisableUser:
                    return <DisableUser
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={(users) => this.enableUser(users.length > 1, { status: EnableStatus.Disabled })}
                    />

                case Action.SetRole:
                    return <SetRole
                        users={selectedUsers}
                        dep={selectedDep.data}
                        userid={this.userid}
                        onComplete={() => { this.changeAction(Action.None); this.userGrid.updateCurrentPage() }}
                    />

                case Action.ProductLicense:
                    return <ProductLicense
                        users={selectedUsers}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={(userInfo) => { 
                            this.changeAction(Action.None); 
                            if(userInfo.length > 1) {
                               this.updateGrid()
                            }else {
                               this.userGrid.updateUserInfo(userInfo?.[0])
                            }
                         }}
                    />

                case Action.ManagePwd:
                    return <PwdManage
                        selectedUser={selectedUsers[0]}
                        userid={this.userid}
                        onRequestConfirm={() => this.changeAction(Action.None)}
                        onRequestCancel={() => this.changeAction(Action.None)}
                    />

                case Action.ImportOrg:
                    return <ExportImportOrganize
                        users={selectedUsers}
                        dep={selectedDep}
                        userid={this.userid}
                        status={true}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={() => { this.changeAction(Action.None); this.orgTree.initOrgTree() }}
                    />

                case Action.ImportDomain:
                    return <ImportDomainUser
                        departmentId={selectedDep.data.id}
                        onRequestSuccess={this.importDomainUser}
                        onRequestCancel={() => this.changeAction(Action.None)}
                        doRedirectDomain={this.doRedirectDomain}
                    />

                case Action.ImportThirdOrg:
                    return <ImportOrganization
                        departmentId={selectedDep.data.id}
                        userid={this.userid}
                        onCancel={() => this.changeAction(Action.None)}
                        onSuccess={() => { this.changeAction(Action.None); this.orgTree.initOrgTree() }}
                        onComplete={() => { this.changeAction(Action.None); this.orgTree.initOrgTree() }}
                    />

                case Action.FreezeUser:
                    return <SetUsersFreezeStatus
                        userid={this.userid}
                        freezeStatus={true}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={() => { this.changeAction(Action.None); this.userGrid.updateCurrentPage() }}
                    />
                case Action.UnfreezeUser:
                    return <SetUsersFreezeStatus
                        userid={this.userid}
                        freezeStatus={false}
                        onComplete={() => this.changeAction(Action.None)}
                        onSuccess={() => { this.changeAction(Action.None); this.userGrid.updateCurrentPage() }}
                    />

                default:
                    return null
            }
        }

        return null
    }
}