#### 何时使用

以标签的形式展示用户选择或者输入的内容时，并包含验证格式

#### 基本使用

```jsx
const Button = require('../Button').default;
initialState =  { value: ['标签1','标签2','标签3'] ,validateState:'noraml'};

<div>
    <ValidateComboArea
        value={state.value}
        placeholder={'请选择一个标签'}
        onChange={(value) => setState(()=>({value}))}
        validateState={state.validateState}
        validateMessages={{
            ['empty']:'输入不允许为空',
        }}
    />

    <div style={{ marginTop: '10px'}}>
        <Button
            width={120} 
            onClick={() => setState(() => ({
                        value: [
                            ...state.value,
                            `标签${state.value.length + 1}`
                        ],
                        validateState:'normal',
            }))}>
            {'添加新标签'}
        </Button>
        <Button 
            onClick={() => {
                    if(state.value.length===0){
                        setState({validateState:'empty'})
                    }else{
                        setState({validateState:'normal'})
                        alert('所有项检查合法，要查看错误状态验证，请删除文本框中所有的标签')
                    }
            }}
        >
        {'确定'}
        </Button>
    </div>
</div>
```
