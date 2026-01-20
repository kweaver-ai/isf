import React, { useState } from "react"
import styles from "./styles.css"
import { Button } from "antd"
import intl from "react-intl-universal"
import { SetRole } from "../SetRole"
import { OperationTypeEnum, RoleClassEnum } from "../types"
export const PermView = (curRole) => {
    const [showRoleAuth, setShowRoleAuth] = useState(false)

    return (
        <div className={styles["perm-view"]}>
            <div className={styles["header"]}>
                <div>{intl.get("perm.view.tip", {count: 22})}</div>
                {
                    curRole.source !==  RoleClassEnum.Business &&
                    <Button onClick={() => {
                        setShowRoleAuth(true)
                    }}>{intl.get("edit")}</Button>
                }
            </div>
            {
                showRoleAuth && (
                    <SetRole
                        operationType={OperationTypeEnum.EditRolePerm}
                        curRole={curRole}
                        onCancel={() =>{
                            setShowRoleAuth(false)
                        }}
                        onEditRoleInfoSuccess={(roleInfo) => {
                            console.info(roleInfo)
                        }}
                    />
                )
            }
        </div>
    )
}