import React from 'react';
import Menu from '../Menu/ui.desktop';
import DropBox from '../DropBox/ui.desktop';
import SelectOption from '../Select.Option/ui.desktop';
import SelectBase from './ui.base';

export default class Select extends SelectBase {
    static Option = SelectOption;

    render() {
        const { role, width: menuWidth = 220, maxHeight: menuMaxHeight = 160 } = this.props.menu;

        return (
            <DropBox
                role={role}
                className={this.props.className}
                value={this.state.value}
                active={this.state.active}
                disabled={this.props.disabled}
                formatter={() => this.state.text}
                width={this.props.width}
            >
                <Menu width={menuWidth} maxHeight={menuMaxHeight} onMouseDown={this.props.onMouseDown}>
                    {
                        React.Children.map(this.props.children, (option) => React.cloneElement(option, {
                            selected: this.state.value === option.props.value,
                            onSelect: this.onOptionSelected.bind(this),
                        }))
                    }
                </Menu>
            </DropBox>
        )
    }
}
