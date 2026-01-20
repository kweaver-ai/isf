
import { ExpandStatus } from './helper';

declare namespace UI {
    namespace AttrTextField {
        type Props = {
            // 文件属性数据
            attr: any,
        }

        type State = {
            // 文件属性数据状态
            status: ExpandStatus,
        }

        interface Base {
            props: Props;

            state: State;
        }
    }
}