declare namespace UI {
    namespace Toast {
        interface Props extends UI.UIIcon.Props {

            closable?: boolean

            onClose?: () => void

        }
    }
}