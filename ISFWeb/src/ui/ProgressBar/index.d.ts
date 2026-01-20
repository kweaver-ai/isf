declare namespace UI {
    namespace ProgressBar {
        interface Props extends React.Props<any> {
            /**
             * 值
             */
            value: string | number;

            /**
             * 进度框宽度
             */
            width?: string | number;

            /**
             * 进度框高度
             */
            height?: string | number;

            /**
             * 进度框边框
             */
            border?: string;

            /**
             * 进度框背景色
             */
            containerBackground?: string;

            /**
             * 进度条背景色
             */
            progressBackground?: string;

            /**
             * 进度提示位置
             */
            textAlign?: 'left' | 'center' | 'right';

            /**
             * 渲染值（当不需要渲染进度提示文字的时候传入renderValue={noop}覆盖）
             */
            renderValue?: (value: number) => string | number;
        }
    }
}
