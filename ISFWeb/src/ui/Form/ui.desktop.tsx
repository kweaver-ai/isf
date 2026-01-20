import React from 'react';
import { noop } from 'lodash';
import classnames from 'classnames';
import FormRow from '../Form.Row/ui.desktop';
import FormLabel from '../Form.Label/ui.desktop';
import FormField from '../Form.Field/ui.desktop';
import styles from './styles.desktop';

const Form: UI.Form.StatelessComponent = function Form({ onSubmit = noop, children, className, ...props }) {
    return (
        <form
            {...props}
            className={classnames(styles['form'], className)}
            onSubmit={(event) => onSubmit(event) && event.preventDefault()}
        >
            {
                children
            }
        </form>
    )
} as UI.Form.StatelessComponent

Form.Row = FormRow;

Form.Label = FormLabel;

Form.Field = FormField;

export default Form