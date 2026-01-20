declare namespace UI {
    namespace SelectMenu {
        interface Props extends UI.PopMenu.Props {
            value?: any

            defaultValue?: any;

            onChange?: (value: any) => any
        }
    }
}