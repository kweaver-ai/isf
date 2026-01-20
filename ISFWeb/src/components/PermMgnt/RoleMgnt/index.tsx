import React, { useState, useEffect, useContext } from "react";
import { Empty, Radio, Space, Spin, Modal, Button, Dropdown, Menu, message, Input } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import intl from "react-intl-universal";
import InfiniteScroll from 'react-infinite-scroll-component';
import { getRoles, deleteRole } from "@/core/apis/console/authorization";
import { OperationTypeEnum, RoleClassEnum } from "../types";
import AppConfigContext from "@/core/context/AppConfigContext";
import { MembersMgnt } from "../MembersMgnt";
import { RoleDetails } from "../RoleDetails";
import { getUserRolem } from "@/core/thrift/sharemgnt";
import { SystemRoleType, getRoleFuntional, getRoleName } from "@/core/role/role";
import loadFailedIcon from "../../../icons/loadFailed.png";
import { getErrorMessage } from "@/core/exception";
import { defaultModalParams } from "@/util/modal"
import AddIcon from "../../../icons/add.svg"
import OperationIcon from "../../../icons/operation.svg";
import EmptyIcon from "../../../icons/empty.png";
import SearchEmptyIcon from "../../../icons/searchEmpty.png";
import { SetRole } from "../SetRole";
import styles from "./styles.css";
import { trim } from "lodash";

const { info, confirm } = Modal

