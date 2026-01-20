declare namespace UI {
    namespace ProgressCircle {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;
            /**
             * 是否显示Mask
             */
            showMask?: boolean;

            /**
             * loading图片下面的文字详情
             */
            detail?: string;

            /**
             * loading图标的选择，深色背景的图标 or 浅色背景的图标
             */
            theme?: 'dark' | 'light';

            /**
             * 是否需要平铺整个页面
             */
            fixedPositioned?: boolean;
        }
    }
}