import React, { useState, useContext } from "react"
import DomainManage from "../DomainManage/component.view"
import ThirdConfig from "../ThirdConfig/component.view"
import __ from './locale'
import styles from './styles';
import { Tabs } from "antd";
import AppConfigContext from "@/core/context/AppConfigContext";
import AccessControl from "../AccessControl";
import { SystemRoleType } from "@/core/role/role";

export enum TabEnum {
    LoginAuth = "login-auth",
    DomainAuth = "domain-auth",
    ThridAuth = "thrid-auth"
}

function CertifictionMgnt () {
    const { prefix, config: { userInfo } } = useContext(AppConfigContext)

    const params = new URLSearchParams(location.search)
    const [urlParams, setUrlParams] = useState(params)

    const onChangeTab = (tabId) => {
        const url = new URL(window.location.href);
        url.searchParams.set('tab', tabId);
        window.history.pushState({}, '', url)
        setUrlParams(url.searchParams)

    }

    const isShowLoginAuth = () => {
        return userInfo.user.roles.some((role) => role.id === SystemRoleType.Supper || role.id === SystemRoleType.Securit)
    }

    const isShowDomainAuthAndThirdAuth = () => {
        return userInfo.user.roles.some((role) => role.id === SystemRoleType.Supper || role.id === SystemRoleType.Admin)
    }

    let tabs = []
    tabs = [...(isShowLoginAuth() ? [{ id: TabEnum.LoginAuth, label: __('登录认证') }] : [])]
    tabs = [...tabs, ...(isShowDomainAuthAndThirdAuth() ? [{ id: TabEnum.DomainAuth, label: __('域认证') }, { id: TabEnum.ThridAuth, label: __('第三方认证') }] : [])]

    return (
        <div className={styles['cerifiction-mgnt']}>
            <Tabs destroyOnHidden={true} activeKey={urlParams?.get("tab") || tabs?.[0]?.id} onChange={onChangeTab} >
                {
                    tabs.map((tab) => (
                        <Tabs.TabPane key={tab.id} tab={tab.label} className={styles['tab']}>
                            {
                                tab.id === TabEnum.LoginAuth ? 
                                    <AccessControl /> : 
                                    tab.id === TabEnum.DomainAuth ? 
                                        <DomainManage /> : 
                                        tab.id === TabEnum.ThridAuth ? 
                                            <ThirdConfig 
                                                swf={'/res/libs/webuploader/Uploader.swf'}
                                                server={`${prefix}/isfweb/api/update_thirdparty/update_thirdparty/`}
                                            /> : null
                            }
                        </Tabs.TabPane>
                    ))
                }
            </Tabs>
        </div>
    )
}

export default CertifictionMgnt;