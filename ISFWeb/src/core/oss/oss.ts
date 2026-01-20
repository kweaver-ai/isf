import __ from './locale';

export enum LocationType {

    /**
     * 用户
     */
    User = 'user',

    /**
     * 组织
     */
    Organization = 'organization',

    /**
     * 部门
     */
    Department = 'department',

    /**
     * 文档库
     */
    DocLib = 'docLib',
}

/**
 * 格式化对象存储展示的信息
 * @param ossInfo 对象存储
  */
export function displayUserOssInfo(ossInfo, locationType?: LocationType) {
    const { ossId } = ossInfo;
    return !ossId ?
        (locationType === LocationType.DocLib ? __('未指定（跟随文件上传者的指定存储位置）') : __('未指定（使用默认存储）'))
        :
        ossInfo.ossName
}