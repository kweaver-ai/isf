declare namespace UI {
    namespace List {
        interface Props extends React.Props<any> {
            /**
             * 从外部传入的数据
             */
            data: Array<any>;

            /**
             * 鼠标点击触发
             */
            onMouseDown: (event: MouseEvent) => any;

            /**
             * 数据格式模板
             */
            template: any;

            /**
             * 当前选中项下标
             */
            selectIndex: number;

            /**
             * 选中项改变时触发
             */
            onSelectionChange: (index: number) => any;

            /**
             * list可视区域高度
             */
            viewHeight?: number;
        }

        interface State {
            /**
             * 当前选中项下标
             */
            selectIndex: number;
        }
    }
}