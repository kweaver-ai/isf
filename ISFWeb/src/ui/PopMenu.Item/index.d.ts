declare namespace UI {
    namespace PopMenuItem {
        interface Props {

            className?: string

            /**
             * 菜单图标
             */
            icon?: React.ReactElement<any> | string
            /**
             * 菜单标题
             */
            label: string

            /**
             * 挂载到DOM时触发
             */
            onDOMNodeMount?: (node: HTMLElement) => any
        }
    }
}