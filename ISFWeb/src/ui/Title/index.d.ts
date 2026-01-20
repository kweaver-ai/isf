declare namespace UI {
    namespace Title {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;

            inline?: boolean

            /** title内容 */
            content: any

            /** 延迟显示 */
            timeout?: number

            /**
             * className
             */
            className?: string;
        }
    }
}