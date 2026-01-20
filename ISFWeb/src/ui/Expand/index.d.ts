declare namespace UI {
    namespace Expand {
        interface Props extends React.Props<any> {
            /**
             * 是否展开
             */
            open: boolean
        }

        interface State {
            /**
             * 子元素的负marginTop值， 通过负marginTop实现折叠效果
             */
            marginTop: number
            /**
             * 组件是否mount
             */
            loaded: boolean
            /**
             * 动画效果
             */
            animation: boolean
        }
    }
}