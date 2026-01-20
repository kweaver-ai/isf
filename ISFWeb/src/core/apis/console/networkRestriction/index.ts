import { consolehttp } from '../../../openapiconsole'

/**
 * 获取访问者网段绑定功能状态
 */
export const getState: Core.APIs.Console.NetworkRestriction.GetState = ({ name }, options?) => {
    return consolehttp('get', ['policy-management', 'v1', 'general'], undefined, { name }, options)
}

/**
 * 启用和关闭访问者网段绑定功能
 */
export const setState: Core.APIs.Console.NetworkRestriction.SetState = ({ name, payload }, options?) => {
    return consolehttp('put', ['policy-management', 'v1', 'general', name, 'value'], payload, {}, options)
}

/**
 * 获取网段列表
 */
export const getNetworkList: Core.APIs.Console.NetworkRestriction.GetNetworkList = ({ key_word, offset, limit }, options?) => {
    return consolehttp('get', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network'], undefined, { key_word, offset, limit }, options)
}

/**
 * 新增网段
 */
export const addNetwork: Core.APIs.Console.NetworkRestriction.AddNetwork = ({ name, start_ip, end_ip, ip_address, netmask, net_type, ip_type }, options?) => {
    return consolehttp('post', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network'], { name, start_ip, end_ip, ip_address, netmask, net_type, ip_type }, {}, { ...options, resHeader: true })
}

/**
 * 获取网段信息
 */
export const getNetworkInfo: Core.APIs.Console.NetworkRestriction.GetNetworkInfo = ({ id }, options?) => {
    return consolehttp('get', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network', id], undefined, {}, options)
}

/**
 * 修改网段
 */
export const editNetwork: Core.APIs.Console.NetworkRestriction.EditNetwork = ({ id, name, start_ip, end_ip, ip_address, netmask, net_type, ip_type }, options?) => {
    return consolehttp('put', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network', id], { id, name, start_ip, end_ip, ip_address, netmask, net_type, ip_type }, {}, options)
}

/**
 * 删除网段
 */
export const deleteNetwork: Core.APIs.Console.NetworkRestriction.DeleteNetwork = ({ id }, options?) => {
    return consolehttp('delete', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network', id], undefined, {}, options)
}

/**
 * 获取网段绑定的访问者列表
 */
export const getAccessorsByNetwork: Core.APIs.Console.NetworkRestriction.GetAccessorsByNetwork = ({ id, key_word, offset, limit }, options?) => {
    return consolehttp('get', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network', id, 'accessor'], undefined, { key_word, offset, limit }, options)
}

/**
 * 向绑定的网段新增访问者
 */
export const addAccessorsByNetwork: Core.APIs.Console.NetworkRestriction.AddAccessorsByNetwork = ({ id, accessorsList }, options?) => {
    return consolehttp('post', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network', id, 'accessor'], [...accessorsList], {}, options)
}

/**
 * 删除网段绑定的访问者
 */
export const deleteAccessorByNetwork: Core.APIs.Console.NetworkRestriction.DeleteAccessorByNetwork = ({ id, accessor_id }, options?) => {
    return consolehttp('delete', ['policy-management', 'v1', 'user-login', 'network-restriction', 'network', id, 'accessor', accessor_id], undefined, {}, options)
}

/**
 * 获取访问者已绑定的网段
 */
export const getNetworkListByAccessor: Core.APIs.Console.NetworkRestriction.GetNetworkListByAccessor = ({ id, offset, limit }, options?) => {
    return consolehttp('get', ['policy-management', 'v1', 'user-login', 'network-restriction', 'accessor', id, 'network'], undefined, { offset, limit }, options)
}