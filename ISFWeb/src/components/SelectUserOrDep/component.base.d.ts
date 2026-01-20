declare namespace Components {
    namespace SelectUserOrDepBase {
        interface Props {

            // 登录用的id
            userid: string;

            // 需要加载的类型   ORGANIZATION=0：组织 DEPARTMENT=1： 部门 USER=2： 用户
            selectedType?: Array<number>

            // 选择事件
            onSelected: (value: any) => void;
        }

        interface State {

        }
    }
}