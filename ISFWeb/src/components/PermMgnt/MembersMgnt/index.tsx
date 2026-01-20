import React, { useState, useEffect, useRef, useContext } from "react";
import intl from "react-intl-universal";
import { getMembers, updateMembers} from "@/core/apis/console/authorization";
import { RoleClassEnum, RoleInfoType } from "../types";
import { Table, Button, Dropdown, Menu, message, Modal, Input,} from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { apis, components } from "@dip/components/dist/dip-components.min.js";
import { AccessorTypeEnum, fromatItem, getRoleType } from "../util";
import { formatDirectDeptInfo } from "@/components/UserOrg/helper";
import { checkMemberExist, deleteUserRolemMember, getUserRolemMember, setUserRolemMember } from "@/core/thrift/sharemgnt";
import AppConfigContext from "@/core/context/AppConfigContext";
import { SystemRoleType, getRoleName } from "@/core/role/role";
import { SystemRoleMember } from "./SystemRoleMember";
import { Level, ManagementOps, manageLog } from "@/core/log";
import { getErrorMessage } from "@/core/exception";
import { defaultModalParams } from "@/util/modal";
import EmptyIcon from "../../../icons/empty.png";
import loadFailedIcon from "../../../icons/loadFailed.png";
import OperationIcon from "../../../icons/operation.svg";
import AddIcon from "../../../icons/add.svg";
import RemoveIcon from "../../../icons/remove.svg";
import SearchEmptyIcon from "../../../icons/searchEmpty.png";
import { UserManagementErrorCode } from "@/core/apis/openapiconsole/errorcode";
import { notExisted } from "@/core/apis/console/usermanagement/types";
import styles from "./styles.css";
import { getIcon } from "../index";
import { AuthorizationErrorCodeEnum, RoleType } from "@/core/apis/console/authorization/type";
import { trim, isEqual, debounce } from 'lodash'

const { confirm, info } = Modal;

