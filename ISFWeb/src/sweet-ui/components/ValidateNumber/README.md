### 这个控件叫什么

数字验证框

### 何时使用

* 需要对数字输入内容进行校验时

### 使用示例

#### 1.输入发生变化时验证

* 输入过程中触发验证。
* 当输入的内容不符合验证规则时，输入框显示为错误样式。当鼠标聚焦或者悬浮到输入框的时候，气泡提示错误详情。当再次输入的内容符合验证规则时候，输入框错误样式消失。

```jsx
initialState = { value: '', validateState: 'normal'};

<ValidateNumber
    value={state.value}
    onValueChange={({detail}) => {
        setState({
            validateState: detail===''?'empty':isNaN(detail)?'invalid':'normal',
            value: detail
        })
    }}
    validateState={state.validateState}
    validateMessages={{
        ['empty']:'输入不允许为空',
        ['invalid']:'请输入合法数字',
  }}
/>
```

#### 2.失焦时验证

* 文本框失去焦点时触发验证。
* 当失去焦点时，输入框输入的内容不符合验证规则，输入框显示为错误样式，当鼠标聚焦或者悬浮到输入框的时候，气泡提示错误详情。下一次失去焦点时检查，若符合验证规则输入框的错误样式消失，不符合，错误样式显示。

```jsx
initialState = { value: '', validateState: 'normal'};

<ValidateNumber
    value={state.value}
    max={99}
    onValueChange={({detail}) => {
    setState({value: detail})}}
    onBlur={(event,value)=> setState({validateState:value===''?'empty':isNaN(value)?'invalid':'normal'})}
    validateState={state.validateState}
    validateMessages={{
        ['empty']:'输入不允许为空',
        ['invalid']:'不允许输入非法字符',
  }}
/>
```

#### 3.外部操作时验证

* 当执行外部操作时候触发验证。
* 如点击确定按钮时，检查输入框中的内容是否符合规范，不符合则输入框显示为错误样式，当鼠标聚焦或者悬浮到输入框的时候，气泡提示错误详情。当文本框输入内容时候，输入框显示为错误样式消失（可控，也可设置为下一次外部操作时改变）。

```jsx
const Button = require('../Button').default;

initialState = { value: '', validateState: '100'};

<div>
    <ValidateNumber
        value={state.value}
        onValueChange={({detail}) => {
            setState({
                validateState:'normal',
                value:detail
            })
        }}
        validateState={state.validateState}
        validateMessages={{
            ['empty']:'输入不允许为空',
            ['invalid']:'不允许输入非法字符',
        }}
    />
    <div style={{ marginTop: '10px'}}>
        <Button onClick={() => {
            setState({
                validateState:state.value===''?'empty':isNaN(state.value)?'invalid':'normal',
            })
        }}>{'确定'}</Button>
    </div>
</div>
```
