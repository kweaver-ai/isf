import * as React from 'react'
import classnames from 'classnames'
import { getDepParentPathById } from '@/core/thrift/sharemgnt/sharemgnt'
import { Title } from '@/ui/ui.desktop'
import { SpecialDep } from '../../helper'
import styles from './styles.view';
import __ from './locale'

interface PathTitleProsp {
    /**
     * 用户信息
     */
    record: Core.ShareMgnt.ncTUsrmGetUserInfo;

    /**
     * 选中的部门
     */
    selectedDep: Core.ShareMgnt.ncTDepartmentInfo;

    /**
     * 新增路径
     */
    onRequestUpdatePath: (depPaths: string[]) => any;
}

export default function PathTitle({ record, selectedDep, onRequestUpdatePath }: PathTitleProsp) {
    const { user: { departmentIds = [], departmentNames = [], status = 0 } = {}, depPath } = record

    const [content, setContent] = React.useState<string[]>([''])

    /**
     * 鼠标悬浮
     */
    const handleMouseEnter = async (): Promise<void> => {
        const { id } = selectedDep

        let paths: string[] = []

        if (id !== SpecialDep.Unassigned) {
            if (depPath) {
                paths = depPath
            } else {
                try {
                    const data = await getDepParentPathById(departmentIds)
                    const parentPaths = data.map(({ parentPath }) => parentPath)

                    departmentIds.map((depId, index) => {
                        if (depId === SpecialDep.Unassigned) {
                            paths = [...paths, __('未分配组')]
                        } else {
                            paths = [...paths, parentPaths[index] ? `${departmentNames[index]}-${parentPaths[index]}` : departmentNames[index]]
                        }
                    })
                } catch (ex) { }
            }
        } else {
            paths = [__('未分配组')]
        }

        onRequestUpdatePath(paths)
        setContent(paths)
    }

    return (
        <div onMouseEnter={handleMouseEnter}>
            <Title
                content={
                    <div className={styles['title']}>
                        {content.map((path) => <div className={styles['item']} key={path}>{path}</div>)}
                    </div>
                }
                role={'ui-title'}
            >
                <span
                    className={classnames(styles['title-text'], { [styles['gray-text']]: status })}
                >
                    {departmentIds.includes(SpecialDep.Unassigned) ? __('未分配组') : departmentNames.length ? departmentNames.join(',') : '---'}
                </span>
            </Title>
        </div>
    )
}