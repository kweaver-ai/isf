import * as React from 'react'
import { SearchBox, Text, UIIcon } from '@/ui/ui.desktop'
import { DataGrid,  SweetIcon, Switch, Select, Message2, PopMenu, Button } from '@/sweet-ui'
import { Icon, Title } from '@/ui/ui.desktop';
import __ from './locale'
import styles from './styles'
import * as loading from './assets/loading.gif';
import { ListTipStatus } from '../../ListTipComponent/helper';
import ListTipComponent from '../../ListTipComponent/component.view';
import { getDepartments, searchDepartments } from '@/core/apis/console/usermanagement';
import { editDepartment, editOrganization, getOrgDepartmentById } from '@/core/thrift/sharemgnt/sharemgnt';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { displayUserOssInfo } from '@/core/oss/oss';
import { Action, formatDepInfo } from '../helper';
import EditDepartment from '../../EditDepartment/component.view'
import session from '@/util/session';
import CreateDepartment from '../../CreateDepartment/component.view';
import DeleteDepartment from '../../DeleteDepartment/component.view';
import CreateOrganization from '../../CreateOrganization/component.view';
import EditOrganization from '../../EditOrganization/component.view';
import DeleteOrganization from '../../DeleteOrganization/component.view';
import { UserRole, getRoleTypes } from '@/core/role/role';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { getConfidentialConfig } from '@/core/apis/eachttp/config/config';
import { getErrorMessage } from '@/core/exception';

const { useState, useEffect, useRef } = React

interface Department {
    id: string;
    name: string;
    code: number;
    level: number;
    enabled: boolean;
    manager: {id: string; name: string;type: string }|null;
    expanded?: boolean;
    checked?: boolean;
    loading?: boolean;
    pathIds?: string[];
    isLoadMore?: boolean;
    parent_deps?: {id: string; code: number; name: string;type: string }[];
    depart_existed?: boolean;

}

