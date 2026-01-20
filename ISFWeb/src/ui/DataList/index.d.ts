declare namespace UI {
    namespace DataList {
        interface Props extends React.Props<any> {
            /**
             * 选中项数组
            */
            selections?: Array<any>

            /**
             * 允许多选
             */
            multiple?: boolean

            /**
             * 列表className
             */
            className?: string

            /**
             * 选择状态
             */
            selecting?: boolean

            /**
             * 选中项改变
             * @param selections 选中项数组
             * @param multiple 是否是多选
             */
            onSelectionChange?: (selections: Array<any>, multiple: boolean) => any

            /**
             * 双击选中项
             * @param e 右击事件
             * @param data 选中项数据
             */
            onClick?: (e: React.SyntheticEvent<any>, data: any, index?: number) => any

            /**
             * 双击选中项
             * @param e 右击事件
             * @param data 选中项数据
             */
            onDoubleClick?: (e: React.SyntheticEvent<any>, data: any, index?: number) => any

            /**
             * 右击选中项
             * @param e 右击事件
             * @param data 选中项数据
             */
            onContextMenu?: (e: React.SyntheticEvent<any>, data: any, index?: number) => any
        }

        interface State {
            /**
             * 选中项数组
            */
            selections: Array<any>

            /** 列表Item数组 */
            items
        }
    }
}