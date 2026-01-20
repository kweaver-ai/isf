### 这个控件叫什么

带校验功能的包装器，一般用来包装支持输入的组件

### 何时使用

* 需要对输入内容进行校验时

* 常用于验证单项内容；如果需要验证多项，推荐使用表单

* 限制项：子组件`props`必须有`value`项，当验证触发方式为`onChange`时还需要有`onValueChange`接口；当验证触发方式为`onBlur`时,需要有`onBlur`接口

### 使用示例

#### 1.输入发生变化时验证（默认）

* `rules`数组项先后顺序决定验证规则优先级，排在前面的验证规则优先验证

* 输入时校验结果为错误，直接在后方显示错误提示；删除不合法内容提示消失，否则一直显示错误提示

```jsx
const TextBox = require('../TextBox').default;

initialState = { value: '', validateResult: true };
<Validator 
    rules={[
        {
            message: '输入不允许为空',
            required: true
        },
        {
            message: '不允许输入非法字符',
            validator: (value) => !/[#\\/:*?"<>|]/.test(value)
        },
        {
            message: '不允许输入字符？',
            validator: (value) => !/[?]/.test(value)
        },
    ]}
    afterValidate={({detail}) => setState({validateResult: detail})}
>

    <TextBox
        value={state.value}
        onValueChange={({detail}) => setState({value: detail})}
        onBlur={console.log}
        status={state.validateResult ? 'normal' : 'error'}
    />
</Validator>

```

#### 2.失焦时验证

* 失焦时校验结果为错误，直接在后方显示错误内容。更改内容，错误提示消失，再次失焦时执行下一次校验

```jsx
const TextArea = require('../TextArea').default;
initialState = { value: '', validateResult: true };

<Validator
    validateTrigger={'onBlur'}
    rules={[
        {
            message: '输入不允许为空',
            required: true
        },
        {
            message: '不允许输入非法字符',
            validator: (value) => !/[#\\/:*?"<>|]/.test(value)
        },
    ]}
    afterValidate={({detail}) => setState({validateResult: detail})}
    placement={'rightTop'}
>
    <TextArea 
        value={state.value} 
        onValueChange={(event) => {setState({value: event.detail})} }
        placeholder={'请输入内容'}
        onBlur={() => console.log('textarea onblur')}
        onPressEnter={console.log}
        status={state.validateResult ? 'normal' : 'error'}
/>

</Validator>


```


