declare namespace UI {
    namespace InlineButton {
        interface Props {
            /**
             * 角色
             */
            role?: string;

            /**
             * 按钮尺寸
             */
            size?: number;

            /**
             * 字体图标代码
             */
            code: string;

            /**
             * 字体图标对应png图片资源
             */
            fallback?: string;

            /**
             * 按钮悬浮标题
             */
            title?: string;

            /**
             * 按钮图标尺寸
             */
            iconSize?: number;

            /**
             * 是否禁用
             */
            disabled?: boolean;

            /**
             * className
             */
            className?: any;

            /**
             * type
             */
            type?: string;

            /**
             * 点击事件
             * @param event 鼠标点击事件对象
             */
            onClick?: (event: React.MouseEvent<HTMLDivElement>) => any;

            /**
             * 挂载元素
             */
            element?: any;
        }
    }
}