const Limit = 100
const rootId = '00000000-0000-0000-0000-000000000000'
function DepartmentGrid() {
    const userid = session.get('isf.userid')
    const roleTypes = getRoleTypes()
    const [departments, setDepartments] = useState<Department[]>([])
    const [curDept, setCurDept] = useState<Department | undefined>()
    const [listTipStatus, setListTipStatus] = useState<ListTipStatus>(ListTipStatus.Loading)
    const [searchField, _setSearchField] = useState('name')
    const [statusFiled, setStatusFiled] = useState('enable')
    const [searchKey, _setSearchKey] = useState<string | boolean>()
    const [actionType, setActionType] = useState(Action.None)
    const [searchTotal, setSearchTotal] = useState(0)
    const [searchPage, setSearchPage] = useState(1)
    const [selection, _setSelection] = useState<Department | undefined>()
    const [isKjzDisabled, setKjzStatus] = useState(true)
    const searchFilter = [{ key:'name', label: __('部门名称') }, { key:'code', label: __('部门编码') }, { key:'direct_department_code', label: __('上级部门编码') }, { key:'manager_name', label: __('部门负责人') }, { key:'remark', label: __('备注') }, { key:'enabled', label: __('状态') } ]
    const statusFilter = [{ key:'enable', label: __('启用') }, { key:'disable', label: __('禁用') }]
    const getRole = (roleTypes) => {
        if(roleTypes.includes(UserRole.Super)) {
            return UserRole.Super
        }else if(roleTypes.includes(UserRole.Admin)) {
            return UserRole.Admin
        }else if(roleTypes.includes(UserRole.Security)) {
            return UserRole.Security
        }else if(roleTypes.includes(UserRole.OrgManager)) {
            return UserRole.OrgManager
        }else {
            return UserRole.Super
        }
    }
    const gridRef = useRef<HTMLDivElement | any>()

    // 解决更改筛选项时searchKey不正确问题
    const searchKeyRef = useRef(searchKey)

    const setSearchKey = (value) => {
        searchKeyRef.current = value
        _setSearchKey(value)
    }

    const searchFieldRef = useRef(searchField)

    const setSearchField = (value) =>{
        searchFieldRef.current = value
        _setSearchField(value)
    }

    const selectionRef = useRef(selection)
    const setSelection = (value) => {
        selectionRef.current = value
        _setSelection(value)
    }
    // 更改搜索值
    const changeSearchKey = (keyWord: string) => {
        if(keyWord !== searchKey) {
            resetParams()
            setSearchKey(keyWord)
        }
    }
    // 更改筛选类型
    const changeFilter = (detail: string) => {
        resetParams()
        const keyWord = detail === 'enabled' ? true : searchKeyRef.current
        // 由筛选状态切换到其他筛选项时清空searchValue
        if(detail !== 'status' && (searchFieldRef.current === 'enabled' || searchFieldRef.current === "disabled")) {
            setSearchField(detail)
            setSearchKey('')
            getDp('').then((dps) => {
                loadDp(dps)
            })
            return
        }

        setSearchField(detail)
        setSearchKey(keyWord)
        if(typeof keyWord === 'string' && keyWord || typeof keyWord === 'boolean') {
            getDp(keyWord, 0, detail).then((dps) =>{
                loadDp(dps)
            })
        }
    }
    // 更改状态选项
    const changeStatusFilter = (detail: string) => {
        resetParams()
        setStatusFiled(detail)
        setSearchKey(detail === 'enable')
        getDp(detail === 'enable', 0, searchField).then((dps) => {
            loadDp(dps)
        })
    }

    const getColumns = () => {
        const columns = [
            {
                key: 'name',
                title: __('部门名称'),
                width: '45%',
                renderCell: (name, record) => (
                    <div className={styles['depart-name']} style={{ display: 'flex', overflow: 'hidden', marginLeft: `${record.level * 20}px` }}>
                        <div className={styles['icon']}>
                            {
                                !record.isLoadMore && record.depart_existed ? (
                                    <span className={styles['expand-icon']} onClick={() => toggleExpand(record)}>
                                        {
                                            record.loading ?
                                                <Icon url={loading} size={16}/> :
                                                <SweetIcon
                                                    name={record.expanded ? 'arrowDown' : 'arrowRight'}
                                                />
                                        }
                                    </span>
                                ) : null
                            }
                            {
                                !record.isLoadMore ?
                                    <div style={{ marginLeft: record.depart_existed || typeof searchKey === 'string' && searchKey || typeof searchKey === 'boolean' ? 0 : 20 }}>
                                        <UIIcon
                                            role={'ui-uiicon'}
                                            size={14}
                                            code={(record.is_root  || !record.parent_deps.length ) ? '\uf008': '\uf009'}
                                        />
                                    </div>
                                    : null
                            }
                        </div>
                        {
                            record.isLoadMore ?
                                <span className={styles['load-more']} onClick={() => loadMore(record)}>
                                    {__('加载更多')}
                                    <SweetIcon
                                        name={'arrowDown'}
                                    />
                                </span>
                                :
                                <Text className={styles['text']} role={'ui-text'}>{record.name || '---'}</Text>
                        }
                    </div>
                ),
            },
            {
                key: 'email',
                title: __('部门邮箱'),
                width: '15%',
                renderCell: (email, record) => {
                    return !record.isLoadMore ?<Text className={styles['text']} role={'ui-text'}>{record.email|| '---'}</Text>:null
                },
            },
            {
                key: 'direct_department_code',
                title: __('部门编码'),
                width: '15%',
                renderCell: (code, record) => {
                    return !record.isLoadMore ?<Text className={styles['text']} role={'ui-text'}>{record.code|| '---'}</Text>:null
                },
            },
            {
                key: 'direct_department',
                title: __('上级部门'),
                width: '15%',
                renderCell: (direct_department, record) => {
                    return !record.isLoadMore ?
                        <span className={styles['text']} title={formatDepInfo(record.parent_deps).departmentPaths}>
                            {formatDepInfo(record.parent_deps).name || '---'}
                        </span>: null
                },
            },
            {
                key: 'directCode',
                title: __('上级部门编码'),
                width: '15%',
                renderCell: (directCode, record) => {
                    return !record.isLoadMore ?<Text className={styles['text']} role={'ui-text'}>{formatDepInfo(record.parent_deps).code || '---'}</Text>:null
                },
            },
            {
                key: 'managerName',
                title: __('部门负责人'),
                width: '15%',
                renderCell: (managerName, record) => {
                    return !record.isLoadMore ?<Text className={styles['text']} role={'ui-text'}>{record.manager && record.manager.name || '---'}</Text>:null
                },
            },
            {
                key: 'remark',
                title: __('备注'),
                width: '15%',
                renderCell: (remark, record) => {
                    return !record.isLoadMore ? <Text className={styles['text']} role={'ui-text'}>{record.remark|| '---'}</Text> :null
                },
            },
            {
                key: 'enabled',
                title: __('状态'),
                width: '10%',
                renderCell: (enabled, record) => {
                    return !record.isLoadMore ?
                        (
                            <Title content={__(`点此${record.enabled ? '禁用' : '启用'}${(record.is_root  || !record.parent_deps.length ) ? '组织' : '部门'}`)} role={'sweetui-title'}>
                                <Switch
                                    role={'sweetui-switch'}
                                    checked={record.enabled}
                                    onChange={({ event, detail }) => { event.stopPropagation(); changeStatus(record, detail) }}
                                    disabled={!isKjzDisabled || (getRole(roleTypes) === UserRole.OrgManager && (record.is_root  || !record.parent_deps.length )) || getRole(roleTypes) === UserRole.Security}
                                />
                            </Title>
                        ) : null
                },
            },
            {
                key: 'operation',
                title: __('操作'),
                width: '10%',
                renderCell: (operation, record) => {
                    return !record.isLoadMore ?
                        <div className={styles['operation']}>
                            <SweetIcon
                                size={16}
                                name={'edit'}
                                color={!isKjzDisabled || (getRole(roleTypes) === UserRole.OrgManager && (record.is_root  || !record.parent_deps.length )) || getRole(roleTypes) === UserRole.Security ? 'rgba(0, 0, 0, 0.3)': 'rgba(0, 0, 0, 0.85)'}
                                className={styles['edit']}
                                title={__('编辑')}
                                onClick={() => {
                                    if(!(!isKjzDisabled || (getRole(roleTypes) === UserRole.OrgManager && (record.is_root  || !record.parent_deps.length )) || getRole(roleTypes) === UserRole.Security)) {
                                        editHandle(record)
                                    }
                                }}
                            />
                            <SweetIcon
                                color={!isKjzDisabled || (getRole(roleTypes) === UserRole.OrgManager && (record.is_root  || !record.parent_deps.length )) || getRole(roleTypes) === UserRole.Security ? 'rgba(0, 0, 0, 0.3)': 'rgba(0, 0, 0, 0.85)'}
                                size={16}
                                name={'bin'}
                                title={__('删除')}
                                className={styles['delete']}
                                onClick={() => {
                                    if(!(!isKjzDisabled || (getRole(roleTypes) === UserRole.OrgManager && (record.is_root  || !record.parent_deps.length )) || getRole(roleTypes) === UserRole.Security)) {
                                        deleteHandle(record)
                                    }
                                }}
                            />
                        </div> : null
                },
            },
        ]
        return columns
    }

    // 加载子部门
    const loadChildren = (dept: Department, offset = 0 ) => {
        const getDpsPromise = getDepartments({ role: getRole(roleTypes), department_id: dept.id, offset }).then(({ departments: { entries, total_count } }) => {
            const datas = entries.map((cur) => {
                return { ...cur, level: dept.level + 1, expanded: false, isLoadMore: false, pathIds: dept.pathIds ? [...dept.pathIds, dept.id]: [dept.id] }
            })
            if(offset + datas.length < total_count) {
                const moreData = { id: dept.id, isLoadMore: true, level: dept.level + 1, offset, expanded: false, pathIds: dept.pathIds ? [...dept.pathIds, dept.id]: [dept.id] }
                return [ ...datas, moreData]
            }else {
                return datas
            }
        }).catch(async(error) =>{
            handleError(error)
        })
        return getDpsPromise
    }
    // 加载更多
    const loadMore = (dept) => {
        getDepartments({ role: getRole(roleTypes), department_id: dept.id, offset: dept.offset + Limit }).then(({ departments: { entries, total_count } }) => {
            const datas = entries.map((cur) => {
                return { ...cur, level: dept.level, expanded: false, isLoadMore: false, pathIds: dept.pathIds }
            })
            const dataIds = datas.map((cur) => cur.id)
            // 若接口返回数据包含新建时插入的数据，进行过滤后再更新
            const newDps = departments.filter((cur) => !dataIds.includes(cur.id))
            const curIndex = newDps.findIndex((cur) => cur.id === dept.id && cur.isLoadMore)
            let newDatas = datas
            if(dept.offset + Limit + datas.length < total_count) {
                const moreData = { id: dept.id, isLoadMore: true, level: dept.level, offset: dept.offset + Limit, expanded: false, pathIds: [...dept.pathIds, dept.id] }
                newDatas = [...datas, moreData]
            }
            if(curIndex !== -1 && datas && datas.length) {
                const newDepartments = newDps.filter((cur) => !(cur.id === dept.id && cur.isLoadMore))
                // eslint-disable-next-line no-restricted-properties
                newDepartments.splice(curIndex, 0, ...newDatas)
                setDepartments([...newDepartments ])
            }
        }).catch(async(error) => {
            handleError(error)
        })
    }

    // 渲染数据
    const loadDp = (dps) => {
        const { entries, total_count, keyWord } = dps
        if(typeof keyWord === 'string' && keyWord || typeof keyWord === 'boolean') {
            setSearchTotal(total_count)
            const newformatEntries = entries.map((cur) => ({ ...cur, level: 0, isLoadMore: false, expanded: false, pathIds: [''], depart_existed: false }))
            setDepartments([...newformatEntries])
            setListTipStatus(newformatEntries.length < 1 ?  ListTipStatus.NoSearchResults: ListTipStatus.None)
        }else {
            const datas = entries.map((cur) => {
                return { ...cur, level: 0, expanded: false, isLoadMore: false, pathIds: [rootId] }
            })
            const moreData = { id: rootId, isLoadMore: true, level: 0, offset: 0, expanded: false, pathIds: [rootId], is_root: true }
            const newDepartments = datas.length < total_count ? [...datas, moreData] : datas
            setDepartments([...newDepartments])
            setListTipStatus(datas.length ? ListTipStatus.None : ListTipStatus.Empty)
        }
    }

    // 获取数据
    const getDp = (keyWord, offset = 0, detail = searchField) => {
        setListTipStatus(ListTipStatus.Loading)
        if(typeof keyWord === 'string' && keyWord || typeof keyWord === 'boolean') {
            const searchhPromise = searchDepartments({ role: getRole(roleTypes), [detail || searchField]: keyWord, offset }).then(({ entries, total_count }) =>{
                return { entries, total_count,  keyWord }
            }).catch(async(error) =>{
                setListTipStatus(ListTipStatus.LoadFailed)
                handleError(error)
            })
            return searchhPromise
        }else {
            const getPromise = getDepartments({ department_id: rootId, role: getRole(roleTypes), offset }).then(({ departments: { entries, total_count } }) =>{
                return { entries, total_count, keyWord }
            }).catch(async(error) =>{
                setListTipStatus(ListTipStatus.LoadFailed)
                handleError(error)
            })
            return getPromise
        }
    }

    // 移除子部门
    const removeChildren = (departments, dept) => {
        const result = departments.filter((cur) => !cur.pathIds.includes(dept.id))
        return result
    };

    // 更新当前部门信息
    const updateCurDp = (departments, dept) => {
        const datas = departments.map((cur) => {
            if(cur.id === dept.id && !cur.isLoadMore) {
                return { ...dept }
            }
            return cur
        })
        return datas
    }

    // 自动展开父目录更新数据处理
    const autoExpandHandle = (dept, createCur, index, datas) => {
        const curCreate = datas.find((cur) => cur.id === createCur.id)
        const newDepartments = updateCurDp(departments, { ...dept, expanded: true, depart_existed: true })
        setDepartments([...newDepartments ])
        if(curCreate) {
            // eslint-disable-next-line no-restricted-properties
            newDepartments.splice(index + 1, 0, ...datas)
            setDepartments([...newDepartments ])
            scrollToIndex(index)
            setSelection(curCreate)
        }else {
            // eslint-disable-next-line no-restricted-properties
            newDepartments.splice(index, 0, createCur, ...datas)
            setDepartments([...newDepartments ])
            scrollToIndex(index + 2)
            setSelection(createCur)
        }
    }

    // 展开、收起
    const toggleExpand = (dept, createCur: Department | null = null) => {
        if(dept.expanded) {
            dept.expanded= false
            const removeMore = departments.filter((cur) => !(cur.id === dept.id && cur.isLoadMore))
            const newDepartments = removeChildren(removeMore, dept);
            setDepartments([...newDepartments])
        }else {
            // 更新loading状态
            setDepartments([...updateCurDp(departments, { ...dept, loading: true }) ])
            // 获取数据
            loadChildren(dept).then((datas) => {
                setDepartments([...updateCurDp(departments, { ...dept, loading: false }) ])
                const curIndex = departments.findIndex((cur) => cur.id === dept.id && !cur.isLoadMore)
                // 新建时自动展开父部门并插入新建数据
                if(createCur) {
                    if(curIndex !== -1) {
                        autoExpandHandle(dept, createCur, curIndex, datas)
                    }
                    return
                }
                if(curIndex !== -1 && datas && datas.length) {
                    dept.expanded = true
                    // eslint-disable-next-line no-restricted-properties
                    departments.splice(curIndex + 1, 0, ...datas)
                    setDepartments([...departments ])
                }else {
                    dept.expanded = false
                    dept.depart_existed = false
                    const newDp = updateCurDp(departments, dept)
                    setDepartments([...newDp ])
                }
            }).catch((ex) => {
                const newDp = updateCurDp(departments, { ...dept, expanded: false, loading: false })
                setDepartments([...newDp ])
            })
        }
    }

    const updateCurAndChild = (departments, dep) => {
        const newDepartments = departments.map((cur) => {
            if(cur.id === dep.id && !cur.isLoadMore) {
                cur ={ ...dep, enabled: dep.status, manager: dep.managerInfo ? dep.managerInfo[0] : null }
            }

            if(cur.pathIds && cur.pathIds.includes(dep.id)) {
                if(cur.parent_deps) {
                    cur.parent_deps = cur.parent_deps.map((cur) => {
                        if(cur.id === dep.id) {
                            cur.name = dep.name
                            cur.code = dep.code
                        }
                        return cur
                    })
                }
                cur.enabled = dep.status
            }
            return cur
        })
        setDepartments([...newDepartments])
    }

    const scrollToIndex = (index) => {
        if (gridRef.current) {
            gridRef.current.tableRoot.contentView.scrollTop = index * 50
        }
    }

    // 新建成功更新数据
    const insertNodeHandle = (dep) => {
        resetAction()
        if(typeof searchKey === 'string' && searchKey || typeof searchKey === 'boolean') {
            setSearchKey('')
            setSearchField('name')
            // 初始化
            getDp('').then((dps) =>{
                loadDp(dps)
            })
            return
        }
        if(dep.is_root) {
            const curDp = { ...dep, pathIds: [rootId], level: 0, enabled: dep.status, manager: dep.managerInfo ? dep.managerInfo[0] : null, parent_deps: [], type: 'department' }
            // eslint-disable-next-line no-restricted-properties
            departments.splice(0, 0, curDp)
            setDepartments([...departments])
            setSelection(curDp)
            scrollToIndex(0)
            return
        }
        const parentIndex = departments.findIndex((cur) => cur.id === dep.parentId && !cur.isLoadMore)
        const parent =  departments.find((cur) => cur.id === dep.parentId && !cur.isLoadMore)
        if(parentIndex !== -1 && parent) {
            const curInfo = { ...dep, pathIds: parent.pathIds ? [...parent.pathIds, dep.parentId] : [], level: parent.level + 1, enabled: dep.status,  manager: dep.managerInfo ? dep.managerInfo[0] : null, is_root: false, depart_existed: false, type: 'department',  parent_deps: parent.parent_deps ? [...parent.parent_deps, { id: parent.id, name: parent.name, type: 'department', code: parent.code }]:[{ id: parent.id, name: parent.name, type: 'department', code: parent.code }] }
            if(parent.expanded) {
                // eslint-disable-next-line no-restricted-properties
                departments.splice(parentIndex + 1, 0, curInfo)
                setDepartments([...departments])
                setSelection(curInfo)
                scrollToIndex(parentIndex + 1)
                return
            }
            if(!parent.expanded) {
                toggleExpand({ ...parent, expanded: false, depart_existed: true }, curInfo)
            }
        }
    }

    // 更改部门状态
    const changeStatus = async (dept, detail) => {
        try {
            let { ossInfo, email, managerID, code, managerDisplayName, remark, departmentName } = await getOrgDepartmentById(dept.id)
            const editParma = {
                ncTEditDepartParam: {
                    departId: dept.id,
                    departName: departmentName,
                    managerID,
                    code,
                    remark,
                    status: detail,
                    ossId: ossInfo.ossId || '',
                    email,
                },
            }
            if(dept.is_root || !dept.parent_deps.length ) {
                await editOrganization([editParma]);
                manageLog(
                    ManagementOps.SET,
                    __('编辑 组织 “${name}” 成功', { name: departmentName.replace(/\.+$/, '') }),
                    __('原组织名 “${oldName}”；组织编码 “${code}”；组织负责人 “${managerDisplayName}”；备注 “${remark}”；邮箱地址 “${emailAddress}”；存储位置 “${ossName}”；状态 “${status}”；', {
                        oldName: departmentName,
                        emailAddress: email,
                        code,
                        managerDisplayName,
                        remark,
                        status: status ? __('启用') : __('禁用'),
                        ossName: displayUserOssInfo(ossInfo),
                    }),
                    Level.INFO,
                )
            }else {
                await editDepartment([editParma]);
                manageLog(
                    ManagementOps.SET,
                    __('编辑部门 “${depName}” 成功', { depName: departmentName
                        .replace(/\.+$/, '') }),
                    __('原部门名 “${oldName}”；部门编码 “${code}”； 部门负责人“${managerDisplayName}”；备注 “${remark}”；邮箱地址 “${emailAddress}”；存储位置 “${ossName}”；状态 “${status}”；', {
                        oldName: departmentName,
                        emailAddress: email,
                        code,
                        managerDisplayName,
                        remark,
                        status: detail ? __('启用') : __('禁用'),
                        ossName: displayUserOssInfo(ossInfo),
                    }),
                    Level.INFO,
                )
            }
            if(detail) {
                setDepartments([...updateCurDp(departments, { ...dept, enabled: detail })])
            }else {
                updateCurAndChild(departments, { ...dept, status: detail })
            }

        } catch(error) {
            if (error) {
                switch (error.errID) {
                    case ErrorCode.OrgNameNotExist:
                        await Message2.alert({ message: __('组织 “${orgName}” 不存在。', { orgName: dept.name }) })
                        deleteNodeHandle(dept)
                        break;
                    case ErrorCode.ParentDepartmentNotExist:
                    case ErrorCode.DepNameNotExist:
                        await Message2.info({ message: __('部门不存在。') })
                        deleteNodeHandle(dept)
                        break;

                    default:
                        handleError(error)
                        break;
                }
            }
        }
    }

    // 编辑部门
    const editHandle = (dept) => {
        setCurDept(dept)
        setActionType(dept.is_root || !dept.parent_deps.length ? Action.EditOrg : Action.EditDep)
    }

    // 编辑成功更新数据
    const editNodeHandle = (dep) => {
        resetAction()
        updateCurAndChild(departments, dep)
    }

    // 删除成功更新数据
    const deleteNodeHandle = (dept) => {
        resetAction()
        const newDepartments = removeChildren(departments, dept).filter((cur) => cur.id !== dept.id);
        const parent = newDepartments.find((cur) => dept.pathIds.length && cur.id === dept.pathIds[dept.pathIds.length - 1])
        let newResult = newDepartments
        if(parent) {
            const hasChild = newDepartments.find((cur) => cur.pathIds.includes(parent.id))
            newResult = hasChild ? newDepartments : updateCurDp(newDepartments, { ...parent, depart_existed: false, expanded: false })
        }
        setDepartments(newResult)
    }

    // 删除部门
    const deleteHandle = (dept) => {
        setCurDept(dept)
        setActionType(dept.is_root || !dept.parent_deps.length ? Action.DelOrg :Action.DelDep)
    }

    // 加载异常
    const loadFailed = (ex: any) => {
        resetParams()
        setDepartments([])
        setListTipStatus(ListTipStatus.LoadFailed)
    }

    const handleError = (error) => {
        if(error) {
            if(error.description) {
                Message2.info({ message: error.description })
            }else {
                Message2.info({ message: getErrorMessage(error.code) })
            }
        }
    }

    // 获取搜索结果
    const getSearchResult = (offset = searchPage ? searchPage * Limit - Limit : 0) => {
        getDp(searchKey, offset,  searchField).then((dps) =>{
            loadDp(dps)
        })
    }

    // 翻页
    const handlePageChange = (page) => {
        setDepartments(departments)
        setSearchPage(page)
        getSearchResult( page ? page * Limit - Limit : 0)
    }

    const resetParams = () => {
        setSearchPage(1)
        setSearchTotal(0)
        if (gridRef.current) {
            gridRef.current.tableRoot.contentView.scrollTop = 0
        }
    }

    // 渲染操作弹框
    const resetAction = () =>{
        setActionType(Action.None)
    }

    const getKjzStatus = async() => {
        const res = await getConfidentialConfig('kjz_disabled')
        setKjzStatus(typeof res === 'undefined' || res)
    }

    useEffect(() => {
        getKjzStatus()
        // 初始化
        getDp('').then((dps) =>{
            loadDp(dps)
        })
        // eslint-disable-next-line
    }, [])

    const renderAction = () => {
        switch(actionType) {
            case Action.CreateOrg:
                return <CreateOrganization dep={undefined} userid={userid} onRequestCancelCreateOrg={resetAction} onCreateOrgSuccess={insertNodeHandle}/>
            case Action.EditOrg:
                return <EditOrganization dep={curDept as any} userid={userid} onRequestCancelEditOrg={resetAction} onEditOrgSuccess={editNodeHandle} onRequestDelOrg={deleteNodeHandle} />
            case Action.DelOrg:
                return <DeleteOrganization dep={curDept as any} userid={userid} onRequestCancelDeleteOrg={resetAction} onDeleteOrgSuccess={deleteNodeHandle}/>
            case Action.CreateDep:
                return <CreateDepartment sourcePage={'dep-tree'} dep={undefined} userid={userid} onRequestCancelCreateDep={resetAction} onCreateDepSuccess={insertNodeHandle} onRequestDelDep={deleteNodeHandle}/>
            case Action.EditDep:
                return <EditDepartment dep={curDept as any} parentName={curDept && curDept.parent_deps? formatDepInfo(curDept.parent_deps).name : ''} userid={userid} onRequestCancelEditDep={resetAction} onEditDepSuccess={editNodeHandle} onRequestDelDep={deleteNodeHandle}/>
            case Action.DelDep:
                return <DeleteDepartment dep={curDept as any} userid={userid} onRequestCancelDeleteDep={resetAction} onDeleteDepSuccess={deleteNodeHandle}/>
            default:
                return null
        }
    }

    return (
        <div className={styles['dp-grid']}>
            <div className={styles['grid-header']}>
                {
                    isKjzDisabled && getRole(roleTypes) !== UserRole.Security ?
                        getRole(roleTypes) === UserRole.OrgManager ?
                            <Button
                                role={'sweetui-button'}
                                width={'auto'}
                                icon={'add'}
                                size={13}
                                theme={'oem'}
                                onClick={() =>{
                                    setActionType(Action.CreateDep)
                                }}
                            >
                                {__('新建部门')}
                            </Button>
                            :
                            <PopMenu role={'sweetui-popmenu'} key={'create'} freeze={false} triggerEvent={'hover'} onRequestCloseWhenClick={(close) => close()}
                                trigger={({ setPopupVisibleOnMouseEnter, setPopupVisibleOnMouseLeave }) =>
                                    <Button
                                        role={'sweetui-button'}
                                        width={'auto'}
                                        key={'create'}
                                        theme={'oem'}
                                        onMouseEnter={setPopupVisibleOnMouseEnter}
                                        onMouseLeave={setPopupVisibleOnMouseLeave}
                                    >
                                        <div className={styles['pop-button']}>
                                            <SweetIcon
                                                size={16}
                                                name={'add'}
                                            />
                                            <div className={styles['menu-name']}>{__('新建')}</div>
                                            <UIIcon
                                                role={'ui-uiicon'}
                                                code={'\uf00b'}
                                                size={12}
                                                className={styles['menu-icon']}
                                            />
                                        </div>
                                    </Button>
                                }>
                                <PopMenu.Item className={styles['action']} role={'sweetui-popmenu.item'} key={'createOrg'} label={__('新建组织')} onClick={() => {
                                    setActionType(Action.CreateOrg)
                                }}/>
                                <PopMenu.Item className={styles['action']} role={'sweetui-popmenu.item'} key={'createDp'} label={__('新建部门')} onClick={() => {
                                    setActionType(Action.CreateDep)
                                }}/>
                            </PopMenu>
                        : null
                }
                <div style={{ display: 'flex' }}>
                    <Select
                        className={styles['select']}
                        width={120}
                        maxMenuHeight={190}
                        value={searchField}
                        onChange={({ detail }) => changeFilter(detail)}
                    >
                        {
                            searchFilter.map(({ key, label }) => {
                                return <Select.Option key={key} value={key}>{label}</Select.Option>
                            })
                        }
                    </Select>
                    {
                        searchField === 'enabled' ?
                            <Select
                                width={120}
                                maxMenuHeight={190}
                                value={statusFiled}
                                onChange={({ detail }) => changeStatusFilter(detail)}
                            >
                                {
                                    statusFilter.map(({ key, label }) => {
                                        return <Select.Option key={key} value={key}>{label}</Select.Option>
                                    })
                                }
                            </Select>
                            :
                            <SearchBox
                                className={styles["search"]}
                                role={'ui-searchbox'}
                                width={220}
                                placeholder={__('请输入搜索内容')}
                                value={searchKey}
                                onChange={changeSearchKey}
                                loader={getDp}
                                onLoad={loadDp}
                                onLoadFailed={loadFailed}
                            />
                    }
                </div>
            </div>
            <div className={styles['grid-main']}>
                <div className={styles['grid-content']}>
                    <DataGrid
                        role={'sweetui-datagrid'}
                        height={'100%'}
                        ref={gridRef}
                        data={departments}
                        showBorder={false}
                        enableSelect={true}
                        enableMultiSelect={false}
                        isResizable={true}
                        onSelectionChange={({ detail }) => setSelection(detail)}
                        DataGridHeader={{ enableSelectAll: false }}
                        selection={selectionRef.current}
                        refreshing={listTipStatus !== ListTipStatus.None}
                        RefreshingComponent={
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                            />
                        }
                        DataGridPager={ typeof searchKey === 'string' && searchKey || typeof searchKey === 'boolean' ? {
                            size: Limit,
                            total: searchTotal,
                            page: searchPage,
                            onPageChange: ({ detail: { page } }) => handlePageChange(page),
                        }: undefined}
                        contentViewClassName={typeof searchKey === 'string' && searchKey || typeof searchKey === 'boolean' ? styles['search-page'] : styles['department-page']}
                        columns={getColumns()}
                    />
                </div>
            </div>
            {renderAction()}
        </div>
    )
}
export default DepartmentGrid