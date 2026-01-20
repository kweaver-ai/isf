import React, { useState, useEffect } from 'react';
import intl from 'react-intl-universal';
import { Button, Checkbox, Col, message, Modal, Row } from 'antd';
import { getLicenseInfo, updateAuthorizedProducts } from '@/core/apis/console/license';
import styles from './styles.css'
import { defaultModalParams } from '@/util/modal';
import { PublicErrorCode } from '@/core/apis/openapiconsole/errorcode';

export const ProductLicense = ({users, onComplete, onSuccess}) => {
    const [products, setProducts] = useState([])
    const [ values, setValues ] = useState([])
    const { confirm, info } = Modal

    const updateProductLicense = async(modalConfirm = undefined) => {
       try{
            const data = users.map((item) => ({
                id: item.id,
                type: 'user',
                products: values,
            }))
            const userInfo = users?.map((item) => ({
                ...item,
                user: {
                    ...item.user,
                    products: values,
                },
            }))
            await updateAuthorizedProducts(data)
            modalConfirm?.destroy()
            onSuccess(userInfo)
        }catch(err) {
            if(err?.code === PublicErrorCode.NotFound) {
                const notExistUser = users.filter((item) => err.detail.ids?.includes(item.id))
                const modalInfo = info({
                    ...defaultModalParams,
                    title: intl.get("user.deleted"),
                    content: (
                        <div className={styles["user-list-tip"]}>
                            {
                                notExistUser.map((cur) =>(
                                    <div key={cur.id} className={styles["item"]} title={cur.user.displayName}>{cur.user.displayName}</div>
                                ))
                            }
                        </div>
                    ),
                    onOk: () => {
                        modalInfo?.destroy();
                        modalConfirm?.destroy();
                        onSuccess(users)
                    },
                    getContainer: document.getElementById('isf-web-plugins'), 
                })
            }else {
                if(err?.description || err?.message) {
                    info({
                        ...defaultModalParams,
                        title: intl.get('tip'),
                        content: (
                            <div style={{whiteSpace: 'pre-line'}}>
                                {
                                    err?.description || err.message
                                }
                            </div>
                        ),
                        getContainer: document.getElementById('isf-web-plugins'), 
                    })
                }
            }
        }
    }

    const changeProductLicense = async() => {
        if(values.length === 0) {
            const modalConfirm = confirm({
                ...defaultModalParams, 
                closable: false,
                title: intl.get('product.license.cancel'),
                content: intl.get('product.license.cancel.content'),
                footer:() => (
                    <div style={{textAlign: 'right'}}>
                        <Button type="primary" onClick={() => updateProductLicense(modalConfirm)}>
                            {intl.get('ok')}
                        </Button>
                        <Button
                            onClick={() => {
                                modalConfirm.destroy();
                            }}
                        >
                            {intl.get('cancel')}
                        </Button>
                    </div>
                ),
            onClose: () => {
                modalConfirm.destroy();
            },
                getContainer: document.getElementById('isf-web-plugins'), 
            })
        }else {
            updateProductLicense()
        }
    }

    const getProductLicense = async() => {
        try{
            const data = await getLicenseInfo()
            if(!data.length) {
                message.info(intl.get("no.use.product"))
            }
            setProducts(data)
        
            const availableProducts = data.map(item => item.product)
            
            const userProducts = users?.[0]?.user?.products || []
            const initValues = users?.length <= 1 
                ? userProducts.filter(product => availableProducts.includes(product))
                : []
            
            setValues(initValues)
        }catch(err) {
            if(err?.description) {
                message.info(err.description)
            }
        }
    }

    useEffect(() => {
       getProductLicense()
    }, [])

    return (
        <Modal
            centered
            maskClosable={false}
            open={true}
            onCancel={onComplete}
            title={intl.get('product.license')}
            footer={[
                <Button type="primary" disabled={!products.length} onClick={changeProductLicense}>
                    {intl.get('ok')}
                </Button>,
                <Button
                    onClick={onComplete}
                >
                    {intl.get('cancel')}
                </Button>,
            ]}
            getContainer={document.getElementById("isf-web-plugins") as HTMLElement}
        >
           <div className={styles["product-license"]}>
                <div className={styles["title"]}>
                    {users?.length > 1 ? intl.get('product.license.batch', {count: users.length}) : intl.get('product.license.single', {name: users?.[0]?.user?.displayName})}
                    
                </div>
                <div className={styles["content"]}>
                    <Checkbox.Group 
                        style={{ width: '100%' }} 
                        value={values} 
                        onChange={(checkedValues) => {
                            setValues(checkedValues)
                        }}
                    >
                        <Row>
                            {products.map((item) => (
                                <Col key={item.product} span={8}>
                                    <Checkbox key={item.product} value={item.product} style={{ verticalAlign: 'middle' }}>
                                        <span className={styles["product-name"]} title={item.product}>{item.product}</span>
                                    </Checkbox>
                                </Col>
                            ))}
                        </Row>
                    </Checkbox.Group>
                </div>
                {
                    users?.length > 1 &&
                    <div className={styles["batch-tip"]}>{intl.get('product.license.batch.tip')}</div>
                }
           </div>
        </Modal>
    )
}