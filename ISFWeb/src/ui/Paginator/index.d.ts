declare namespace UI {
    namespace Paginator {
        interface Props extends React.Props<any> {
            /**
             * 当前页码
             */
            page: number;

            /**
             * 分页大小
             */
            limit: number;

            /**
             * 总页数
             */
            total: number;

            /**
             * 翻页时触发
             */
            onChange?: (page: number, limit: number) => any;
        }

        interface State {
            /**
             * 当前页码
             */
            page: number;

            /**
             * 总页数
             */
            pages: number;

            /**
             * 页码框输入值
             */
            pageInput: number | '';
        }

    }
}