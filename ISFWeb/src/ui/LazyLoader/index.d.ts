declare namespace UI {
    namespace LazyLoader {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;

            /**
             * 每页条数
             */
            limit?: number;

            /**
             * 滚动到某个位置
             */
            scroll?: number;

            /**
             * 触发滚动翻页的位置，使用小于1的数表示百分比位置
             */
            trigger?: number;

            /**
             * 滚动时触发
             * @param scroll 滚动位置
             */
            onScroll?: (scroll: number) => any;

            /**
             * 滚动翻页
             * @param page 页码
             * @param limit 每页条数
             */
            onChange?: (page: number, limit: number) => any;
        }

        interface State {
            /**
             * 当前所在页
             */
            page: number;
        }
    }
}