declare namespace Console {
    namespace DisplayManager {
        interface Props extends React.Props<void> {
            /**
             * 部门Id
             */
            departmentId: string;

            /**
             * 部门名
             */
            departmentName: string;

            /**
             * 当前用户id
             */
            userid: string;

            /**
             * 是否有权限修改
             */
            hasPermission: boolean;

            /**
             * 结束设置管理员操作
             */
            onComplete: () => void;
        }

        interface State {
            /**
             * 当前部门的组织管理员
             */
            manager: ReadonlyArray<string>;

            /**
             * 时候编辑
             */
            // TODO isEditing
            isEdited: boolean;
        }
    }
}