export const MembersMgnt = ({currentRole, roleClass, updateRoleList}: {currentRole: RoleType; roleClass: RoleClassEnum, updateRoleList?: () => void }) => {
    const { config: { userInfo} } = useContext(AppConfigContext);
    const [curPageSize, setPageSize] = useState(50);
    const [curPage, setPage] = useState(1);
    const [roleMembers, setRoleMembers] = useState([]);
    const [roleMembersCount, setRoleMembersCount] = useState(0);
    const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);
    const [selections, setSelections] = useState([]);
    const [isLoading, setIsLoading] = useState(false);  
    const [isError, setIsError] = useState(false);
    const [searchValue, setSearchValue] = useState('');
    const [filterValue, setFilterValue] = useState([]);
    const accessorRef = useRef(null);
    const selectmembersRef = useRef(null);
    const selectRangesRef = useRef(null);
    const [showSelectMembers, setShowSelectMembers] = useState(false);
    const [type, setType]= useState('create')
    const businessRoleTop = currentRole?.description ? '310px' : '286px'
    const systemRoleTop = currentRole?.description ? '244px': '220px'
    const topHeight = roleClass === RoleClassEnum.Business ? businessRoleTop : systemRoleTop

    const onEditMembers = async() => {
        setType("edit")
        setShowSelectMembers(true)
    };

    const getLogMessage = (memberInfo) => {
        switch (currentRole.id) {
            case SystemRoleType.OrgAudit:
            case SystemRoleType.OrgManager:
                return intl.get("member.range.log", {
                    name: memberInfo.displayName,
                    memberRange: memberInfo.manageDeptInfo ? memberInfo.manageDeptInfo.departmentNames.join(intl.get("quota")) : '',
                })
            default:
                return ''
        }
    }

    const updateCurrentPage = (deleteNumber) => {
        const totalPage = Math.ceil((roleMembersCount - deleteNumber) / curPageSize)
        if(totalPage < curPage && totalPage !== 0) {
            setPage(totalPage)
            getRoleMemebers({ id: currentRole.id, offset: (totalPage - 1) * curPageSize, limit: curPageSize, type: filterValue, keyword: searchValue });
        }else {
            getRoleMemebers({ id: currentRole.id, offset: (curPage - 1) * curPageSize, limit: curPageSize, type: filterValue, keyword: searchValue });
        }
    }

    const handleDelete = async(selections, deleteModal) => {
        try {
            if(roleClass === RoleClassEnum.Business) {
                const members = selections.map(item => {
                    return {id: item.id, type:item.type}
                });
                await updateMembers({method: "DELETE",id: currentRole.id, members})
            }else {
                for(let member of selections) {
                    await deleteUserRolemMember([userInfo.id, currentRole.id, member.userId])
                    if (currentRole.id === SystemRoleType.OrgManager) {
                        manageLog(
                            ManagementOps.DELETE,
                            intl.get("cancel.member.org.log", { userName: member.displayName, departmentName: member.manageDeptInfo.departmentNames.join('”，“') }),
                            null,
                            Level.INFO,
                        )
                    } else {
                        manageLog(
                            ManagementOps.DELETE,
                            intl.get("cancel.member.log", { roleName: currentRole.name, userName: member.displayName }),
                            getLogMessage(member),
                            Level.INFO,
                        )
                    }
                }
            }
            message.success(intl.get("remove.success"));
            deleteModal?.destroy();
            if(roleClass === RoleClassEnum.Business) {
                updateCurrentPage(selections?.length)
            }else {
                const selectIds = selections.map(cur => cur.userId)
                const newRoleMembers = roleMembers.filter((cur) => !selectIds.includes(cur.userId))
                setRoleMembers(newRoleMembers)
                initData()
            }
            
        } catch (e) {
            handleError(e)
            deleteModal?.destroy()
        }
    }

    const onDeleteMembers = (selections) => {
        const deleteModal = confirm({
            ...defaultModalParams,
            title: intl.get('delete.member.title'),
            content:
              selections.length > 1
                  ? intl.get('delete.batch.member.tip', { count: selections.length })
                  : intl.get('delete.member.tip'),
            footer: () => (
                <div style={{ textAlign: 'right' }}>
                    <Button key="delete" type="primary" danger onClick={() => handleDelete(selections, deleteModal)}>
                        {intl.get('remove')}
                    </Button>
                    <Button
                        key="back"
                        onClick={() => {
                            deleteModal?.destroy();
                        }}
                    >
                        {intl.get('cancel')}
                    </Button>
                </div>
            ),
            onClose: () => {
                deleteModal?.destroy();
            },
            getContainer: document.getElementById('isf-web-plugins')
        });
    };

    // 使用useRef保存防抖函数实例，避免每次渲染创建新的函数实例
    const debouncedSearch = useRef(
        debounce(({ value, role, pageSize, filterValue }) => {
            getRoleMemebers({ id: role.id, offset: 0, limit: pageSize, type: filterValue, keyword: trim(value)})
        }, 300)
    )

    const handleSearchMemebers = (e) => {
        const value = e.target.value; 
        setSearchValue(value);
        setPage(1)
        debouncedSearch.current?.({ value, role: currentRole, pageSize: curPageSize, filterValue })
    }

    const getSelections = (recordKeys: string[]) => {
        const selections = roleMembers.filter((item: any) => recordKeys.includes(item.id || item.userId));
        setSelections(selections);
    };
  
    const onRow = (record) => ({
        onClick: () => {
            const id = record.id || record.userId
            setSelectedRowKeys([id]);
            getSelections([id]);
        },
    });

    const rowSelection = {
        selectedRowKeys,
        onChange: (selectedKeys: string[]) => {
            setSelectedRowKeys(selectedKeys);
            getSelections(selectedKeys);
        },
        getCheckboxProps: (record) => ({
            disabled: false,
            key: record.id || record.userId,
        }),
    };

    const getRoleMemebers = async ({ id, offset = 0, limit = 50, type = [], keyword = "" }) => {
        try {
            setIsLoading(true);
            setIsError(false);  
            const query = {
                id, 
                offset, 
                limit,
                ...(type.length ? { type } : {}),
                ...(keyword ? { keyword } : {}),
            }
            const { entries, total_count } = await getMembers(query);
            setIsLoading(false);   
            setRoleMembers(entries);
            setRoleMembersCount(total_count);
        } catch (e) {
            setIsLoading(false);
            setIsError(true);
            setRoleMembers([])   
        }
    };

    const getSystemRoleMembers = async () => {
        try {
            setIsLoading(true);
            setIsError(false);   
            const allMember = await getUserRolemMember([userInfo.id, currentRole.id])
            setIsLoading(false);  
            setRoleMembers(allMember)
            setRoleMembersCount(allMember.length);
        } catch(e) {
            setIsLoading(false);
            setIsError(true);   
            setRoleMembers([])  
        }
    }

    const addHandle = async (datas, unmount) => {
        const members = datas.map((cur) => fromatItem(cur))
        try {
            await updateMembers({ method: "POST", id: currentRole.id, members})
            message.success(intl.get("add.success"));
            unmount();
            setFilterValue([])
            setSearchValue('')
            setPage(1)
            getRoleMemebers({ id: currentRole.id, offset: 0, limit: curPageSize });
        } catch (e) {
            handleError(e, datas, unmount)
        }
    };

    const setSystemRole = async (visitors, unmount) => {
        try {
            let existMember = []
            for (let member of visitors) {
                const existResult = await checkMemberExist([currentRole.id, member.id])
                if (existResult) {
                    existMember = [...existMember, member]
                }
            }

            if(existMember.length) {
                info({
                    ...defaultModalParams,
                    title: intl.get("member.existed"),
                    content: (
                        <div className={styles["user-list-tip"]}>
                            {
                                existMember.map((cur) =>(
                                    <div key={cur.userId} className={styles["item"]} title={cur.name}>{cur.name}</div>
                                ))
                            }
                        </div>
                    ),
                    getContainer: document.getElementById('isf-web-plugins')
                })
                return
            }

            for (let member of visitors) {
                await setUserRolemMember([userInfo.id, currentRole.id, {
                    userId: member.id,
                    displayName: member.name,
                }])

                manageLog(
                    ManagementOps.SET,
                    intl.get("set.member.log", { roleName: getRoleName(currentRole), userName: member.name }),
                    "",
                    Level.INFO,
                )
            }
            message.success(intl.get("add.success"));
            unmount()
            getSystemRoleMembers()
        }catch(e) {
            handleError(e)
        }
    }

    const addSystemRoleMember = () => {
        if(currentRole.id === SystemRoleType.OrgManager || currentRole.id === SystemRoleType.OrgAudit) {
            setShowSelectMembers(true)
        }else {
            const unmount = apis.mountComponent(
                components.AccessorPicker,
                {
                    title: intl.get("add.member"),
                    tabs: ["organization",],
                    range: ["user"],
                    isAdmin: true,
                    role: getRoleType(userInfo?.user?.roles),
                    onSelect: (data) => {
                        const newSelections = data.map(item => {
                            return fromatItem(item, true);
                        });
                        if(newSelections.length) {
                            setSystemRole(newSelections, unmount)
                        }
                    },
                    onCancel: () => {
                        unmount();
                    },
                },
                accessorRef.current,
            );
        }
    }

    const addVisitorHandle = () => {
        initData()
        if(roleClass === RoleClassEnum.Business) {
            const unmount = apis.mountComponent(
                components.AccessorPicker,
                {
                    title: intl.get("add.member"),
                    tabs: ["organization", "group", "app"],
                    range: ["user", "department", "group", "app"],
                    isAdmin: true,
                    role: getRoleType(userInfo?.user?.roles),
                    onSelect: (data) => {
                        const newSelections = data.map(item => {
                            return fromatItem(item, true);
                        });
                        if(newSelections.length) {
                            addHandle(newSelections, unmount)
                        }
                    },
                    onCancel: () => {
                        unmount();
                    },
                },
                accessorRef.current,
            );
        }else {
            addSystemRoleMember()
        }
    };

    const getErrorTitle = (code) => {
        switch(code) {
            case UserManagementErrorCode.GroupMemberNotExisted:
                return intl.get("user.deleted")
            case UserManagementErrorCode.DepartmentNotExisted:
                return intl.get("department.deleted")
            case UserManagementErrorCode.UserGroupNotFound:
                return intl.get("usergroup.deleted")
            case UserManagementErrorCode.AppAccountNotFound:
                return intl.get("appaccount.deleted")
            default:
                return ""
        }
    }

    const handleError = (e, datas = [], unmount = null) => {
        if(e.code && e?.detail?.ids && notExisted.includes(e.code)) {
            const users = datas.filter((cur) => e.detail.ids.includes(cur.id))
            info({
                ...defaultModalParams,
                title: getErrorTitle(e.code),
                content: (
                    <div className={styles["user-list-tip"]}>
                        {
                            users.map((cur) =>(
                                <div key={cur.id} className={styles["item"]} title={cur.name}>{cur.name}</div>
                            ))
                        }
                    </div>
                ),
                getContainer: document.getElementById('isf-web-plugins')
            })
        }else if(e?.code === AuthorizationErrorCodeEnum.RoleNotFound) {
            info({ 
                ...defaultModalParams, 
                closable: false,
                content: intl.get("role.not.exist"), 
                getContainer: document.getElementById('isf-web-plugins'), 
                onOk: () => {
                    updateRoleList?.()
                    unmount?.()
                }
            })
        }else {
            const msg = roleClass === RoleClassEnum.Business ? e?.description || "" : e?.error?.errID ? getErrorMessage(e?.error?.errID) || "" : ""
            msg && info({ ...defaultModalParams, content: msg, getContainer: document.getElementById('isf-web-plugins')})
        }
    }

    const baseCoumns = [
        {
            title: intl.get("members"),
            dataIndex: "name",
            key: "name",
            width: currentRole && (currentRole.id === SystemRoleType.OrgManager || currentRole.id === SystemRoleType.OrgAudit) ? "20%" : "30%",
            filters: roleClass === RoleClassEnum.Business ? [
                {
                    text: intl.get("department"),
                    value: "department"
                },
                {
                    text: intl.get("user"),
                    value: "user"
                },
                {
                    text: intl.get("group"),
                    value: "group"
                },
                {
                    text: intl.get("app"),
                    value: "app"
                }
            ]: undefined,
            filterMultiple: true,
            filteredValue: filterValue || null,
            render: (name: string, cur: any) => {
                return (
                    <div className={styles["item"]}>
                        <div className={styles["icon"]}>
                            {getIcon(cur.type || "user", 14)}
                        </div>
                        <div className={styles["name"]} title={cur.name || cur.displayName}>{cur.name || cur.displayName}</div>
                    </div>
                )
            }
        },
        {
            title: intl.get("operation"),
            dataIndex: "operation",
            key: "operation",
            width: "12%",
            render: (operation: any, cur: RoleInfoType) => (
                <Dropdown
                    trigger={['click']}
                    dropdownRender={() => {
                        return (
                            <Menu>
                                {currentRole && (currentRole.id === SystemRoleType.OrgManager || currentRole.id === SystemRoleType.OrgAudit) && 
                                <Menu.Item onClick={() => onEditMembers()}>
                                    {intl.get("edit")}
                                </Menu.Item>
                                }
                                <Menu.Item onClick={() => {
                                    onDeleteMembers([cur])
                                }}>
                                    {intl.get("remove")}
                                </Menu.Item>
                            </Menu>
                        );
                    }}
                >
                    <Button
                        type="text"
                        size={"small"}
                        icon={<OperationIcon style={{ width: "16px", height: "16px" }} />}
                    />
                </Dropdown>
            ),
        },
        {
            title: intl.get("from"),
            dataIndex: "parent_deps",
            key: "parent_deps",
            width: currentRole && (currentRole.id === SystemRoleType.OrgManager || currentRole.id === SystemRoleType.OrgAudit) ? "20%" : "58%",
            render: (parent_deps: any[], cur: any) => {
                if(roleClass === RoleClassEnum.Business) {
                    const fromNames = cur.type === AccessorTypeEnum.User || cur.type === AccessorTypeEnum.Department ? 
                        formatDirectDeptInfo(parent_deps).departmentNames : cur.type === AccessorTypeEnum.Group ? 
                            [intl.get("group")] : cur.type === AccessorTypeEnum.App ?
                                [intl.get("app")] : [];
                    const fromPaths = cur.type === AccessorTypeEnum.User || cur.type === AccessorTypeEnum.Department ?
                        formatDirectDeptInfo(parent_deps).departmentPaths : cur.type === AccessorTypeEnum.Group ? 
                            [intl.get("group")] : cur.type === AccessorTypeEnum.App ?
                                [intl.get("app")] : [];
                    return (
                        <div className={styles["ellipsis-item"]} title={fromPaths.join('\n')}>
                            { fromNames.join(',') }
                        </div>
                    )
                }else {
                    return (
                        <div className={styles["ellipsis-item"]} title={cur?.departmentNames?.join('\n')}>{cur?.departmentNames?.join(',')}</div>
                    )
                }
            }
        },
    ];

    const columns = currentRole && (currentRole.id === SystemRoleType.OrgManager || currentRole.id === SystemRoleType.OrgAudit) ?
        [
            ...baseCoumns,
            {
                title: intl.get("range"),
                dataIndex: "range",
                key: "range",
                width: "48%",
                render: (range: any, cur: any) => {
                    return (
                        <div className={styles["ellipsis-item"]} title={cur?.manageDeptInfo?.departmentNames.join('\n')}>
                            {cur?.manageDeptInfo?.departmentNames?.join(',')}
                        </div>
                    )
                }
            },
        ]: baseCoumns;

    const initData = () => {
        setSelectedRowKeys([])
        setSelections([])
        setType('create')
    }

    useEffect(() => {
        initData()
        if (currentRole?.id) {
            setPage(1)
            setRoleMembers([])
            setSearchValue('')
            setFilterValue([])
            if(roleClass === RoleClassEnum.Business) {
                getRoleMemebers({ id: currentRole.id, offset: 0, limit: curPageSize });
            }else {
                getSystemRoleMembers()
            }
        }
    }, [currentRole?.id]);

    useEffect(() => {
        if(roleMembers) {
            getSelections(selectedRowKeys)
        }
    }, [roleMembers])

    useEffect(() => {
        setSelectedRowKeys([]);
        setSelections([]);
    }, [searchValue, curPage, curPageSize]);

    return (
        <div className={styles["members"]}>
            <div className={styles["header"]}>
                <div className={styles["menu"]}>
                    <Button
                        className={styles["btn"]}
                        type={"primary"}
                        onClick={addVisitorHandle}
                        icon={<AddIcon style={{width: "14px", height: "14px"}}/>}
                    >
                        {intl.get("add.member")}
                    </Button>
                    <Button 
                        className={styles["btn"]} 
                        disabled={!selections.length} 
                        onClick={() => onDeleteMembers(selections)}
                        icon={<RemoveIcon style={{width: "14px", height: "14px"}}/>}
                    >
                        {intl.get("remove")}
                    </Button>
                </div>
                {
                    roleClass === RoleClassEnum.Business && 
                    <Input 
                        allowClear
                        placeholder={intl.get("search.name")}
                        className={styles["search"]}
                        prefix={<SearchOutlined />}
                        value={searchValue}
                        onChange={handleSearchMemebers}
                    />
                }
            </div>
            <div className={styles["table"]}>
                <Table
                    size="small"
                    loading={isLoading}
                    tableLayout="fixed"
                    scroll={{y: `calc(100vh - ${topHeight})`}}
                    columns={columns}
                    dataSource={roleMembers}
                    rowSelection={rowSelection}
                    onRow={onRow}
                    rowKey={(record) => record.id || record.userId}
                    locale={{
                        emptyText: (
                            <div className={styles["empty"]}>
                                <img
                                    src={isError ? loadFailedIcon : searchValue ? SearchEmptyIcon : EmptyIcon}
                                    alt=""
                                    width={128}
                                    height={128}
                                />
                                <span>{intl.get(isError ? "loadFailed" : searchValue ? "no.search.result" : "no.member")}</span>
                            </div>
                        ),
                    }}
                    onChange={(_, filters) => {
                        const type = filters.name || [];
                        if(isEqual(type, filterValue)) return
                        setFilterValue(type)
                        getRoleMemebers({ id: currentRole.id, offset: 0, limit: curPageSize, type, keyword: searchValue})
                    }}
                    pagination={roleClass === RoleClassEnum.Business ? {
                        current: curPage, 
                        pageSize: curPageSize,
                        total: roleMembersCount,
                        showSizeChanger: true,
                        showQuickJumper: false,
                        showTotal: (total) => {
                            return intl.get("list.total.tip", { total });
                        },
                        onChange: (page, pageSize) => {
                            initData()
                            if(page === curPage && pageSize === curPageSize) return
                            setPage(pageSize !== curPageSize ? 1 : page)
                            setPageSize(pageSize);
                            getRoleMemebers({
                                id: currentRole.id,
                                offset: ((pageSize !== curPageSize ? 1 : page) - 1) * pageSize,
                                limit: pageSize,
                                type: filterValue,
                                keyword: searchValue
                            });
                        },
                    }: false}
                />
            </div>
            <div ref={accessorRef}></div>
            <div ref={selectmembersRef}></div>
            <div ref={selectRangesRef}></div>
            {
                showSelectMembers && <SystemRoleMember title={intl.get(type === 'create' ? "add.member": "edit.member")} currentRole={currentRole} selectmembersRef={selectmembersRef} selectRangesRef={selectRangesRef} type={type} selections={selections} onCancel={() => setShowSelectMembers(false)} onSuccess={() =>{
                    getSystemRoleMembers()
                }}/>
            }
        </div>
    )
}