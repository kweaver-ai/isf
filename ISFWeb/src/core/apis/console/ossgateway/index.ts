import { currify } from '@/util/currify';
import { privateAPI } from '../../../privateAPI/privateAPI'

/**
 * 获取当前站点的默认对象存储
 */
export const getDefaultStorage: Core.APIs.Console.OssGateWay.GetDefaultStorage = currify(privateAPI, 'ossgateway', 'default-storage', 'get')

/**
 * 获取当前站点可界面管理的对象存储
 * @param app: 产品名，如 'as'
 * @param enable: 是否开启
 */
export const getObjectStorageInfoByApp: Core.APIs.Console.OssGateWay.GetObjectStorageInfoByApp = currify(privateAPI, 'ossgateway', 'objectstorageinfo', 'get', '')

/**
 * 获取指定对象存储服务的信息
 */
export const getObjectStorageInfoById: Core.APIs.Console.OssGateWay.GetObjectStorageInfoById = currify(privateAPI, 'ossgateway', 'objectstorageinfo', 'get')