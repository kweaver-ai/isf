import React from 'react'
import UserOrgMgnt from '@/components/UserOrg/component.view'
import { bootstrap, getMount, unmount } from '@/core/lifecycle'

const mount = getMount(<UserOrgMgnt />)
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