declare namespace UI {
    namespace Gallery {
        type Props = {
            /**
             * 子节点
             */
            children: any | Array<any>;

            /**
             * 当期分组
             */
            groupIndex: number;

            /**
             * 分组大小
             */
            groupSize?: number;
        }

        type State = {
            /**
             * 当前分组
             */
            currentGroupIndex: number;
        }

        interface Base {
            props: Props;

            state: State;
        }
    }
}