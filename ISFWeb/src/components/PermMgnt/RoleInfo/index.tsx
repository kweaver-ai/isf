import React, { useState, useEffect, useRef }  from "react"
import intl from "react-intl-universal";
import { getRoleInfo, getResources } from "@/core/apis/console/authorization";
import { Form, Input, Modal, Select, } from "antd"
import { trim } from "lodash"
import styles from "./styles.css"
import { defaultModalParams } from "@/util/modal";
import { AuthorizationErrorCodeEnum } from "@/core/apis/console/authorization/type";
import { isNormalShortName } from "@/util/validators";

const { info } = Modal

export const RoleInfo = ({curRole = null, getRoleInfoValue, nameTip, setNameTip, descriptionTip, setDescriptionTip, updateRoleList }) => {
    const [roleName, setRoleName] = useState("")
    const [description, setDescription] = useState("")
    const [range, setRange] = useState([])
    const roleInfoRef = useRef(null)
    const [options, setOptions] = useState([])

    const handleError = (error) => {
        if(error?.code === AuthorizationErrorCodeEnum.RoleNotFound) {
            info({ 
                ...defaultModalParams, 
                closable: false,
                content: intl.get("role.not.exist"), 
                getContainer: document.getElementById('isf-web-plugins'), 
                onOk: () => {
                    updateRoleList?.()
                }
            })
        }else {
            const msg = error?.description || ""
            msg && info({ ...defaultModalParams, content: msg, getContainer: document.getElementById('isf-web-plugins')})
        }
    }

    const getAllResources = async({offset = 0, limit = 1000}) => {
        try {
            const { entries } = await getResources({ offset, limit })
            const options = entries.map(cur => {
                return {value: cur.id, label: cur.name}
            })
            setOptions(options)
        }catch(e) {
            handleError(e)
        }
    }

    const initEditInfo = async(roleId) => {
        try {
            const {name, description, resource_type_scopes} = await getRoleInfo({id: roleId})
            const {unlimited, types} = resource_type_scopes
            setRoleName(name)
            setDescription(description)
            const range = types.map(cur => {
                return {value: cur.id, label: cur.name}
            })
            setRange(unlimited ? [] : range)
            getRoleInfoValue({roleName: name, description, range})
        }catch(e) {
            handleError(e)
        }
    }

    const handleNameChange = (e) => {
        setRoleName(e.target.value)
        getRoleInfoValue({roleName: e.target.value, description, range})
        setNameTip("")
    }

    const handleNameBlur = (e) => {
        if(!trim(e.target.value)) {
            setNameTip(intl.get("role.name.placeholder"))
        }else if(!isNormalShortName(trim(e.target.value))) {
            setNameTip(intl.get("role.name.invalid"))
        }
    }

    const handleDescriptionChange = (e) => {
        setDescription(e.target.value)
        getRoleInfoValue({roleName, description: e.target.value, range})
        setDescriptionTip("")
    }

    const handleDescriptionBlur = (e) => {
        if(e.target.value && trim(e.target.value).length > 300) {
            setDescriptionTip(intl.get("role.description.invalid.length"))
        }
    }

    const handleRangeChange = (value) => {
        setRange(value)
        getRoleInfoValue({roleName, description, range: value})
    }

    useEffect(() => {
        if(curRole) {
            initEditInfo(curRole?.id)
        }

        getAllResources({offset: 0, limit: 1000})
    }, [])
    
    return (
        <div className={styles["role-info"]} ref={roleInfoRef}>
            <Form layout="vertical">
                <Form.Item 
                    required
                    className={styles['required-label']} 
                    layout="vertical" 
                    label={intl.get("role.name")}
                    validateStatus={nameTip ? 'error' : ''}
                    help={nameTip ? (<div>{nameTip}</div>) : ''}
                >
                    <Input 
                        placeholder={intl.get("role.name.placeholder")}
                        value={roleName}
                        onChange={handleNameChange}
                        onBlur={handleNameBlur}
                    />
                </Form.Item>
                <Form.Item 
                    className={styles['label']} 
                    layout="vertical" 
                    label={intl.get("role.description")}
                    validateStatus={descriptionTip ? 'error' : ''}
                    help={descriptionTip ? (<div>{descriptionTip}</div>) : ''}
                >
                    <Input 
                        placeholder={intl.get("role.description.placeholder")}
                        value={description}
                        onChange={handleDescriptionChange}
                        onBlur={handleDescriptionBlur}
                    />
                </Form.Item>
                <Form.Item 
                    className={styles['label']} 
                    layout="vertical" 
                    label={intl.get("role.range")}
                >
                    <Select
                        placeholder={"不限"}
                        mode="multiple"
                        allowClear={true}
                        labelInValue={true}
                        optionFilterProp="label"
                        getPopupContainer={() => roleInfoRef.current}
                        options={options}
                        value={range}
                        onChange={handleRangeChange}
                        
                    />
                </Form.Item>
            </Form>
        </div>
    )
}