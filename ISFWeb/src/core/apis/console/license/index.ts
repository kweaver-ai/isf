import { consolehttp } from "../../../openapiconsole";

/**
 * 获取许可证信息
 */
export const getLicenseInfo = () => {
    return consolehttp(
        "get",
        ["license", "v1", "console", "licenses"],
        null,
        null
    )
}

/**
 * 查询授权对象授权产品信息
 */
export const getAuthorizedProducts = ({
    user_ids,
}: {
    user_ids: string[];
}) => {
    return consolehttp(
        "post",
        ["license", "v1", "console", "query-authorized-products"],
        { method: 'GET', user_ids },
        null,
    )
}

/**
 * 修改授权对象授权产品信息
 */
export const updateAuthorizedProducts = (data : {
    id: string;
    type: string;
    products: string[];
}[]) => {
    return consolehttp(
        "put",
        ["license", "v1", "console", "authorized-products"],
        data,
        null
    )
}


