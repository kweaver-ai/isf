import React, { useState, useEffect } from "react";
import intl from "react-intl-universal";
import { Button, Dropdown, Input, Menu, Modal, Table, message } from "antd";
import { getUseAccountMgnt,  generateToken, delUseAccountMgnt } from "@/core/apis/console/useaccountmgnt";
import OperationIcon from "../../icons/operation.svg"
import AddIcon from "../../icons/add.svg"
import CopyIcon from "../../icons/copy.svg"
import EmptyIcon from "../../icons/empty.png";
import loadFailedIcon from "../../icons/loadFailed.png";
import SearchEmptyIcon from "../../icons/searchEmpty.png";
import { CreateOrEditAccount } from "./CreateOrEditAccount";
import { AccountItemType } from "./types";
import { trim } from "lodash"
import styles from "./styles.css";
import { SearchOutlined } from "@ant-design/icons";
import { UserManagementErrorCode } from "@/core/apis/openapiconsole/errorcode";

const { confirm } = Modal
export const UseAccountMgnt = () => {
    const [isLoading, setLoading] = useState(false)
    const [isError, setError] = useState(false)
    const [data, setData] = useState<AccountItemType[]>([])
    const [total, setTotal] = useState<number>(0)
    const [curPageSize, setPageSize] = useState<number>(50);
    const [curPage, setPage] = useState<number>(1);
    const [showCreateAccount, setCreateAccount] = useState<boolean>(false)
    const [token, setToken] = useState<string>("")
    const [showToken, setShowToken] = useState<boolean>(false)
    const [isEdit, setEdit] = useState<boolean>(false)
    const [cur, setCur] = useState<AccountItemType | null>(null)
    const [searchKey, setSearchKey] = useState<string>("")

    const initData = () => {
        setToken("")
        setEdit(false)
        setCur(null)
    }

    const getUserAccount = async({offset = 0, limit = 50, keyword = ""}) => {
        try {
            setLoading(true)
            const {entries, total_count} = await getUseAccountMgnt({ offset, limit, direction: 'desc', sort: 'date_created', keyword})
            setLoading(false)
            setData(entries)
            setTotal(total_count)
        }catch(e) {
            setLoading(false)
            setError(true)
        }
    }

    const handleError = (error) => {
        if(error?.description || error?.message) {
            message.info(error?.description || error?.message)
        }

        if(error?.code && error.code === UserManagementErrorCode.AppAccountNotFound) {
            updateCurrentPage(1)
        }
    }

    const reGenerateToken = async (id) => {
        try {
            const { token } = await generateToken({ id })
            setToken(token)
            setShowToken(true)
        }catch(error) {
            handleError(error)
        }
    }

    const createSuccess = async(id: string) => {
        setSearchKey("")
        updateList()
        await reGenerateToken(id)
    }

    const editSuccess = (cur) => {
        const newData = data.map(item => {
            if(item.id === cur.id) {
                return cur
            }else {
                return item
            }
        })
        setData(newData)
    }

    const onRename = (cur) => {
        setCreateAccount(true)
        setEdit(true)
        setCur(cur)
    }

    const onRegenerateToken = async(cur) => {
        const modal = confirm({
            centered: true,
            closable: true,
            title: intl.get("generate.new.token"),
            content: intl.get("regenerate.token.tip", {account: cur.name}),
            footer: () => (
                <div style={{ textAlign: 'right' }}>
                    <Button type="primary" onClick={() => {
                        setEdit(true)
                        reGenerateToken(cur.id)
                        modal.destroy();
                    }}>
                        {intl.get('ok')}
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

    const updateList = (keyword = "") => {
        setPage(1)
        getUserAccount({ offset: 0, limit: curPageSize, keyword })
    }

    const updateCurrentPage = (deleteNumber, keyword = "") => {
        const totalPage = Math.ceil((total - deleteNumber) / curPageSize)
        if(totalPage < curPage && totalPage !== 0) {
            setPage(totalPage)
            getUserAccount({ offset: (totalPage - 1) * curPageSize, limit: curPageSize, keyword });
        }else {
            getUserAccount({ offset: (curPage - 1) * curPageSize, limit: curPageSize, keyword });
        }
    }

    const onDeleteAccount = async(cur) => {
        try {
            await delUseAccountMgnt({id: cur.id})
            message.success(intl.get("delete.success"))
            updateCurrentPage(1)
        }catch(error) {
            handleError(error)
        }
    }

    const onDelete = async(cur) => {
        const modal = confirm({
            centered: true,
            closable: true,
            title: intl.get("delete.account"),
            content: intl.get("delete.account.tip", {account: cur.name}),
            footer: () => (
                <div style={{ textAlign: 'right' }}>
                    <Button type="primary" danger onClick={() => {
                        onDeleteAccount(cur)
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

    const onCopy = async(token) => {
        try{
            await navigator.clipboard.writeText(token);
            message.success(intl.get("copy.success"))
        }catch(e) {
            console.info("paste failed")
        }
    }

    useEffect(() => {
        getUserAccount({offset: 0, limit: 50})
    },[])

    const columns=[
        {
            title: intl.get("account.name"),
            dataIndex: "name",
            key: "name",
            width: "38%",
            render: (name, cur) => {
                return <div className={styles["item"]} title={cur.name}>{cur.name}</div>
            }
        },
        {
            title: intl.get("operation"),
            dataIndex: "operation",
            key: "operation",
            width: "12%",
            render: (operation, cur) => {
                return (
                    <Dropdown
                        trigger={['click']}
                        dropdownRender={() => {
                            return (
                                <Menu>
                                    <Menu.Item onClick={() => onRename(cur)}>
                                        {intl.get("rename")}
                                    </Menu.Item>
                                    <Menu.Item onClick={() => onRegenerateToken(cur)}>
                                        {intl.get("generate.new.token")}
                                    </Menu.Item>
                                    <Menu.Item onClick={() => onDelete(cur)}>
                                        {intl.get("delete")}
                                    </Menu.Item>
                                </Menu>
                            )
                        }}
                    >
                        <Button
                            type="text"
                            size={"small"}
                            icon={<OperationIcon style={{ width: "16px", height: "16px" }} />}
                        />
                    </Dropdown>
                )
            }
        },
        {
            title: intl.get("account.id"),
            dataIndex: "id",
            key: "id",
            width: "50%",
            render: (id, cur) => {
                return <div className={styles["account-id"]}>
                    <div className={styles["item"]} title={id}>{id}</div>
                    <Button className={styles["copy-icon"]} type="text" icon={<CopyIcon style={{width: "16px", height: "16px"}}/>} onClick={() => onCopy(id)}/>
                </div>
            }
        }
    ]

    const changeSearchKey = (keyWord: string) => {
        if(keyWord !== searchKey) {
            setSearchKey(keyWord)
            updateList(trim(keyWord))
        }
    }

    return <div className={styles["user-account"]}>
        <div className={styles["header"]}>
            <Button
                className={styles["btn"]}
                type={"primary"}
                icon={<AddIcon style={{width: "14px", height: "14px"}}/>}
                onClick={() => {
                    initData()
                    setCreateAccount(true)
                }}
            >
                {intl.get("create.account")}
            </Button>
            <Input
                placeholder={intl.get("search.name")}
                prefix={<SearchOutlined />}
                size="middle"
                value={searchKey}
                onChange={(e) => {
                    changeSearchKey(trim(e.target.value))
                }}
                allowClear
            />
        </div>
        <div className={styles["table"]}>
            <Table 
                size="small"
                tableLayout="fixed"
                loading={isLoading}
                columns={columns}
                dataSource={data}
                scroll={{y: 'calc(100vh - 280px)'}}
                locale={{
                    emptyText: (
                        <div className={styles["empty"]}>
                            <img
                                src={isError ? loadFailedIcon : searchKey ? SearchEmptyIcon : EmptyIcon}
                                alt=""
                                width={128}
                                height={128}
                            />
                            <span>{intl.get(isError ? "loadFailed" : searchKey ? "no.search.result":"list.empty")}</span>
                        </div>
                    ),
                }}
                pagination={
                    {
                        current: curPage, 
                        pageSize: curPageSize,
                        total,
                        showSizeChanger: true,
                        showQuickJumper: false,
                        showTotal: (total) => {
                            return intl.get("list.total.tip", { total });
                        },
                        onChange: (page, pageSize) => {
                            setPage(pageSize !== curPageSize ? 1 : page)
                            setPageSize(pageSize);
                            getUserAccount({
                                offset: ((pageSize !== curPageSize ? 1 : page) - 1) * pageSize,
                                limit: pageSize,
                                keyword: searchKey
                            });
                        },
                    }
                }
            />
        </div>
        {
            showCreateAccount && (
                <CreateOrEditAccount
                    isEdit={isEdit}
                    onCancel={() =>{
                        setCreateAccount(false)
                    }}
                    createSuccess={createSuccess}
                    editSuccess={editSuccess}
                    updateList={updateList}
                    handleError={handleError}
                    cur={cur}
                />
            )
        }
        {
            showToken && (
                <Modal
                    centered
                    open={true}
                    maskClosable={false}
                    title={intl.get( isEdit ? "generate.new.token" : "create.app.account")}
                    width={454}
                    onCancel={() => {
                        setShowToken(false)
                    }}
                    footer={[
                        <Button type="primary" key="ok" onClick={() => onCopy(token)}>{intl.get("copy")}</Button>

                    ]}
                    getContainer={document.getElementById("isf-web-plugins") as HTMLElement}
                >
                    <div className={styles["generate-token"]}>
                        <div className={styles["success-tip"]}>{intl.get("generate.token.success.tip")}</div>
                        <div className={styles["token-tip"]}>{intl.get("token.tip")}</div>
                        <Input value={token} />
                    </div>
                </Modal>
            )
        }
    </div>
}