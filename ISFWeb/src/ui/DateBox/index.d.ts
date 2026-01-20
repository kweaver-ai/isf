declare namespace UI {
    namespace DateBox {
        interface Props extends React.Props<void> {

            /**
             * 初始日期
             */
            value?: Date;

            /**
             * 是否显示日历控件
             */
            active?: boolean;

            /**
             * 控件宽度
             */
            width?: number;

            /**
             * 未选择日期时的默认文字
             */
            placeholder?: string;

            /**
             * 对齐方式
             * @example
             * ```
             * <DateBox dropAlign="left bottom" />
             * ```
             */
            dropAlign?: string;

            /**
             * 日期选择范围
             */
            selectRange?: [Date, Date];

            /**
             * 选择改变时触发
             * @param date 日期对象
             */
            onChange?: (date: Date) => any;

            /**
             * 日期对象显示格式
             * @example
             * ```
             * <DateBox format="yyyy/MM/dd" />
             * ```
             */
            format: string;

            /**
             * 当为空时显示空白或'---'
             */
            shouldShowblankStatus?: boolean;

            /**
             * 选中的日期对象从当天00:00:00开始
             */
            startsFromZero: false;

            /**
             * 日历面板点击事件
             */
            onDatePickerClick: () => void;

            element?: HTMLElement;
        }

        interface State {
            /**
             * 当前选中的日期对象
             */
            value?: Date;

            /**
             * 当前是否显示日历控件
             */
            active?: boolean;
        }

        interface Component extends React.ReactElement<Props> {
        }
    }
}