import { NetInfo as VisitorNetInfo } from '../VisitorNetBind/helper'
import { ListTipStatus } from '../../ListTipComponent/helper';
declare namespace Console {
    namespace BindingQuery {
        interface Props extends React.Props<void> {
            /**
             * 用户id
             */
            userid: string;
        }

        interface State {
            /**
             * 按什么类型查询
             */
            searchField: number;

            /**
             * 访问者网段信息
             */
            visitorNetInfos: ReadonlyArray<VisitorNetInfo>;

            /**
             * 设备绑定信息
             */
            deviceInfos: ReadonlyArray<{
                /**
                 * 设备识别码
                 */
                udid: string;

                /**
                 * 设备类型
                 */
                osType: number;

                /**
                 * 绑定状态
                 */
                bindFlag: boolean;
            }>;

            /**
             * 绑定该设备识别码的用户
             */
            deviceIdUsers: ReadonlyArray<string>;

            /**
             * 搜索结果
             */
            searchResults: ReadonlyArray<any>;

            /**
             * 搜索关键字
             */
            searchKey: string;

            /**
             * 列表提示状态
             */
            listTipStatus: ListTipStatus;

            /**
             * 设备列表提示状态
             */
            deviceListStatus: ListTipStatus;
        }
    }
}