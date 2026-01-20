import React, { useRef, useEffect, useState } from "react";
import styles from "./styles.css"
import { Button, Drawer, Tabs } from "antd";
import intl from "react-intl-universal";
import { MembersMgnt } from "../MembersMgnt";
import { PermView } from "../PermView"
import RoleBgImg from "../../../icons/rolebg.png"
import EditIcon from "../../../icons/edit.svg"
import { SetRole } from "../SetRole";
import { OperationTypeEnum, RoleClassEnum } from "../types";

export const RoleDetails = ({onCancel, roleClass, currentRole}) => {
    const headerRef = useRef(null)
    const [showRoleInfo, setShowRoleInfo]= useState(false)
    const [headerHeight, setHeaderHeight] = useState<number | undefined>(undefined)
    const [scopes, setScopes] = useState<string | undefined>()

    useEffect(() => {
        if(currentRole?.resource_type_scopes) {
            const scopes = currentRole.resource_type_scopes.map((item) => {
                return item.name
            })
            setScopes(scopes.join("、"))
        }
    }, [currentRole])
    
    useEffect(() => {
        const updateHeight = () => {
            if (headerRef.current) {
                const height = headerRef.current.clientHeight;
                setHeaderHeight(height);
            }
        };

        // 使用 requestAnimationFrame 确保在浏览器重绘前执行
        const animationFrameId = requestAnimationFrame(() => {
            updateHeight();
        });

        window.addEventListener('resize', updateHeight);

        return () => {
            // 清除动画帧
            cancelAnimationFrame(animationFrameId);
            window.removeEventListener('resize', updateHeight);
        };
    }, []);

    return (
        <Drawer
            className={styles["role-details-drawer"]}
            title={intl.get("role.details")} 
            open={true}
            maskClosable={true} 
            destroyOnClose={true}
            destroyOnHidden={true} 
            width={"50vw"}
            onDrawerClose={onCancel} 
            onClose={() => {
                onCancel()
            }}
        >
            <div className={styles["role-details"]} style={{backgroundSize: `100% ${headerHeight + 46}px`, backgroundImage: `url(${RoleBgImg})`, pointerEvents: "all",}}>
                <div className={styles["header"]} ref={headerRef}>
                    <div className={styles["name-info"]}>
                        <div className={styles["name"]} title={currentRole.name}>{currentRole.name}</div>
                        {
                            currentRole.source !== RoleClassEnum.Business &&
                            <Button type={"text"} icon={<EditIcon style={{width: "16px", height: "16px"}} onClick={() => {
                                setShowRoleInfo(true)
                            }}/>} />
                        }
                    </div>
                    <div className={styles["description"]}>{currentRole.description}</div>
                    <div className={styles["scope"]}>
                        <span>{intl.get("scope-tip")}</span>
                        <span>{scopes}</span>
                    </div>
                </div>
                <Tabs destroyOnHidden={true}>
                    <Tabs.TabPane tab={intl.get("permission")} key="permission">
                        <PermView curRole={currentRole}/>
                    </Tabs.TabPane>
                    <Tabs.TabPane tab={intl.get("members")} key="members">
                        <MembersMgnt roleClass={roleClass} currentRole={currentRole} />
                    </Tabs.TabPane>
                </Tabs>
            </div>
            {
                showRoleInfo &&(
                    <SetRole
                        operationType={OperationTypeEnum.EditRoleInfo}
                        curRole={currentRole}
                        onCancel={() =>{
                            setShowRoleInfo(false)
                        }}
                        onEditRoleInfoSuccess={(roleInfo) => {
                            console.info(roleInfo)
                        }}
                    /> 
                )
            }
        </Drawer>
    )
}