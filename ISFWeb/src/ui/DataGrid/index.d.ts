declare namespace UI {
    namespace DataGrid {

        /**
         * 数据状态
         */
        interface RecordStatus {
            /**
             * 复选框状态
             * 如果返回false，则不出现复选框，否则返回复选框状态
             */
            checkbox?: CheckBoxStatus & boolean;

            /**
             * 选中
             * @todo 暂不生效，使用this.props.getDefaultSelection获取选中项
             */
            selected?: boolean;

            /**
             * 禁用
             */
            disabled?: boolean;
        }

        /**
         * 复选框状态
         */
        interface CheckBoxStatus {
            disabled?: boolean;
        }

        /**
         * 翻页
         */
        interface Paginator {
            // 初始页
            page: number;

            // 每页条数
            limit: number;

            // 总条数
            total: number;
        }

        interface LazyLoad {
            // 每页条数
            limit?: number;
        }

        interface Props extends React.Props<any> {

            /**
             * 数据
             */
            data: Array<Object>;

            /**
             * 表格高度
             */
            height?: string | number;

            className?: string;

            rowHoverClass?: string;

            /**
             * 双击数据行触发
             */
            onDblClickRow?: Function;

            /**
             * 点击行
             */
            onClickRow?: Function;

            /**
             * 获取默认选中项目
             */
            getDefaultSelection?: (nextData: Array<any>, prevData: Array<any>) => Array<any>;

            /**
             * 获取数据状态
             * @param record 数据对象
             * @param index 下标
             */
            getRecordStatus?: (record: Object, index: number) => RecordStatus;

            /**
             * 选中项改变时触发
             */
            onSelectionChange?: (selection: Object | Array<Object>) => any;

            /**
             * 翻页时触发
             * @param page 页码
             * @param limit 每页条数
             */
            onPageChange?: (page: number, limit: number) => any;

            /**
             * 是否允许选中 | 是否允许多选 && 是否必须要选中一项
             */
            select?: boolean | { multi: boolean, required: boolean };

            /**
             * 是否显示表头
             */
            headless?: boolean;

            /**
             * 表头高度
             */
            headHeight?: string | number;

            /**
             * 自动定位到某一条数据
             * @return 返回该条数据的下标
             */
            locator?: (data: Array<any>) => number;

            /**
             * 允许翻页
             */
            paginator?: Paginator & boolean;

            /**
             * 滚动加载
             */
            lazyLoad?: LazyLoad & boolean;

            /**
             * 显示奇偶交错的斑马线样式
             */
            strap?: boolean;

            /**
             * 获取React repeat时的key属性，用于优化大数据更新性能
             */
            getKey?(data): string | number;

            /**
             * 列表为空时显示的内容
             */
            EmtpyComponent: React.ReactNode;
        }

        interface State {
            /**
             * 当前所在位置
             */
            scroll?: number;

            /**
             * 当前选中项
             */
            selection?: Object | Array<Object> | null;

            /**
             * 允许翻页
             */
            paginator?: Paginator | false;

            /**
             * 内容是否溢出
             */
            isOverflow: boolean;
        }
    }
}