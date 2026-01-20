import * as React from 'react'
import { noop } from 'lodash'
import { Select } from '@/sweet-ui'
import { Title, Text } from '@/ui/ui.desktop'
import { getDepParentPathById } from '@/core/thrift/sharemgnt/sharemgnt'
import { Props, DepInfo } from './type'
import styles from './styles.view';
import __ from './locale'

const { useState, useEffect, useCallback } = React

const DepartmentsOfUserSelector: React.FC<Props> = ({
    userInfo,
    width = 200,
    dep,
    onSelectionChange = noop,
}) => {
    // 所选部门
    const [selectedDep, setSelectedDep] = useState<DepInfo>({ id: '', name: '' })
    // 用户所属部门列表
    const [depList, setDepList] = useState<ReadonlyArray<DepInfo>>([])

    // 获取用户所属部门列表
    const getDepList = useCallback(async () => {
        let depList: ReadonlyArray<DepInfo> = []

        try {
            const { user: { departmentIds, departmentNames } } = userInfo

            // 判断是否只属于一个部门
            if (departmentIds.length === 1) {
                const [{ parentPath }] = await getDepParentPathById(departmentIds)
                const path = `${parentPath ? parentPath + '/' : ''}${departmentNames[0]}`

                depList = [{ id: departmentIds[0], name: departmentNames[0], path, is_root: !parentPath }]
            } else {
                // 用户所属部门中是否有{dep}部门
                let isHaveDep = false
                const parentPathList = await getDepParentPathById([...departmentIds, dep.id])
                // 当前(进入)部门的路径
                const depPath = `${parentPathList[parentPathList.length - 1].parentPath ? parentPathList[parentPathList.length - 1].parentPath + '/' : ''}${dep.name}`

                depList = parentPathList.slice(0, parentPathList.length - 1).reduce((prev, { departmentId, parentPath }, index) => {
                    const path = `${parentPath ? parentPath + '/' : ''}${departmentNames[index]}`
                    const is_root = !parentPath

                    if (!isHaveDep) {
                        isHaveDep = departmentId === dep.id

                        if (isHaveDep || path.startsWith(depPath)) {
                            return [{
                                id: departmentId,
                                name: departmentNames[index],
                                path,
                                is_root,
                            }, ...prev]
                        }
                    }

                    return [...prev, {
                        id: departmentId,
                        name: departmentNames[index],
                        path,
                        is_root,
                    }]
                }, [])
            }
        } catch {
            depList = []
        }
        
        // 统一和父组件的选中项
        onSelectionChange(depList[0]);
        setSelectedDep(depList[0])
        setDepList(depList)

    }, [userInfo, dep, onSelectionChange])

    useEffect(() => {
        getDepList()
    }, [getDepList])

    return (
        <>
            {
                depList.length === 1
                    ? <div className={styles['dep-name']}>
                        {__('“')}
                        <span className={styles['dep-name-content']}>
                            <Text>{depList[0].name}</Text>
                        </span>
                        {__('”')}
                    </div >
                    : depList.length > 1
                        ? (
                            <Select
                                value={selectedDep.id}
                                onChange={({ detail }) => { 
                                    const cur = depList.find((cur) => cur.id === detail)
                                    if(cur) {
                                        onSelectionChange(cur);
                                        setSelectedDep(cur)
                                    }
                                }}
                                selectorWidth={width}
                            >
                                {
                                    depList.map((item) => {
                                        const { id, name, path } = item

                                        return (
                                            <Select.Option
                                                key={id}
                                                value={id}
                                                className={styles['select-option']}
                                            >
                                                <Title role={'ui-title'} content={path}>
                                                    <div className={styles['select-ellipsis']}>{name}</div>
                                                </Title>
                                            </Select.Option>
                                        )
                                    })
                                }
                            </Select>
                        )
                        : null
            }
        </>
    )
}

export default DepartmentsOfUserSelector;