declare namespace Components {
    namespace SetUsersFreezeStatus {
        namespace ShowSelected {
            interface Props extends React.Props<void> {

                // 选中的用户
                selecteds: Array<any>

                // 确定事件
                onDelete: (value) => void;
            }
        }
    }

}