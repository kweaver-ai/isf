declare namespace UI {
    namespace DataListItem {
        interface Props extends React.Props<any> {
            className?: string

            /**
             * 是否可选
             */
            selectable?: boolean

            /**
             * 是否显示复选框
             */
            checkbox?: boolean

            /**
             * 数据
             */
            data?: any

            /**
             * 选中状态
             */
            selected?: boolean

            /**
             * 是否可展开
             */
            expandable?: boolean

            /**
             * 展开内容
             */
            expandContent?: any

            /**
             * 展开状态
             */
            expanded?: boolean

            /**
             * 选择状态
             */
            selecting?: boolean

            /**
             * 单击
             * @param e 单击事件
             */
            onClick?: (e: React.MouseEvent<HTMLDivElement>) => any

            /**
             * 右击
             * @param e 右击事件
             */
            onContextMenu?: (e: React.MouseEvent<HTMLDivElement>) => any

            /**
             * 双击
             * @param e 双击事件
             */
            onDoubleClick?: (e: React.MouseEvent<HTMLDivElement>) => any

            /**
             * 切换选中状态
             * @param e 切换选中状态
             */
            onToggleSelect?: () => any

            /**
             * 复选框自定义样式
             */
            checkBoxClassName?: string;
        }
    }
}