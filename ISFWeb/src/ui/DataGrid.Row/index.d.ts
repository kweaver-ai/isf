declare namespace UI {
    namespace DataGridRow {
        interface Props extends React.Props<void> {
            // 是否选中
            selected?: boolean;

            // 是否显示斑马线样式
            strap?: boolean;

            // 是否显示复选框
            checkbox?: boolean;

            /**
             * 单击行触发
             */
            onClick: (e: React.SyntheticEvent<HTMLTableRowElement>) => any;

            /**
             * 双击行触发
             */
            onDoubleClick: (e: React.SyntheticEvent<HTMLTableRowElement>) => any;

            /**
             * 点击复选框触发
             */
            onCheckChange: ({ multi }: { multi: boolean }) => any;
        }
    }
}