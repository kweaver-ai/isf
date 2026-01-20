import * as React from 'react';
import classNames from 'classnames';
import { DataGrid, Button, Select, Trigger } from '@/sweet-ui';
import { Text, SearchBox, Title, UIIcon } from '@/ui/ui.desktop';
import { NodeType } from '@/core/organization';
import OrganizationPicker from '../../OrganizationPicker/component.view';
import { ListTipStatus, ListTipMessage } from '../../ListTipComponent/helper';
import ListTipComponent from '../../ListTipComponent/component.view';
import { MemberType, Limit } from '../helper';
import { SearchType, SearchTypeText } from './type';
import MemberGridBase from './component.base';
import __ from './locale';
import styles from './styles.view';

const listTipMessage = {
    ...ListTipMessage,
    [ListTipStatus.Empty]: __('请点击左上角【添加成员】，添加用户组成员'),
}

const NodeTypeIcon: Record<MemberType, string> = {
    [MemberType.Dep]: '\uf009',
    [MemberType.User]: '\uf007',
}

export default class MemberGrid extends MemberGridBase {
    render() {
        const {
            data: {
                members,
                total,
                page,
            },
            searchKey,
            isShowAdd,
            selections,
            listTipStatus,
            searchBoxIsOnFocus,
            searchType,
        } = this.state

        return (
            <div className={classNames(
                styles['member-grid'],
            )}>
                <div className={styles['header']}>
                    <Button
                        role={'sweetui-button'}
                        theme={'oem'}
                        icon={'add'}
                        size={'auto'}
                        disabled={!this.props.selectedGroup}
                        onClick={() => this.changePickerStatus(true)}
                    >
                        {__('添加成员')}
                    </Button>
                    {
                        selections.length ?
                            <Button
                                role={'sweetui-button'}
                                className={styles['btn']}
                                icon={'bin'}
                                size={'auto'}
                                disabled={!this.props.selectedGroup}
                                onClick={this.deleteMembers}
                            >
                                {__('删除成员')}
                            </Button>
                            : null
                    }
                    <div className={styles['search-area']}>
                        <Select
                            width={134}
                            maxMenuHeight={190}
                            className={styles['search-scope']}
                            value={searchType}
                            onChange={({ detail }) => this.changeSearchType(detail)}
                        >
                            {
                                [SearchType.UserGroupMemberName, SearchType.UserDisplayNmae].map((item) => (
                                    <Select.Option
                                        key={item}
                                        value={item}
                                    >
                                        {SearchTypeText[item]}
                                    </Select.Option>
                                ))
                            }
                        </Select>
                        <SearchBox
                            role={'ui-searchbox'}
                            ref={(searchBox) => this.searchBox = searchBox}
                            className={classNames(styles['search-box'], { [styles['focus']]: searchBoxIsOnFocus })}
                            width={170}
                            placeholder={searchType === SearchType.UserGroupMemberName ? __('搜索用户组成员 ') : __('请输入用户显示名')}
                            disabled={!this.props.selectedGroup}
                            value={searchKey}
                            onChange={this.changeSearchKey}
                            loader={this.getMembers}
                            onLoad={this.loadMembers}
                            onLoadFailed={this.loadFailed}
                            onFocus={() => this.setState({ searchBoxIsOnFocus: true })}
                            onBlur={() => this.setState({ searchBoxIsOnFocus: false })}
                        />
                    </div>
                </div>
                <div className={
                    classNames(
                        styles['grid-wrapper'],
                    )
                }>
                    <DataGrid
                        role={'swwetui-datagrid'}
                        ref={(dataGrid) => this.dataGrid = dataGrid}
                        height={'100%'}
                        data={members}
                        showBorder={true}
                        enableSelect={true}
                        enableMultiSelect={true}
                        enableVirtualScroll={false}
                        DataGridHeader={{ enableSelectAll: true }}
                        selection={selections}
                        refreshing={listTipStatus !== ListTipStatus.None}
                        RefreshingComponent={
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                                listTipMessage={listTipMessage}
                            />
                        }
                        onSelectionChange={({ detail }) => this.changeSelection(detail)}
                        columns={this.getColumns(searchType)}
                        DataGridPager={
                            {
                                size: Limit,
                                total,
                                page,
                                onPageChange: ({ detail: { page } }) => this.handlePageChange(page),
                            }
                        }
                    />
                </div>
                {
                    isShowAdd ?
                        <OrganizationPicker
                            zIndex={20}
                            isCascadeTree={true}
                            convererOut={this.convererOutData}
                            selectType={[
                                NodeType.USER,
                                NodeType.DEPARTMENT,
                                NodeType.ORGANIZATION,
                            ]}
                            onCancel={() => this.changePickerStatus(false)}
                            onConfirm={this.confirmAdd}
                            title={__('添加用户或部门')}
                        />
                        : null
                }
            </div >
        )
    }

    /**
     * 转出数据时转换数据格式
     */
    protected convererOutData = ({ id, type, name }: { id: string; type: number; name: string }): { id: string; type: string; name: string } => {
        return {
            id,
            name,
            type: type === NodeType.USER ? MemberType.User : MemberType.Dep,
        }
    }

