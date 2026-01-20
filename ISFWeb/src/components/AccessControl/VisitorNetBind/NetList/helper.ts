import { IpVersion, NetInfo as OutNetInfo } from '../../NetSegment/helper';
import { NetInfo, NetType } from '../helper';
import __ from './locale';

/**
 * 所有外网网段
 */
export const PubliceNet = {
    name: __('所有外网网段'),
    netInfo: '0.0.0.0-9.255.255.255' + __('，') +
        '11.0.0.0-172.15.255.255' + __('，') +
        '172.32.0.0-192.167.255.255' + __('，') +
        '192.169.0.0-255.255.255.255',
}

/**
 * 数据key类型
 */
export const NetSegmentKey = {
    id: 'id',
    ip: 'ip',
    originIP: 'originIP',
    endIP: 'endIP',
    mask: 'mask',
    ip_type: 'ipVersion',
    net_type: 'single',
}

/**
 * 格式化NetSegment要接收的数据
 */
export const setNetSegment = (netInfo: NetInfo): OutNetInfo => {
    return {
        name: netInfo.name,
        id: netInfo.id,
        ip: netInfo.ip,
        originIP: netInfo.originIP,
        endIP: netInfo.endIP,
        ipVersion: netInfo.ipVersion || IpVersion.Ipv4,
        mask: netInfo.ipVersion === IpVersion.Ipv4
            ? netInfo.mask
            : '',
        prefix: netInfo.ipVersion === IpVersion.Ipv6
            ? netInfo.mask
            : '',
        single: netInfo.netType === NetType.Mask,
        ...(() => (
            netInfo.netType === NetType.Mask
                ? {
                    originIP: '',
                    endIP: '',
                }
                : { ip: '' }
        ))(),
    }
}

/**
 * 格式化接收的NetSegment数据
 */
export const getNetSegment = (netInfo: OutNetInfo): NetInfo => {
    return {
        name: netInfo.name,
        id: netInfo.id,
        ip: netInfo.ip,
        originIP: netInfo.originIP,
        endIP: netInfo.endIP,
        ipVersion: netInfo.ipVersion,
        netType: netInfo.single
            ? NetType.Mask
            : NetType.Range,
        mask: netInfo.ipVersion === IpVersion.Ipv4
            ? netInfo.mask
            : netInfo.prefix,
    }
}

