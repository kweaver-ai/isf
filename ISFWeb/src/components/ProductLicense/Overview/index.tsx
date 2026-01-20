import React, { useContext, useEffect, useState } from "react"
import styles from "./styles.css"
import classNames from "classnames"
import AppConfigContext from "@/core/context/AppConfigContext";
import { getLicenseInfo } from "@/core/apis/console/license";
import intl from "react-intl-universal";
import { message } from "antd";

interface ProductLicenseType {
    product: string;
    total_user_quota: number;
    authorized_user_count: number;
}

export const ProductLicenseOverview = () => {
    const { lang } = useContext(AppConfigContext);
    const [products, setProducts] = useState<ProductLicenseType[]>([])
    
    // 判断是否授权数超标
    function isQuotaExceeded(item: ProductLicenseType): boolean {
        return item.authorized_user_count >= item.total_user_quota && item.total_user_quota !== -1
    }

    // 格式化数值显示
    function formatNumber(value: number, isEnglish: boolean): string {
        if (value === -1 || value === Infinity) {
            return intl.get('unlimited');
        }

        // 数值截断到一位小数的辅助函数
        const truncateToOneDecimal = (num: number): string => {
            return (Math.floor(num * 10) / 10).toFixed(1);
        };

        if (isEnglish) {
            if (value < 1000000) {
                return value.toString();
            } else if (value < 1000000000) {
                return `${truncateToOneDecimal(value / 1000000)}M`;
            } else {
                return `${truncateToOneDecimal(value / 1000000000)}B`;
            }
        } else {
            if (value < 10000) {
                return value.toString();
            } else if (value < 100000000) {
                return `${truncateToOneDecimal(value / 10000)}${intl.get('thousand')}`;
            } else {
                return `${truncateToOneDecimal(value / 100000000)}${intl.get('million')}`;
            }
        }
    }

    const getProducts = async() => {
        try{
            const data = await getLicenseInfo()
            if(!data.length) {
                message.info(intl.get("no.use.product"))
            }
            setProducts(data)
        }catch(err){
            if(err?.description) {
                message.info(err.description)
            }
        }
    }

    useEffect(() => {
       getProducts()
    }, [])
    
    return (
        <div className={styles.container}>
            {
                products.map((item) => (
                    <div 
                        key={item.product} 
                        className={classNames(styles.item, {
                            [styles.singleItem]: products.length === 1
                        })}
                    >
                        <div className={styles.product} title={item.product}>{item.product}</div>
                        <div 
                            className={classNames(styles.quota, {
                                [styles.exceeded]: isQuotaExceeded(item)
                            })}
                        >
                            <span>{intl.get('authorized.user.count')}：</span>
                            <span 
                                className={classNames({
                                    [styles.exceeded]: isQuotaExceeded(item)
                                })}
                                title={`${item.authorized_user_count}/${item.total_user_quota === -1 ? intl.get('unlimited') : item.total_user_quota}`}
                            >
                                {formatNumber(item.authorized_user_count, lang === 'en-us')}/
                                {formatNumber(item.total_user_quota, lang === 'en-us')}
                            </span>
                        </div>
                    </div>
                ))
            }
        </div>
    )
}