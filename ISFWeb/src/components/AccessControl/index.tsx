import React, { useState, useContext } from "react"
import { Button, Divider } from "antd";
import session from '@/util/session';
import styles from './styles';
import intl from "react-intl-universal";
import LoginSecurityPolicy from "./LoginSecurityPolicy/component.view";
import VisitorNetBind from "./VisitorNetBind/component.view";
import DeviceBind from "./DeviceBind/component.view";
import BindingQuery from "./BindingQuery/component.view";
import Tabs from '@/ui/Tabs/ui.desktop'

export enum TabEnum {
    LoginSecurity = "login-security",
    VisitorSubnetBinding = "visitor-subnet-binding",
    DeviceBinding = "device-binding",
    BindingQuery = "binding-query",
}
function AccessControl () { 
    const tabs = [{key: TabEnum.LoginSecurity, label: intl.get("login.security.policy")}, {key: TabEnum.VisitorSubnetBinding, label: intl.get("visitor.subnet.binding")}, {key: TabEnum.DeviceBinding, label: intl.get("device.binding")}, {key: TabEnum.BindingQuery, label: intl.get("binding.query")}]

    return (
        <div className={styles['access-control']}>
            <Tabs>
                <Tabs.Navigator key={'navigator'}>
                    {
                        tabs.map((tab, index) => (
                            <Tabs.Tab key={tab.key} className={styles['tab']}>
                                {tab.label}
                                {(index !== tabs.length - 1) && <span key={tab + 'line'} className={styles['line']}></span>}
                            </Tabs.Tab>
                        ))
                    }
                </Tabs.Navigator>
                <Tabs.Main className={styles['tabs-main']}>
                    {
                        tabs.map(tab => {
                            switch(tab.key) {
                                case TabEnum.LoginSecurity:
                                    return <Tabs.Content key={tab.key} className={styles['tabs-content']}><LoginSecurityPolicy /></Tabs.Content>
                                case TabEnum.VisitorSubnetBinding:
                                    return <Tabs.Content key={tab.key} className={styles['tabs-content']}><VisitorNetBind /></Tabs.Content>
                                case TabEnum.DeviceBinding:
                                    return <Tabs.Content key={tab.key} className={styles['tabs-content']}><DeviceBind /></Tabs.Content>
                                case TabEnum.BindingQuery:
                                    return <Tabs.Content key={tab.key} className={styles['tabs-content']}><BindingQuery  userid={session.get('isf.userid')}/></Tabs.Content>
                            }
                        })
                    }
                </Tabs.Main>
            </Tabs>
        </div>
    )
}

export default AccessControl;