    /**
     * 通过成员类型分类渲染成员
     */
    private renderMembersByType = (memberType: MemberType, groupMembers) => {
        const filteredMembers = groupMembers.filter(({ type }) => type === memberType)

        return filteredMembers.length ? (
            <>
                <UIIcon
                    className={styles['node-icon']}
                    size={16}
                    code={NodeTypeIcon[memberType]}
                />
                {
                    filteredMembers.map(({ name, id, parent_deps }, index) => (
                        <Title
                            key={id}
                            inline={true}
                            content={parent_deps.length && memberType === MemberType.Dep ? parent_deps.map((dep) => dep.map(({ name }) => name).join('/').concat(`/${name}`)).join('\n') : name}
                        >
                            {`${name}${index !== filteredMembers.length - 1 ? __('，') : '\u00A0\u00A0'}`}
                        </Title>
                    ))
                }
            </>
        ) : null
    }

    /**
     * 获取columns
     */
    protected getColumns = (type: SearchType): ReadonlyArray<any> => {
        const { searchKey } = this.state
        return type === SearchType.UserDisplayNmae && !!searchKey
            ? [
                {
                    title: __('用户名'),
                    key: 'name',
                    width: 25,
                    renderCell: (name, record) => (
                        <div className={styles['icon-text']}>
                            <UIIcon
                                className={styles['node-icon']}
                                size={16}
                                code={NodeTypeIcon[record.type]}
                            />
                            <Text role={'ui-text'}>
                                {name}
                            </Text>
                        </div>
                    ),
                },
                {
                    title: __('所在用户组成员'),
                    key: 'group_members',
                    width: 35,
                    renderCell: (group_members, record) => (
                        <Trigger
                            triggerEvent={'click'}
                            anchorOrigin={['left', 'bottom']}
                            alignOrigin={['left', 'top']}
                            freeze={false}
                            renderer={({ setPopupVisibleOnClick }) => (
                                <span
                                    key={'accountDropMenuTrigger'}
                                    className={styles['info']}
                                    onClick={(e) => {e.stopPropagation(); setPopupVisibleOnClick && setPopupVisibleOnClick()}}
                                >
                                    {__('查看')}
                                </span>
                            )}
                        >
                            <div className={styles['departments']}>
                                <div className={styles['depts-header']}>{__('所在用户组成员')}</div>
                                <div className={styles['depts-info']}>
                                    {
                                        group_members && this.renderMembersByType(MemberType.User, group_members)
                                    }
                                    {
                                        group_members && this.renderMembersByType(MemberType.Dep, group_members)
                                    }
                                </div>
                            </div>
                        </Trigger>
                    ),
                },
                {
                    title: __('直属部门'),
                    key: 'parent_deps',
                    width: 40,
                    renderCell: (parent_deps, record) => (
                        <Trigger
                            triggerEvent={'click'}
                            anchorOrigin={['left', 'bottom']}
                            alignOrigin={['left', 'top']}
                            freeze={false}
                            renderer={({ setPopupVisibleOnClick }) => (
                                <span
                                    key={'accountDropMenuTrigger'}
                                    className={styles['info']}
                                    onClick={(e) => {e.stopPropagation(); setPopupVisibleOnClick && setPopupVisibleOnClick()}}
                                >
                                    {__('查看')}
                                </span>
                            )}
                        >
                            <div className={styles['departments']}>
                                <div className={styles['depts-header']}>{__('直属部门')}</div>
                                <div className={styles['depts-info']}>
                                    {
                                        parent_deps && parent_deps.length ?
                                            parent_deps.map((dep, index) => (
                                                <Title
                                                    key={dep.id}
                                                    inline={true}
                                                    content={dep.map(({ name }) => name).join('/')}
                                                >
                                                    {`${dep[dep.length - 1].name}${index !== parent_deps.length - 1 ? __('，') : ''}`}
                                                </Title>
                                            ))
                                            :
                                            (
                                                <Text role={'ui-text'}>{__('未分配组')}</Text>
                                            )
                                    }
                                </div>
                            </div>
                        </Trigger>
                    ),
                },
            ]
            : [
                {
                    title: __('用户组成员'),
                    key: 'name',
                    width: 50,
                    renderCell: (name, record) => (
                        <div className={styles['icon-text']}>
                            <UIIcon
                                className={styles['node-icon']}
                                size={16}
                                code={NodeTypeIcon[record.type]}
                            />
                            <div className={styles['name']}>
                                <Text role={'ui-text'}>
                                    {name}
                                </Text>
                            </div>
                        </div>
                    ),
                },
                {
                    title: __('直属部门'),
                    key: 'department_names',
                    width: 50,
                    renderCell: (department_names, record) => (
                        <Trigger
                            triggerEvent={'click'}
                            anchorOrigin={['left', 'bottom']}
                            alignOrigin={['left', 'top']}
                            freeze={false}
                            renderer={({ setPopupVisibleOnClick }) => (
                                <span
                                    key={'accountDropMenuTrigger'}
                                    className={styles['info']}
                                    onClick={(e) => {e.stopPropagation(); setPopupVisibleOnClick && setPopupVisibleOnClick()}}
                                >
                                    {__('查看')}
                                </span>
                            )}
                        >
                            <div className={styles['departments']}>
                                <div className={styles['depts-header']}>{__('直属部门')}</div>
                                <div className={styles['depts-info']}>
                                    {
                                        department_names && department_names.length ? department_names.map((depName, index) =>(
                                            <Title
                                                key={record.parent_deps[index].slice(-1)[0].id}
                                                inline={true}
                                                content={record.parent_deps[index].map(({ name }) => name).join('/')}
                                            >
                                                {`${depName}${index !== department_names.length - 1 ? __('，') : ''}`}
                                            </Title>))
                                            : record.type === MemberType.Dep ?
                                                record.name
                                                : __('未分配组')
                                    }
                                </div>
                            </div>
                        </Trigger>
                    )
                    ,
                },
            ]
    }
}