declare namespace UI {
    namespace Form {
        interface Props extends React.Props<any> {
            /**
             * 表单提交时触发
             * @param event 事件对象
             */
            onSubmit?: (event: Event) => any;
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface StatelessComponent extends React.FunctionComponent<Props> {
            (props: Props): Element;
            Row: UI.FormRow.Component;
            Label: UI.FormLabel.Component;
            Field: UI.FormField.Component;
        }
    }
}