const limit = 50
function RoleMgnt() {
    const { oemColor, config: { userInfo} } = useContext(AppConfigContext);
    const [roles, setRoles] = useState([]);
    const [roleClass, setRoleClass] = useState(RoleClassEnum.Business);
    const [currentRole, setCurrentRole] = useState(null);
    const [searchValue, setSearchValue] = useState("");
    const [hasMore, setHasMore] = useState(false);
    const [hoverIndex, setHoverIndex] = useState(0)
    const [showDetails, setShowDetails] = useState(false);
    const [operationType, setOperationType] = useState("")
    const [showRoleDrawer, setRoleDrawer]= useState(false)
    const [loading, setLoading]= useState(false)
    const [isError, setIsError] = useState(false)
    const scrollRef = React.useRef<HTMLElement>(null);

    const getBusinessRole = async ({ isClear = true, isSelectFirst = true, offset = 0, limit, keyword = "", source = ["business", "user"] }) => {
        let scrollTop = 0;
        if (!isSelectFirst && scrollRef.current) {
            scrollTop = scrollRef.current.scrollTop;
        }
        try {
            isClear && initData()
            const { entries, total_count } = await getRoles({ offset, limit, keyword, source });
            setHasMore(total_count > offset + limit); 
            setLoading(false)
            const newRoles = isClear ? entries : roles.concat(entries);
            setRoles(newRoles);
            if(!newRoles.length) {
                setCurrentRole(null)
            }
            if (isSelectFirst && newRoles.length) {
                setCurrentRole(newRoles[0]);
            }
            if (!isSelectFirst && scrollRef.current) {
                setTimeout(() => {
                    scrollRef.current.scrollTop = scrollTop;
                }, 0);
            }
        } catch (e) {
            handleError(e)
        }
    };

    const getSystemRole = async () => {
        try {
            initData()
            const allRoles = (await getUserRolem([userInfo.id])).map((value) => ({
                ...value,
                name: getRoleName(value),
                description: getRoleFuntional(value)
            }))
            setLoading(false)
            setRoles(allRoles);
            if (allRoles.length) {
                setCurrentRole(allRoles[0]);
            }
        }catch(e) {
            handleError(e)
        } 
    }

    const onChangeRadio = (e) => {
        setRoles([])
        setRoleClass(e.target.value);
    };

    const changeSelect = (role) => {
        setCurrentRole(role);
    };

    const handleError = (e) => {
        setLoading(false)
        setIsError(true)
        const msg = roleClass === RoleClassEnum.Business ? e?.description || "" : e?.error?.errID ? getErrorMessage(e?.error?.errID) || "" : ""
        msg && info({ ...defaultModalParams, content: msg, getContainer: document.getElementById('isf-web-plugins')})
    }

    const initData = () => {
        setLoading(true)
        setIsError(false)
    }

    const onAddRole = () => {
        setRoleDrawer(true)
        setOperationType(OperationTypeEnum.AddRole)
    }

    const onAddRoleSuccess = (role) => {
        const formatRole = { ...role, source: RoleClassEnum.User }
        setSearchValue("")
        getBusinessRole({ isClear: true, isSelectFirst: false, offset: 0, limit });
        setCurrentRole(formatRole)
    }

    const onEditRoleInfoSuccess = (role) => {
        setRoleDrawer(false)
        setCurrentRole({...currentRole, ...role})
        const newRoles = roles.map(cur => {
            if(cur.id === role.id) {
                return {...cur,...role}
            }else{
                return cur
            }
        })
        setRoles(newRoles)
    }

    const onEditRole = (role) => {
        setCurrentRole(role)
        setRoleDrawer(true)
        setOperationType(OperationTypeEnum.EditRoleInfo)
    }

    const onEditRolePerm = (role) => {
        setCurrentRole(role)
        setRoleDrawer(true)
        setOperationType(OperationTypeEnum.EditRolePerm)
    }

    const onDeleteRole = async(role) => {
        try {
            await deleteRole(role.id)
            message.success(intl.get("delete.success"))
            const index = roles.findIndex(cur => cur.id === role.id)
            const newRoles = roles.filter(cur => cur.id !== role.id)
            setRoles(newRoles)
            setCurrentRole(newRoles[index > 0 ? index - 1 : 0]);
        }catch(e) {
            handleError(e)
        }
    }

    const onDeleteHandle = async(role) => {
        const modal = confirm({
            centered: true,
            closable: true,
            title: intl.get("delete.role"),
            content: intl.get("delete.role.tip"),
            footer: () => (
                <div style={{ textAlign: 'right' }}>
                    <Button type="primary" danger onClick={() => {
                        onDeleteRole(role)
                        modal.destroy()
                    }}>
                        {intl.get('delete')}
                    </Button>
                    <Button
                        key="back"
                        onClick={() => {
                            modal.destroy();
                        }}
                    >
                        {intl.get('cancel')}
                    </Button>
                </div>
            ),
            onClose: () => {
                modal.destroy();
            },
            getContainer: document.getElementById('isf-web-plugins')
        })
    }

    const updateRoleList = () => {
        getBusinessRole({ isClear: true, isSelectFirst: true, offset: 0, limit });
        setRoleDrawer(false)
    }

    useEffect(() => {
        getBusinessRole({ offset: 0, limit });
    }, []);

    useEffect(() => {
        setSearchValue("")
        if (roleClass === RoleClassEnum.Business) {
            getBusinessRole({ offset: 0, limit });
        }else{
            getSystemRole();
        }
    }, [roleClass]);

    return (
        <div className={styles["role-mgnt"]}>
            <div className={styles["left"]}>
                <div className={styles["header"]}>
                    {
                        !(userInfo?.user?.roles?.length === 1 && userInfo.user.roles[0]?.id === SystemRoleType.Admin) &&
                        <Space>
                            <Radio.Group value={roleClass} onChange={onChangeRadio}>
                                <Radio.Button value={RoleClassEnum.Business}>
                                    {intl.formatMessage({ id: "business.roles" })}
                                </Radio.Button>
                                <Radio.Button value={RoleClassEnum.System}>
                                    {intl.formatMessage({ id: "system.roles" })}
                                </Radio.Button>
                            </Radio.Group>
                        </Space>
                    }
                    {
                        roleClass === RoleClassEnum.Business && (
                            <Input
                                className={styles["search-role"]}
                                allowClear
                                placeholder={intl.get("search.name")}
                                prefix={<SearchOutlined />}
                                value={searchValue}
                                onChange={(e) => {
                                    const value = e.target.value;
                                    setSearchValue(trim(value));
                                    getBusinessRole({ isClear: true, isSelectFirst: true, offset: 0, limit, keyword: trim(value) })
                                }} 
                            />
                        )
                    }
                </div>
                <div className={styles["content"]}  id={"scrollable-role"} style={{height: !(userInfo?.user?.roles?.length === 1 && userInfo.user.roles[0]?.id === SystemRoleType.Admin) ? roleClass === RoleClassEnum.Business ? "calc(100% - 172px)" : "calc(100% - 64px)" : "calc(100% - 52px)"}}>
                    {
                        isError ?
                            <Empty image={loadFailedIcon} description={intl.get("loadFailed")}/>
                            : loading ? 
                                <div className={styles["loading"]}><Spin /></div>
                                : !roles.length ?
                                    <Empty image={searchValue ? SearchEmptyIcon : EmptyIcon} description={intl.get(searchValue ? "no.match.search.result" : "list.empty")}/>
                                    :
                                    <InfiniteScroll
                                        dataLength={roles.length}
                                        next={() => {
                                            if(roleClass === RoleClassEnum.System) return 
                                            getBusinessRole({ isClear: false, isSelectFirst: false, offset: roles.length, limit, keyword: searchValue})
                                        }}
                                        hasMore={roleClass === RoleClassEnum.System ? false : hasMore}
                                        loader={
                                            <div className={styles['loading']}>
                                                <Spin size="small" />
                                            </div>
                                        }
                                        scrollableTarget={"scrollable-role"}
                                        scrollThreshold={0.9}
                                    >
                                        { roles.map((role, index) => {
                                            return (
                                                <div
                                                    key={role.id}
                                                    className={styles["item"]}
                                                    style={{
                                                        backgroundColor:
                                                        currentRole.id === role.id || index === hoverIndex
                                                            ? oemColor.colorPrimaryBg
                                                            : "transparent",
                                                    }}
                                                    title={role.name}
                                                    onClick={() => changeSelect(role)}
                                                    onMouseEnter={() => {
                                                        setHoverIndex(index)
                                                    }}
                                                    onMouseLeave={() => {
                                                        setHoverIndex(undefined)
                                                    }}
                                                >
                                                    <div className={styles["info"]}>
                                                        <div className={styles["name"]}>{role.name}</div>
                                                        {
                                                            role.source ===  RoleClassEnum.Business &&
                                                            <div 
                                                                className={styles["inner"]}
                                                                style={{color: oemColor.colorPrimary, backgroundColor: oemColor.colorPrimaryBg}}
                                                                color="primary" 
                                                                variant="filled" 
                                                                autoInsertSpace={false}
                                                            >
                                                                {intl.get("built.in")}
                                                            </div>
                                                        }
                                                    </div>
                                                    {  
                                                        role.source !==  RoleClassEnum.Business && roleClass === RoleClassEnum.Business && index === hoverIndex &&
                                                        <div className={styles["icon"]}>
                                                            <Dropdown
                                                                trigger={['click']}
                                                                dropdownRender={() => {
                                                                    return (
                                                                        <Menu>
                                                                            <Menu.Item onClick={() => onEditRole(role)}>
                                                                                {intl.get("edit.role.info")}
                                                                            </Menu.Item>
                                                                            <Menu.Item onClick={() => onEditRolePerm(role)}>
                                                                                {intl.get("edit.role.perm")}
                                                                            </Menu.Item>
                                                                            <Menu.Item onClick={()=> onDeleteHandle(role)}>
                                                                                {intl.get("delete")}
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
                                                        </div>
                                                    }
                                                </div>
                                            );
                                        })}
                                    </InfiniteScroll>
                    }
                </div>
                {
                    roleClass === RoleClassEnum.Business && (
                        <div  className={styles["add-role-btn"]}>
                            <Button 
                                className={styles["btn"]}
                                type={"default"} 
                                icon={<AddIcon style={{width: "14px", height: "14px"}}/>} 
                                onClick={onAddRole}
                            >
                                {intl.get("add.role")}
                            </Button>
                        </div>
                    )
                }
            </div>
            <div className={styles["right"]}>
                <div className={styles["header"]}>
                    <div className={styles["title"]}>
                        <div className={styles["introduction"]}>
                            {currentRole && (
                                <span className={styles["name"]}>{currentRole.name}</span>
                            )}
                            {/* {
                                roleClass === RoleClassEnum.Business && <a className={styles["details"]} onClick={() => setShowDetails(true)}>{intl.get("role.details")}</a>
                            } */}
                        </div>
                        {currentRole && (
                            <div className={styles["description"]} title={currentRole.description}>
                                {currentRole.description}
                            </div>
                        )}
                    </div>
                </div>
                {
                    currentRole && (
                        <div className={styles["members"]}>
                            <MembersMgnt currentRole={currentRole} roleClass={roleClass} updateRoleList={updateRoleList}/>
                        </div>
                    )
                }
            </div>
            {
                showDetails && (
                    <RoleDetails roleClass={roleClass} currentRole={currentRole} onCancel={() => setShowDetails(false)} />
                )
            }
            {
                showRoleDrawer && (
                    <SetRole 
                        operationType={operationType}
                        curRole={currentRole}
                        onCancel={()=> {
                            setRoleDrawer(false)
                        }}
                        onAddRoleSuccess={onAddRoleSuccess}
                        onEditRoleInfoSuccess={onEditRoleInfoSuccess}
                        updateRoleList={updateRoleList}
                    />
                )
            }
        </div>
    );
}

export default RoleMgnt;
