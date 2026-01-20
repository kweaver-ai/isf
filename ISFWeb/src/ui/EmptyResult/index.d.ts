declare namespace UI {
    namespace EmptyResult {
        interface Props extends React.Props<any> {
            /**
             * 图片
             */
            picture?: string;

            /**
             * 字体图标Unicode码
             */
            code?: string;

            /**
             * 结果为空时的提示内容
             */
            details?: string;

            /**
             * 图片大小（默认64*64）|字体图标大小
             */
            size?: number | string,

            /**
             * 提示文字字体家族名称（默认AnyShare）
             */
            font?: string;

            /**
             * 提示文字字体大小（默认13px）
             */
            fontSize?: number | string;

            /**
            * 提示文字颜色（默认#868686）
            */
            color?: string;
        }
    }
}