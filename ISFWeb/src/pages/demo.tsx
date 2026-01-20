import React, {useRef, useState} from 'react'
import intl from 'react-intl-universal'
import { bootstrap, getMount, unmount } from '@/core/lifecycle'
import { Button } from "antd"
import { apis, components } from '@dip/components/dist/dip-components.min.js';
import { TextArea } from '@/sweet-ui';

const PermissionMgnt = () => {
    const accessPickerContainerRef = useRef<HTMLDivElement>(null)
    const defaultValue = 
        {
            title: "标题",
            resource: { id: "menu", name: "菜单", type: "menu"},
            pickerParams: {
                isAdmin: true,
                tabs: ['organization', 'group', 'app', 'role'],
                range: ['user', 'department', 'group', 'app', 'role'],
                role:'super_admin',
            }
        }
    
    const [value, setValue] = useState(JSON.stringify(defaultValue))

    const showAccessPicker = () => {
        const unmount = apis.mountComponent(
            components.PermConfig,
            {
                ...JSON.parse(value),
                onCancel: () => {
                    unmount();
                },
            },
            accessPickerContainerRef.current
        );
    }

    return (
        <div >
            <h1>权限配置组件测试</h1>
            <div>1.填写组件参数</div>
            <TextArea value={value} width={400} height={300} onValueChange={({detail}) => {
                console.info({value: detail})
                setValue(detail)
            }}/>
            <Button type='primary' style={{marginTop: 10}} onClick={showAccessPicker}>{intl.get('ok')}</Button>
            <div ref={accessPickerContainerRef}></div>
        </div>
    )
}

const mount = getMount(<PermissionMgnt />)

if (!window.__POWERED_BY_QIANKUN__) {
    mount({
        getToken: () => 'ory_at_meG9MJpP8eDdoD1-1HmnddAFmvK1EHLC2eASKasgnK4.dWxoWBreEnmAcmBwuby3AepV1TGAwMOYRfWTkAxm9K0',
        protocol: 'https:',
        host: location.hostname,
        port: '443',
        prefix: '',
        lang: 'zh-cn',
    })
}

export {
    bootstrap,
    mount,
    unmount,
}