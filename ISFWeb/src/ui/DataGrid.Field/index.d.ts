declare namespace UI {
    namespace DataGridField {

        /**
         * field配置
         */
        interface Field extends React.Props<any> {
            /**
             * 字段键名
             */
            field: string;

            /**
             * 样式
             */
            className?: string;

            /**
             * 字段名称
             */
            label: string;

            /**
             * 水平显示
             */
            align?: 'left'|'center'|'right';

            /**
             * 单元格宽度
             */
            width?: string | number;

            /**
             * 格式化函数
             * @prop value 字段值
             * @prop record 整条数据
             * @prop status 该条数据的状态
             */
            formatter?: (value: any, record: Object, status: UI.DataGrid.RecordStatus) => any;
        }
    }
}