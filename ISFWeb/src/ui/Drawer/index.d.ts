declare namespace UI {
    namespace Drawer {

        type DrawerPosition = 'top' | 'right' | 'bottom' | 'left'

        interface Props extends React.Props<void> {
            /**
             * 遮罩层
             */
            mask?: boolean

            /**
             * 展开
             */
            open?: boolean

            /**
             * 位置
             */
            position?: DrawerPosition

            className?: any

            /**
             * 点击遮罩层
             */
            onClickMask?: (e: React.SyntheticEvent<any>) => void
        }
    }
}