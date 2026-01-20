declare namespace UI {
    namespace Divider {
        interface Props extends React.Props<void> {
            /**
             * 分割线颜色
             */
            color?: string;

            /**
             * 分割线缩进
             */
            inset?: number | string | Array<number | string>
        }
    }
}