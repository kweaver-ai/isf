import React, { useState, useEffect } from "react"
import intl from "react-intl-universal"
import { Button, Form, Input, Modal, message } from "antd"
import { trim } from "lodash"
import { isNormalName } from "@/util/validators"
import { createUseAccountMgnt, setUseAccountMgnt } from "@/core/apis/console/useaccountmgnt"
import { UserManagementErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { AccountItemType } from "../types"
import styles from "./styles.css"

const { info } = Modal
export const CreateOrEditAccount = ({isEdit = false, onCancel, handleError, createSuccess, cur = null, editSuccess, updateList}:{isEdit?: boolean; onCancel: () => void, handleError: (error: any) => void; createSuccess?: (id: string) => void; cur?: AccountItemType | null; editSuccess?: (cur) => void; updateList?:() => void}) => {
    const [accountName, setAccountName] = useState("")
    const [nameTip, setNameTip] = useState("")

    const preCheck = () => {
        if(!accountName) {
            setNameTip(intl.get("account.name.placeholder"))
            return false
        }

        if(!isNormalName(accountName)){
            setNameTip(intl.get("invalid.account.name"))
            return false
        }

        return true
    }

    const handleErrorModal = (content) => {
        const modal = info({
            centered: true,
            closable: true,
            title: intl.get("tip"),
            content,
            footer: () => (
                <div style={{ textAlign: 'right' }}>
                    <Button type="primary" onClick={() => {
                        updateList()
                        modal?.destroy()
                        onCancel()
                    }}>
                        {intl.get('ok')}
                    </Button>
                    <Button
                        key="back"
                        onClick={() => {
                            modal.destroy()
                        }}
                    >
                        {intl.get('cancel')}
                    </Button>
                </div>
            ),
            onClose: () => {
                modal.destroy()
            },
            getContainer: document.getElementById('isf-web-plugins')
        })
    }

    const onCreate = async () => {
        try {
            if(preCheck()) {
                if(isEdit) {
                    await setUseAccountMgnt({id: cur.id, fields: "name", name: accountName, })
                    editSuccess({...cur, name: accountName})
                    message.success(intl.get("edit.success"))
                }else {
                    const { id } = await createUseAccountMgnt({name: accountName, password: ""})
                    createSuccess(id)
                }
                onCancel()
            }
        }catch(error) {
            if(error?.code) {
                const { code, detail} = error
                switch(code) {
                    case UserManagementErrorCode.AppConflict:
                        if(detail) {
                            if(detail.type === "user") {
                                setNameTip(intl.get("same.account.name"))
                            }else {
                                setNameTip(intl.get("duplicate.account.name"))
                            }
                        }
                        break
                    case UserManagementErrorCode.AppNotFound:
                        await handleErrorModal(intl.get("account.not.found"))
                        updateList()
                        break
                    default:
                        handleError(error)
                } 
            }else {
                handleError(error)
            }
        }
    }

    useEffect(() => {
        if(cur) {
            setAccountName(cur.name)
        }
    },[cur])

    return (
        <Modal
            centered
            open={true}
            maskClosable={false}
            title={intl.get(isEdit ? "rename" : "create.app.account")}
            width={454}
            onCancel={onCancel}
            footer={[
                <Button type="primary" key="ok" onClick={onCreate}>{intl.get("ok")}</Button>,
                <Button key="cancel" onClick={onCancel}>{intl.get("cancel")}</Button>
            ]}
            getContainer={document.getElementById("isf-web-plugins") as HTMLElement}
        >
            <Form className={styles["create-app-account-form"]}>
                <Form.Item 
                    label={intl.get("account.name")} 
                    required
                    validateStatus={nameTip ? "error" : ""}
                    help={nameTip ? <div>{nameTip}</div>: ""}
                >
                    <Input 
                        placeholder={intl.get("account.name.placeholder")} 
                        value={accountName}
                        onChange={(e) => {
                            setNameTip("")
                            const value = trim(e.target.value)
                            setAccountName(value)
                        }}
                        onBlur={(e) => {
                            const value = trim(e.target.value)
                            if(value && !isNormalName(value)){
                                setNameTip(intl.get("invalid.account.name"))
                            }
                        }}
                    />
                </Form.Item>
            </Form>
        </Modal>
    )
}