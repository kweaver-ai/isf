import React, { useState } from "react"
import intl from "react-intl-universal"
import { createRoles, EditRole } from "@/core/apis/console/authorization";
import { Button, Drawer, Steps, message, Modal } from "antd"
import { RoleInfo } from "../RoleInfo";
import { Authorization } from "../Authorization";
import styles from "./styles.css"
import { OperationTypeEnum } from "../types";
import { AuthorizationErrorCodeEnum } from "@/core/apis/console/authorization/type";
import { defaultModalParams } from "@/util/modal";
import { trim } from "lodash";
import { isNormalShortName } from "@/util/validators";

const { info } = Modal

export const SetRole = ({operationType = OperationTypeEnum.AddRole, curRole, onCancel, updateRoleList, onAddRoleSuccess, onEditRoleInfoSuccess}:{ operationType?: OperationTypeEnum; curRole?: { id: string; name: string }; onCancel: () => void; updateRoleList?:() => void; onAddRoleSuccess?: (roleInfo) => void; onEditRoleInfoSuccess?:(roleInfo) => void }) => {
    const [step, setStep] = useState(0)
    const [roleInfo, setRoleInfo] = useState({name: "", description: "", resource_type_scope: {types: [], unlimited: true} })
    const [nameTip, setNameTip] = useState("")
    const [descriptionTip, setDescriptionTip] = useState("")
    const [roleId, setRoleId] = useState("")
    
    const getRoleInfoValue = ({roleName, description, range }) => {
        const types = range.map(cur => ({ id: cur.value }))
        setRoleInfo({
            name: trim(roleName),
            description: trim(description),
            resource_type_scope: {
                types,
                unlimited: !range.length
            }
        })
    }

    const preCheckRoleInfo = () => {
        let result = true
        if(!trim(roleInfo.name)) {
            setNameTip(intl.get("role.name.placeholder"))
            result = false
        }
        if(trim(roleInfo.name) && !isNormalShortName(trim(roleInfo.name))) {
            setNameTip(intl.get("role.name.invalid"))
            result = false
        }
        if(trim(roleInfo.description).length > 300) {
            setDescriptionTip(intl.get("role.description.invalid.length"))
            result = false
        }
        return result
    }

    const handleError = (error) => {
        if(error?.code === AuthorizationErrorCodeEnum.RoleNotFound) {
            info({ 
                ...defaultModalParams, 
                closable: false,
                content: intl.get("role.not.exist"), 
                getContainer: document.getElementById('isf-web-plugins'), 
                onOk: () => {
                    updateRoleList?.()
                }
            })
        }else {
            const msg = error?.description || ""
            msg && info({ ...defaultModalParams, content: msg, getContainer: document.getElementById('isf-web-plugins')})
        }
    }

    const onSave = async() => {
        try {
            if(preCheckRoleInfo()) {
                if(operationType === OperationTypeEnum.AddRole) {
                    const {id} = await createRoles(roleInfo)
                    message.success(intl.get("role.create.success"))
                    setStep(step + 1)
                    setRoleId(id)
                    onAddRoleSuccess?.({...roleInfo, id})
                }else {
                    await EditRole({id: curRole?.id, ...roleInfo})
                    message.success(intl.get("edit.success"))
                    onEditRoleInfoSuccess?.({...roleInfo, id: curRole?.id })
                }
               
            }
        }catch(e) {
            if(e?.code === AuthorizationErrorCodeEnum.RoleNameConflict) {
                setNameTip(intl.get("role.name.duplicate"))
                return
            }
            handleError(e)
        }
    }

    const getDrawerTitle = () => {
        switch(operationType) {
            case OperationTypeEnum.AddRole: 
                return intl.get("add.role")
            case OperationTypeEnum.EditRoleInfo: 
                return intl.get("edit.role.info")
            case OperationTypeEnum.EditRolePerm: 
                return <div className={styles["custom-drawer-title"]} title={intl.get("edit.role.perm") + ` - ${curRole?.name}`}>{intl.get("edit.role.perm") + ` - ${curRole?.name}`}</div>
        }
    }

    return (
        <Drawer
            title={getDrawerTitle()}
            open={true}
            closable={true}
            maskClosable={true}
            destroyOnHidden={true} 
            destroyOnClose={true}
            styles={{ wrapper: {
                width: operationType === OperationTypeEnum.EditRoleInfo ? "45vw" : "60vw",
                minWidth: operationType === OperationTypeEnum.EditRoleInfo ? "600px" : "900px"
            } }}
            onDrawerClose={onCancel}
            onClose={onCancel}
            footer={step !== 1 && operationType !== OperationTypeEnum.EditRolePerm?  
                <div className={styles["drawer-footer"]}>
                    <Button key="ok" type="primary" className={styles["btn"]} onClick={() => onSave()}>{intl.get("ok")}</Button>
                    <Button key="cancel" onClick={onCancel}>{intl.get("cancel")}</Button>
                </div>
                : null
            }
        >
            <div className={styles["create-role"]}>
                {
                    operationType === OperationTypeEnum.AddRole ? (
                        <Steps current={step} className={styles["step"]} items={[{title: intl.get("role.info")}, {title: intl.get("role.auth")}]}/>
                    ) : null
                }
                {
                    (step === 0 || operationType === OperationTypeEnum.EditRoleInfo) && operationType !== OperationTypeEnum.EditRolePerm ? 
                        <div className={styles["content"]}>
                            <RoleInfo curRole={operationType === OperationTypeEnum.EditRoleInfo ?  curRole : null} getRoleInfoValue={getRoleInfoValue} nameTip={nameTip} setNameTip={setNameTip} descriptionTip={descriptionTip} setDescriptionTip={setDescriptionTip} updateRoleList={updateRoleList}/>
                        </div>
                        : <Authorization curRole={roleId ? {id: roleId, name: roleInfo?.name} : curRole} topMargin={operationType === OperationTypeEnum.AddRole ? "300px" : "260px"} updateRoleList={updateRoleList}/>
                }
            </div>
        </Drawer>
    )
}