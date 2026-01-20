import React from 'react';

export default class StackBarBase extends React.Component<UI.StackBar.Props, any> implements UI.StackBar.Base {
    static defaultProps = {
        width: '100%',
    }
}