declare namespace UI {
    namespace Calendar {
        type Props = {
            /**
             * 选中日期触发
             */
            onSelect: (date: Date) => void;

            /**
             * 每周第一天是周几
             */
            firstOfDay: number;

            /**
             * 允许选择的日期范围
             */
            selectRange: Array<Date>;

            /**
             * 选中的日期对象从当天00:00:00开始, 默认从当天23:59:59开启
             */
            startsFromZero?: boolean;

            /**
             * 当前年份，可以是日期对象、字符串或时间戳
             */
            year?: Date | string | number;

            /**
             * 当前月份，可以是日期对象、字符串或时间戳
             */
            month?: Date | string | number;

            /**
             * 高亮的日期，可以是日期对象、字符串或时间戳
             */
            date?: Date | string | number;

            /**
             * 当前选中的时间，支持选择时间时有效
             */
            time?: string;

            /**
             * 点击时间处理函数
             */
            onClickTime?: () => void;

        }

        type State = {
            /**
             * 高亮的日期，可以是日期对象、字符串或时间戳
             */
            selected: Date | string | number;

            /**
             * 一个月中的所有日期对象
             */
            weeks: Array<Array<Date>>;
        }

        interface Base {
            props: Props;

            state: State;
        }
    }
}