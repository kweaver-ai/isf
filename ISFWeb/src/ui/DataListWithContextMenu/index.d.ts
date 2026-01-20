declare namespace UI {
    namespace DataListWithContextMenu {

        interface ContextMenuProps extends React.Props<any> {
            open: boolean
            selections: Array<any>
            position: [number, number]
            onClose: () => void
        }

        interface Props extends UI.DataList.Props {
            contextMenu: (props: ContextMenuProps) => React.ComponentElement<UI.PopMenu.Props, any> | null
        }

        interface State {
            /**
             * 选中项数组
            */
            selections: Array<any>

            open: boolean

            contextMenuPosition: [number, number]
        }
    }
}