import * as React from 'react'
import { ModalDialog2, SweetIcon } from '@/sweet-ui'
import { useForceUpdate } from '@/core/hooks'
import OrgAndGroupPick from '../OrgAndGroupPick/component.view'
import { OrgAndGroupPickProps } from '../OrgAndGroupPick/component.base'
import __ from './locale'
import { noop } from 'lodash';

const { useState } = React

interface Props extends OrgAndGroupPickProps {
    /**
     * 标题
     */
    title: string;

    /**
     * 默认已选
     */
    defaultSelections?: ReadonlyArray<any>;

    element?: HTMLElement;

    /**
     * 弹窗的zIndex
     */
    zIndex?: number;

    /**
     * 确定
     */
    onRequestConfirm: (data: ReadonlyArray<any>) => void;

    /**
     * 取消
     */
    onRequestCancel: () => void;

    /**
     * 是否为知识管理员
     */
    isKnowledger?: boolean;
    /**
      * 是否请求普通用户的组织架构
      */
    isRequestNormal?: boolean;
}

const OrgAndGroupPicker: React.FC<Props> = ({
    title = '',
    defaultSelections = [],
    element = null,
    zIndex = 51,
    onRequestConfirm = noop,
    onRequestCancel = noop,
    ...otherProps
}) => {
    const [selections, setSelections] = useState(defaultSelections)
    const forceUpdate = useForceUpdate() // 强制更新方法（为什么？OrgAndGroupPick内异步判断tab的隐藏/显示，导致了宽度会变，故在OrgAndGroupPick didMount之后需更新Dialog，使Dialog重新定位）

    return (
        <ModalDialog2
            zIndex={zIndex}
            title={title}
            icons={[{
                icon: <SweetIcon role={'sweetui-sweeticon'} name={'x'} size={16} />,
                onClick: onRequestCancel,
            }]}
            element={element}
            buttons={[
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: () => onRequestConfirm(selections),
                    disabled: !selections.length,
                },
                {
                    text: __('取消'),
                    theme: 'regular',
                    onClick: onRequestCancel,
                },
            ]}
        >
            <OrgAndGroupPick
                isShowPadding={false}
                isShowCheckBox={false}
                selections={selections}
                onRequestSelectionsChange={(data) => setSelections(data)}
                onDidMount={forceUpdate}
                {...otherProps}
            />
        </ModalDialog2>
    )
}

export default OrgAndGroupPicker