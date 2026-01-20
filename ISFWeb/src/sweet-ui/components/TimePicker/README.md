### 这个控件叫什么

时间选择器

### 何时使用

当用户需要输入一个时间，可以点击标准输入框，弹出时间面板进行选择。

### 使用示例

#### 1. 基本使用

* 需要导入moment：JavaScript日期处理类库。

* 默认格式 `format: 'HH:mm:ss'`，显示时分秒。`defaultValue`设置默认显示时间。

```jsx
const moment = require('moment');

<TimePicker 
    placeholder={'请选择时间'}
    defaultValue={moment('00:00:00', 'HH:mm:ss')}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```

* 不设置`defaultValue`，当没有选择时间时显示提示文字。

```jsx
<TimePicker 
    placeholder={'请选择时间'}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```

* 可以通过`defaultOpenValue`设置面板打开时选中的时间，默认时间00:00:00（格式：'HH:mm:ss'）。

```jsx
const moment = require('moment');

<TimePicker 
    placeholder={'请选择时间'}
    defaultOpenValue={moment('08:00:00', 'HH:mm:ss')}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```

#### 2.指定时间值`value`（受控组件）

* value 和 onValueChange 需要配合使用。

```jsx
const moment = require('moment');

initialState={value: moment('18:05', 'HH:mm')};

<TimePicker 
    placeholder={'请选择时间'}
    value={state.value}
    format={'HH:mm'}
    onValueChange={({detail: {time, timeString}}) => {setState({value: time}); console.log(timeString)}}
/>

```

```jsx
const moment = require('moment');

initialState={value: moment({hour: 18, minute: 5, second: 15})};

<TimePicker 
    placeholder={'请选择时间'}
    value={state.value}
    onValueChange={({detail: {time, timeString}}) => {setState({value: time}); console.log(timeString)}}
/>

```

#### 3. 选择时分

* 时间面板中的列会随着 `format` 变化，当略去 `format` 中的某部分时，浮层中对应的列也会消失。指定`format: 'HH:mm'`，只显示时分。

```jsx
<TimePicker 
    placeholder={'请选择时间'}
    format={'HH:mm'}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```

#### 4. 禁选部分时间

* 指定 `disabledHours` 禁止选择部分小时

```jsx
<TimePicker 
    placeholder={'请选择时间'}
    disabledHours={() => [0,1,2,3,4]}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```
* 指定 `disabledMinutes` 禁止选择部分分钟

```jsx
<TimePicker 
    placeholder={'请选择时间'}
    disabledMinutes={(hour) => {console.log(hour); return [1,2,3,4]}}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```

#### 5. 步长选项

* 通过`hourStep`、`minuteStep`、`secondStep`分别指定小时、分钟、秒选项的间隔。

```jsx
<TimePicker 
    placeholder={'请选择时间'}
    hourStep={2}
    minuteStep={15}
    secondStep={10}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```

#### 6. 禁用

```jsx
<TimePicker 
    placeholder={'请选择时间'}
    disabled={true}
    onValueChange={({detail: {time, timeString}}) => console.log(timeString)}
/>

```
