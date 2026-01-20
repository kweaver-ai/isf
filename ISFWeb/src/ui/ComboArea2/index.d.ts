declare namespace UI {
    namespace ComboArea2 {
        interface Props extends React.Props<any> {
            /**
             * className
             */
            className?: string;

            /**
             * 宽度，含padding和border
             */
            width?: number | string;

            /**
             * 高度，含padding和border
             */
            height?: number | string;

            /**
             * 最小高度
             */
            minHeight?: number;

            /**
             * 最大允许高度
             */
            maxHeight?: number;

            /**
             * 只读
             */
            readOnly?: boolean;

            /**
             * 是否不可编辑
             * 不可编辑状态下仍然可以删除已有项
             */
            uneditable?: boolean;

            /**
             * 禁用
             */
            disabled?: boolean;

            /**
             * 占位文本
             */
            placeholder?: string;

            /**
             * 自动创建Chip的分割字符
             */
            spliter?: Array<string>;

            /**
             * 初始数据
             */
            value?: Array<any>

            /**
             * 格式化函数
             */
            formatter?: (o: any) => string;

            /**
             * 验证是否可以自动创建Chip
             * @params input 输入的值
             * @params data 已生成的数据
             * @return 返回验证是否通过
             */
            validator?: (input: string, data: Array<any>) => boolean;

            /**
             * 数据改变时触发
             */
            onChange?: (data: Array<any>) => any;
        }

        interface State {
            placeholder: string;
        }
    }
}