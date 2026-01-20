import React from 'react';
import { formatTime } from '@/util/formatters';
import { DropBox } from '@/sweet-ui';
import DatePicker from '../DatePicker/ui.desktop';
import DateBoxBase from './ui.base';

export default class DateBox extends DateBoxBase {
    render() {
        return (
            <DropBox
                iconLabel={'date'}
                dropAlign={this.props.dropAlign}
                value={this.state.value}
                formatter={(date) => this.props.shouldShowblankStatus ? this.props.placeholder : formatTime(date, this.props.format)}
                active={this.state.active}
                width={this.props.width}
                onActive={(active) => { this.props.onActive(active) }}
                disabled={this.props.disabled}
                onBeforePopupClose={this.handleBeforePopupClose}
                element={this.props.element}
            >
                {({ close, open }) =>
                    this.props.disabled ?
                        null
                        :
                        <DatePicker
                            value={this.state.value}
                            selectRange={this.props.selectRange}
                            onChange={(value) => this.select(value, close)}
                            startsFromZero={this.props.startsFromZero}
                            onDatePickerClick={() => this.props.onDatePickerClick(open)}
                        />
                }

            </DropBox>
        )
    }
}